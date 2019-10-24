package aliIOT

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliIOT/model"
)

type Client struct {
	c         mqtt.Client
	containOf *model.Manager
}

func (sf *Client) Publish(topic string, payload interface{}) error {
	return sf.c.Publish(topic, 1, false, payload).Error()
}

func (sf *Client) UnderlyingClient() interface{} {
	return sf.c
}
func (sf *Client) ContainerOf() *model.Manager {
	return sf.containOf
}

func (sf *Client) Subscribe(topic string, streamFunc model.ProcDownStreamFunc) error {
	return sf.c.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
		_ = streamFunc(sf.containOf, message.Topic(), message.Payload())
	}).Error()
}

func New(productKey, deviceName string, c mqtt.Client) *model.Manager {
	sf := model.New(productKey, deviceName)
	return sf.SetCon(&Client{c: c, containOf: sf})
}
