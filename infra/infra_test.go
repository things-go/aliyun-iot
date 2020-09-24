package infra

import (
	"testing"
)

func BenchmarkCalcSign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalcSign("hmacsha1", MetaTriad{
			ProductKey:   "11",
			DeviceName:   "22",
			DeviceSecret: "333",
		}, 124134134)
	}
}
