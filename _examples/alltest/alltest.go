package main

import (
	"fmt"
	"time"

	"github.com/thinkgos/cache-go"
)

func main() {
	ts := cache.New(cache.DefaultExpiration, time.Second)
	ts.Set("haha", "bbb", time.Second)
	ts.SetDefault("delete", "ccc")
	ts.OnEvicted(func(s string, i interface{}) {
		fmt.Println(s)
	})
	ts.OnDeleted(func(s string, i interface{}) {
		fmt.Println(s)
	})
	ts.Delete("delete")
	time.Sleep(time.Second * 3)
}
