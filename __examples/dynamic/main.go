package main

import (
	"log"

	"github.com/thinkgos/aliIOT/dynamic"
	"github.com/thinkgos/aliIOT/infra"
)

func main() {
	meta := dynamic.MetaInfo{
		ProductKey:    "a1iJcssSlPC",
		ProductSecret: "lw3QzKHNfh7XvOxO",
		DeviceName:    "dyncreg",
	}
	err := dynamic.Register2Cloud(&meta, infra.CloudRegionShangHai)
	if err != nil {
		panic(err)
	}
	log.Println(meta)
}
