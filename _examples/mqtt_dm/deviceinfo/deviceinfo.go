package main

import (
	"log"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init(mock.MetaTriad)
	DeviceInfoTest(client) // done
}

func DeviceInfoTest(client *aiot.MQTTClient) {
	err := client.LinkThingDeviceInfoUpdate(mock.ProductKey, mock.DeviceName,
		[]aiot.DeviceInfoLabel{
			{AttrKey: "attrKey", AttrValue: "attrValue"},
		}, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute * 1)
	err = client.LinkThingDeviceInfoDelete(mock.ProductKey, mock.DeviceName,
		[]aiot.DeviceLabelKey{{AttrKey: "attrKey"}}, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
}
