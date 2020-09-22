package main

import (
	"github.com/thinkgos/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init(mock.MetaTriad)

	mock.ThingEventPropertyPost(client) // done
}
