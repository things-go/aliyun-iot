package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliIOT/model"
	"github.com/thinkgos/aliIOT/sign"
)

type wrapper struct {
	Client mqtt.Client
}

func (sf *wrapper) Publish(topic string, payload interface{}) error {
	return sf.Client.Publish(topic, 1, false, payload).Error()
}

func main() {
	signs, err := sign.NewMQTTSign().Generate(&sign.MetaInfo{
		ProductKey:    "a1QR3GD1Db3",
		ProductSecret: "mvngTYBlX9Z9l1V0",
		DeviceName:    "MPA19GT010070140",
		DeviceSecret:  "CsC7Gmb6EvDLOm8V40HLOQwFPdc3KCHT",
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
	wp := &wrapper{Client: client}
	manage := &model.Manager{
		Conn:       wp,
		ProductKey: "a1QR3GD1Db3",
		DeviceName: "MPA19GT010070140",
	}

	client.Subscribe("/sys/a1QR3GD1Db3/MPA19GT010070140/thing/event/property/post_reply", 1, func(client mqtt.Client, message mqtt.Message) {
		_ = model.ProcThingEventPostReply(manage, message.Topic(), message.Payload())
	})
	client.Connect().Wait()
	for {
		err = manage.UpstreamThingEventPropertyPost(map[string]interface{}{
			"AcTm": rand.Intn(255),
			"AeTm": rand.Intn(255),
		})
		if err != nil {
			log.Printf("error: %#v", err)
		} else {
			log.Printf("success")
		}
		time.Sleep(time.Second * 10)
	}
}
