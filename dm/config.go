package dm

import (
	"time"
)

// 缓存默认值
const (
	DefaultCacheExpiration      = time.Second * 10
	DefaultCacheCleanupInterval = time.Second * 30
)

// 当前的工作方式
const (
	workOnMQTT = iota
	workOnCOAP
	workOnHTTP
)

// Config
type Config struct {
	productKey   string
	deviceName   string
	deviceSecret string

	uriOffset int
	workOnWho byte

	cacheExpiration      time.Duration
	cacheCleanupInterval time.Duration
	hasCache             bool

	hasNTP      bool
	hasRawModel bool
	hasDesired  bool
	hasExtRRPC  bool
	hasGateway  bool
}

// NewConfig 创建新配置
func NewConfig(productKey, deviceName, deviceSecret string) *Config {
	return &Config{
		productKey:   productKey,
		deviceName:   deviceName,
		deviceSecret: deviceSecret,

		uriOffset: 0,
		workOnWho: workOnMQTT,

		cacheExpiration:      DefaultCacheExpiration,
		cacheCleanupInterval: DefaultCacheCleanupInterval,
	}
}

// Valid 校验配置有效,无效条目将采用相应默认值
func (sf *Config) Valid() *Config {
	return sf
}

// MetaInfo 获取设备mata info
func (sf *Config) MetaInfo() (productKey, deviceName, deviceSecret string) {
	return sf.productKey, sf.deviceName, sf.deviceSecret
}

// SetEnableCache 使能消息缓存
func (sf *Config) SetEnableCache(enable bool) *Config {
	sf.hasCache = enable
	return sf
}

// SetCacheTimeout 设备消息缓存超时时间
func (sf *Config) SetCacheTimeout(expiration, cleanupInterval time.Duration) *Config {
	sf.cacheExpiration = expiration
	sf.cacheCleanupInterval = cleanupInterval
	return sf
}

// EnableCOAP 采用COAP
func (sf *Config) EnableCOAP() *Config {
	sf.workOnWho = workOnCOAP
	sf.uriOffset = 1
	return sf
}

// EnableHTTP 采用HTTP
func (sf *Config) EnableHTTP() *Config {
	sf.workOnWho = workOnHTTP
	sf.uriOffset = 1
	return sf
}

// EnableNTP 使能NTP
func (sf *Config) EnableNTP() *Config {
	sf.hasNTP = true
	return sf
}

// EnableModelRaw 使能透传
func (sf *Config) EnableModelRaw() *Config {
	sf.hasRawModel = true
	return sf
}

// EnableDesired 使能期望属性
func (sf *Config) EnableDesired() *Config {
	sf.hasDesired = true
	return sf
}

// EnableExtRRPC 使能扩展RRPC功能
func (sf *Config) EnableExtRRPC() *Config {
	sf.hasExtRRPC = true
	return sf
}

// EnableGateway 使能网关功能
func (sf *Config) EnableGateway() *Config {
	sf.hasGateway = true
	return sf
}
