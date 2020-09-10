package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/thinkgos/aliyun-iot/_examples/testmeta"
	"github.com/thinkgos/aliyun-iot/ahttp"
	"github.com/thinkgos/aliyun-iot/dm"
)

// 采用物模型测试
func main() {
	var err error

	client := ahttp.New(ahttp.WithDeviceMetaInfo(testmeta.ProductKey, testmeta.DeviceName, testmeta.DeviceSecret))
	client.LogMode(true)

	uri := dm.URICOAPHTTPPrePrefix + dm.URIService(dm.URISysPrefix, dm.URIThingEventPropertyPost, testmeta.ProductKey, testmeta.DeviceName)
	bPayload, err := json.Marshal(
		dm.Request{
			ID:      rand.Int(),
			Version: dm.Version,
			Params: map[string]interface{}{
				"Temp":         rand.Intn(200),
				"Humi":         rand.Intn(100),
				"switchStatus": rand.Intn(2),
			},
			Method: dm.MethodEventPropertyPost,
		})
	for {
		err = client.Publish(uri, bPayload)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
