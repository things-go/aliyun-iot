package main

import (
	"fmt"

	"github.com/thinkgos/aliyun-iot/_examples/testmeta"
	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
)

func main() {
	meta := &infra.MetaInfo{
		ProductKey:    testmeta.ProductKey,
		ProductSecret: testmeta.ProductSecret,
		DeviceName:    testmeta.DeviceName,
		DeviceSecret:  testmeta.DeviceSecret,
	}
	crd := infra.CloudRegionDomain{
		Region: infra.CloudRegionShangHai,
	}

	info, err := sign.Generate(meta, crd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", info)
}
