// Package dm imp aliyun dm
package dm

import (
	"encoding/json"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/thinkgos/go-core-package/lib/logger"

	"github.com/thinkgos/aliyun-iot/infra"
)

// 缓存默认值
const (
	DefaultCacheExpiration      = time.Second * 10
	DefaultCacheCleanupInterval = time.Second * 30
)

// 当前的工作方式
const (
	WorkOnMQTT = iota
	WorkOnCOAP
	WorkOnHTTP
)

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

	workOnWho byte

	// 选项功能
	isGateway   bool
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

	log logger.Logger
}

// Version 平台通信版本
var Version = "1.0"

// New 创建一个物管理客户端
func New(meta infra.MetaTriad, conn Conn, opts ...Option) *Client {
	c := &Client{
		requestID: rand.Uint32(),
		tetrad:    meta,

		workOnWho: WorkOnMQTT,

		cacheExpiration:      DefaultCacheExpiration,
		cacheCleanupInterval: DefaultCacheCleanupInterval,

		DevMgr: NewDevMgr(meta),
		Conn:   conn,
		cb:     NopCb{},
		gwCb:   NopGwCb{},
		log:    logger.NewDiscard(),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.workOnWho != WorkOnHTTP {
		c.msgCache = cache.New(c.cacheExpiration, c.cacheCleanupInterval)
	}

	return c
}

// NewSubDevice 创建一个子设备
func (sf *Client) NewSubDevice(meta infra.MetaTetrad) error {
	if sf.isGateway {
		return sf.Create(meta)
	}
	return ErrNotSupportFeature
}

// RequestID 获得下一个requestID,协程安全
func (sf *Client) RequestID() uint {
	return uint(atomic.AddUint32(&sf.requestID, 1))
}

// SendRequest 发送请求,API内部已实现json序列化
// uriService 唯一定位服务器或(topic)
// requestID: 请求ID
// method: 方法
// params: 消息体Request的params
func (sf *Client) SendRequest(uriService string, requestID uint, method string, params interface{}) error {
	out, err := json.Marshal(&Request{requestID, Version, params, method})
	if err != nil {
		return err
	}
	return sf.Publish(uriService, 1, out)
}

// SendResponse 发送回复
// uriService 唯一定位服务器或(topic)
// responseID: 回复ID
// code: 回复code
// Data: 数据域
// API内部已实现json序列化
func (sf *Client) SendResponse(uriService string, responseID uint, code int, data interface{}) error {
	out, err := json.Marshal(&Response{ID: responseID, Code: code, Data: data})
	if err != nil {
		return err
	}
	return sf.Publish(uriService, 1, out)
}

// Connect 将订阅所有相关主题,主题有config配置
func (sf *Client) Connect() error {
	return sf.SubscribeAllTopic(sf.tetrad.ProductKey, sf.tetrad.DeviceName, false)
}

// SubDeviceConnect 子设备连接注册并添加到网关拓扑关系
func (sf *Client) SubDeviceConnect(pk, dn string) error {
	node, err := sf.SearchDevNode(pk, dn)
	if err != nil {
		return err
	}
	if node.DeviceSecret() == "" { // 需要注册
		// 子设备注册
		if err := sf.LinKitGwSubRegister(pk, dn); err != nil {
			return err
		}
	}

	// 子设备添加到拓扑
	return sf.LinkKitGwTopoAdd(pk, dn)
}

func (sf *Client) LinKitGwSubRegister(pk, dn string) error {
	entry, err := sf.ThingGwSubRegister(pk, dn)
	if err != nil {
		return err
	}
	_, err = entry.Wait(time.Second)
	return err
}

func (sf *Client) LinkKitGwTopoAdd(pk, dn string) error {
	entry, err := sf.ThingGwTopoAdd(pk, dn)
	if err != nil {
		return err
	}
	_, err = entry.Wait(time.Second)
	return err
}

// LinkKitGwTopoDelete 删除网关与子设备的拓扑关系
func (sf *Client) LinkKitGwTopoDelete(pk, dn string) error {
	if !sf.isGateway {
		return ErrNotSupportFeature
	}
	entry, err := sf.ThingGwTopoDelete(pk, dn)
	if err != nil {
		return err
	}
	_, err = entry.Wait(time.Second)
	return err
}

func (sf *Client) LinkKitExtCombineLogin(pk, dn string) error {
	if !sf.isGateway {
		return ErrNotSupportFeature
	}
	entry, err := sf.ExtCombineLogin(pk, dn)
	if err != nil {
		return err
	}
	_, err = entry.Wait(time.Second)
	return err
}

func (sf *Client) LinkKitExtCombineLogout(pk, dn string) error {
	if !sf.isGateway {
		return ErrNotSupportFeature
	}
	entry, err := sf.ExtCombineLogout(pk, dn)
	if err != nil {
		return err
	}

	_, err = entry.Wait(time.Second)
	return err
}
