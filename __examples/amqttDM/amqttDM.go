package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliIOT"
	"github.com/thinkgos/aliIOT/dm"
	"github.com/thinkgos/aliIOT/infra"
	"github.com/thinkgos/aliIOT/model"
	"github.com/thinkgos/aliIOT/sign"
)

const (
	productKey    = "a1QR3GD1Db3"
	productSecret = ""
	deviceName    = "MPA19GT010070140"
	deviceSecret  = "CsC7Gmb6EvDLOm8V40HLOQwFPdc3KCHT"
)

var manage *aliIOT.MQTTClient

func main() {
	signs, err := sign.NewMQTTSign().
		SetSDKVersion(infra.IOTSDKVersion).
		Generate(&sign.MetaInfo{
			ProductKey:    productKey,
			ProductSecret: productSecret,
			DeviceName:    deviceName,
			DeviceSecret:  deviceSecret}, sign.CloudRegionShangHai)
	if err != nil {
		panic(err)
	}
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("%s:%d", signs.HostName, signs.Port)).
		SetClientID(signs.ClientID).
		SetUsername(signs.UserName).
		SetPassword(signs.Password).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(cli mqtt.Client) {
			log.Println("mqtt client connection success")
		}).
		SetConnectionLostHandler(func(cli mqtt.Client, err error) {
			log.Println("mqtt client connection lost, ", err)
		})
	client := mqtt.NewClient(opts)

	dmopt := model.NewOption(productKey, deviceName, deviceSecret).
		SetEnableCache(true).
		Valid()
	manage = aliIOT.NewWithMQTT(dmopt, client)
	manage.LogMode(true)

	client.Connect().Wait()
	if err = manage.Connect(); err != nil {
		panic(err)
	}

	//go DslTemplateTest()
	//go ConfigTest()
	//go DeviceInfoTest()
	//go NTPTest()
	EventPostTest()
}

func EventPostTest() {
	go func() {
		for {
			err := manage.UpstreamThingEventPost(model.DevLocal, "tempAlarm", map[string]interface{}{
				"high": 1,
			})
			if err != nil {
				log.Printf("error: %#v", err)
			}
			time.Sleep(time.Second * 30)
		}

	}()

	for {
		err := manage.UpstreamThingEventPropertyPost(model.DevLocal, map[string]interface{}{
			"Temp":         rand.Intn(200),
			"Humi":         rand.Intn(100),
			"switchStatus": rand.Intn(1),
		})
		if err != nil {
			log.Printf("error: %#v", err)
		}
		time.Sleep(time.Second * 30)
	}
}

func DeviceInfoTest() {
	if err := manage.UpstreamThingDeviceInfoUpdate(model.DevLocal,
		[]dm.DevInfoLabelUpdate{
			{AttrKey: "attrKey", AttrValue: "attrValue"},
		}); err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute * 1)
	if err := manage.UpstreamThingDeviceInfoDelete(model.DevLocal,
		[]dm.DevInfoLabelDelete{
			{AttrKey: "attrKey"},
		}); err != nil {
		log.Println(err)
		return
	}

}

func ConfigTest() {
	err := manage.UpstreamThingConfigGet(model.DevLocal)
	if err != nil {
		log.Println(err)
		return
	}
}

func DslTemplateTest() {
	err := manage.UpstreamThingDsltemplateGet(model.DevLocal)
	if err != nil {
		log.Println(err)
		return
	}
}

func NTPTest() {
	err := manage.UpstreamExtNtpRequest()
	if err != nil {
		log.Println(err)
		return
	}
}

//
//type UserProc struct {
//	model.DevUserProc
//}
