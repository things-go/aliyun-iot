// Package sign 实现MQTT设备签名
package sign

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"sort"
	"strings"

	"github.com/thinkgos/aliyun-iot/infra"
)

// default defined
const (
	fixedTimestamp = "2524608000000"
)

// sign method MQTT设备签名只支持以下签名方法
const (
	signMethodHMACSHA256 = "hmacsha256"
	signMethodHMACSHA1   = "hmacsha1"
)

// all secure mode define
const (
	modeTLSGuider      = "-1"
	modeTLSDirect      = "2"
	modeTCPDirectPlain = "3"
	modeITLSDNSID2     = "8"
)

// MQTTSignInfo 签名后的信息
type MQTTSignInfo struct {
	HostName string
	Port     uint16
	ClientID string
	UserName string
	Password string
}

// SecureMode 安全模式
type SecureMode byte

// SecureMode 支持的安全模型
const (
	SecureModeTLSGuider SecureMode = iota
	SecureModeTLSDirect
	SecureModeTCPDirectPlain
	SecureModeITLSDNSID2
)

// MQTTSign MQTT签名主要设置
type MQTTSign struct {
	enableTLS   bool
	deviceModel bool
	clientIDkv  map[string]string
	hfc         func() hash.Hash
}

// NewMQTTSign 新建一个签名,默认不支持PreAUTH也不支持TLS(即安全模式为SecureModeTcpDirectPlain)
// 默认支持物模型,默认hmacsha256签名加密
func NewMQTTSign() *MQTTSign {
	return &MQTTSign{
		deviceModel: true,
		clientIDkv: map[string]string{
			"timestamp":  fixedTimestamp,
			"securemode": modeTCPDirectPlain,
			"signmethod": signMethodHMACSHA256,
			"lan":        "Golang",
			"v":          infra.IOTAlinkVersion,
		},
		hfc: sha256.New,
	}
}

// SetSignMethod 设置签名方法,目前只支持hmacsha1,hmacsha256, see package infra
func (sf *MQTTSign) SetSignMethod(method string) *MQTTSign {
	if method == signMethodHMACSHA1 {
		sf.clientIDkv["signmethod"] = signMethodHMACSHA1
		sf.hfc = sha1.New
	} else {
		sf.clientIDkv["signmethod"] = signMethodHMACSHA256
		sf.hfc = sha256.New
	}
	return sf
}

// SetSupportSecureMode 设置支持的安全模式
func (sf *MQTTSign) SetSupportSecureMode(mode SecureMode) *MQTTSign {
	switch mode {
	case SecureModeTLSGuider:
		sf.enableTLS = true
		sf.clientIDkv["securemode"] = modeTLSGuider
	case SecureModeTLSDirect:
		sf.enableTLS = true
		sf.clientIDkv["securemode"] = modeTLSDirect
	case SecureModeTCPDirectPlain:
		sf.enableTLS = false
		sf.clientIDkv["securemode"] = modeTCPDirectPlain
	case SecureModeITLSDNSID2:
		sf.enableTLS = true
		sf.clientIDkv["securemode"] = modeITLSDNSID2
	default:
		panic("invalid secure mode")
	}
	return sf
}

// SetSupportDeviceModel 设置支持物模型
func (sf *MQTTSign) SetSupportDeviceModel(enable bool) *MQTTSign {
	if enable {
		sf.clientIDkv["v"] = infra.IOTAlinkVersion
		delete(sf.clientIDkv, "gw")
		delete(sf.clientIDkv, "ext")
	} else {
		sf.clientIDkv["gw"] = "0"
		sf.clientIDkv["ext"] = "0"
		delete(sf.clientIDkv, "v")
	}
	return sf
}

// SetSupportExtRRPC 支持扩展RRPC 仅物模型下支持
func (sf *MQTTSign) SetSupportExtRRPC() *MQTTSign {
	if _, ok := sf.clientIDkv["v"]; ok {
		sf.clientIDkv["ext"] = "1"
	}

	return sf
}

// SetSDKVersion 设备SDK版本
func (sf *MQTTSign) SetSDKVersion(ver string) *MQTTSign {
	sf.clientIDkv["_v"] = ver
	return sf
}

// AddCustomKV 添加一个用户的键值对,键值对将被添加到clientID上
func (sf *MQTTSign) AddCustomKV(key, value string) *MQTTSign {
	sf.clientIDkv[key] = value
	return sf
}

// DeleteCustomKV 删除一个用户的键值对
func (sf *MQTTSign) DeleteCustomKV(key string) *MQTTSign {
	delete(sf.clientIDkv, key)
	return sf
}

// generateClientID 根据deviceID生成clientID
func (sf *MQTTSign) generateClientID(deviceID string) string {
	builder := new(strings.Builder)
	builder.WriteString(deviceID)
	builder.WriteString("|")

	sKey := make([]string, 0, len(sf.clientIDkv))
	for k := range sf.clientIDkv {
		sKey = append(sKey, k)
	}
	// 对键进行排序
	sort.Strings(sKey)

	for _, Value := range sKey {
		builder.WriteString(Value)
		builder.WriteString("=")
		builder.WriteString(sf.clientIDkv[Value])
		builder.WriteString(",")
	}

	return strings.TrimRight(builder.String(), ",") + "|"
}

// Generate 根据MetaInfo和region生成签名
func (sf *MQTTSign) Generate(meta *infra.MetaInfo, crd infra.CloudRegionDomain) (*MQTTSignInfo, error) {
	signOut := &MQTTSignInfo{}

	/* setup ClientID */
	deviceID := fmt.Sprintf("%s.%s", meta.ProductKey, meta.DeviceName)

	signOut.ClientID = sf.generateClientID(deviceID)

	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%s",
		deviceID, meta.DeviceName, meta.ProductKey, fixedTimestamp)
	/* setup Password */
	h := hmac.New(sf.hfc, []byte(meta.DeviceSecret))
	if _, err := h.Write([]byte(signSource)); err != nil {
		return nil, err
	}
	signOut.Password = hex.EncodeToString(h.Sum(nil))

	/* setup HostName */
	if crd.Region == infra.CloudRegionCustom {
		if crd.CustomDomain == "" {
			return nil, errors.New("custom domain invalid")
		}
		signOut.HostName = crd.CustomDomain
	} else {
		signOut.HostName = fmt.Sprintf("%s.%s", meta.ProductKey, infra.MQTTCloudDomain[crd.Region])
	}
	/* setup UserName */
	signOut.UserName = fmt.Sprintf("%s&%s", meta.DeviceName, meta.ProductKey)

	/* setup Port */
	if sf.enableTLS {
		signOut.Port = 443
	} else {
		signOut.Port = 1883
	}

	return signOut, nil
}
