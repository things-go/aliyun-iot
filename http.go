package aiot

import (
	"github.com/thinkgos/aliyun-iot/ahttp"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/infra"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	*dm.Client
}

// NewWithHTTP 新建HTTP客户端
func NewWithHTTP(meta infra.MetaTriad) *HTTPClient {
	return &HTTPClient{
		dm.New(meta, ahttp.New(meta), dm.WithWork(dm.WorkOnHTTP)),
	}
}

// Underlying 底层客户端
func (sf *HTTPClient) Underlying() *ahttp.Client {
	return sf.Conn.(*ahttp.Client)
}
