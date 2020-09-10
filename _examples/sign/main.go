package main

import (
	"fmt"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
)

// just for test
const (
	testProductKey    = "a1QR3GD1Db3"
	testProductSecret = "mvngTYBlX9Z9l1V0"
	testDeviceName    = "dynamic"
	testDeviceSecret  = "9690f9da431078f105b7969b23e05762"
)

func main() {
	info, err := sign.NewMQTTSign().
		Generate(
			&infra.MetaInfo{
				ProductKey:    testProductKey,
				ProductSecret: testProductSecret,
				DeviceName:    testDeviceName,
				DeviceSecret:  testDeviceSecret,
			},
			infra.CloudRegionDomain{
				Region: infra.CloudRegionShangHai,
			})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", info)
}
