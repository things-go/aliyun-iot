// Package dm imp aliyun dm
package dm

import (
	"encoding/json"
	"math/rand"
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

// Version 平台通信版本
var Version = "1.0"

// New 创建一个物管理客户端
func New(triad infra.MetaTriad, conn Conn, opts ...Option) *Client {
	c := &Client{
		requestID: rand.Uint32(),
		tetrad:    triad,

		workOnWho: WorkOnMQTT,

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
	if c.workOnWho != WorkOnHTTP {
		c.msgCache = cache.New(c.cacheExpiration, c.cacheCleanupInterval)
	}
	return c
}

// Connect 将订阅所有相关主题,主题有config配置
func (sf *Client) Connect() error {
	return sf.SubscribeAllTopic(sf.tetrad.ProductKey, sf.tetrad.DeviceName, false)
}

// NewSubDevice 创建一个子设备
func (sf *Client) NewSubDevice(meta infra.MetaTetrad) error {
	if sf.isGateway {
		return sf.Create(meta)
	}
	return ErrNotSupportFeature
}

// SubDeviceConnect 子设备连接注册并添加到网关拓扑关系
func (sf *Client) SubDeviceConnect(pk, dn string, timeout time.Duration) error {
	node, err := sf.SearchAvail(pk, dn)
	if err != nil {
		return err
	}
	if node.DeviceSecret() == "" { // 需要注册
		// 子设备注册
		if err := sf.LinkThingSubRegister(pk, dn, timeout); err != nil {
			return err
		}
	}
	// 子设备添加到拓扑
	return sf.LinkThingTopoAdd(pk, dn, timeout)
}
