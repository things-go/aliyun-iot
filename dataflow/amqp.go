// Package dataflow 阿里云iot服务器订阅数据流
package dataflow

// Properties amqp properties
type Properties struct {
	GenerateTime int64  `json:"generateTime"` //
	MessageID    int64  `json:"messageId"`
	Qos          int    `json:"qos"`
	Topic        string `json:"topic"`
}
