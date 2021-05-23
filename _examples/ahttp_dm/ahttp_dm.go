package main

import (
	"fmt"
	"log"
	"os"
	"time"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/_examples/mock"
	"github.com/things-go/aliyun-iot/ahttp"
	"github.com/things-go/aliyun-iot/logger"
)

func main() {
	l := logger.New(log.New(os.Stdout, "mqtt --> ", log.LstdFlags), logger.WithEnable(true))
	client := aiot.New(mock.MetaTriad, ahttp.New(mock.MetaTriad), aiot.WithMode(aiot.ModeHTTP), aiot.WithLogger(l))
	for {
		_, err := client.ThingEventPropertyPost(mock.ProductKey, mock.DeviceName, mock.InstanceValue())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
