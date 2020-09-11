package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/clog"
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

// EventProc 事件接口
type EventProc interface {
	// 上行应答
	EvtThingModelUpRawReply(c *Client, productKey, deviceName string, payload []byte) error
	EvtThingEventPropertyPostReply(c *Client, err error, productKey, deviceName string) error
	EvtThingEventPostReply(c *Client, err error, eventID, productKey, deviceName string) error
	EvtThingEventPropertyPackPostReply(c *Client, err error, productKey, deviceName string) error
	EvtThingDeviceInfoUpdateReply(c *Client, err error, productKey, deviceName string) error
	EvtThingDeviceInfoDeleteReply(c *Client, err error, productKey, deviceName string) error
	EvtThingDesiredPropertyGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	EvtThingDesiredPropertyDeleteReply(c *Client, err error, productKey, deviceName string) error
	EvtThingDsltemplateGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	EvtThingDynamictslGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	EvtExtNtpResponse(c *Client, productKey, deviceName string, payload NtpResponse) error
	EvtThingConfigGetReply(c *Client, err error, productKey, deviceName string, data ConfigParamsData) error
	// 下行
	// 透传请求,需要用户自己处理及应答
	EvtThingModelDownRaw(c *Client, productKey, deviceName string, payload []byte) error
	// 配置推送,已做默认回复
	EvtThingConfigPush(c *Client, productKey, deviceName string, params ConfigParamsData) error
	// 设置设备属性, 需用户自行做回复
	EvtThingServicePropertySet(c *Client, productKey, deviceName string, payload []byte) error
	// 设备服务调用,需用户自行做回复
	EvtThingServiceRequest(c *Client, srvID, productKey, deviceName string, payload []byte) error
	// 系统RRPC调用, 仅支持设备端Qos = 0的回复,需用户自行做回复
	EvtRRPCRequest(c *Client, messageID, productKey, deviceName string, payload []byte) error
	// 自定义RRPC调用,仅支持设备端Qos = 0的回复, 需用户自行做回复
	EvtExtRRPCRequest(c *Client, messageID, topic string, payload []byte) error
}

// EventGwProc 网关事件接口
type EventGwProc interface {
	// 520错误已做自动登陆回复
	EvtExtErrorResponse(c *Client, err error, productKey, deviceName string) error
	EvtThingGwSubTopoGetReply(c *Client, err error, params []GwSubTopoGetData) error
	EvtThingListFoundReply(c *Client, err error) error
	EvtThingTopoAddNotify(c *Client, params []GwSubTopoAddNotifyParams) error
	EvtThingTopoChange(c *Client, params GwTopoChangeParams) error
	EvtThingDisable(c *Client, productKey, deviceName string) error
	EvtThingEnable(c *Client, productKey, deviceName string) error
	EvtThingDelete(c *Client, productKey, deviceName string) error
}
