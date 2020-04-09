package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/thinkgos/aiot"
	"github.com/thinkgos/aiot/dm"
)

const (
	productKey    = "a1QR3GD1Db3"
	productSecret = ""
	deviceName    = "MPA19GT010070140"
	deviceSecret  = "ld9Xf2BtKGfdEC7G9nSMe1wYfgllvi3Q"
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
