package aliIOT

import (
	"bytes"
	"errors"

	"github.com/go-ocf/go-coap"
	"github.com/thinkgos/aliIOT/clog"
	"github.com/thinkgos/aliIOT/dm"
)

// COAPClient COAP客户端
type COAPClient struct {
	c *coap.ClientConn
	*dm.Client
	log *clog.Clog
}

// Publish 实现dm.Conn接口
func (sf *COAPClient) Publish(topic string, qos byte, payload interface{}) error {
	var buf *bytes.Buffer
	switch v := payload.(type) {
	case string:
		buf = bytes.NewBufferString(v)
	case []byte:
		buf = bytes.NewBuffer(v)
	default:
		return errors.New("payload must be string or []byte")
	}

	// TODO
	_, _ = sf.c.Post(topic, coap.AppJSON, buf)
	return nil
}

// Subscribe 实现dm.Conn接口
func (sf *COAPClient) Subscribe(topic string, streamFunc dm.ProcDownStreamFunc) error {
	return nil
}

// UnSubscribe 实现dm.Conn接口
func (sf *COAPClient) UnSubscribe(topic ...string) error {
	return nil
}

// LogProvider 实现dm.Conn接口
func (sf *COAPClient) LogProvider() clog.LogProvider {
	return sf.log
}

// LogMode 实现dm.Conn接口
func (sf *COAPClient) LogMode(enable bool) {
	sf.log.LogMode(enable)
}

// UnderlyingClient 获得底层的Client
func (sf *COAPClient) UnderlyingClient() *coap.ClientConn {
	return sf.c
}

// NewWithMQTT 新建MQTTClient
func NewWithCOAP(config *dm.Config, c *coap.ClientConn) *COAPClient {
	m := dm.New(config)
	cli := &COAPClient{c, m, clog.NewWithPrefix("mqtt --> ")}
	m.SetConn(cli)
	return cli
}
