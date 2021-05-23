package main

import (
	"log"
	"math/rand"
	"time"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init(mock.MetaTriad)
	ThingEventPost(client) // done
}

// done
func ThingEventPost(client *aiot.MQTTClient) {
	for {
		err := client.LinkThingEventPost(mock.ProductKey, mock.DeviceName, "alarm",
			mock.Alarm{
				rand.Intn(2),
			}, time.Second*5)
		if err != nil {
			log.Printf("error: %#v", err)
		}
		time.Sleep(time.Second * 15)
	}
}
