package main

import (
	"log"
	"time"

	"github.com/thinkgos/aliyun-iot/_examples/mock"
)

// TODO: 子设备如何上传属性
func main() {
	client := mock.Init()
	err := client.Add(mock.SensorTriad)
	if err != nil {
		panic(err)
	}
	err = client.SubDeviceConnect(mock.SensorProductKey, mock.SensorDeviceName, true, 5*time.Second)
	if err != nil {
		panic(err)
	}

	for {
		err := client.LinkThingEventPropertyPost(mock.SensorProductKey, mock.SensorDeviceName, mock.SensorInstanceValue(), time.Second)
		if err != nil {
			log.Printf("error: %#v", err)
		}
		time.Sleep(time.Second * 30)
	}
}
