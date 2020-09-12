// Package sign 实现MQTT设备签名
// see https://help.aliyun.com/document_detail/73742.html?spm=a2c4g.11186623.6.599.76216eebzbvrYq
package sign

import (
	"crypto/hmac"
	"crypto/md5"
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
	alinkVersion   = "20"

	// sign method MQTT设备签名只支持以下签名方法
	hmacsha256 = "hmacsha256"
	hmacsha1   = "hmacsha1"
	hmacmd5    = "hmacmd5"
)

// all secure mode define
const (
	modeTLSGuider      = "-1"
	modeTLSDirect      = "2"
	modeTCPDirectPlain = "3"
	modeITLSDNSID2     = "8"
)

// SecureMode 安全模式
type SecureMode byte

// SecureMode 支持的安全模型
const (
	SecureModeTLSGuider SecureMode = iota
	SecureModeTLSDirect
	SecureModeTCPDirectPlain
	SecureModeITLSDNSID2
)

// SignInfo 签名后的信息
type SignInfo struct {
	HostName string
	Port     uint16
	ClientID string
	UserName string
	Password string
}

// Addr address like host:port
func (sf *SignInfo) Addr() string {
	return fmt.Sprintf("%s:%d", sf.HostName, sf.Port)
}

// Sign MQTT签名主要设置
type Sign struct {
	enableTLS   bool
	deviceModel bool
	clientIDkv  map[string]string
	hfc         func() hash.Hash
}

// AlinkSDKVersion alink sdk version
var AlinkSDKVersion = "sdk-golang-v0.0.1"

// New 新建一个签名,默认不支持PreAUTH也不支持TLS(即安全模式为SecureModeTcpDirectPlain)
// 默认支持物模型,默认hmacmd5签名加密
// TODO: 支持tls
func New(opts ...Option) *Sign {
	ms := &Sign{
		deviceModel: true,
		clientIDkv: map[string]string{
			"timestamp":  fixedTimestamp,
			"securemode": modeTCPDirectPlain,
			"signmethod": hmacmd5,
			"lan":        "Golang",
			"v":          alinkVersion,
		},
		hfc: md5.New,
	}
	for _, opt := range opts {
		opt(ms)
	}
	return ms
}

// Generate 根据MetaInfo和region生成签名
func (sf *Sign) Generate(meta *infra.MetaInfo, crd infra.CloudRegionDomain) (*SignInfo, error) {
	if crd.Region == infra.CloudRegionCustom && crd.CustomDomain == "" {
		return nil, errors.New("invalid custom domain")
	}

	// setup ClientID
	clientID := fmt.Sprintf("%s.%s", meta.ProductKey, meta.DeviceName)

	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%s",
		clientID, meta.DeviceName, meta.ProductKey, fixedTimestamp)
	// setup Password
	h := hmac.New(sf.hfc, []byte(meta.DeviceSecret))
	if _, err := h.Write([]byte(signSource)); err != nil {
		return nil, err
	}
	pwd := hex.EncodeToString(h.Sum(nil))

	signOut := &SignInfo{
		Port:     1883,
		ClientID: generateClientID(sf.clientIDkv, clientID),
		UserName: fmt.Sprintf("%s&%s", meta.DeviceName, meta.ProductKey),
		Password: pwd,
	}

	// setup HostName
	if crd.Region == infra.CloudRegionCustom {
		signOut.HostName = crd.CustomDomain
	} else {
		signOut.HostName = fmt.Sprintf("%s.%s", meta.ProductKey, infra.MQTTCloudDomain[crd.Region])
	}
	// setup Port
	if sf.enableTLS {
		signOut.Port = 443
	}
	return signOut, nil
}

// generateClientID 根据deviceID生成clientID
func generateClientID(clientIDkv map[string]string, deviceID string) string {
	builder := new(strings.Builder)
	builder.WriteString(deviceID)
	builder.WriteString("|")

	sKey := make([]string, 0, len(clientIDkv))
	for k := range clientIDkv {
		sKey = append(sKey, k)
	}
	// 对键进行排序
	sort.Strings(sKey)

	for _, Value := range sKey {
		builder.WriteString(Value)
		builder.WriteString("=")
		builder.WriteString(clientIDkv[Value])
		builder.WriteString(",")
	}
	return strings.TrimRight(builder.String(), ",") + "|"
}
