// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package sign 实现MQTT设备签名
// see https://help.aliyun.com/document_detail/73742.html?spm=a2c4g.11186623.6.599.76216eebzbvrYq
package sign

import (
	"errors"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/thinkgos/aliyun-iot/infra"
)

// default defined
const (
	fixedTimestamp = 2524608000000
	alinkVersion   = "20"

	// config method MQTT设备签名只支持以下签名方法
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

// Sign 签名后的信息
type Sign struct {
	HostName  string
	Port      uint16
	ClientID  string
	extParams string // clientID的扩展参数
	UserName  string // deviceName & productKey
	Password  string
}

// Addr address like host:port
func (ms *Sign) Addr() string {
	return net.JoinHostPort(ms.HostName, strconv.FormatUint(uint64(ms.Port), 10))
}

// ClientIDWithExt 登录使用的 client + "| 扩展参数 |"
func (ms *Sign) ClientIDWithExt() string {
	return ms.ClientID + ms.extParams
}

// config MQTT签名主要配置
type config struct {
	enableTLS bool              // 使能tls
	enableDM  bool              // 使能物模型
	extParams map[string]string // clientID扩展参数
}

// New 新建一个签名

// Generate 根据MetaInfo和region生成签名
// 默认不支持PreAUTH
// 默认安全模式为SecureModeTcpDirectPlain)
// 默认使能物模型
// 默认hmacmd5签名加密
// 默认sdk版本为 SDKVersion
// TODO: 支持tls
func Generate(triad infra.MetaTriad, crd infra.CloudRegionDomain, opts ...Option) (*Sign, error) {
	if crd.Region == infra.CloudRegionCustom && crd.CustomDomain == "" {
		return nil, errors.New("invalid custom domain")
	}
	ms := &config{
		enableDM: true,
		extParams: map[string]string{
			"timestamp":  strconv.FormatUint(fixedTimestamp, 10), // 表示当前时间的毫秒值,可以不传递
			"securemode": modeTCPDirectPlain,
			"signmethod": hmacsha256,
			"lan":        "Golang",
			"v":          alinkVersion,
		},
	}
	for _, opt := range opts {
		opt(ms)
	}

	// setup ClientID,Password
	clientID, pwd := infra.CalcSign(ms.extParams["signmethod"], triad, fixedTimestamp)
	info := &Sign{
		Port:      1883,
		ClientID:  clientID,
		extParams: encodeExtParam(ms.extParams),
		UserName:  triad.DeviceName + "&" + triad.ProductKey,
		Password:  pwd,
	}

	domain := crd.CustomDomain
	if crd.Region != infra.CloudRegionCustom {
		domain = infra.MQTTCloudDomain[crd.Region]
	}
	// setup HostName
	info.HostName = triad.ProductKey + "." + domain
	// setup Port
	if ms.enableTLS {
		info.Port = 443
	}
	return info, nil
}

// encodeExtParam 根据extParams编码扩展字符串
func encodeExtParam(extParams map[string]string) string {
	if len(extParams) == 0 {
		return ""
	}

	l := 0
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
