package main

import (
	"log"
	"math/rand"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/dmd"
	"github.com/thinkgos/aliyun-iot/infra"
)

func main() {
	client := mock.Init()

	// go DiagPostTest(client) // done
	// go ConfigLogGetTest(client) // done
	// go dynamictslTest(client)
	// go DslTemplateTest(client) // done
	// go DesiredGetTest(client) // done
	// go DesiredDeleteTest(client) // done
	// go ConfigTest(client) // done
	go NTPTest(client) // done
	// go DeviceInfoTest(client) // done
	// go ThingEventPost(client) // done

	mock.ThingEventPropertyPost(client)
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

// done
func DeviceInfoTest(client *aiot.MQTTClient) {
	err := client.LinkThingDeviceInfoUpdate(mock.ProductKey, mock.DeviceName,
		[]dmd.DeviceInfoLabel{
			{AttrKey: "attrKey", AttrValue: "attrValue"},
		}, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute * 1)
	err = client.LinkThingDeviceInfoDelete(mock.ProductKey, mock.DeviceName,
		[]dmd.DeviceLabelKey{{AttrKey: "attrKey"}}, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
}

// done
func ConfigTest(client *aiot.MQTTClient) {
	cpd, err := client.LinkThingConfigGet(mock.ProductKey, mock.DeviceName, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("config: %+v", cpd)
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
	_, err := client.LinkThingDynamictslGet(mock.ProductKey, mock.DeviceName, time.Second*5)
	if err != nil {
		log.Println(err)
	}
}

// done
func NTPTest(client *aiot.MQTTClient) {
	err := client.ExtNtpRequest()
	if err != nil {
		log.Println(err)
		return
	}
}

func DesiredGetTest(client *aiot.MQTTClient) {
	data, err := client.LinkThingDesiredPropertyGet(mock.ProductKey, mock.DeviceName, []string{"temp", "humi"}, time.Second*5)
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

func ConfigLogGetTest(client *aiot.MQTTClient) {
	data, err := client.LinkThingConfigLogGet(mock.ProductKey, mock.DeviceName, dm.ConfigLogParam{}, time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", data)
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
