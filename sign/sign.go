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
	HostName  string
	Port      uint16
	ClientID  string
	extParams string // clientID的扩展参数
	UserName  string // deviceName & productKey
	Password  string
}

// Addr address like host:port
func (sf *SignInfo) Addr() string {
	return fmt.Sprintf("%s:%d", sf.HostName, sf.Port)
}

// ClientIDWithExt 登录使用的 client + "| 扩展参数 |"
func (sf SignInfo) ClientIDWithExt() string {
	return sf.ClientID + sf.extParams
}

// Sign MQTT签名主要设置
type Sign struct {
	enableTLS bool              // 使能tls
	enableDM  bool              // 使能物模型
	extParams map[string]string // clientID扩展参数
	hfc       func() hash.Hash
}

// SDKVersion alink sdk version
var SDKVersion = "sdk-golang-v0.0.1"

// New 新建一个签名
// 默认不支持PreAUTH
// 默认安全模式为SecureModeTcpDirectPlain)
// 默认使能物模型
// 默认hmacmd5签名加密
// 默认sdk版本为 SDKVersion
// TODO: 支持tls
func New(opts ...Option) *Sign {
	ms := &Sign{
		enableDM: true,
		extParams: map[string]string{
			"timestamp":  fixedTimestamp,
			"securemode": modeTCPDirectPlain,
			"signmethod": hmacmd5,
			"lan":        "Golang",
			"v":          alinkVersion,
			"_v":         SDKVersion,
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

	signInfo := &SignInfo{
		Port:      1883,
		ClientID:  clientID,
		extParams: generateExtParam(sf.extParams),
		UserName:  meta.DeviceName + "&" + meta.ProductKey,
		Password:  pwd,
	}

	domain := crd.CustomDomain
	if crd.Region != infra.CloudRegionCustom {
		domain = infra.MQTTCloudDomain[crd.Region]
	}
	// setup HostName
	signInfo.HostName = meta.ProductKey + "." + domain
	// setup Port
	if sf.enableTLS {
		signInfo.Port = 443
	}
	return signInfo, nil
}

// generateExtParam 根据deviceID生成clientID
func generateExtParam(extParams map[string]string) string {
	var l int

	if len(extParams) == 0 {
		return ""
	}

	keys := make([]string, 0, len(extParams))
	for k, v := range extParams {
		keys = append(keys, k)
		l += len(keys) + len(v) + 2 // key=value, key=value,
	}
	l-- // 减去多的那个','
	// sort key
	sort.Strings(keys)

	builder := new(strings.Builder)
	builder.Grow(2 + l)
	builder.WriteString("|")
	l = 0
	for _, key := range keys {
		if l == 0 {
			l = 1
		} else {
			builder.WriteString(",")
		}
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(extParams[key])
	}
	builder.WriteString("|")
	return builder.String()
}
