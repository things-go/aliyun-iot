// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package coap

import (
	"bytes"
	"errors"

	"github.com/go-ocf/go-coap"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/57697.html?spm=a2c4g.11186623.6.606.5d7a12e0FGY05a

// 确保 NopCb 实现 dm.Conn 接口
var _ aiot.Conn = (*coapClient)(nil)

// COAPClient COAP客户端
type coapClient struct {
	c *coap.ClientConn
}

// NewCOAP new coap con
func NewCOAP(c *coap.ClientConn) aiot.Conn {
	return &coapClient{c}
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
func (*coapClient) Subscribe(string, aiot.ProcDownStream) error { return nil }

// UnSubscribe 实现dm.Conn接口
func (sf *coapClient) UnSubscribe(...string) error { return nil }

// Close 实现dm.Conn接口
func (sf *coapClient) Close() error { return sf.c.Close() }
