package sign

import (
	"testing"

	"github.com/thinkgos/aliIOT/infra"
)

const (
	productKey   = "a1QR3GD1Db3"
	deviceName   = "MPA19GT010070140"
	deviceSecret = "CsC7Gmb6EvDLOm8V40HLOQwFPdc3KCHT"
	customRegion = "http"
)

func TestMQTTSign(t *testing.T) {
	t.Run("MQTT sign almost all", func(t *testing.T) {
		signout, err := NewMQTTSign().
			SetSignMethod(infra.SignMethodHMACSHA256).
			SetSDKVersion("SDK-Golang-v1.13.3").
			SetSupportExtRRPC().
			AddCustomKV("testKey", "testValue").
			DeleteCustomKV("deleteKey").
			Generate(&MetaInfo{
				ProductKey:   productKey,
				DeviceName:   deviceName,
				DeviceSecret: deviceSecret,
			}, CloudRegionShangHai)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v", signout)
	})

	t.Run("MQTT sign custom cloud region", func(t *testing.T) {
		signout, err := NewMQTTSign().
			SetSignMethod(infra.SignMethodHMACSHA1).
			Generate(&MetaInfo{
				ProductKey:   productKey,
				DeviceName:   deviceName,
				DeviceSecret: deviceSecret,
				CustomDomain: "iot.custom.com",
			}, CloudRegionCustom)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v", signout)
	})

	t.Run("MQTT sign empty custom cloud region", func(t *testing.T) {
		_, err := NewMQTTSign().
			Generate(&MetaInfo{
				ProductKey:   productKey,
				DeviceName:   deviceName,
				DeviceSecret: deviceSecret,
				CustomDomain: "",
			}, CloudRegionCustom)
		if err == nil {
			t.Fatal("must be error")
		}
	})

}
