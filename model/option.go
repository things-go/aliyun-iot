package model

import (
	"time"
)

// 默认值
const (
	DefaultExpiration      = time.Second * 10
	DefaultCleanupInterval = time.Second * 30
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
	enableCache          bool

	hasNTP      bool
	hasRawModel bool
	hasDesired  bool
	hasExtRRPC  bool
	hasGateway  bool
}

// NewOption 创建选项
func NewOption(productKey, deviceName, deviceSecret string) *Config {
	return &Config{
		productKey:   productKey,
		deviceName:   deviceName,
		deviceSecret: deviceSecret,

		cacheExpiration:      DefaultExpiration,
		cacheCleanupInterval: DefaultCleanupInterval,
	}
}

// Valid 校验消息有效,无效采用相应默认值
func (sf *Config) Valid() *Config {
	return sf
}

// MetaInfo 获取设备mata info
func (sf *Config) MetaInfo() (productKey, deviceName, deviceSecret string) {
	return sf.productKey, sf.deviceName, sf.deviceSecret
}

// SetEnableCache 使能消息缓存
func (sf *Config) SetEnableCache(enable bool) *Config {
	sf.enableCache = enable
	return sf
}

// SetCacheTimeout 设备消息缓存超时时间
func (sf *Config) SetCacheTimeout(expiration, cleanupInterval time.Duration) *Config {
	sf.cacheExpiration = expiration
	sf.cacheCleanupInterval = cleanupInterval
	return sf
}

// EnableCOAP 采用COAP
func (sf *Config) EnableCOAP(enable bool) *Config {
	if enable {
		sf.workOnWho = workOnCOAP
		sf.uriOffset = 1
	} else {
		sf.workOnWho = workOnMQTT
		sf.uriOffset = 0
	}
	return sf
}

// EnableCOAP 采用HTTP
func (sf *Config) EnableHTTP(enable bool) *Config {
	if enable {
		sf.workOnWho = workOnHTTP
		sf.uriOffset = 1
	} else {
		sf.workOnWho = workOnHTTP
		sf.uriOffset = 0
	}
	return sf
}

func (sf *Config) EnableNTP() *Config {
	sf.hasNTP = true
	return sf
}

// 使能透传 EnableModelRaw
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

func (sf *Config) EnableGateway() *Config {
	sf.hasGateway = true
	return sf
}
