// package 实现阿里云设备签名
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
)

const (
	IotxAlinkVersion = "20"
	IotxSDKVersion   = "sdk-golang-3.0.1"
	fixedTimestamp   = "2524608000000"
)

/* all secure mode define */
const (
	modeTLSGuider      = "-1"
	modeTLSDirect      = "2"
	modeTcpDirectPlain = "3"
	modeITLSDNSID2     = "8"
)

// signmethod 签名方法
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

// http 域名
var httpDomain = []string{
	"iot-auth.cn-shanghai.aliyuncs.com",    /* Shanghai */
	"iot-auth.ap-southeast-1.aliyuncs.com", /* Singapore */
	"iot-auth.ap-northeast-1.aliyuncs.com", /* Japan */
	"iot-auth.us-west-1.aliyuncs.com",      /* America */
	"iot-auth.eu-central-1.aliyuncs.com",   /* Germany */
}

// MQTTCloudRegion MQTT云端地域
type MQTTCloudRegion byte

// 云平台地域
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
	ProductKey, ProductSecret string
	DeviceName, DeviceSecret  string
	CustomDomain              string
}

// MQTTSignInfo 签名后的信息
type MQTTSignInfo struct {
	HostName string
	port     uint16
	clientID string
	UserName string
	password string
}

// SecureMode 安全模式
type SecureMode byte

// SecureMode 支持的安全模型
const (
	SecureModeTLSGuider SecureMode = iota
	SecureModeTLSDirect
	SecureModeTcpDirectPlain
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
			"_v":         IotxSDKVersion,
			"securemode": modeTcpDirectPlain,
			"signmethod": SignMethodSHA256,
			"lan":        "Golang",
			"v":          IotxAlinkVersion,
		},
		hfc: sha256.New,
	}
}

// SetSignMethod 设置签名方法
func (this *MQTTSign) SetSignMethod(method string) {
	if method == SignMethodSHA1 {
		this.clientIDkv["signmethod"] = SignMethodSHA1
		this.hfc = sha1.New
	} else {
		this.clientIDkv["signmethod"] = SignMethodSHA256
		this.hfc = sha256.New
	}
}

// SetSupportSecureMode 设置支持的安全模式
func (this *MQTTSign) SetSupportSecureMode(mode SecureMode) {
	switch mode {
	case SecureModeTLSGuider:
		this.enableTLS = true
		this.clientIDkv["securemode"] = modeTLSGuider
	case SecureModeTLSDirect:
		this.enableTLS = true
		this.clientIDkv["securemode"] = modeTLSDirect
	case SecureModeTcpDirectPlain:
		this.enableTLS = false
		this.clientIDkv["securemode"] = modeTcpDirectPlain
	case SecureModeITLSDNSID2:
		this.enableTLS = true
		this.clientIDkv["securemode"] = modeITLSDNSID2
	default:
		panic("invalid secure mode")
	}
}

// SetSupportDeviceModel 设置支持物模型
func (this *MQTTSign) SetSupportDeviceModel(enable bool) {
	if enable {
		this.clientIDkv["v"] = IotxAlinkVersion
		delete(this.clientIDkv, "gw")
		delete(this.clientIDkv, "ext")
	} else {
		this.clientIDkv["gw"] = "0"
		this.clientIDkv["ext"] = "0"
		delete(this.clientIDkv, "v")
	}
}

// AddCustomKV 添加一个用户的键值对,键值对将被添加到clientID上
func (this *MQTTSign) AddCustomKV(key, value string) {
	this.clientIDkv[key] = value
}

// DeleteCustomKV 删除一个用户的键值对
func (this *MQTTSign) DeleteCustomKV(key string) {
	delete(this.clientIDkv, key)
}

// generateClientID 根据deviceID生成clientID
func (this *MQTTSign) generateClientID(deviceID string) string {
	builder := new(strings.Builder)
	builder.WriteString(deviceID)
	builder.WriteString("|")

	sKey := make([]string, 0, len(this.clientIDkv))
	for k := range this.clientIDkv {
		sKey = append(sKey, k)
	}

	sort.Strings(sKey)

	for _, kValue := range sKey {
		builder.WriteString(kValue)
		builder.WriteString("=")
		builder.WriteString(this.clientIDkv[kValue])
		builder.WriteString(",")
	}

	return strings.TrimRight(builder.String(), ",") + "|"
}

// Generate 根据MetaInfo和region生成签名
func (this *MQTTSign) Generate(meta *MetaInfo, region MQTTCloudRegion) (*MQTTSignInfo, error) {
	signOut := &MQTTSignInfo{}

	/* setup clientID */
	deviceID := fmt.Sprintf("%s.%s", meta.ProductKey, meta.DeviceName)

	signOut.clientID = this.generateClientID(deviceID)
	/* setup password */
	h := hmac.New(this.hfc, []byte(meta.DeviceSecret))
	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%s",
		deviceID, meta.DeviceName, meta.ProductKey, fixedTimestamp)
	if _, err := h.Write([]byte(signSource)); err != nil {
		return nil, err
	}
	signOut.password = hex.EncodeToString(h.Sum(nil))

	/* setup HostName */
	if region == CloudRegionCustom {
		if len(meta.CustomDomain) == 0 {
			return nil, errors.New("custom domain invalid")
		}
		signOut.HostName = meta.CustomDomain
	} else {
		signOut.HostName = fmt.Sprintf("%s.%s", meta.ProductKey, mqttDomain[region])
	}
	/* setup UserName */
	signOut.UserName = fmt.Sprintf("%s&%s", meta.DeviceName, meta.ProductKey)

	/* setup port */
	if this.enableTLS {
		signOut.port = 443
	} else {
		signOut.port = 1883
	}

	return signOut, nil
}
