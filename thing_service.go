// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import "github.com/thinkgos/aliyun-iot/uri"

// ProcThingServiceRequest 处理设备服务调用(异步)
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/service/[{tsl.service.identifier},property/set]
// response:  /sys/{productKey}/{deviceName}/thing/service/[{tsl.service.identifier}_reply,property/set_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/service/[+,#]
func ProcThingServiceRequest(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	pk, dn := uris[1], uris[2]
	serviceID := uris[5]
	if serviceID == property && len(uris) >= 7 && uris[6] == "set" {
		c.Log.Debugf("thing.service.property.set")
		return c.cb.ThingServicePropertySet(c, pk, dn, payload)
	}
	c.Log.Debugf("thing.service.%s", serviceID)
	return c.cb.ThingServiceRequest(c, serviceID, pk, dn, payload)
}
