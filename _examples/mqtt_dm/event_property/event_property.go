package main

import (
	"github.com/thinkgos/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init()

	mock.ThingEventPropertyPost(client) // done
}
