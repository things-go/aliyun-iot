package main

import (
	"fmt"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
)

func main() {
	info, err := sign.NewMQTTSign().Generate(&infra.MetaInfo{
		ProductKey:    "a1iJcssSlPC",
		ProductSecret: "lw3QzKHNfh7XvOxO",
		DeviceName:    "dynamic",
	}, infra.CloudRegionDomain{infra.CloudRegionShangHai, ""})
	if err != nil {
		panic(err)
	}

	fmt.Println(info)
}
