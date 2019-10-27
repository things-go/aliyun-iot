package main

import (
	"encoding/json"
	"fmt"
)

type Name struct {
	A int `json:"a,string"`
}

func main() {
	name := Name{}
	json.Unmarshal([]byte(`{"a":"12000"}`), name)
	fmt.Println(name)
}
