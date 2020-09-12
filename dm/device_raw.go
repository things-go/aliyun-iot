package dm

import (
	"github.com/thinkgos/aliyun-iot/infra"
)

// ThingModelUpRaw 上传透传数据
// request: /sys/{productKey}/{deviceName}/thing/model/up_raw
// response: /sys/{productKey}/{deviceName}/thing/model/up_raw_reply
func (sf *Client) ThingModelUpRaw(devID int, payload interface{}) error {
	if !sf.hasRawModel {
		return ErrNotSupportFeature
	}
	if devID < 0 {
		return ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return err
	}
	uri := infra.URI(infra.URISysPrefix, infra.URIThingModelUpRaw, node.ProductKey(), node.DeviceName())
	if err = sf.Publish(uri, 1, payload); err != nil {
		return err
	}
	sf.log.Debugf("upstream thing <model>: up raw")
	return nil
}

// ProcThingModelUpRawReply 处理透传上行的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/model/up_raw
// response: /sys/{productKey}/{deviceName}/thing/model/up_raw_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/model/up_raw_reply
func ProcThingModelUpRawReply(c *Client, rawURI string, payload []byte) error {
	uris := infra.URISpilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	c.log.Debugf("downstream thing <model>: up raw reply")
	pk, dn := uris[1], uris[2]
	return c.cb.ThingModelUpRawReply(c, pk, dn, payload)
}

// ProcThingModelDownRaw 处理透传下行数据
// 下行
// request: /sys/{productKey}/{deviceName}/thing/model/down_raw
// response: /sys/{productKey}/{deviceName}/thing/model/down_raw_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/model/down_raw
func ProcThingModelDownRaw(c *Client, rawURI string, payload []byte) error {
	uris := infra.URISpilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	c.log.Debugf("downstream thing <model>: down raw request")
	pk, dn := uris[1], uris[2]
	return c.cb.ThingModelDownRaw(c, pk, dn, payload)
}
