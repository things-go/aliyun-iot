// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sign

import (
	"strconv"
)

// Option option
type Option func(*config)

// WithSignMethod 设置签名方法,目前只支持hmacsha1,hmacmd5,hmacsha256(默认)
func WithSignMethod(method string) Option {
	return func(ms *config) {
		switch method {
		case hmacsha1:
			ms.extParams["signmethod"] = hmacsha1
		case hmacmd5:
			ms.extParams["signmethod"] = hmacmd5
		case hmacsha256:
			fallthrough
		default:
			ms.extParams["signmethod"] = hmacsha256
		}
	}
}

// WithSecureMode 设置支持的安全模式
func WithSecureMode(mode SecureMode) Option {
	return func(ms *config) {
		switch mode {
		case SecureModeTLSGuider:
			ms.enableTLS = true
			ms.extParams["securemode"] = modeTLSGuider
		case SecureModeTLSDirect:
			ms.enableTLS = true
			ms.extParams["securemode"] = modeTLSDirect
		case SecureModeITLSDNSID2:
			ms.enableTLS = true
			ms.extParams["securemode"] = modeITLSDNSID2
		case SecureModeTCPDirectPlain:
			fallthrough
		default:
			ms.enableTLS = false
			ms.extParams["securemode"] = modeTCPDirectPlain
		}
	}
}

// WithEnableDeviceModel 设置是否支持物模型
func WithEnableDeviceModel(enable bool) Option {
	return func(ms *config) {
		if enable {
			ms.extParams["v"] = alinkVersion
			delete(ms.extParams, "gw")
			delete(ms.extParams, "ext")
		} else {
			ms.extParams["gw"] = "0"
			ms.extParams["ext"] = "0"
			delete(ms.extParams, "v")
		}
	}
}

// WithExtRRPC 支持扩展RRPC 仅物模型下支持
func WithExtRRPC() Option {
	return func(ms *config) {
		if _, ok := ms.extParams["v"]; ok {
			ms.extParams["ext"] = "1"
		}
	}
}

// WithSDKVersion 设备SDK版本
func WithSDKVersion(ver string) Option {
	return func(ms *config) {
		ms.extParams["_v"] = ver
	}
}

// WithExtParamsKV 添加一个扩展参数的键值对,键值对将被添加到clientID的扩展参数上
func WithExtParamsKV(key, value string) Option {
	return func(ms *config) {
		ms.extParams[key] = value
	}
}

// WithTimestamp 添加当前时间的毫秒值,可以不传递,默认为一个固定的值
func WithTimestamp(t uint64) Option {
	return WithExtParamsKV("timestamp", strconv.FormatUint(t, 10))
}
