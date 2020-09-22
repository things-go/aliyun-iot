package main

import (
	"fmt"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
)

func main() {
	dmClient := aiot.NewWithHTTP(mock.MetaTriad)
	for {
		_, err := dmClient.ThingEventPropertyPost(mock.ProductKey, mock.DeviceName, mock.InstanceValue())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
