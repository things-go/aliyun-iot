package main

import (
	"log"
	"time"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init(mock.MetaTriad)
	dynamictslTest(client)  // done
	DslTemplateTest(client) // done
}

func DslTemplateTest(client *aiot.MQTTClient) {
	data, err := client.LinkThingDsltemplateGet(mock.ProductKey, mock.DeviceName, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", string(data))
}

func dynamictslTest(client *aiot.MQTTClient) {
	data, err := client.LinkThingDynamictslGet(mock.ProductKey, mock.DeviceName, time.Second*5)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", string(data))
}
