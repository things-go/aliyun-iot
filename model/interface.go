package model

// Conn conn接口
type Conn interface {
	// Publish will publish a message with the specified QoS and content
	Publish(topic string, qos byte, payload interface{}) error
	UnderlyingClient() interface{}
	Subscribe(topic string, streamFunc ProcDownStreamFunc) error
	ContainerOf() *Manager
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
