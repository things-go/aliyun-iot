package main

import (
	"log"

	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
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
