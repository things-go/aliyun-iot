package infra

import (
	"time"
)

func Millisecond(tm time.Time) int64 {
	return tm.Unix()*1000 + int64(tm.Nanosecond())/1000000
}

func Time(msec int64) time.Time {
	return time.Unix(msec/1000, (msec%1000)*1000000)
}
