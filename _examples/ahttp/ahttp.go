package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/thinkgos/aliyun-iot/_examples/mock"
	"github.com/thinkgos/aliyun-iot/ahttp"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/infra"
	uri2 "github.com/thinkgos/aliyun-iot/uri"
)

// 采用物模型测试
func main() {
	var err error

	client := ahttp.New(mock.MetaInfo())

	uri := uri2.URI(uri2.SysPrefix, uri2.ThingEventPropertyPost, mock.ProductKey, mock.DeviceName)
	bPayload, err := json.Marshal(
		dm.Request{
			ID:      uint(rand.Int()),
			Version: dm.Version,
			Params: map[string]interface{}{
				"Temp":         rand.Intn(200),
				"Humi":         rand.Intn(100),
				"switchStatus": rand.Intn(2),
			},
			Method: infra.MethodEventPropertyPost,
		})
	for {
		err = client.Publish(uri, 1, bPayload)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
