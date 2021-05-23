package main

import (
	"log"
	"time"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/_examples/mock"
	"github.com/things-go/aliyun-iot/infra"
)

func main() {
	client := mock.Init(mock.MetaTriad)
	DiagPostTest(client) // done
}

func DiagPostTest(client *aiot.MQTTClient) {
	err := client.LinkThingDiagPost(mock.ProductKey, mock.DeviceName, aiot.P{
		Wifi: aiot.Wifi{
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
