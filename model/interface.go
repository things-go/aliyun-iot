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
	ContainerOf() *Manager
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
	DownstreamThingModelUpRawReply(productKey, deviceName string, payload []byte) error
	DownstreamThingEventPropertyPostReply(rsp *Response) error
	DownstreamThingEventPostReply(eventID string, rsp *Response) error
	DownstreamThingDeviceInfoUpdateReply(rsp *Response) error
	DownstreamThingDeviceInfoDeleteReply(rsp *Response) error
	DownstreamThingDesiredPropertyGetReply(rsp *Response) error
	DownstreamThingDesiredPropertyDeleteReply(rsp *Response) error
	DownstreamThingDsltemplateGetReply(rsp *Response) error
	DownstreamThingDynamictslGetReply(rsp *Response) error
	DownstreamExtNtpResponse(rsp *NtpResponsePayload) error
	DownstreamThingConfigGetReply(rsp *ConfigGetResponse) error
	DownstreamThingConfigPush(rsp *ConfigPushRequest) error

	DownstreamExtErrorResponse(rsp *Response) error
	DownstreamThingModelDownRaw(productKey, deviceName string, payload []byte) error
	DownstreamThingServicePropertyGet(productKey, deviceName string, payload []byte) error
	DownstreamThingServiceRequest(productKey, deviceName, srvID string, payload []byte) error
	DownstreamThingServicePropertySet(payload []byte) error
}
