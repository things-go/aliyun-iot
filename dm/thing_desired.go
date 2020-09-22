package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/109807.html?spm=a2c4g.11186623.6.707.31c552ceZhSvWp

// ThingDesiredPropertyGet 获取期望属性值
// request:  /sys/{productKey}/{deviceName}/thing/property/desired/get
// response: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
func (sf *Client) ThingDesiredPropertyGet(pk, dn string, params []string) (*Token, error) {
	if !sf.hasDesired {
		return nil, ErrNotSupportFeature
	}
	if !sf.IsActive(pk, dn) {
		return nil, ErrNotActive
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingDesiredPropertyGet, pk, dn)
	return sf.SendRequest(_uri, infra.MethodDesiredPropertyGet, params)
}

// ThingDesiredPropertyDelete 清空期望属性值
// request:  /sys/{productKey}/{deviceName}/thing/property/desired/delete
// response: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
func (sf *Client) ThingDesiredPropertyDelete(pk, dn string, params interface{}) (*Token, error) {
	if !sf.hasDesired {
		return nil, ErrNotSupportFeature
	}
	if !sf.IsActive(pk, dn) {
		return nil, ErrNotActive
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingDesiredPropertyDelete, pk, dn)
	return sf.SendRequest(_uri, infra.MethodDesiredPropertyDelete, params)
}

// ProcThingDesiredPropertyGetReply 处理获取期望属性值的应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/property/desired/get
// response:  /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
func ProcThingDesiredPropertyGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	rsp := &ResponseRawData{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, cloneJSONRawMessage(rsp.Data), err})

	c.Log.Debugf("thing.property.desired.get.reply @%d", rsp.ID)

	pk, dn := uris[1], uris[2]
	return c.cb.ThingDesiredPropertyGetReply(c, err, pk, dn, rsp.Data)
}

// ProcThingDesiredPropertyDeleteReply 处理清空期望属性值的应答
// request:   /sys/{productKey}/{deviceName}/thing/property/desired/delete
// response:  /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
func ProcThingDesiredPropertyDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	rsp := &Response{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, nil, err})

	c.Log.Debugf("thing.property.desired.delete.reply @%d", rsp.ID)

	pk, dn := uris[1], uris[2]
	return c.cb.ThingDesiredPropertyDeleteReply(c, err, pk, dn)
}
