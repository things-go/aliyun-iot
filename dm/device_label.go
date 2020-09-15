package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// ThingDeviceInfoUpdate 设备信息上传(如厂商、设备型号等，可以保存为设备标签)
// request:  /sys/{productKey}/{deviceName}/thing/deviceinfo/update
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
func (sf *Client) ThingDeviceInfoUpdate(pk, dn string, params interface{}) (*Token, error) {
	id := sf.nextRequestID()
	_uri := uri.URI(uri.SysPrefix, uri.ThingDeviceInfoUpdate, pk, dn)
	if err := sf.Request(_uri, id, infra.MethodDeviceInfoUpdate, params); err != nil {
		return nil, err
	}
	sf.log.Debugf("thing <deviceInfo>: update, @%d", id)
	return sf.putPending(id), nil
}

// ThingDeviceInfoDelete 设备信息删除
// request:  /sys/{productKey}/{deviceName}/thing/deviceinfo/delete
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
func (sf *Client) ThingDeviceInfoDelete(pk, dn string, params interface{}) (*Token, error) {
	id := sf.nextRequestID()
	_uri := uri.URI(uri.SysPrefix, uri.ThingDeviceInfoDelete, pk, dn)
	if err := sf.Request(_uri, id, infra.MethodDeviceInfoDelete, params); err != nil {
		return nil, err
	}
	sf.log.Debugf("thing <deviceInfo>: delete, @%d", id)
	return sf.putPending(id), nil
}

// ProcThingDeviceInfoUpdateReply 处理设备信息更新应答
// request:   /sys/{productKey}/{deviceName}/thing/deviceinfo/update
// response:  /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
func ProcThingDeviceInfoUpdateReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.log.Debugf("thing <deviceInfo>: update reply, @%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, nil, err})
	pk, dn := uris[1], uris[2]
	return c.cb.ThingDeviceInfoUpdateReply(c, err, pk, dn)
}

// ProcThingDeviceInfoDeleteReply 处理设备信息删除的应答
// request:   /sys/{productKey}/{deviceName}/thing/deviceinfo/delete
// response:  /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
func ProcThingDeviceInfoDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, nil, err})
	pk, dn := uris[1], uris[2]
	c.log.Debugf("thing <deviceInfo>: delete reply, @%d", rsp.ID)
	return c.cb.ThingDeviceInfoDeleteReply(c, err, pk, dn)
}
