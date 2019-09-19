// package 实现阿里云设备签名
package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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

type MQTTSignOut struct {
	HostName string
	port     uint16
	clientID string
	UserName string
	password string
}

// SecureMode
type SecureMode byte

// SecureMode 支持的安全模型
const (
	SecureModeTLSGuider SecureMode = iota
	SecureModeTLSDirect
	SecureModeTcpDirectPlain
	SecureModeITLSDNSID2
)

type Sign struct {
	deviceModel bool
	clientIDkv  map[string]string
}

// New 新建一个签名,默认不支持PreAUTH也不支持TLS(即安全模式为SecureModeTcpDirectPlain)
// 支持物模型,hmacsha256签名加密
func New() *Sign {
	sig := &Sign{
		deviceModel: true,
		clientIDkv: map[string]string{
			"timestamp":  fixedTimestamp,
			"_v":         IotxSDKVersion,
			"securemode": modeTcpDirectPlain,
			"signmethod": "hmacsha256",
			"lan":        "Golang",
			"v":          IotxAlinkVersion,
		},
	}
	return sig
}

// SetSupportSecureMode 设置支持的安全模式
func (this *Sign) SetSupportSecureMode(mode SecureMode) {
	switch mode {
	case SecureModeTLSGuider:
		this.clientIDkv["securemode"] = modeTLSGuider
	case SecureModeTLSDirect:
		this.clientIDkv["securemode"] = modeTLSDirect
	case SecureModeTcpDirectPlain:
		this.clientIDkv["securemode"] = modeTcpDirectPlain
	case SecureModeITLSDNSID2:
		this.clientIDkv["securemode"] = modeTcpDirectPlain
	default:
		panic("invalid secure mode")
	}
}

// SetSupportDeviceModel 设置支持物模型
func (this *Sign) SetSupportDeviceModel(enable bool) {
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

// AddCustomKV
func (this *Sign) AddCustomKV(key, value string) {
	this.clientIDkv[key] = value
}

// DeleteCustomKV
func (this *Sign) DeleteCustomKV(key string) {
	delete(this.clientIDkv, key)
}

func (this *Sign) generateClientID(deviceID string) string {
	builder := new(strings.Builder)
	builder.WriteString(deviceID)
	builder.WriteString("|")

	for k, v := range this.clientIDkv {
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(v)
		builder.WriteString(",")
	}

	return strings.TrimRight(builder.String(), ",") + "|"
}

func (this *Sign) Generate(meta *MetaInfo, region MQTTCloudRegion) (*MQTTSignOut, error) {
	signOut := &MQTTSignOut{}

	/* setup clientID */
	deviceID := fmt.Sprintf("%s.%s", meta.ProductKey, meta.DeviceName)

	signOut.clientID = this.generateClientID(deviceID)
	/* setup password */
	h := hmac.New(sha256.New, []byte(meta.DeviceSecret))
	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%s",
		deviceID, meta.DeviceName, meta.ProductKey, fixedTimestamp)
	h.Write([]byte(signSource))
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
	signOut.port = 1883

	return signOut, nil
}
