package main

import (
	"log"

	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/dynamic"
	"github.com/thinkgos/aliyun-iot/infra"
)

func main() {
	meta := &infra.MetaInfo{
		ProductKey:    mock.ProductKey,
		ProductSecret: mock.ProductSecret,
		DeviceName:    mock.DeviceName,
	}
	crd := infra.CloudRegionDomain{
		Region:       infra.CloudRegionShangHai,
		CustomDomain: "127.0.0.1:8080",
	}

	dclient := dynamic.New()

	err := dclient.RegisterCloud(meta, crd)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", meta)
}
