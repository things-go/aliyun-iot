package main

import (
	"log"
	"time"

	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/dm"
)

func main() {
	client := mock.Init()
	err := client.Create(mock.SensorTriad)
	if err != nil {
		panic(err)
	}
	err = client.LinkThingSubRegister(mock.SensorProductKey, mock.SensorDeviceName, time.Second*5)
	if err != nil {
		log.Println(err)
	}
	ds, err := client.DeviceSecret(mock.SensorProductKey, mock.SensorDeviceName)
	if err != nil {
		log.Println(err)
	}
	log.Println(ds)

	err = client.LinkThingTopoAdd(mock.SensorProductKey, mock.SensorDeviceName, time.Second*5)
	if err != nil {
		log.Println(err)
	}

	err = client.LinkExtCombineLogin(dm.CombinePair{
		mock.SensorProductKey,
		mock.SensorDeviceName,
		true,
	}, time.Second*5)
	if err != nil {
		log.Println(err)
	}

	err = client.LinkExtCombineLogin(dm.CombinePair{
		mock.SensorProductKey,
		mock.SensorDeviceName,
		true,
	}, time.Second*5)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(time.Second * 30)
}
