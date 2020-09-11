package aiot

import (
	"github.com/thinkgos/aliyun-iot/ahttp"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/infra"
)

// 确保 NopEvt 实现 dm.Conn 接口
var _ dm.Conn = (*httpClient)(nil)

// httpClient HTTP客户端
type httpClient struct {
	c *ahttp.Client
}

// Publish 实现dm.Conn接口
func (sf *httpClient) Publish(topic string, _ byte, payload interface{}) error {
	return sf.c.Publish(topic, payload)
}

// Subscribe 实现dm.Conn接口
func (*httpClient) Subscribe(string, dm.ProcDownStream) error { return nil }

// UnSubscribe 实现dm.Conn接口
func (*httpClient) UnSubscribe(...string) error { return nil }

// HTTPClient HTTP客户端
type HTTPClient struct {
	c *ahttp.Client
	*dm.Client
}

// NewWithHTTP 新建HTTP客户端
func NewWithHTTP(meta infra.MetaInfo) *HTTPClient {
	c := ahttp.New(meta)
	return &HTTPClient{c, dm.New(meta, &httpClient{c})}
}

// UnderlyingClient 底层客户端
func (sf *HTTPClient) UnderlyingClient() *ahttp.Client { return sf.c }
