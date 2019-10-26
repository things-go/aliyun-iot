package model

import (
	"time"
)

// 默认值
const (
	DefaultExpiration      = time.Second * 10
	DefaultCleanupInterval = time.Second * 30
)

// Options
type Options struct {
	productKey   string
	deviceName   string
	deviceSecret string
	enableCache  bool
	uriOffset    int

	expiration      time.Duration
	cleanupInterval time.Duration
}

// NewOption 创建选项
func NewOption(productKey, deviceName, deviceSecret string) *Options {
	return &Options{
		productKey:   productKey,
		deviceName:   deviceName,
		deviceSecret: deviceSecret,

		expiration:      DefaultExpiration,
		cleanupInterval: DefaultCleanupInterval,
	}
}

// Valid 校验消息有效,无效采用相应默认值
func (sf *Options) Valid() *Options {
	return sf
}

func (sf *Options) MetaInfo() (productKey, deviceName, deviceSecret string) {
	return sf.productKey, sf.deviceName, sf.deviceSecret
}

// SetEnableCache 使能消息缓存
func (sf *Options) SetEnableCache(enable bool) *Options {
	sf.enableCache = enable
	return sf
}

// SetCacheTimeout 设备消息缓存超时时间
func (sf *Options) SetCacheTimeout(expiration, cleanupInterval time.Duration) *Options {
	sf.expiration = expiration
	sf.cleanupInterval = cleanupInterval
	return sf
}

// EnableCOAP 采用COAP
func (sf *Options) EnableCOAP(enable bool) *Options {
	if enable {
		sf.uriOffset = 1
	} else {
		sf.uriOffset = 0
	}
	return sf
}

func (sf *Options) EnableHTTP(enable bool) *Options {
	sf.EnableCOAP(enable)
	return sf
}
