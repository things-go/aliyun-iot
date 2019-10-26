package aliIOT

import (
	"github.com/thinkgos/aliIOT/ahttp"
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

func (sf *httpClient) ContainerOf() *model.Manager {
	return sf.containOf
}

func (sf *httpClient) Subscribe(topic string, streamFunc model.ProcDownStreamFunc) error {
	return nil
}

func NewWithHTTP(options *model.Options) *model.Manager {
	client := ahttp.New().
		SetDeviceMetaInfo(options.EnableHTTP(true).MetaInfo())
	sf := model.New(options)
	return sf.SetConn(&httpClient{c: client, containOf: sf})
}
