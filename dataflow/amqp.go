// Package dataflow
package dataflow

// Properties amqp properties
type Properties struct {
	GenerateTime int64  `json:"generateTime"` //
	MessageId    int64  `json:"messageId"`
	Qos          int    `json:"qos"`
	Topic        string `json:"topic"`
}
