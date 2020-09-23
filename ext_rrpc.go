// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import (
	"strings"

	"github.com/thinkgos/aliyun-iot/uri"
)

// RRPCResponse rrcpc 回复
// response: /sys/${YourProductKey}/${YourDeviceName}/rrpc/response/${messageId}
func (sf *Client) RRPCResponse(pk, dn, messageID string, rsp Response) error {
	_uri := uri.URI(uri.SysPrefix, uri.RRPCResponse, pk, dn, messageID)
	return sf.Response(_uri, rsp)
}

// ExtRRPCResponse ext rrpc回复
// response:  /ext/rrpc/${messageId}/${topic}
// see https://help.aliyun.com/document_detail/90570.html?spm=a2c4g.11186623.6.656.64076175n5VFHO#title-0r5-s8c-t1c
func (sf *Client) ExtRRPCResponse(messageID, topic string, payload interface{}) error {
	_uri := uri.ExtRRPC(messageID, topic)
	return sf.Publish(_uri, 0, payload)
}

// ProcRRPCRequest 处理RRPC请求
// request:   /sys/${YourProductKey}/${YourDeviceName}/rrpc/request/${messageId}
// response:  /sys/${YourProductKey}/${YourDeviceName}/rrpc/response/${messageId}
// subscribe: /sys/${YourProductKey}/${YourDeviceName}/rrpc/request/+
func ProcRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	pk, dn := uris[1], uris[2]
	messageID := uris[5]
	c.Log.Debugf("rrpc.request.%s", messageID)
	return c.cb.RRPCRequest(c, messageID, pk, dn, payload)
}

// ProcExtRRPCRequest 处理扩展RRPC请求
// ${topic} 不为空,设备建立要求clientID传ext = 1
// request:   /ext/rrpc/${messageId}/${topic}
// response:  /ext/rrpc/${messageId}/${topic}
// subscribe: /ext/rrpc/+/${topic}
// 			  /ext/rrpc/#
func ProcExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := strings.SplitN(strings.TrimLeft(rawURI, uri.Sep), uri.Sep, 4)
	if len(uris) < 3 {
		return ErrInvalidParameter
	}
	messageID, topic := uris[2], uris[3]
	c.Log.Debugf("ext.rrpc.%s -- topic: %s", messageID, topic)
	return c.cb.ExtRRPCRequest(c, messageID, topic, payload)
}
