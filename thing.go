// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import (
	"encoding/json"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
)

// @see https://help.aliyun.com/document_detail/89301.html?spm=a2c4g.11186623.6.706.570f3f69J3fW5z

const (
	property = "property"
)

// ProcDownStream 处理下行数据
type ProcDownStream func(c *Client, rawURI string, payload []byte) error

// Conn conn接口
type Conn interface {
	// Publish will publish a Message with the specified QoS and content
	Publish(topic string, qos byte, payload interface{}) error
	Subscribe(topic string, callback ProcDownStream) error
	UnSubscribe(topic ...string) error
}

// Callback 事件回调接口
type Callback interface {
	// 透传应答
	ThingModelUpRawReply(c *Client, productKey, deviceName string, payload []byte) error
	// 透传请求,需要用户自己处理及应答
	ThingModelDownRaw(c *Client, productKey, deviceName string, payload []byte) error
	// event
	ThingEventPropertyPostReply(c *Client, err error, productKey, deviceName string) error
	ThingEventPostReply(c *Client, err error, eventID, productKey, deviceName string) error
	ThingEventPropertyPackPostReply(c *Client, err error, productKey, deviceName string) error
	ThingEventPropertyHistoryPostReply(c *Client, err error, productKey, deviceName string) error
	// device info
	ThingDeviceInfoUpdateReply(c *Client, err error, productKey, deviceName string) error
	ThingDeviceInfoDeleteReply(c *Client, err error, productKey, deviceName string) error
	// desired property
	ThingDesiredPropertyGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	ThingDesiredPropertyDeleteReply(c *Client, err error, productKey, deviceName string) error
	// template
	ThingDsltemplateGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	ThingDynamictslGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	// config
	ThingConfigGetReply(c *Client, err error, productKey, deviceName string, data ConfigParamsData) error
	// 配置推送,已做默认回复
	ThingConfigPush(c *Client, productKey, deviceName string, params ConfigParamsData) error

	// Log
	ThingConfigLogGetReply(c *Client, err error, productKey, deviceName string, data ConfigLogParamData) error
	ThingConfigLogPush(c *Client, productKey, deviceName string, param ConfigLogParamData) error
	ThingLogPostReply(c *Client, err error, productKey, deviceName string) error

	// diag
	ThingDialPostReply(c *Client, err error, productKey, deviceName string) error

	// service
	// 设置设备属性, 需用户自行做回复
	ThingServicePropertySet(c *Client, productKey, deviceName string, payload []byte) error
	// 设备服务调用,需用户自行做回复
	ThingServiceRequest(c *Client, srvID, productKey, deviceName string, payload []byte) error

	// ntp
	ExtNtpResponse(c *Client, productKey, deviceName string, exact time.Time) error

	// 系统RRPC调用, 仅支持设备端Qos = 0的回复,需用户自行做回复
	RRPCRequest(c *Client, messageID, productKey, deviceName string, payload []byte) error
	// 自定义RRPC调用,仅支持设备端Qos = 0的回复, 需用户自行做回复
	ExtRRPCRequest(c *Client, messageID, topic string, payload []byte) error
	// ota
	OtaUpgrade(c *Client, productKey, deviceName string, rsp *OtaFirmwareResponse) error
	ThingOtaFirmwareGetReply(c *Client, productKey, deviceName string, data OtaFirmwareData) error
}

// GwCallback 网关事件接口
type GwCallback interface {
	// 520错误已做自动登陆回复
	ExtErrorResponse(c *Client, err error, productKey, deviceName string) error
	ThingTopoGetReply(c *Client, err error, params []infra.MetaPair) error
	ThingListFoundReply(c *Client, err error) error
	ThingTopoAddNotify(c *Client, params []infra.MetaPair) error
	ThingTopoChange(c *Client, params TopoChangeParams) error
	ThingDisable(c *Client, productKey, deviceName string) error
	ThingEnable(c *Client, productKey, deviceName string) error
	ThingDelete(c *Client, productKey, deviceName string) error
}
