package main

import (
	"log"
	"time"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init(mock.MetaTriad)
	ConfigTest(client) // done
	time.Sleep(time.Minute * 5)
}

// TODO: 配置推送不正常
func ConfigTest(client *aiot.MQTTClient) {
	cpd, err := client.LinkThingConfigGet(mock.ProductKey, mock.DeviceName, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("config: %+v", cpd)
}
