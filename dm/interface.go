package dm

import (
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

// DevUserProc 设备用户回调
type DevUserProc interface {
	DownstreamThingModelUpRawReply(m *Client, productKey, deviceName string, payload []byte) error
	DownstreamThingEventPropertyPostReply(m *Client, rsp *Response) error
	DownstreamThingEventPostReply(m *Client, eventID string, rsp *Response) error
	DownstreamThingDeviceInfoUpdateReply(m *Client, rsp *Response) error
	DownstreamThingDeviceInfoDeleteReply(m *Client, rsp *Response) error
	DownstreamThingDesiredPropertyGetReply(m *Client, rsp *Response) error
	DownstreamThingDesiredPropertyDeleteReply(m *Client, rsp *Response) error
	DownstreamThingDsltemplateGetReply(m *Client, rsp *Response) error
	DownstreamThingDynamictslGetReply(m *Client, rsp *Response) error
	DownstreamExtNtpResponse(m *Client, rsp *NtpResponsePayload) error
	DownstreamThingConfigGetReply(m *Client, rsp *ConfigGetResponse) error
	DownstreamExtErrorResponse(m *Client, rsp *Response) error
	// 透传请求
	DownstreamThingModelDownRaw(m *Client, productKey, deviceName string, payload []byte) error
	// 推送,已做默认回复,覆盖本接口并不覆盖默认回复
	DownstreamThingConfigPush(m *Client, rsp *ConfigPushRequest) error
	// 设置设备属性,已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	DownstreamThingServicePropertySet(m *Client, topic string, payload []byte) error
	// 设备服务调用,已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	DownstreamThingServiceRequest(m *Client, productKey, deviceName, srvID string, payload []byte) error
	// 系统RRPC调用, 仅支持设备端Qos = 0的回复. 已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	DownStreamRRPCRequest(m *Client, productKey, deviceName, messageID string, payload []byte) error
	// 自定义RRPC调用,仅支持设备端Qos = 0的回复, 已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	DownStreamExtRRPCRequest(m *Client, rawURI string, payload []byte) error
}

// GatewayUserProc 网关用户处理回调
type GatewayUserProc interface {
	DownstreamGwExtSubDevRegisterReply(m *Client, rsp *GwSubDevRegisterResponse) error
	DownstreamGwExtSubDevCombineLoginReply(m *Client, rsp *Response) error
	DownstreamGwExtSubDevCombineLogoutReply(m *Client, rsp *Response) error
	DownstreamGwSubDevThingDisable(m *Client, productKey, deviceName string) error
	DownstreamGwSubDevThingEnable(m *Client, productKey, deviceName string) error
	DownstreamGwSubDevThingDelete(m *Client, productKey, deviceName string) error
	DownstreamGwThingTopoAddReply(m *Client, rsp *Response) error
	DownstreamGwThingTopoDeleteReply(m *Client, rsp *Response) error
	DownstreamGwThingTopoGetReply(m *Client, rsp *GwTopoGetResponse) error
}
