package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/89304.html?spm=a2c4g.11186623.6.710.31c552ceVRAsmU

// ThingDeviceInfoUpdate 设备信息上传(如厂商、设备型号等，可以保存为设备标签)
// request:  /sys/{productKey}/{deviceName}/thing/deviceinfo/update
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
func (sf *Client) ThingDeviceInfoUpdate(pk, dn string, params interface{}) (*Token, error) {
	_uri := uri.URI(uri.SysPrefix, uri.ThingDeviceInfoUpdate, pk, dn)
	return sf.SendRequest(_uri, infra.MethodDeviceInfoUpdate, params)
}

// ThingDeviceInfoDelete 删除标签信息
// request:  /sys/{productKey}/{deviceName}/thing/deviceinfo/delete
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
func (sf *Client) ThingDeviceInfoDelete(pk, dn string, params interface{}) (*Token, error) {
	_uri := uri.URI(uri.SysPrefix, uri.ThingDeviceInfoDelete, pk, dn)
	return sf.SendRequest(_uri, infra.MethodDeviceInfoDelete, params)
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

	c.Log.Debugf("thing.deviceinfo.update.reply @%d", rsp.ID)
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
	c.Log.Debugf("thing.deviceinfo.delete.reply @%d", rsp.ID)
	return c.cb.ThingDeviceInfoDeleteReply(c, err, pk, dn)
}
