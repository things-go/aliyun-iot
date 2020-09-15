package main

import (
	"fmt"

	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
)

func main() {
	crd := infra.CloudRegionDomain{
		Region: infra.CloudRegionShangHai,
	}

	info, err := sign.Generate(mock.MetaTriad, crd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", info)
}
