package main

import (
	"fmt"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/ahttp"
)

func main() {
	client := aiot.New(mock.MetaTriad, ahttp.New(mock.MetaTriad), aiot.WithMode(aiot.ModeHTTP))
	for {
		_, err := client.ThingEventPropertyPost(mock.ProductKey, mock.DeviceName, mock.InstanceValue())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
