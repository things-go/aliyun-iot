// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package infra

import (
	"time"
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

// Millisecond time.Time 转为 毫秒
func Millisecond(tm time.Time) int64 {
	return tm.Unix()*1000 + int64(tm.Nanosecond())/1000000
}

// Time 毫秒转time.Time
func Time(msec int64) time.Time {
	return time.Unix(msec/1000, (msec%1000)*1000000)
}
