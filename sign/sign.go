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

	"github.com/thinkgos/aliIOT/infra"
)

// default defined
const (
	fixedTimestamp = "2524608000000"
)

// all secure mode define
const (
	modeTLSGuider      = "-1"
	modeTLSDirect      = "2"
	modeTCPDirectPlain = "3"
	modeITLSDNSID2     = "8"
)

// sign method 签名方法
const (
	SignMethodSHA256 = "hmacsha256"
	SignMethodSHA1   = "hmacsha1"
)

// mqtt 域名
var mqttDomain = []string{
	"iot-as-mqtt.cn-shanghai.aliyuncs.com",    /* Shanghai */
	"iot-as-mqtt.ap-southeast-1.aliyuncs.com", /* Singapore */
	"iot-as-mqtt.ap-northeast-1.aliyuncs.com", /* Japan */
	"iot-as-mqtt.us-west-1.aliyuncs.com",      /* America */
	"iot-as-mqtt.eu-central-1.aliyuncs.com",   /* Germany */
}

// MQTTCloudRegion MQTT云端地域
type MQTTCloudRegion byte

// 云平台地域定义
const (
	CloudRegionShangHai MQTTCloudRegion = iota
	CloudRegionSingapore
	CloudRegionJapan
	CloudRegionAmerica
	CloudRegionGermany
	CloudRegionCustom
)

// MetaInfo 产品与设备三元组
type MetaInfo struct {
	ProductKey    string
	ProductSecret string
	DeviceName    string
	DeviceSecret  string
	CustomDomain  string // 如果使用CloudRegionCustom,需要定义此字段
}

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
// 支持物模型,默认hmacsha256签名加密
func NewMQTTSign() *MQTTSign {
	return &MQTTSign{
		deviceModel: true,
		clientIDkv: map[string]string{
			"timestamp":  fixedTimestamp,
			"_v":         infra.IOTSDKVersion,
			"securemode": modeTCPDirectPlain,
			"signmethod": SignMethodSHA256,
			"lan":        "Golang",
			"v":          infra.IOTAlinkVersion,
		},
		hfc: sha256.New,
	}
}

// SetSignMethod 设置签名方法
func (sf *MQTTSign) SetSignMethod(method string) *MQTTSign {
	if method == SignMethodSHA1 {
		sf.clientIDkv["signmethod"] = SignMethodSHA1
		sf.hfc = sha1.New
	} else {
		sf.clientIDkv["signmethod"] = SignMethodSHA256
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

// AddCustomKV 添加一个用户的键值对,键值对将被添加到clientID上
func (sf *MQTTSign) AddCustomKV(key, value string) *MQTTSign {
	sf.clientIDkv[key] = value
	return sf
}

// DeleteCustomKV 删除一个用户的键值对
func (sf *MQTTSign) DeleteCustomKV(key string) {
	delete(sf.clientIDkv, key)
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
func (sf *MQTTSign) Generate(meta *MetaInfo, region MQTTCloudRegion) (*MQTTSignInfo, error) {
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
	if region == CloudRegionCustom {
		if meta.CustomDomain == "" {
			return nil, errors.New("custom domain invalid")
		}
		signOut.HostName = meta.CustomDomain
	} else {
		signOut.HostName = fmt.Sprintf("%s.%s", meta.ProductKey, mqttDomain[region])
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
