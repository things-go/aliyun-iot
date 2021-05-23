// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import (
	"encoding/json"

	"github.com/things-go/aliyun-iot/uri"
)

// URIGateway 获得本设备网关URI
func (sf *Client) URIGateway(prefix, name string) string {
	return uri.URI(prefix, name, sf.tetrad.ProductKey, sf.tetrad.DeviceName)
}

func dupJSONRawMessage(d json.RawMessage) json.RawMessage {
	v := make(json.RawMessage, len(d))
	copy(v, d)
	return v
}
