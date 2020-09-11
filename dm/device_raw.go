package dm

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
	err = sf.Publish(sf.URIService(URISysPrefix, URIThingModelUpRaw, node.ProductKey(), node.DeviceName()), 1, payload)
	if err != nil {
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
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.log.Debugf("downstream thing <model>: up raw reply")
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	return c.cb.ThingModelUpRawReply(c, pk, dn, payload)
}

// ProcThingModelDownRaw 处理透传下行数据
// 下行
// request: /sys/{productKey}/{deviceName}/thing/model/down_raw
// response: /sys/{productKey}/{deviceName}/thing/model/down_raw_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/model/down_raw
func ProcThingModelDownRaw(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.log.Debugf("downstream thing <model>: down raw request")
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	return c.cb.ThingModelDownRaw(c, pk, dn, payload)
}
