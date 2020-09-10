package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/thinkgos/aliyun-iot/ahttp"
	"github.com/thinkgos/aliyun-iot/dm"
)

const (
	productKey    = "a1QR3GD1Db3"
	productSecret = "mvngTYBlX9Z9l1V0"
	deviceName    = "dynamic"
	deviceSecret  = "9690f9da431078f105b7969b23e05762"
)

// 透传
func main() {
	var err error

	client := ahttp.New(ahttp.WithDeviceMetaInfo(productKey, deviceName, deviceSecret))
	// client.LogMode(true)

	uri := dm.URICOAPHTTPPrePrefix + fmt.Sprintf(dm.URISysPrefix, productKey, deviceName) + dm.URIThingEventPropertyPost
	bPayload, err := json.Marshal(dm.Request{
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
