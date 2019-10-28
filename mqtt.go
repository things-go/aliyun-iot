package aliIOT

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliIOT/clog"
	"github.com/thinkgos/aliIOT/model"
)

type MQTTClient struct {
	c mqtt.Client
	*model.Manager
	log *clog.Clog
}

func (sf *MQTTClient) Publish(topic string, qos byte, payload interface{}) error {
	return sf.c.Publish(topic, qos, false, payload).Error()
}

func (sf *MQTTClient) Subscribe(topic string, streamFunc model.ProcDownStreamFunc) error {
	return sf.c.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
		if message.Duplicate() {
			return
		}
		if err := streamFunc(sf.Manager, message.Topic(), message.Payload()); err != nil {
			sf.log.Error("topic: %s, Error: %+v", message.Topic(), err)
		}
	}).Error()
}

func (sf *MQTTClient) UnSubscribe(topic ...string) error {
	return sf.c.Unsubscribe(topic...).Error()
}

func (sf *MQTTClient) LogProvider() clog.LogProvider {
	return sf.log
}

func (sf *MQTTClient) LogMode(enable bool) {
	sf.log.LogMode(enable)
}

func (sf *MQTTClient) UnderlyingClient() interface{} {
	return sf.c
}

func NewWithMQTT(config *model.Config, c mqtt.Client) *MQTTClient {
	m := model.New(config)
	mqttCli := &MQTTClient{c, m, clog.NewWithPrefix("mqtt --> ")}
	m.SetConn(mqttCli)
	return mqttCli
}
