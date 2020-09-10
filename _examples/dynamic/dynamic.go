package main

import (
	"log"

	"github.com/thinkgos/aliyun-iot/dynamic"
	"github.com/thinkgos/aliyun-iot/infra"
)

// just for test
const (
	testProductKey    = "a1QR3GD1Db3"
	testProductSecret = "mvngTYBlX9Z9l1V0"
	testDeviceName    = "dynamic"
)

func main() {
	meta := infra.MetaInfo{
		ProductKey:    testProductKey,
		ProductSecret: testProductSecret,
		DeviceName:    testDeviceName,
	}
	crd := infra.CloudRegionDomain{
		Region:       infra.CloudRegionShangHai,
		CustomDomain: "127.0.0.1:8080",
	}

	dclient := dynamic.New()

	err := dclient.RegisterCloud(&meta, crd)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", meta)
}
