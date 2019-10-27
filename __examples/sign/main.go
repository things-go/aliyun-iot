package main

import (
	"fmt"

	"github.com/thinkgos/aliIOT/infra"
	"github.com/thinkgos/aliIOT/sign"
)

const (
	productKey    = "a1QR3GD1Db3"
	productSecret = ""
	deviceName    = "MPA19GT010070140"
	deviceSecret  = "CsC7Gmb6EvDLOm8V40HLOQwFPdc3KCHT"
)

func main() {
	sig := sign.NewMQTTSign().SetSignMethod(infra.SignMethodHMACSHA256)
	signout, err := sig.Generate(&sign.MetaInfo{
		ProductKey:   productKey,
		DeviceName:   deviceName,
		DeviceSecret: deviceSecret,
	}, sign.CloudRegionShangHai)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", signout)
}
