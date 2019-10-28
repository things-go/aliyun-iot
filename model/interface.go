package model

import (
	"github.com/thinkgos/aliIOT/clog"
)

// Conn conn接口
type Conn interface {
	// Publish will publish a message with the specified QoS and content
	Publish(topic string, qos byte, payload interface{}) error
	UnderlyingClient() interface{}
	Subscribe(topic string, streamFunc ProcDownStreamFunc) error
	// 目志调试
	LogProvider() clog.LogProvider
	LogMode(enable bool)
}

type GatewayUserProc interface {
	DownstreamGwExtSubDevRegisterReply(m *Manager, rsp *GwSubDevRegisterResponse) error
	DownstreamGwExtSubDevCombineLoginReply(m *Manager, rsp *Response) error
	DownstreamGwExtSubDevCombineLogoutReply(m *Manager, rsp *Response) error
	DownstreamGwSubDevThingDisable(m *Manager, productKey, deviceName string) error
	DownstreamGwSubDevThingEnable(m *Manager, productKey, deviceName string) error
	DownstreamGwSubDevThingDelete(m *Manager, productKey, deviceName string) error
	DownstreamGwThingTopoAddReply(m *Manager, rsp *Response) error
	DownstreamGwThingTopoDeleteReply(m *Manager, rsp *Response) error
	DownstreamGwThingTopoGetReply(m *Manager, rsp *GwTopoGetResponse) error
}

type DevUserProc interface {
	DownstreamThingModelUpRawReply(m *Manager, productKey, deviceName string, payload []byte) error
	DownstreamThingEventPropertyPostReply(m *Manager, rsp *Response) error
	DownstreamThingEventPostReply(m *Manager, eventID string, rsp *Response) error
	DownstreamThingDeviceInfoUpdateReply(m *Manager, rsp *Response) error
	DownstreamThingDeviceInfoDeleteReply(m *Manager, rsp *Response) error
	DownstreamThingDesiredPropertyGetReply(m *Manager, rsp *Response) error
	DownstreamThingDesiredPropertyDeleteReply(m *Manager, rsp *Response) error
	DownstreamThingDsltemplateGetReply(m *Manager, rsp *Response) error
	DownstreamThingDynamictslGetReply(m *Manager, rsp *Response) error
	DownstreamExtNtpResponse(m *Manager, rsp *NtpResponsePayload) error
	DownstreamThingConfigGetReply(m *Manager, rsp *ConfigGetResponse) error
	DownstreamExtErrorResponse(m *Manager, rsp *Response) error
	// 透传请求
	DownstreamThingModelDownRaw(m *Manager, productKey, deviceName string, payload []byte) error
	// 推送,已做默认回复,覆盖本接口并不覆盖默认回复
	DownstreamThingConfigPush(m *Manager, rsp *ConfigPushRequest) error
	// 设置设备属性,已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	DownstreamThingServicePropertySet(m *Manager, topic string, payload []byte) error
	// 设备服务调用,已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	DownstreamThingServiceRequest(m *Manager, productKey, deviceName, srvID string, payload []byte) error
	// 系统RRPC调用, 仅支持设备端Qos = 0的回复. 已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	DownStreamRRPCRequest(m *Manager, productKey, deviceName, messageID string, payload []byte) error
	// 自定义RRPC调用,仅支持设备端Qos = 0的回复, 已做默认回复,覆盖本接口覆盖默认回复,需用户自行做回复
	DownStreamExtRRPCRequest(m *Manager, rawURI string, payload []byte) error
}
