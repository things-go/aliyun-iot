package main

import (
	"log"

	"github.com/things-go/aliyun-iot/_examples/mock"
	"github.com/things-go/aliyun-iot/dynamic"
	"github.com/things-go/aliyun-iot/infra"
)

var cli = dynamic.New()

func main() {
	var meta = infra.MetaTetrad{
		mock.ProductKey,
		mock.ProductSecret,
		"dynamic",
		"",
	}

	crd := infra.CloudRegionDomain{
		Region:       infra.CloudRegionShangHai,
		CustomDomain: "127.0.0.1:8080",
	}

	err := cli.RegisterCloud(&meta, crd)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", meta)
}
