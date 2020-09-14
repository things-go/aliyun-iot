package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/go-core-package/lib/logger"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/dmd"
	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
)

var dmClient *aiot.MQTTClient

func main() {
	meta := mock.MetaInfo()
	signs, err := sign.Generate(&meta, infra.CloudRegionDomain{Region: infra.CloudRegionShangHai})
	if err != nil {
		panic(err)
	}
	opts :=
		mqtt.NewClientOptions().
			AddBroker(signs.Addr()).
			SetClientID(signs.ClientIDWithExt()).
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

	dmClient = aiot.NewWithMQTT(
		mock.MetaInfo(),
		mqtt.NewClient(opts),
		dm.WithLogger(logger.New(log.New(os.Stdout, "mqtt --> ", log.LstdFlags), logger.WithEnable(true))),
	)

	dmClient.Underlying().Connect().Wait()
	if err = dmClient.Connect(); err != nil {
		panic(err)
	}

	//go DslTemplateTest()
	//go ConfigTest()
	//go DeviceInfoTest()
	//go NTPTest()

	for {
		time.Sleep(time.Second * 10)
		entry, err := dmClient.ThingEventPropertyPost(dm.DevNodeLocal,
			mock.Instance{
				rand.Intn(200),
				rand.Intn(100),
				rand.Intn(2),
			})
		if err != nil {
			log.Printf("error: %#v", err)
			continue
		}
		id, err := entry.WaitID(time.Second)
		if err != nil {
			log.Printf("error: %#v", err)
			continue
		}
		log.Printf("wait and got id: %d", id)
	}
}

func ThingEventPost() {
	for {
		_, err := dmClient.ThingEventPost(dm.DevNodeLocal, "tempAlarm", map[string]interface{}{
			"high": 1,
		})
		if err != nil {
			log.Printf("error: %#v", err)
		}
		time.Sleep(time.Second * 10)
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
