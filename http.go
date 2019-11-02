package aliIOT

import (
	"github.com/thinkgos/aliIOT/ahttp"
	"github.com/thinkgos/aliIOT/clog"
	"github.com/thinkgos/aliIOT/dm"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	c *ahttp.Client
	*dm.Client
}

// 确保 NopEvt 实现 dm.Conn 接口
var _ dm.Conn = (*HTTPClient)(nil)

//Publish 实现dm.Conn接口
func (sf *HTTPClient) Publish(topic string, _ byte, payload interface{}) error {
	return sf.c.Publish(topic, payload)
}

//Subscribe 实现dm.Conn接口
func (sf *HTTPClient) Subscribe(_ string, _ dm.ProcDownStreamFunc) error {
	return nil
}

// UnSubscribe 实现dm.Conn接口
func (sf *HTTPClient) UnSubscribe(_ ...string) error {
	return nil
}

// LogProvider 实现dm.Conn接口
func (sf *HTTPClient) LogProvider() clog.LogProvider {
	return sf.c.Clog
}

// LogMode 实现dm.Conn接口
func (sf *HTTPClient) LogMode(enable bool) {
	sf.c.LogMode(enable)
}

// UnderlyingClient 底层客户端
func (sf *HTTPClient) UnderlyingClient() *ahttp.Client {
	return sf.c
}

// NewWithHTTP 新建HTTP客户端
func NewWithHTTP(config *dm.Config) *HTTPClient {
	client := ahttp.New().
		SetDeviceMetaInfo(config.
			EnableHTTP().
			MetaInfo())

	sf := dm.New(config)
	cli := &HTTPClient{c: client, Client: sf}
	sf.SetConn(cli)
	return cli
}
