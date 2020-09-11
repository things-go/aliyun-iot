package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/testmeta"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/dmd"
	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
)

const (
	productKey    = "a1QR3GD1Db3"
	productSecret = ""
	deviceName    = "MPA19GT010070140"
	deviceSecret  = "ld9Xf2BtKGfdEC7G9nSMe1wYfgllvi3Q"
)

var dmClient *aiot.MQTTClient

func main() {
	signs, err :=
		sign.NewMQTTSign(sign.WithSDKVersion(sign.AlinkSDKVersion)).
			Generate(&infra.MetaInfo{
				ProductKey:    productKey,
				ProductSecret: productSecret,
				DeviceName:    deviceName,
				DeviceSecret:  deviceSecret},
				infra.CloudRegionDomain{Region: infra.CloudRegionShangHai})
	if err != nil {
		panic(err)
	}
	opts :=
		mqtt.NewClientOptions().
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

	dmClient = aiot.NewWithMQTT(testmeta.MetaInfo(), mqtt.NewClient(opts))
	dmClient.LogMode(true)

	dmClient.UnderlyingClient().Connect().Wait()
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
			_, err := dmClient.ThingEventPost(dm.DevNodeLocal, "tempAlarm", map[string]interface{}{
				"high": 1,
			})
			if err != nil {
				log.Printf("error: %#v", err)
			}
			time.Sleep(time.Second * 30)
		}

	}()

	for {
		_, err := dmClient.ThingEventPropertyPost(dm.DevNodeLocal, map[string]interface{}{
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
	if _, err := dmClient.ThingDeviceInfoUpdate(dm.DevNodeLocal,
		[]dmd.DevInfoLabelUpdate{
			{AttrKey: "attrKey", AttrValue: "attrValue"},
		}); err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute * 1)
	if _, err := dmClient.ThingDeviceInfoDelete(dm.DevNodeLocal,
		[]dmd.DevInfoLabelDelete{
			{AttrKey: "attrKey"},
		}); err != nil {
		log.Println(err)
		return
	}

}

func ConfigTest() {
	_, err := dmClient.ThingConfigGet(dm.DevNodeLocal)
	if err != nil {
		log.Println(err)
		return
	}
}

func DslTemplateTest() {
	_, err := dmClient.ThingDsltemplateGet(dm.DevNodeLocal)
	if err != nil {
		log.Println(err)
		return
	}
}

func dynamictslTest() {
	_, err := dmClient.ThingDynamictslGet(dm.DevNodeLocal)
	if err != nil {
		panic(err)
	}
}

func NTPTest() {
	err := dmClient.ExtNtpRequest()
	if err != nil {
		log.Println(err)
		return
	}
}
