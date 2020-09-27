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

// WithSignMethod 设置签名方法,目前只支持hmacsha1,hmacmd5,hmacsha256(默认)
func WithSignMethod(method string) Option {
	return func(c *config) {
		switch method {
		case hmacsha1:
			c.extParams["signmethod"] = hmacsha1
		case hmacmd5:
			c.extParams["signmethod"] = hmacmd5
		case hmacsha256:
			fallthrough
		default:
			c.extParams["signmethod"] = hmacsha256
		}
	}
}

// WithSecureMode 设置支持的安全模式
func WithSecureMode(mode SecureMode) Option {
	return func(c *config) {
		switch mode {
		case SecureModeTLSGuider:
			c.enableTLS = true
			c.extParams["securemode"] = modeTLSGuider
		case SecureModeTLSDirect:
			c.enableTLS = true
			c.extParams["securemode"] = modeTLSDirect
		case SecureModeITLSDNSID2:
			c.enableTLS = true
			c.extParams["securemode"] = modeITLSDNSID2
		case SecureModeTCPDirectPlain:
			fallthrough
		default:
			c.enableTLS = false
			c.extParams["securemode"] = modeTCPDirectPlain
		}
	}
}

// WithEnableDeviceModel 设置是否支持物模型
func WithEnableDeviceModel(enable bool) Option {
	return func(c *config) {
		if enable {
			c.extParams["v"] = alinkVersion
			delete(c.extParams, "gw")
			delete(c.extParams, "ext")
		} else {
			c.extParams["gw"] = "0"
			c.extParams["ext"] = "0"
			delete(c.extParams, "v")
		}
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

// WithTimestamp 添加当前时间的毫秒值,可以不传递,默认: fixedTimestamp 单位ms
func WithTimestamp() Option {
	return func(c *config) {
		c.timestamp = infra.Millisecond(time.Now())
	}
}
