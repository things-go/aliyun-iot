package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/_examples/mock"
	"github.com/things-go/aliyun-iot/http"
	"github.com/things-go/aliyun-iot/infra"
	"github.com/things-go/aliyun-iot/logger"
	"github.com/things-go/aliyun-iot/uri"
)

// 采用物模型测试
func main() {
	var err error
	l := logger.New(log.New(os.Stdout, "mqtt --> ", log.LstdFlags), logger.WithEnable(true))

	client := http.New(mock.MetaTriad, http.WithLogger(l))

	_uri := uri.URI(uri.SysPrefix, uri.ThingEventPropertyPost, mock.ProductKey, mock.DeviceName)
	bPayload, err := json.Marshal(aiot.Request{
		ID:      uint(rand.Int()),
		Version: aiot.DefaultVersion,
		Params:  mock.InstanceValue(),
		Method:  infra.MethodEventPropertyPost,
	})
	for {
		err = client.Publish(_uri, 1, bPayload)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 15)
	}
}
