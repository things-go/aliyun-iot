package main

import (
	"fmt"
	"math/rand"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/_examples/testmeta"
	"github.com/thinkgos/aliyun-iot/dm"
)

func main() {
	dmClient := aiot.NewWithHTTP(testmeta.MetaInfo())
	dmClient.LogMode(true)

	for {
		err := dmClient.AlinkReport(dm.MsgTypeEventPropertyPost, dm.DevNodeLocal, map[string]interface{}{
			"Temp":         rand.Intn(200),
			"Humi":         rand.Intn(100),
			"switchStatus": rand.Intn(1),
		})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}
