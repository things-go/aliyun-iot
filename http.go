package aliIOT

import (
	"github.com/thinkgos/aliIOT/ahttp"
	"github.com/thinkgos/aliIOT/clog"
	"github.com/thinkgos/aliIOT/model"
)

type HTTPClient struct {
	c *ahttp.Client
	*model.Manager
}

func (sf *HTTPClient) Publish(topic string, qos byte, payload interface{}) error {
	return sf.c.Publish(topic, payload)
}

func (sf *HTTPClient) Subscribe(topic string, streamFunc model.ProcDownStreamFunc) error {
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

func NewWithHTTP(config *model.Config) *HTTPClient {
	client := ahttp.New().
		SetDeviceMetaInfo(config.
			EnableHTTP(true).
			MetaInfo())

	sf := model.New(config)
	httpClient := &HTTPClient{c: client, Manager: sf}
	sf.SetConn(httpClient)
	return httpClient
}
