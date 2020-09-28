// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sign

import (
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
)

// Option option
type Option func(*config)

// WithDeviceToken 设置device token 只在SecureModeNoPreRegistration时需设置
func WithDeviceToken(deviceToken string) Option {
	return func(c *config) {
		c.deviceToken = deviceToken
	}
}

// WithPort 设置端口, 默认1883
func WithPort(port uint16) Option {
	return func(c *config) {
		c.port = port
	}
}

// WithSignMethod 设置签名方法,目前只支持hmacsha1,hmacmd5,hmacsha256(默认)
func WithSignMethod(method string) Option {
	return func(c *config) {
		c.method = method
	}
}

// WithSecureMode 设置安全模式, 默认SecureModeTCPDirectPlain,tcp直连
func WithSecureMode(secureMode string) Option {
	return func(c *config) {
		c.secureMode = secureMode
	}
}

// WithEnableDeviceModel 设置是否支持物模型,默认为物模型
func WithEnableDeviceModel(enable bool) Option {
	return func(c *config) {
		c.enableDM = enable
	}
}

// WithExtRRPC 支持扩展RRPC 仅物模型下支持
func WithExtRRPC() Option {
	return func(c *config) {
		if _, ok := c.extParams["v"]; ok {
			c.extParams["ext"] = "1"
		}
	}
}

// WithSDKVersion 设备SDK版本
func WithSDKVersion(ver string) Option {
	return func(c *config) {
		c.extParams["_v"] = ver
	}
}

// WithExtParamsKV 添加一个扩展参数的键值对,键值对将被添加到clientID的扩展参数上
func WithExtParamsKV(key, value string) Option {
	return func(c *config) {
		c.extParams[key] = value
	}
}

// WithTimestamp 添加当前时间的毫秒值,可以不传递,默认: fixedTimestamp,单位ms
func WithTimestamp() Option {
	return func(c *config) {
		c.timestamp = infra.Millisecond(time.Now())
	}
}
