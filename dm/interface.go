package dm

import (
	"encoding/json"
)

// Conn conn接口
type Conn interface {
	// Publish will publish a message with the specified QoS and content
	Publish(topic string, qos byte, payload interface{}) error
	Subscribe(topic string, callback ProcDownStream) error
	UnSubscribe(topic ...string) error
}

// Callback 事件回调接口
type Callback interface {
	// 上行应答
	ThingModelUpRawReply(c *Client, productKey, deviceName string, payload []byte) error
	ThingEventPropertyPostReply(c *Client, err error, productKey, deviceName string) error
	ThingEventPostReply(c *Client, err error, eventID, productKey, deviceName string) error
	ThingEventPropertyPackPostReply(c *Client, err error, productKey, deviceName string) error
	ThingDeviceInfoUpdateReply(c *Client, err error, productKey, deviceName string) error
	ThingDeviceInfoDeleteReply(c *Client, err error, productKey, deviceName string) error
	ThingDesiredPropertyGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	ThingDesiredPropertyDeleteReply(c *Client, err error, productKey, deviceName string) error
	ThingDsltemplateGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	ThingDynamictslGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error
	ExtNtpResponse(c *Client, productKey, deviceName string, payload NtpResponse) error
	ThingConfigGetReply(c *Client, err error, productKey, deviceName string, data ConfigParamsData) error
	// 下行
	// 透传请求,需要用户自己处理及应答
	ThingModelDownRaw(c *Client, productKey, deviceName string, payload []byte) error
	// 配置推送,已做默认回复
	ThingConfigPush(c *Client, productKey, deviceName string, params ConfigParamsData) error
	// 设置设备属性, 需用户自行做回复
	ThingServicePropertySet(c *Client, productKey, deviceName string, payload []byte) error
	// 设备服务调用,需用户自行做回复
	ThingServiceRequest(c *Client, srvID, productKey, deviceName string, payload []byte) error
	// 系统RRPC调用, 仅支持设备端Qos = 0的回复,需用户自行做回复
	RRPCRequest(c *Client, messageID, productKey, deviceName string, payload []byte) error
	// 自定义RRPC调用,仅支持设备端Qos = 0的回复, 需用户自行做回复
	ExtRRPCRequest(c *Client, messageID, topic string, payload []byte) error
}

// GwCallback 网关事件接口
type GwCallback interface {
	// 520错误已做自动登陆回复
	ExtErrorResponse(c *Client, err error, productKey, deviceName string) error
	ThingGwTopoGetReply(c *Client, err error, params []GwTopoGetData) error
	ThingGwListFoundReply(c *Client, err error) error
	ThingGwTopoAddNotify(c *Client, params []GwTopoAddNotifyParams) error
	ThingGwTopoChange(c *Client, params GwTopoChangeParams) error
	ThingGwDisable(c *Client, productKey, deviceName string) error
	ThingGwEnable(c *Client, productKey, deviceName string) error
	ThingGwDelete(c *Client, productKey, deviceName string) error
}
