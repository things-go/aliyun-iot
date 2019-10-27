package aliIOT

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliIOT/clog"
	"github.com/thinkgos/aliIOT/model"
)

type mqttClient struct {
	c         mqtt.Client
	containOf *model.Manager
	log       *clog.Clog
}

func (sf *mqttClient) Publish(topic string, qos byte, payload interface{}) error {
	return sf.c.Publish(topic, qos, false, payload).Error()
}

func (sf *mqttClient) UnderlyingClient() interface{} {
	return sf.c
}

func (sf *mqttClient) Subscribe(topic string, streamFunc model.ProcDownStreamFunc) error {
	return sf.c.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
		if message.Duplicate() {
			return
		}
		if err := streamFunc(sf.containOf, message.Topic(), message.Payload()); err != nil {
			sf.log.Error("topic: %s, Error: %+v", message.Topic(), err)
		}
	}).Error()
}

func (sf *mqttClient) ContainerOf() *model.Manager {
	return sf.containOf
}

func (sf *mqttClient) LogProvider() clog.LogProvider {
	return sf.log
}

func (sf *mqttClient) LogMode(enable bool) {
	sf.log.LogMode(enable)
}

func NewWithMQTT(options *model.Options, c mqtt.Client) *model.Manager {
	sf := model.New(options)
	return sf.SetConn(&mqttClient{c: c, containOf: sf, log: clog.NewWithPrefix("mqtt --> ")})
}
