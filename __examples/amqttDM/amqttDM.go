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

var manage *model.Manager

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

	//_ = manage.Subscribe(manage.URIServiceSelf(model.URISysPrefix, model.URIThingServiceRequestMultiWildcard2),
	//	model.ProcThingServiceRequest)
	//_ = manage.Subscribe(manage.URIServiceSelf(model.URISysPrefix, model.URIRRPCRequestSingleWildcard),
	//	model.ProcRRPCRequest)

	//go DslTemplateTest()
	//go ConfigTest()
	//go DeviceInfoTest()
	go NTPTest()
	EventPostTest()
}

func EventPostTest() {
	_ = manage.Subscribe(manage.URIServiceSelf(model.URISysPrefix, model.URIThingEventPostReplySingleWildcard),
		model.ProcThingEventPostReply)

	go func() {
		for {
			err := manage.UpstreamThingEventPost(model.DevSelf, "tempAlarm", map[string]interface{}{
				"high": 1,
			})
			if err != nil {
				log.Printf("error: %#v", err)
			}
			time.Sleep(time.Second * 30)
		}

	}()

	for {
		err := manage.UpstreamThingEventPropertyPost(model.DevSelf, map[string]interface{}{
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
	// 设备标签
	_ = manage.Subscribe(manage.URIServiceSelf(model.URISysPrefix, model.URIThingDeviceInfoUpdateReply),
		model.ProcThingDeviceInfoUpdateReply)
	_ = manage.Subscribe(manage.URIServiceSelf(model.URISysPrefix, model.URIThingDeviceInfoDeleteReply),
		model.ProcThingDeviceInfoDeleteReply)

	if err := manage.UpstreamThingDeviceInfoUpdate(model.DevSelf,
		[]dm.DevInfoLabelUpdate{
			{AttrKey: "attrKey", AttrValue: "attrValue"},
		}); err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute * 1)
	if err := manage.UpstreamThingDeviceInfoDelete(model.DevSelf,
		[]dm.DevInfoLabelDelete{
			{AttrKey: "attrKey"},
		}); err != nil {
		log.Println(err)
		return
	}

}

func ConfigTest() {
	// 设备配置
	_ = manage.Subscribe(manage.URIServiceSelf(model.URISysPrefix, model.URIThingConfigGetReply),
		model.ProcThingConfigGetReply)
	_ = manage.Subscribe(manage.URIServiceSelf(model.URISysPrefix, model.URIThingConfigPush),
		model.ProcThingConfigPush)

	err := manage.UpstreamThingConfigGet(model.DevSelf)
	if err != nil {
		log.Println(err)
		return
	}
}

func DslTemplateTest() {
	_ = manage.Subscribe(manage.URIServiceSelf(model.URISysPrefix, model.URIThingDslTemplateGetReply),
		model.ProcThingDsltemplateGetReply)
	err := manage.UpstreamThingDsltemplateGet(model.DevSelf)
	if err != nil {
		log.Println(err)
		return
	}
}

func NTPTest() {
	_ = manage.Subscribe(manage.URIServiceSelf(model.URIExtNtpPrefix, model.URINtpResponse),
		model.ProcExtNtpResponse)

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
