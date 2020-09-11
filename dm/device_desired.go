package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
)

// upstreamThingDesiredPropertyGet 获取期望值
// request: /sys/{productKey}/{deviceName}/thing/property/desired/get
// response: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
func (sf *Client) upstreamThingDesiredPropertyGet(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDesiredPropertyGet, node.ProductKey(), node.DeviceName()),
		id, MethodDesiredPropertyGet, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDesiredPropertyGet)
	sf.debugf("upstream thing <desired>: property get,@%d", id)
	return nil
}

// upstreamThingDesiredPropertyDelete 清空期望值
// request: /sys/{productKey}/{deviceName}/thing/property/desired/delete
// response: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
func (sf *Client) upstreamThingDesiredPropertyDelete(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNode(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDesiredPropertyDelete, node.ProductKey(), node.DeviceName()),
		id, MethodDesiredPropertyDelete, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDesiredPropertyDelete)
	sf.debugf("upstream thing <desired>: property delete,@%d", id)
	return nil
}

// ProcThingDesiredPropertyGetReply 处理期望属性获取的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/property/desired/get
// response: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
func ProcThingDesiredPropertyGetReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 7) {
		return ErrInvalidURI
	}
	rsp := ResponseRawData{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debugf("downstream thing <desired>: property get reply,@%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.debugf("downstream thing <desired>: property get reply,@%d", rsp.ID)
	return c.eventProc.EvtThingDesiredPropertyGetReply(c, err, pk, dn, rsp.Data)
}

// ProcThingDesiredPropertyDeleteReply 处理期望属性删除的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/property/desired/delete
// response: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
func ProcThingDesiredPropertyDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 7) {
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

	c.CacheDone(rsp.ID, err)
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.debugf("downstream thing <desired>: property delete reply,@%d", rsp.ID)
	return c.eventProc.EvtThingDesiredPropertyDeleteReply(c, err, pk, dn)
}
