package aiot

import (
	"github.com/thinkgos/aliyun-iot/ahttp"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/infra"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	c *ahttp.Client
	*dm.Client
}

// 确保 NopEvt 实现 dm.Conn 接口
var _ dm.Conn = (*HTTPClient)(nil)

// Publish 实现dm.Conn接口
func (sf *HTTPClient) Publish(topic string, _ byte, payload interface{}) error {
	return sf.c.Publish(topic, payload)
}

// Subscribe 实现dm.Conn接口
func (sf *HTTPClient) Subscribe(_ string, _ dm.ProcDownStream) error {
	return nil
}

// UnSubscribe 实现dm.Conn接口
func (sf *HTTPClient) UnSubscribe(_ ...string) error {
	return nil
}

// UnderlyingClient 底层客户端
func (sf *HTTPClient) UnderlyingClient() *ahttp.Client {
	return sf.c
}

// NewWithHTTP 新建HTTP客户端
func NewWithHTTP(meta infra.MetaInfo) *HTTPClient {
	client := ahttp.New(meta)
	sf := dm.New(meta)
	cli := &HTTPClient{c: client, Client: sf}
	sf.SetConn(cli)
	return cli
}
