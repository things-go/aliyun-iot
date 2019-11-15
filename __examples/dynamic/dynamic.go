package main

import (
	"log"

	"github.com/thinkgos/aliIOT/dynamic"
	"github.com/thinkgos/aliIOT/infra"
)

const (
	productKey    = "a1QR3GD1Db3"
	productSecret = "mvngTYBlX9Z9l1V0"
	deviceName    = "MPA19GT010070140"

//	deviceSecret  = "ld9Xf2BtKGfdEC7G9nSMe1wYfgllvi3Q"
)

func main() {
	meta := infra.MetaInfo{
		ProductKey:    productKey,
		ProductSecret: productSecret,
		DeviceName:    deviceName,
	}
	err := dynamic.Register2Cloud(&meta, infra.CloudRegionShangHai)
	if err != nil {
		panic(err)
	}
	log.Println(meta)
}
