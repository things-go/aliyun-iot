// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package aiot imp aliyun dm
package aiot

import (
	"encoding/json"
	"io"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/things-go/aliyun-iot/infra"
	"github.com/things-go/aliyun-iot/logger"
)

// 缓存默认值
const (
	DefaultCacheExpiration      = time.Second * 10
	DefaultCacheCleanupInterval = time.Second * 30
)

// DefaultVersion 平台通信版本
const DefaultVersion = "1.0"

// Mode 工作模式
type Mode byte

// 当前的工作模式
const (
	ModeMQTT Mode = iota
	ModeCOAP
	ModeHTTP
)

// ProcDownStream 处理下行数据
type ProcDownStream func(c *Client, rawURI string, payload []byte) error

// Conn conn接口
type Conn interface {
	// Publish will publish a Message with the specified QoS and content
	Publish(topic string, qos byte, payload interface{}) error
	Subscribe(topic string, callback ProcDownStream) error
	UnSubscribe(topic ...string) error
	io.Closer
}

// Request 请求
type Request struct {
	ID      uint        `json:"id,string"`
	Version string      `json:"version"`
	Params  interface{} `json:"params"`
	Method  string      `json:"method"`
}

// Response 应答
type Response struct {
	ID      uint        `json:"id,string"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// ResponseRawData 应答, data域为 json.RawMessage
type ResponseRawData struct {
	ID      uint            `json:"id,string"`
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message,omitempty"`
}

// Client 客户端
type Client struct {
	requestID uint32
	tetrad    infra.MetaTriad

	cacheExpiration      time.Duration
	cacheCleanupInterval time.Duration

	mode    Mode
	version string
	// 选项功能
	isGateway   bool
	hasDiag     bool
	hasNTP      bool
	hasRawModel bool
	hasDesired  bool
	hasExtRRPC  bool
	hasOTA      bool

	*DevMgr
	msgCache *cache.Cache
	Conn
	cb   Callback
	gwCb GwCallback
	Log  logger.Logger
}

// New 创建一个物管理客户端
func New(triad infra.MetaTriad, conn Conn, opts ...Option) *Client {
	c := &Client{
		requestID: uint32(time.Now().Nanosecond()),
		tetrad:    triad,

		mode:    ModeMQTT,
		version: DefaultVersion,

		cacheExpiration:      DefaultCacheExpiration,
		cacheCleanupInterval: DefaultCacheCleanupInterval,

		DevMgr: NewDevMgr(triad),
		Conn:   conn,
		cb:     NopCb{},
		gwCb:   NopGwCb{},
		Log:    logger.NewDiscard(),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.mode != ModeHTTP {
		c.msgCache = cache.New(c.cacheExpiration, c.cacheCleanupInterval)
	}
	return c
}

// Connect 将订阅所有相关主题,主题有config配置
func (sf *Client) Connect() error {
	if sf.mode != ModeMQTT {
		return nil
	}
	return sf.SubscribeAllTopic(sf.tetrad.ProductKey, sf.tetrad.DeviceName, false)
}

// AddSubDevice 增加一个一个子设备
func (sf *Client) AddSubDevice(meta infra.MetaTriad) error {
	if sf.isGateway {
		return sf.Add(meta)
	}
	return ErrNotSupportFeature
}

// SubDeviceConnect 子设备连接注册并添加到网关拓扑关系
// 子设备上线流程: (确保网关已接入物联网平台)
//      1. 子设备发起动态注册，返回成功注册的子设备的设备证书(当平台使能动态注册子设备时)
//      2. 子设备身份注册后,通过网关向平台上报网关与子设备的拓扑关系
//      3. 子设备进行上线(此时平台会校验子设备的身份和与网关的拓扑关系。所有校验通过，才会建立并绑定子设备逻辑通道至网关物理通道上)
//      4. 子设备与物联网平台的数据上下行通信与直连设备的通信协议一致，协议上不需要露出网关信息
//      5. 删除拓扑关系后,子设备不能再通过网关上线
func (sf *Client) SubDeviceConnect(pk, dn string, cleanSession bool, timeout time.Duration) error {
	node, err := sf.SearchAvail(pk, dn)
	if err != nil {
		return err
	}
	if node.Status() < DevStatusRegistered || node.DeviceSecret() == "" { // 需要注册
		// 子设备注册
		if _, err := sf.LinkThingSubRegister(pk, dn, timeout); err != nil {
			return err
		}
	}
	// 子设备添加到拓扑
	err = sf.LinkThingTopoAdd(pk, dn, timeout)
	if err != nil {
		return err
	}
	// 上线
	err = sf.LinkExtCombineLogin(CombinePair{pk, dn, cleanSession}, timeout)
	if err != nil {
		return err
	}
	// 订阅
	err = sf.SubscribeAllTopic(pk, dn, true)
	if err != nil {
		return err
	}
	sf.SetDeviceStatus(pk, dn, DevStatusOnline) // nolint: errcheck
	return nil
}
