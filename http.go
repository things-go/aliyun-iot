package aliIOT

import (
	"github.com/thinkgos/aliIOT/ahttp"
	"github.com/thinkgos/aliIOT/clog"
	"github.com/thinkgos/aliIOT/dm"
)

type HTTPClient struct {
	c *ahttp.Client
	*dm.Client
}

func (sf *HTTPClient) Publish(topic string, qos byte, payload interface{}) error {
	return sf.c.Publish(topic, payload)
}

func (sf *HTTPClient) Subscribe(topic string, streamFunc dm.ProcDownStreamFunc) error {
	return nil
}

func (sf *HTTPClient) UnSubscribe(topic ...string) error {
	return nil
}

func (sf *HTTPClient) LogProvider() clog.LogProvider {
	return sf.c.Clog
}

func (sf *HTTPClient) LogMode(enable bool) {
	sf.c.LogMode(enable)
}

func (sf *HTTPClient) UnderlyingClient() interface{} {
	return sf.c
}

func NewWithHTTP(config *dm.Config) *HTTPClient {
	client := ahttp.New().
		SetDeviceMetaInfo(config.
			EnableHTTP(true).
			MetaInfo())

	sf := dm.New(config)
	httpClient := &HTTPClient{c: client, Client: sf}
	sf.SetConn(httpClient)
	return httpClient
}
