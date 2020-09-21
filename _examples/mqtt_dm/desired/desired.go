package main

import (
	"log"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init()

	DesiredGetTest(client)    // done
	DesiredDeleteTest(client) // done
}

func DesiredGetTest(client *aiot.MQTTClient) {
	data, err := client.LinkThingDesiredPropertyGet(mock.ProductKey, mock.DeviceName, []string{"memory_usage", "memory_total"}, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", string(data))
}

func DesiredDeleteTest(client *aiot.MQTTClient) {
	err := client.LinkThingDesiredPropertyDelete(mock.ProductKey, mock.DeviceName, "{}", time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
}
