package main

import (
	"fmt"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/ahttp"
)

func main() {
	dmClient := aiot.New(mock.MetaTriad,
		ahttp.New(mock.MetaTriad),
		aiot.WithWork(aiot.WorkOnHTTP))
	for {
		_, err := dmClient.ThingEventPropertyPost(mock.ProductKey, mock.DeviceName, mock.InstanceValue())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
