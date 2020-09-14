package main

import (
	"encoding/json"
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
		dm.WithEnableNTP(),
		dm.WithEnableDesired(),
		dm.WithCallback(mockCb{}),
		dm.WithLogger(logger.New(log.New(os.Stdout, "mqtt --> ", log.LstdFlags), logger.WithEnable(true))),
	)

	dmClient.Underlying().Connect().Wait()
	if err = dmClient.Connect(); err != nil {
		panic(err)
	}

	//go DslTemplateTest()
	// go DesiredGetTest() // done
	// go DesiredDeleteTest()
	// go ConfigTest() // done
	// NTPTest() // done
	// DeviceInfoTest()  // done
	// ThingEventPost() // done
	for {
		time.Sleep(time.Second * 15)
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
		msg, err := entry.Wait(time.Second)
		if err != nil {
			log.Printf("error: %#v", err)
			continue
		}
		log.Printf("wait and got id: %d", msg.ID)
	}
}

// done
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

// done
func DeviceInfoTest() {
	tk, err := dmClient.ThingDeviceInfoUpdate(dm.DevNodeLocal,
		[]dmd.DeviceInfoLabel{{AttrKey: "attrKey", AttrValue: "attrValue"}})
	if err != nil {
		log.Println(err)
		return
	}
	_, err = tk.Wait(time.Second * 5)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute * 1)
	_, err = dmClient.ThingDeviceInfoDelete(dm.DevNodeLocal,
		[]dmd.DeviceLabelKey{{AttrKey: "attrKey"}})
	if err != nil {
		log.Println(err)
		return
	}
}

// done
func ConfigTest() {
	tk, err := dmClient.ThingConfigGet(dm.DevNodeLocal)
	if err != nil {
		log.Println(err)
		return
	}
	msg, err := tk.Wait(time.Second * 5)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(msg.Data.(dm.ConfigParamsData))
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

// done
func NTPTest() {
	err := dmClient.ExtNtpRequest()
	if err != nil {
		log.Println(err)
		return
	}
}

func DesiredGetTest() {
	tk, err := dmClient.ThingDesiredPropertyGet(dm.DevNodeLocal, []string{"temp", "humi"})
	if err != nil {
		log.Println(err)
		return
	}
	msg, err := tk.Wait(time.Second * 5)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%+v", msg)
	log.Printf("%+v", string(msg.Data.(json.RawMessage)))
}

func DesiredDeleteTest() {
	tk, err := dmClient.ThingDesiredPropertyDelete(dm.DevNodeLocal, "{}")
	if err != nil {
		log.Println(err)
		return
	}
	msg, err := tk.Wait(time.Second * 5)
	if err != nil {
		log.Printf("wait failed, %+v", err)
		return
	}

	log.Printf("%+v", msg.ID)
	log.Printf("%+v %+v", msg.Data)
}

type mockCb struct {
	dm.NopCb
}

func (sf mockCb) RRPCRequest(c *dm.Client, messageID, productKey, deviceName string, payload []byte) error {
	log.Println(string(payload))
	return nil
}
