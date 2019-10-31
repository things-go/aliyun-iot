package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliIOT/clog"
)

// Conn conn接口
type Conn interface {
	// Publish will publish a message with the specified QoS and content
	Publish(topic string, qos byte, payload interface{}) error
	Subscribe(topic string, callback ProcDownStreamFunc) error
	UnSubscribe(topic ...string) error
	// 目志调试
	LogProvider() clog.LogProvider
	LogMode(enable bool)
}

// EventProc 事件回调接口
type EventProc interface {
	EvtThingModelUpRawReply(m *Client, productKey, deviceName string, payload []byte) error
	EvtThingEventPropertyPostReply(m *Client, err error, productKey, deviceName string) error
	EvtThingEventPostReply(m *Client, err error, eventID, productKey, deviceName string) error
	EvtThingDeviceInfoUpdateReply(m *Client, err error, productKey, deviceName string) error
	EvtThingDeviceInfoDeleteReply(m *Client, err error, productKey, deviceName string) error
	EvtThingDesiredPropertyGetReply(m *Client, err error, productKey, deviceName string, data json.RawMessage) error
	EvtThingDesiredPropertyDeleteReply(m *Client, err error, productKey, deviceName string) error
	EvtThingDsltemplateGetReply(m *Client, err error, productKey, deviceName string, data json.RawMessage) error
	EvtThingDynamictslGetReply(m *Client, err error, productKey, deviceName string, data json.RawMessage) error
	EvtExtNtpResponse(m *Client, productKey, deviceName string, payload NtpResponsePayload) error
	EvtThingConfigGetReply(m *Client, err error, productKey, deviceName string, data ConfigParamsAndData) error
	EvtExtErrorResponse(m *Client, rsp *Response) error
	// 透传请求
	EvtThingModelDownRaw(m *Client, productKey, deviceName string, payload []byte) error
	// 推送,已做默认回复,覆盖本接口并不覆盖默认回复
	EvtThingConfigPush(m *Client, productKey, deviceName string, params ConfigParamsAndData) error
	// 设置设备属性,已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	EvtThingServicePropertySet(m *Client, productKey, deviceName string, payload []byte) error
	// 设备服务调用,已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	EvtThingServiceRequest(m *Client, srvID, productKey, deviceName string, payload []byte) error
	// 系统RRPC调用, 仅支持设备端Qos = 0的回复. 已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	EvtRRPCRequest(m *Client, messageID, productKey, deviceName string, payload []byte) error
	// 自定义RRPC调用,仅支持设备端Qos = 0的回复, 已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	EvtExtRRPCRequest(m *Client, rawURI string, payload []byte) error
}
