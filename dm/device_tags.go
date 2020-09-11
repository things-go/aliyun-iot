package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
)

// ThingDeviceInfoUpdate 设备信息上传(如厂商、设备型号等，可以保存为设备标签)
// request:  /sys/{productKey}/{deviceName}/thing/deviceinfo/update
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
func (sf *Client) ThingDeviceInfoUpdate(devID int, params interface{}) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDeviceInfoUpdate, node.ProductKey(), node.DeviceName()),
		id, MethodDeviceInfoUpdate, params)
	if err != nil {
		return nil, err
	}

	sf.debugf("upstream thing <deviceInfo>: update,@%d", id)
	return sf.Insert(id), nil
}

// ThingDeviceInfoDelete 设备信息删除
// request: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
func (sf *Client) ThingDeviceInfoDelete(devID int, params interface{}) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDeviceInfoDelete, node.ProductKey(), node.DeviceName()),
		id, MethodDeviceInfoDelete, params)
	if err != nil {
		return nil, err
	}
	sf.debugf("upstream thing <deviceInfo>: delete,@%d", id)
	return sf.Insert(id), nil
}

// ProcThingDeviceInfoUpdateReply 处理设备信息更新应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/deviceinfo/update
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
func ProcThingDeviceInfoUpdateReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := ResponseRawData{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debugf("downstream thing <deviceInfo>: update reply,@%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.done(rsp.ID, err, nil)
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	return c.eventProc.EvtThingDeviceInfoUpdateReply(c, err, pk, dn)
}

// ProcThingDeviceInfoDeleteReply 处理设备信息删除的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
func ProcThingDeviceInfoDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
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
	c.done(rsp.ID, err, nil)
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.debugf("downstream thing <deviceInfo>: delete reply,@%d", rsp.ID)
	return c.eventProc.EvtThingDeviceInfoDeleteReply(c, err, pk, dn)
}
