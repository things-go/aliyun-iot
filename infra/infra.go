// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package infra

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"math/bits"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// MetaPair 产品名与设备名
type MetaPair struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// MetaTriad 产品与设备三元组
type MetaTriad struct {
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
}

// MetaTetrad 产品与设备四元组
type MetaTetrad struct {
	ProductKey    string `json:"productKey"`
	ProductSecret string `json:"productSecret"`
	DeviceName    string `json:"deviceName"`
	DeviceSecret  string `json:"deviceSecret"`
}

// ClientID to client id like {pk}.{dn}
func ClientID(pk, dn string) string {
	return pk + "." + dn
}

// Millisecond time.Time 转为 毫秒
func Millisecond(tm time.Time) int64 {
	return tm.UnixNano() / int64(time.Millisecond)
}

// Time 毫秒转time.Time
func Time(msec int64) time.Time {
	return time.Unix(msec/1000, (msec%1000)*int64(time.Millisecond))
}

// CalcSign 计算签名, 返回clientID和加签后的值
// sign,mqtt 支持 hmacmd5, hmacsha1, hmacsha256
// http 支持 hmacmd5, hmacsha1
// timestamp: 时间戳,单位: ms
func CalcSign(method string, meta MetaTriad, timestamp int64) (string, string) {
	clientID := ClientID(meta.ProductKey, meta.DeviceName)
	source := "clientId" + clientID +
		"deviceName" + meta.DeviceName +
		"productKey" + meta.ProductKey +
		"timestamp" + strconv.FormatInt(timestamp, 10)
	return clientID, Hmac(method, meta.DeviceSecret, source)
}

// LoadCrt 加载tls cert
// 如果cert有"base64://"前缀,直接解析后面的字符串,否则认为这是个cert文件名
func LoadCrt(cert string) ([]byte, error) {
	if strings.HasPrefix(cert, "base64://") {
		return base64.StdEncoding.DecodeString(cert[len("base64://"):])
	}
	return ioutil.ReadFile(cert)
}

var defaultAlphabet = []byte("QWERTYUIOPLKJHGFDSAZXCVBNMabcdefghijklmnopqrstuvwxyz")

// RandAlphabet rand alpha with give length(只包含字母)
func RandAlphabet(length int) string {
	b := randBytes(length, defaultAlphabet)
	return *(*string)(unsafe.Pointer(&b))
}

func randBytes(length int, alphabets []byte) []byte {
	b := make([]byte, length)
	if _, err := cryptoRand.Read(b); err == nil {
		for i, v := range b {
			b[i] = alphabets[v%byte(len(alphabets))]
		}
		return b
	}

	bn := bits.Len(uint(len(alphabets)))
	mask := int64(1)<<bn - 1
	max := 63 / bn

	// A rand.Int63() generates 63 random bits, enough for alphabets letters!
	for i, cache, remain := 0, rand.Int63(), max; i < length; {
		if remain == 0 {
			cache, remain = rand.Int63(), max
		}
		if idx := int(cache & mask); idx < len(alphabets) {
			b[i] = alphabets[idx]
			i++
		}
		cache >>= bn
		remain--
	}
	return b
}
