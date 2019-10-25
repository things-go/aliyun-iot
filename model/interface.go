package model

// Conn conn接口
type Conn interface {
	// Publish will publish a message with the specified QoS and content
	Publish(topic string, qos byte, payload interface{}) error
	UnderlyingClient() interface{}
	Subscribe(topic string, streamFunc ProcDownStreamFunc) error
	ContainerOf() *Manager
}
