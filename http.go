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

// NewWithHTTP 新建HTTP客户端
func NewWithHTTP(meta infra.MetaInfo) *HTTPClient {
	c := ahttp.New(meta)
	return &HTTPClient{
		c,
		dm.New(meta, c, dm.WithWork(dm.WorkOnHTTP)),
	}
}

// UnderlyingClient 底层客户端
func (sf *HTTPClient) UnderlyingClient() *ahttp.Client { return sf.c }
