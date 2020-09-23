// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dataflow

import "time"

const (
	layout    = "2006-01-02 15:04:05.999"
	utcLayout = "2006-01-02T15:04:05.999Z"
)

// Time local time layout like "2006-01-02 15:04:05.999"
type Time time.Time

// MarshalJSON implemented interface Marshaler
func (sf Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(layout)+2)
	b = append(b, '"')
	b = time.Time(sf).Local().AppendFormat(b, layout)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implemented interface Unmarshaler
func (sf *Time) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation(`"`+layout+`"`, string(data), time.Local)
	*sf = Time(t)
	return err
}

// String implemented interface Stringer
func (sf Time) String() string {
	return time.Time(sf).Format(layout)
}

// UTCtime utc time layout like "2006-01-02T15:04:05.999Z"
type UTCtime time.Time

// MarshalJSON implemented interface Marshaler
func (sf UTCtime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(layout)+2)
	b = append(b, '"')
	b = time.Time(sf).UTC().AppendFormat(b, utcLayout)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implemented interface Unmarshaler
func (sf *UTCtime) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation(`"`+utcLayout+`"`, string(data), time.UTC)
	*sf = UTCtime(t)
	return err
}

// String implemented interface Stringer
func (sf UTCtime) String() string {
	return time.Time(sf).Format(utcLayout)
}

// Unix 时间戳转换
func Unix(msec int64) time.Time {
	return time.Unix(msec/1000, (msec%1000)*1000000)
}
