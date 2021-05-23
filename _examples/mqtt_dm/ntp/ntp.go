package main

import (
	"log"
	"time"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init(mock.MetaTriad)
	NTPTest(client) // done
	time.Sleep(time.Second * 2)
}

func NTPTest(client *aiot.MQTTClient) {
	err := client.ExtNtpRequest()
	if err != nil {
		log.Println(err)
		return
	}
}
