package aiot

import (
	"bytes"
	"errors"

	"github.com/go-ocf/go-coap"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/infra"
)

// COAPClient COAP客户端
type COAPClient struct {
	c *coap.ClientConn
	*dm.Client
}

// 确保 NopEvt 实现 dm.Conn 接口
var _ dm.Conn = (*COAPClient)(nil)

// Publish 实现dm.Conn接口
func (sf *COAPClient) Publish(topic string, _ byte, payload interface{}) error {
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
func (sf *COAPClient) Subscribe(topic string, streamFunc dm.ProcDownStream) error {
	return nil
}

// UnSubscribe 实现dm.Conn接口
func (sf *COAPClient) UnSubscribe(...string) error {
	return nil
}

// UnderlyingClient 获得底层的Client
func (sf *COAPClient) UnderlyingClient() *coap.ClientConn {
	return sf.c
}

// NewWithCOAP 新建MQTTClient
func NewWithCOAP(meta infra.MetaInfo, c *coap.ClientConn, opts ...dm.Option) *COAPClient {
	m := dm.New(meta, append(opts, dm.WithWork(dm.WorkOnCOAP))...)
	cli := &COAPClient{
		c,
		m,
	}
	m.SetConn(cli)
	return cli
}
