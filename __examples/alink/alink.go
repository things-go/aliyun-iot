package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliIOT"
	"github.com/thinkgos/aliIOT/infra"
	"github.com/thinkgos/aliIOT/model"
	"github.com/thinkgos/aliIOT/sign"
)

const (
	productKey    = "a1iJcssSlPC"
	productSecret = "mvngTYBlX9Z9l1V0"
	deviceName    = "dyncreg"
	deviceSecret  = "irqurH8zaIg1ChoeaBjLHiqBXEZnlVq8"
)

func main() {
	signs, err := sign.NewMQTTSign().SetSDKVersion(infra.IOTSDKVersion).Generate(&sign.MetaInfo{
		ProductKey:    productKey,
		ProductSecret: productSecret,
		DeviceName:    deviceName,
		DeviceSecret:  deviceSecret,
	}, sign.CloudRegionShangHai)
	if err != nil {
		panic(err)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", signs.HostName, signs.Port))
	opts.SetClientID(signs.ClientID)
	opts.SetUsername(signs.UserName)
	opts.SetPassword(signs.Password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(func(cli mqtt.Client) {
		log.Println("mqtt client connection success")
	})
	opts.SetConnectionLostHandler(func(cli mqtt.Client, err error) {
		log.Println("mqtt client connection lost, ", err)
	})
	client := mqtt.NewClient(opts)

	dmopt := model.NewOption(productKey, deviceName, deviceSecret).Valid()
	manage := aliIOT.NewWithMQTT(dmopt, client)

	client.Connect().Wait()

	_ = manage.Subscribe(manage.URIServiceItself(model.URISysPrefix, model.URIThingEventPropertyPostReply), model.ProcThingEventPostReply)

	for {
		err = manage.UpstreamThingEventPropertyPost(model.DevItself, map[string]interface{}{
			"Temp":         rand.Intn(200),
			"Humi":         rand.Intn(100),
			"SwitchStatus": rand.Intn(1),
		})
		if err != nil {
			log.Printf("error: %#v", err)
		} else {
			log.Printf("success")
		}
		time.Sleep(time.Second * 10)
	}
}
