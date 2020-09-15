package mock

import (
	"github.com/thinkgos/aliyun-iot/infra"
)

// just for test
const (
	ProductKey    = "a1QR3GD1Db3"
	ProductSecret = "mvngTYBlX9Z9l1V0"
	DeviceName    = "dynamic"
	DeviceSecret  = "9690f9da431078f105b7969b23e05762"
)

var MetaTetrad = infra.MetaTetrad{
	ProductKey,
	ProductSecret,
	DeviceName,
	DeviceSecret,
}

var MetaTriad = infra.MetaTriad{
	ProductKey,
	DeviceName,
	DeviceSecret,
}

type Instance struct {
	Temp         int `json:"temp"`
	Humi         int `json:"humi"`
	SwitchStatus int `json:"switchStatus"`
}

type Alarm struct {
	High int `json:"high"`
}
