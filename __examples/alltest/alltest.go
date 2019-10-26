package main

import (
	"fmt"
	"time"
)

func main() {
	var t = time.Date(2019, 10, 26, 22, 56, 0, 0, time.Local)

	ts := time.Since(t)
	if ts < time.Minute*15 {
		fmt.Println("haha")
	}
	fmt.Println(ts.Minutes())

}
