package model

import (
	"time"

	"github.com/thinkgos/aliIOT/feature"
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
	*feature.Options
}

// NewOption 创建选项
func NewOption(productKey, deviceName, deviceSecret string) *Config {
	return &Config{
		productKey:   productKey,
		deviceName:   deviceName,
		deviceSecret: deviceSecret,

		cacheExpiration:      DefaultExpiration,
		cacheCleanupInterval: DefaultCleanupInterval,
		Options:              feature.New(),
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

func (sf *Config) FeatureOption() *feature.Options {
	return sf.Options
}
