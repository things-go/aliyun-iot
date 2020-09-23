package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/thinkgos/go-core-package/lib/logger"

	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/ahttp"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// 采用物模型测试
func main() {
	var err error
	l := logger.New(log.New(os.Stdout, "mqtt --> ", log.LstdFlags), logger.WithEnable(true))

	client := ahttp.New(mock.MetaTriad, ahttp.WithLogger(l))

	_uri := uri.URI(uri.SysPrefix, uri.ThingEventPropertyPost, mock.ProductKey, mock.DeviceName)
	bPayload, err := json.Marshal(
		dm.Request{
			ID:      uint(rand.Int()),
			Version: dm.DefaultVersion,
			Params:  mock.InstanceValue(),
			Method:  infra.MethodEventPropertyPost,
		})
	for {
		err = client.Publish(_uri, 1, bPayload)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
