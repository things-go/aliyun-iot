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
	DownstreamExtSubDevRegisterReply(m *Manager, rsp *SubDevRegisterResponse) error
	DownstreamExtSubDevCombineLoginReply(m *Manager, rsp *Response) error
	DownstreamExtSubDevCombineLogoutReply(m *Manager, rsp *Response) error
	DownstreamThingTopoAddReply(m *Manager, rsp *Response) error
	DownstreamThingTopoDeleteReply(m *Manager, rsp *Response) error
}
