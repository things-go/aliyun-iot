package dm

import (
	uri "github.com/thinkgos/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/89301.html?spm=a2c4g.11186623.6.706.570f3f69J3fW5z

// ThingModelUpRaw 上传透传数据
// request: /sys/{productKey}/{deviceName}/thing/model/up_raw
// response: /sys/{productKey}/{deviceName}/thing/model/up_raw_reply
func (sf *Client) ThingModelUpRaw(pk, dn string, payload interface{}) error {
	if !sf.hasRawModel {
		return ErrNotSupportFeature
	}
	sf.log.Debugf("thing.model.up.raw")
	_uri := uri.URI(uri.SysPrefix, uri.ThingModelUpRaw, pk, dn)
	return sf.Publish(_uri, 1, payload)
}

// ProcThingModelUpRawReply 处理透传上行的应答
// request: /sys/{productKey}/{deviceName}/thing/model/up_raw
// response: /sys/{productKey}/{deviceName}/thing/model/up_raw_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/model/up_raw_reply
func ProcThingModelUpRawReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	c.log.Debugf("thing.model.up.raw.reply")
	pk, dn := uris[1], uris[2]
	return c.cb.ThingModelUpRawReply(c, pk, dn, payload)
}

// ProcThingModelDownRaw 处理透传下行数据
// 下行
// request: /sys/{productKey}/{deviceName}/thing/model/down_raw
// response: /sys/{productKey}/{deviceName}/thing/model/down_raw_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/model/down_raw
func ProcThingModelDownRaw(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	c.log.Debugf("thing.model.down.raw")
	pk, dn := uris[1], uris[2]
	return c.cb.ThingModelDownRaw(c, pk, dn, payload)
}
