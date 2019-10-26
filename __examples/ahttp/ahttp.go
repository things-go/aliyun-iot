package main

import (
	"fmt"

	"github.com/thinkgos/aliIOT/ahttp"
	"github.com/thinkgos/aliIOT/infra"
)

const (
	productKey    = "a1iJcssSlPC"
	productSecret = "lw3QzKHNfh7XvOxO"
	deviceName    = "dyncreg"
	deviceSecret  = "irqurH8zaIg1ChoeaBjLHiqBXEZnlVq8"
)

func main() {
	client := ahttp.New()
	client.
		SetDeviceMetaInfo(productKey, deviceName, deviceSecret).
		SetSignMethod(infra.SignMethodHMACSHA1)
	err := client.SendAuth()
	if err != nil {
		panic(err)
	}

	fmt.Println("auth success")

}
