// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package sign 实现MQTT设备签名
// see https://help.aliyun.com/document_detail/73742.html?spm=a2c4g.11186623.6.599.76216eebzbvrYq
package sign

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/thinkgos/go-core-package/extcert"

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
	// itls domain
	// itlsDomain = "x509.itls.cn-shanghai.aliyuncs.com"
)

// SecureMode 支持的安全模型
const (
	SecureModeNoPreRegistration = "-2" // 一型一密免预注册
	SecureModeTLSGuider         = "-1"
	SecureModeTLSDirect         = "2" // TLS直连
	SecureModeTCPDirectPlain    = "3" // TCP直连
	SecureModeITLSDNSID2        = "8" // ITLS, ID2方式
)

// Sign 签名后的信息
type Sign struct {
	Addr      string // broker addr
	HostName  string
	Port      uint16
	ClientID  string
	extParams string // clientID的扩展参数
	UserName  string // deviceName & productKey
	Password  string
}

// ClientIDWithExt 登录使用的 client + "| 扩展参数 |"
func (ms *Sign) ClientIDWithExt() string {
	return ms.ClientID + ms.extParams
}

// config MQTT签名主要配置
type config struct {
	secureMode  string            // 安全模式
	deviceToken string            // only use on SecureModeNoPreRegistration
	method      string            // 签名方法
	enableDM    bool              // 使能物模型
	extRRPC     bool              // 物模型下,支持扩展RRPC
	port        uint16            // 端口,默认为1883
	timestamp   int64             // 表示当前时间的毫秒值,可以不传递, 默认 fixedTimestamp
	extParams   map[string]string // clientID扩展参数
}

// Generate 根据MetaTriad和region生成签名
// 默认不支持PreAUTH
// 默认安全模式为SecureModeTcpDirectPlain)
// 默认使能物模型
// 默认固定时间戳
// 默认hmacsha256签名加密
// TODO: itls x509方式
func Generate(triad infra.MetaTriad, crd infra.CloudRegionDomain, opts ...Option) (*Sign, error) {
	if crd.Region == infra.CloudRegionCustom && crd.CustomDomain == "" {
		return nil, errors.New("invalid custom domain")
	}
	c := &config{
		SecureModeTCPDirectPlain,
		"",
		hmacsha256,
		true,
		false,
		1883,
		fixedTimestamp,
		map[string]string{
			"securemode": SecureModeTCPDirectPlain,
			"signmethod": hmacsha256,
			"gw":         "0",
			"ext":        "0",
			"lan":        "Golang",
		},
	}
	for _, opt := range opts {
		opt(c)
	}

	c.extParams["timestamp"] = strconv.FormatInt(c.timestamp, 10)
	if c.enableDM && c.extRRPC {
		c.extParams["ext"] = "1"
	}
	if !c.enableDM {
		c.extParams["v"] = alinkVersion
		delete(c.extParams, "gw")
		delete(c.extParams, "ext")
	}

	var enableTLS bool // 使能tls
	switch c.secureMode {
	case SecureModeNoPreRegistration:
		panic("feature not support")
	case SecureModeTLSGuider, SecureModeTLSDirect, SecureModeITLSDNSID2:
		enableTLS = true
	default: // SecureModeTCPDirectPlain
		c.secureMode = SecureModeTCPDirectPlain
		enableTLS = false
	}
	c.extParams["securemode"] = c.secureMode

	schema := "tcp://"
	if enableTLS {
		schema = "tls://"
	}

	// setup HostName
	hostname := triad.ProductKey + "."
	if crd.Region == infra.CloudRegionCustom {
		hostname += crd.CustomDomain
	} else {
		hostname += infra.MQTTCloudDomain[crd.Region]
	}

	addr := schema + net.JoinHostPort(hostname, strconv.Itoa(int(c.port)))
	username := triad.DeviceName + "&" + triad.ProductKey

	if c.secureMode == SecureModeNoPreRegistration {
		c.extParams["authType"] = "connwl"
		delete(c.extParams, "timestamp")
		delete(c.extParams, "signmethod")
		return &Sign{
			addr,
			hostname,
			c.port,
			infra.ClientID(triad.ProductKey, triad.DeviceName),
			encodeExtParam(c.extParams),
			username,
			c.deviceToken,
		}, nil
	}

	switch c.method {
	case hmacsha1, hmacmd5, hmacsha256:
	default:
		c.method = hmacsha256
	}
	c.extParams["signmethod"] = c.method
	// setup ClientID,Password
	clientID, pwd := infra.CalcSign(c.method, triad, c.timestamp)
	return &Sign{
		addr,
		hostname,
		c.port,
		clientID,
		encodeExtParam(c.extParams),
		username,
		pwd,
	}, nil
}

// encodeExtParam 根据extParams编码扩展字符串
func encodeExtParam(extParams map[string]string) string {
	if len(extParams) == 0 {
		return ""
	}

	// key=value, key=value
	n := len(extParams)*2 - 1 // ','和'='的个数
	keys := make([]string, 0, len(extParams))
	for k, v := range extParams {
		keys = append(keys, k)
		n += len(k) + len(v)
	}
	sort.Strings(keys) // sort key

	builder := strings.Builder{}
	builder.Grow(2 + n)
	builder.WriteString("|")
	for _, key := range keys {
		if builder.Len() > 1 {
			builder.WriteString(",")
		}
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(extParams[key])
	}
	builder.WriteString("|")
	return builder.String()
}

// NewTLSConfig new tls config from ca file
// 如果ca有"base64://"前缀,直接解析后面的字符串,否则认为这是个ca为文件名
func NewTLSConfig(ca string) (*tls.Config, error) {
	bs, err := extcert.LoadCrt(ca)
	if err != nil {
		return nil, err
	}
	return TLSConfig(bs)
}

// TLSConfig tls config
func TLSConfig(cacertPem []byte) (*tls.Config, error) {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(cacertPem)

	//	// Import client certificate/key pair
	//	cert, err := tls.X509KeyPair([]byte(cert_pem), []byte(key_pem))
	//	if err != nil {
	//		return nil, err
	//	}

	//	// Just to print out the client certificate..
	//	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	//	if err != nil {
	//		return nil, err
	//	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		//		Certificates: []tls.Certificate{cert},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
	}, nil
}
