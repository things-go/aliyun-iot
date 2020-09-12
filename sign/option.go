package sign

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
)

// Option option
type Option func(*Sign)

// WithSignMethod 设置签名方法,目前只支持hmacsha1,hmacsha256,hmacmd5(默认)
func WithSignMethod(method string) Option {
	return func(ms *Sign) {
		switch method {
		case hmacsha1:
			ms.clientIDkv["signmethod"] = hmacsha1
			ms.hfc = sha1.New
		case hmacsha256:
			ms.clientIDkv["signmethod"] = hmacsha256
			ms.hfc = sha256.New
		case hmacmd5:
			fallthrough
		default:
			ms.clientIDkv["signmethod"] = hmacmd5
			ms.hfc = md5.New
		}
	}
}

// WithSecureMode 设置支持的安全模式
func WithSecureMode(mode SecureMode) Option {
	return func(ms *Sign) {
		switch mode {
		case SecureModeTLSGuider:
			ms.enableTLS = true
			ms.clientIDkv["securemode"] = modeTLSGuider
		case SecureModeTLSDirect:
			ms.enableTLS = true
			ms.clientIDkv["securemode"] = modeTLSDirect
		case SecureModeITLSDNSID2:
			ms.enableTLS = true
			ms.clientIDkv["securemode"] = modeITLSDNSID2
		case SecureModeTCPDirectPlain:
			fallthrough
		default:
			ms.enableTLS = false
			ms.clientIDkv["securemode"] = modeTCPDirectPlain
		}
	}
}

// WithDeviceModel 设置支持物模型
func WithDeviceModel(enable bool) Option {
	return func(ms *Sign) {
		if enable {
			ms.clientIDkv["v"] = alinkVersion
			delete(ms.clientIDkv, "gw")
			delete(ms.clientIDkv, "ext")
		} else {
			ms.clientIDkv["gw"] = "0"
			ms.clientIDkv["ext"] = "0"
			delete(ms.clientIDkv, "v")
		}
	}
}

// WithExtRRPC 支持扩展RRPC 仅物模型下支持
func WithExtRRPC() Option {
	return func(ms *Sign) {
		if _, ok := ms.clientIDkv["v"]; ok {
			ms.clientIDkv["ext"] = "1"
		}
	}
}

// WithSDKVersion 设备SDK版本
func WithSDKVersion(ver string) Option {
	return func(ms *Sign) {
		ms.clientIDkv["_v"] = ver
	}
}

// WithCustomKV 添加一个用户的键值对,键值对将被添加到clientID上
func WithCustomKV(key, value string) Option {
	return func(ms *Sign) {
		ms.clientIDkv[key] = value
	}
}
