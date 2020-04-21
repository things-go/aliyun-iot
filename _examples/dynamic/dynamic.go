package main

import (
	"log"

	"github.com/thinkgos/aliyun-iot/dynamic"
	"github.com/thinkgos/aliyun-iot/infra"
)

const (
	productKey    = "a1iJcssSlPC"
	productSecret = "lw3QzKHNfh7XvOxO"
	deviceName    = "dynamic"
	//productKey    = "a1iJcssSlPC"
	//productSecret = "lw3QzKHNfh7XvOxO"
	//deviceName    = "1Myx6uC9RjJnucEraO2R"
)

func main() {
	meta := infra.MetaInfo{
		ProductKey:    productKey,
		ProductSecret: productSecret,
		DeviceName:    deviceName,
	}
	err := dynamic.Register2Cloud(&meta, infra.CloudRegionDomain{
		Region:       infra.CloudRegionShangHai,
		CustomDomain: "127.0.0.1:8080",
	})
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", meta)
}
