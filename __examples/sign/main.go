package main

import (
	"fmt"

	"github.com/thinkgos/aliIOT/sign"
)

func main() {
	sig := sign.NewMQTTSign()
	sig.SetSignMethod(sign.SignMethodSHA1)
	signout, err := sig.Generate(&sign.MetaInfo{
		ProductKey:   "a1QR3GD1Db3",
		DeviceName:   "testcar",
		DeviceSecret: "eOHD59KvSI45Vv8HEYpj6ImmqNCEgBEc",
	}, sign.CloudRegionShangHai)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", signout)

}
