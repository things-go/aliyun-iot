// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package dataflow 定义数据流转的数据格式
// 实现iot转发的解析以及相关
// 客户端将获得一个messageBody,解析到message, payload承载着真实的dataflow数据流
// 根据不同的主题,解析不同的payload,see https://help.aliyun.com/document_detail/73736.html?spm=a2c4g.11186623.6.630.1ce25a10TgnylI
package dataflow

// message type 消息类型
const (
	MessageTypeStatus = "status"
	MessageTypeUpload = "upload"
)

// Message mns 消息负载
type Message struct {
	Payload     string `json:"payload"`     // dataflow消息负载,须base64解码
	MessageType string `json:"messagetype"` // 消型类型
	Topic       string `json:"topic"`       // dataflow主题
	MessageID   int64  `json:"messageid"`   // 消息id,平台发送,平台内唯一
	Timestamp   int64  `json:"timestamp"`   // 消息时间戳
}
