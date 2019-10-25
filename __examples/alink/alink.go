package alink

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

func main() {
	signs, err := sign.NewMQTTSign().SetSDKVersion(infra.IOTSDKVersion).Generate(&sign.MetaInfo{
		ProductKey:    "a1iJcssSlPC",
		ProductSecret: "mvngTYBlX9Z9l1V0",
		DeviceName:    "dyncreg",
		DeviceSecret:  "irqurH8zaIg1ChoeaBjLHiqBXEZnlVq8",
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
	manage := aliIOT.NewWithMQTT(
		"a1iJcssSlPC",
		"dyncreg",
		"bOtuBFIZgDdNb2RamMRh7eaOn6VoyurP",
		client)

	client.Connect().Wait()

	_ = manage.Subscribe(manage.URIService(model.URISysPrefix, model.URIThingEventPropertyPostReply), model.ProcThingEventPostReply)

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
