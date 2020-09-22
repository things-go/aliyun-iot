package main

import (
	"log"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/dm"
)

func main() {
	client := mock.Init(mock.MetaTriad)
	ConfigLogGetTest(client) // done
	time.Sleep(time.Minute)
}

func ConfigLogGetTest(client *aiot.MQTTClient) {
	data, err := client.LinkThingConfigLogGet(mock.ProductKey, mock.DeviceName, dm.ConfigLogParam{}, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", data)
}
