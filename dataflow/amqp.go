// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dataflow

// Properties amqp properties
type Properties struct {
	GenerateTime int64  `json:"generateTime"` //
	MessageID    int64  `json:"messageId"`
	Qos          int    `json:"qos"`
	Topic        string `json:"topic"`
}
