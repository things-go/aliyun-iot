package main

import (
	"fmt"
	"math/rand"
	"time"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/dm"
)

const (
	productKey    = "a1QR3GD1Db3"
	productSecret = ""
	deviceName    = "MPA19GT010070140"
	deviceSecret  = "26b136834c0dc1b9f3afc64158f6580d"
)

func main() {
	dmopt := dm.NewConfig(productKey, deviceName, deviceSecret).
		//EnableModelRaw().
		Valid()
	dmClient := aiot.NewWithHTTP(dmopt)
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
