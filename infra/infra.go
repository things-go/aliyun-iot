package infra

import (
	"time"
)

// MetaPair 产品名与设备名
type MetaPair struct {
	ProductKey string
	DeviceName string
}

// MetaTriad 产品与设备三元组
type MetaTriad struct {
	ProductKey   string
	DeviceName   string
	DeviceSecret string
}

// MetaTetrad 产品与设备四元组
type MetaTetrad struct {
	ProductKey    string
	ProductSecret string
	DeviceName    string
	DeviceSecret  string
}

// Millisecond time.Time 转为 毫秒
func Millisecond(tm time.Time) int64 {
	return tm.Unix()*1000 + int64(tm.Nanosecond())/1000000
}

// Time 毫秒转time.Time
func Time(msec int64) time.Time {
	return time.Unix(msec/1000, (msec%1000)*1000000)
}
