package main

import (
	"log"

	"github.com/things-go/aliyun-iot/_examples/mock"
	"github.com/things-go/aliyun-iot/infra"
	"github.com/things-go/aliyun-iot/sign"
)

func main() {
	crd := infra.CloudRegionDomain{
		Region: infra.CloudRegionShangHai,
	}

	info, err := sign.Generate(mock.MetaTriad, crd, sign.WithTimestamp())
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", info)
}
