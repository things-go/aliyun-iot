package main

import (
	"github.com/things-go/aliyun-iot/_examples/mock"
)

func main() {
	client := mock.Init(mock.MetaTriad)

	mock.ThingEventPropertyPost(client) // done
}
