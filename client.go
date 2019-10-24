package aliIOT

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliIOT/model"
)

type mqttClient struct {
	c         mqtt.Client
	containOf *model.Manager
}

func (sf *mqttClient) Publish(topic string, payload interface{}) error {
	return sf.c.Publish(topic, 1, false, payload).Error()
}

func (sf *mqttClient) UnderlyingClient() interface{} {
	return sf.c
}

func (sf *mqttClient) ContainerOf() *model.Manager {
	return sf.containOf
}

func (sf *mqttClient) Subscribe(topic string, streamFunc model.ProcDownStreamFunc) error {
	return sf.c.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
		_ = streamFunc(sf.containOf, message.Topic(), message.Payload())
	}).Error()
}

func NewWithMQTT(productKey, deviceName, deviceSecret string, c mqtt.Client) *model.Manager {
	sf := model.New(productKey, deviceName, deviceSecret)
	return sf.SetCon(&mqttClient{c: c, containOf: sf})
}
