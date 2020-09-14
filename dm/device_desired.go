package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	uri2 "github.com/thinkgos/aliyun-iot/uri"
)

// ThingDesiredPropertyGet 获取期望属性值
// request:  /sys/{productKey}/{deviceName}/thing/property/desired/get
// response: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
func (sf *Client) ThingDesiredPropertyGet(devID int, params interface{}) (*Entry, error) {
	if !sf.hasDesired {
		return nil, ErrNotSupportFeature
	}
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	uri := uri2.URI(uri2.SysPrefix, uri2.ThingDesiredPropertyGet, node.ProductKey(), node.DeviceName())
	if err := sf.SendRequest(uri, id, infra.MethodDesiredPropertyGet, params); err != nil {
		return nil, err
	}
	sf.log.Debugf("upstream thing <desired>: property get,@%d", id)
	return sf.Insert(id), nil
}

// ThingDesiredPropertyDelete 清空期望属性值
// request:  /sys/{productKey}/{deviceName}/thing/property/desired/delete
// response: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
func (sf *Client) ThingDesiredPropertyDelete(devID int, params interface{}) (*Entry, error) {
	if !sf.hasDesired {
		return nil, ErrNotSupportFeature
	}
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	uri := uri2.URI(uri2.SysPrefix, uri2.ThingDesiredPropertyDelete, node.ProductKey(), node.DeviceName())
	err = sf.SendRequest(uri, id, infra.MethodDesiredPropertyDelete, params)
	if err != nil {
		return nil, err
	}
	sf.log.Debugf("upstream thing <desired>: property delete,@%d", id)
	return sf.Insert(id), nil
}

// ProcThingDesiredPropertyGetReply 处理获取期望属性值的应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/property/desired/get
// response:  /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
func ProcThingDesiredPropertyGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri2.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	rsp := ResponseRawData{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.log.Debugf("downstream thing <desired>: property get reply,@%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signal(rsp.ID, err, nil)
	pk, dn := uris[1], uris[2]
	c.log.Debugf("downstream thing <desired>: property get reply,@%d", rsp.ID)
	return c.cb.ThingDesiredPropertyGetReply(c, err, pk, dn, rsp.Data)
}

// ProcThingDesiredPropertyDeleteReply 处理清空期望属性值的应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/property/desired/delete
// response:  /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
func ProcThingDesiredPropertyDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := uri2.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	rsp := ResponseRawData{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signal(rsp.ID, err, nil)
	pk, dn := uris[1], uris[2]
	c.log.Debugf("downstream thing <desired>: property delete reply,@%d", rsp.ID)
	return c.cb.ThingDesiredPropertyDeleteReply(c, err, pk, dn)
}
