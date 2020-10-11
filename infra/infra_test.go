package infra

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func TestClientID(t *testing.T) {
	pk, dn := "pk", "dn"
	require.Equal(t, pk+"."+dn, ClientID(pk, dn))
}

func TestTime(t *testing.T) {
	tm := time.Date(2020, 9, 29, 10, 10, 10, 10000000, time.UTC)
	msec := int64(1601374210010)

	require.Equal(t, msec, Millisecond(tm))
	require.True(t, tm.Equal(Time(msec)))
}
