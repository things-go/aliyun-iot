package main

import (
	"fmt"
	"strings"
)

func main() {
	a := "messageId"
	uris := strings.SplitN(strings.TrimLeft(a, "/"), "/", 2)
	fmt.Println(uris[0])
	fmt.Println(uris[1])
}
