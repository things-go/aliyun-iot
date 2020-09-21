package main

import (
	"log"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/infra"
)

func main() {
	client := mock.Init()
	DiagPostTest(client) // done
}

func DiagPostTest(client *aiot.MQTTClient) {
	err := client.LinkThingDiagPost(mock.ProductKey, mock.DeviceName, dm.P{
		Wifi: dm.Wifi{
			Rssi:     100,
			Snr:      10,
			Per:      2,
			ErrStats: "",
		},
		Time: infra.Millisecond(time.Now()),
	}, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
}