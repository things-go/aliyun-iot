package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aiot"
	"github.com/thinkgos/aiot/dm"
	"github.com/thinkgos/aiot/dmd"
	"github.com/thinkgos/aiot/infra"
	"github.com/thinkgos/aiot/sign"
)

const (
	productKey    = "a1QR3GD1Db3"
	productSecret = ""
	deviceName    = "MPA19GT010070140"
	deviceSecret  = "ld9Xf2BtKGfdEC7G9nSMe1wYfgllvi3Q"
)

var dmClient *aiot.MQTTClient

func main() {
	signs, err := sign.NewMQTTSign().
		SetSDKVersion(infra.IOTSDKVersion).
		Generate(&infra.MetaInfo{
			ProductKey:    productKey,
			ProductSecret: productSecret,
			DeviceName:    deviceName,
			DeviceSecret:  deviceSecret,
		}, infra.CloudRegionDomain{
			Region: infra.CloudRegionShangHai,
		})
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

	dmopt := dm.NewConfig(productKey, deviceName, deviceSecret).
		Valid()
	dmClient = aiot.NewWithMQTT(dmopt, client)
	dmClient.LogMode(true)

	client.Connect().Wait()
	if err = dmClient.AlinkConnect(); err != nil {
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
			err := dmClient.AlinkTriggerEvent(dm.DevNodeLocal, "tempAlarm", map[string]interface{}{
				"high": 1,
			})
			if err != nil {
				log.Printf("error: %#v", err)
			}
			time.Sleep(time.Second * 30)
		}

	}()

	for {
		err := dmClient.AlinkReport(dm.MsgTypeEventPropertyPost, dm.DevNodeLocal, map[string]interface{}{
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
	if err := dmClient.AlinkReport(dm.MsgTypeDeviceInfoUpdate, dm.DevNodeLocal,
		[]dmd.DevInfoLabelUpdate{
			{AttrKey: "attrKey", AttrValue: "attrValue"},
		}); err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute * 1)
	if err := dmClient.AlinkReport(dm.MsgTypeDeviceInfoDelete, dm.DevNodeLocal,
		[]dmd.DevInfoLabelDelete{
			{AttrKey: "attrKey"},
		}); err != nil {
		log.Println(err)
		return
	}

}

func ConfigTest() {
	err := dmClient.AlinkQuery(dm.MsgTypeConfigGet, dm.DevNodeLocal)
	if err != nil {
		log.Println(err)
		return
	}
}

func DslTemplateTest() {
	err := dmClient.AlinkQuery(dm.MsgTypeDsltemplateGet, dm.DevNodeLocal)
	if err != nil {
		log.Println(err)
		return
	}
}

func dynamictslTest() {
	err := dmClient.AlinkQuery(dm.MsgTypeDynamictslGet, dm.DevNodeLocal)
	if err != nil {
		panic(err)
	}
}

func NTPTest() {
	err := dmClient.AlinkQuery(dm.MsgTypeExtNtpRequest, dm.DevNodeLocal)
	if err != nil {
		log.Println(err)
		return
	}
}
