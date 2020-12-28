// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package infra

import (
	"strconv"
	"time"

	"github.com/thinkgos/go-core-package/extime"
	"github.com/thinkgos/go-core-package/lib/algo"
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
	return extime.Millisecond(tm)
}

// Time 毫秒转time.Time
func Time(msec int64) time.Time {
	return extime.Time(msec)
}

// CalcSign 返回clientID和加签后的值
// sign,mqtt 可以采用 hmacmd5,hmacsha1,hmacsha256
// http 支持 hmacmd5,hmacsha1
// timestamp: 时间戳,单位ms
func CalcSign(method string, meta MetaTriad, timestamp int64) (string, string) {
	clientID := ClientID(meta.ProductKey, meta.DeviceName)
	source := "clientId" + clientID +
		"deviceName" + meta.DeviceName +
		"productKey" + meta.ProductKey +
		"timestamp" + strconv.FormatInt(timestamp, 10)
	return clientID, algo.Hmac(method, meta.DeviceSecret, source)
}
