package main

import (
	"log"

	"github.com/thinkgos/aliIOT/dynamic"
	"github.com/thinkgos/aliIOT/infra"
)

const (
	//productKey    = "a1iJcssSlPC"
	//productSecret = "lw3QzKHNfh7XvOxO"
	//deviceName    = "dynamic"
	productKey    = "a1iJcssSlPC"
	productSecret = "lw3QzKHNfh7XvOxO"
	deviceName    = "1Myx6uC9RjJnucEraO2R"

//	deviceSecret  = "ld9Xf2BtKGfdEC7G9nSMe1wYfgllvi3Q"
)

func main() {
	meta := infra.MetaInfo{
		ProductKey:    productKey,
		ProductSecret: productSecret,
		DeviceName:    deviceName,
	}
	err := dynamic.Register2Cloud(&meta, infra.CloudRegionDomain{
		Region:       infra.CloudRegionShangHai,
		CustomDomain: "127.168.1.14:8080",
	})
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", meta)
}
