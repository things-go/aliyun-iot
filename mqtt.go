package aiot

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliyun-iot/clog"
	"github.com/thinkgos/aliyun-iot/dm"
)

// MQTTClient MQTT客户端
type MQTTClient struct {
	c mqtt.Client
	*dm.Client
	log *clog.Clog
}

// 确保 NopEvt 实现 dm.Conn 接口
var _ dm.Conn = (*MQTTClient)(nil)

// Publish 实现dm.Conn接口
func (sf *MQTTClient) Publish(topic string, qos byte, payload interface{}) error {
	return sf.c.Publish(topic, qos, false, payload).Error()
}

// Subscribe 实现dm.Conn接口
func (sf *MQTTClient) Subscribe(topic string, streamFunc dm.ProcDownStreamFunc) error {
	return sf.c.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
		if message.Duplicate() {
			return
		}
		if err := streamFunc(sf.Client, message.Topic(), message.Payload()); err != nil {
			sf.log.Error("topic: %s, Error: %+v", message.Topic(), err)
		}
	}).Error()
}

// UnSubscribe 实现dm.Conn接口
func (sf *MQTTClient) UnSubscribe(topic ...string) error {
	return sf.c.Unsubscribe(topic...).Error()
}

// LogProvider 实现dm.Conn接口
func (sf *MQTTClient) LogProvider() clog.LogProvider {
	return sf.log
}

// LogMode 实现dm.Conn接口
func (sf *MQTTClient) LogMode(enable bool) {
	sf.log.LogMode(enable)
}

// UnderlyingClient 获得底层的Client
func (sf *MQTTClient) UnderlyingClient() mqtt.Client {
	return sf.c
}

// NewWithMQTT 新建MQTTClient
func NewWithMQTT(config *dm.Config, c mqtt.Client) *MQTTClient {
	m := dm.New(config)
	cli := &MQTTClient{
		c,
		m,
		clog.New(clog.WithLogger(clog.NewLogger(log.New(os.Stderr, "mqtt --> ", log.LstdFlags)))),
	}
	m.SetConn(cli)
	return cli
}
