package main

import (
	"fmt"
	"math/rand"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
)

func main() {
	dmClient := aiot.NewWithHTTP(mock.MetaTriad)
	for {
		_, err := dmClient.ThingEventPropertyPost(mock.ProductKey, mock.DeviceName,
			mock.Instance{
				rand.Intn(200),
				rand.Intn(100),
				rand.Intn(2),
			})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
