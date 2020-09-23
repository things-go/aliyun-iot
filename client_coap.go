package aiot

import (
	"bytes"
	"errors"

	"github.com/go-ocf/go-coap"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/57697.html?spm=a2c4g.11186623.6.606.5d7a12e0FGY05a

// 确保 NopCb 实现 dm.Conn 接口
var _ Conn = (*COAPClient)(nil)

// COAPClient COAP客户端
type coapClient struct {
	c *coap.ClientConn
}

// Publish 实现dm.Conn接口
func (sf *coapClient) Publish(_uri string, _ byte, payload interface{}) error {
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
	_, _ = sf.c.Post(uri.TopicPrefix+_uri, coap.AppJSON, buf)
	return nil
}

// Subscribe 实现dm.Conn接口
func (*coapClient) Subscribe(string, ProcDownStream) error { return nil }

// UnSubscribe 实现dm.Conn接口
func (sf *coapClient) UnSubscribe(...string) error { return nil }

// COAPClient COAP客户端
type COAPClient struct {
	*Client
}

// NewWithCOAP 新建MQTTClient
func NewWithCOAP(meta infra.MetaTriad, c *coap.ClientConn, opts ...Option) *COAPClient {
	return &COAPClient{
		New(meta, &coapClient{c}, append(opts, WithWork(WorkOnCOAP))...),
	}
}

// Underlying 获得底层的Client
func (sf *COAPClient) Underlying() *coap.ClientConn {
	return sf.Conn.(*coapClient).c
}
