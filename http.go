package aliIOT

import (
	"github.com/thinkgos/aliIOT/ahttp"
	"github.com/thinkgos/aliIOT/clog"
	"github.com/thinkgos/aliIOT/model"
)

type httpClient struct {
	c         *ahttp.Client
	containOf *model.Manager
}

func (sf *httpClient) Publish(topic string, qos byte, payload interface{}) error {
	return sf.c.Publish(topic, payload)
}

func (sf *httpClient) UnderlyingClient() interface{} {
	return sf.c
}

func (sf *httpClient) Subscribe(topic string, streamFunc model.ProcDownStreamFunc) error {
	return nil
}

func (sf *httpClient) LogProvider() clog.LogProvider {
	return sf.c.Clog
}

func (sf *httpClient) LogMode(enable bool) {
	sf.c.LogMode(enable)
}

func NewWithHTTP(config *model.Config) *Client {
	client := ahttp.New().
		SetDeviceMetaInfo(config.
			EnableHTTP(true).
			MetaInfo())

	sf := model.New(config)

	return &Client{
		sf.SetConn(&httpClient{
			c:         client,
			containOf: sf}),
		config.FeatureOption(),
	}
}
