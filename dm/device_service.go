package dm

import "github.com/thinkgos/aliyun-iot/uri"

// ProcThingServicePropertySet 处理属性设置
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/service/property/set
// response:  /sys/{productKey}/{deviceName}/thing/service/property/set_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/service/[+,#]
func ProcThingServicePropertySet(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	c.log.Debugf("thing <service>: property set request")
	pk, dn := uris[1], uris[2]
	return c.cb.ThingServicePropertySet(c, pk, dn, payload)
}

// ProcThingServiceRequest 处理设备服务调用(异步)
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}
// response:  /sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/service/[+,#]
func ProcThingServiceRequest(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	pk, dn := uris[1], uris[2]
	serviceID := uris[5]
	c.log.Debugf("thing <service>: %s set request", serviceID)
	return c.cb.ThingServiceRequest(c, serviceID, pk, dn, payload)
}
