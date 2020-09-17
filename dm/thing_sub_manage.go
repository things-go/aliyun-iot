package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// ProcThingDisable 禁用子设备
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/disable
// response:  /sys/{productKey}/{deviceName}/thing/disable_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/disable
func ProcThingDisable(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 5 {
		return ErrInvalidURI
	}

	pk, dn := uris[1], uris[2]
	c.log.Debugf("thing.disable -- %s.%s", pk, dn)

	req := &Request{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		return err
	}

	if err = c.SetDeviceAvail(pk, dn, false); err != nil {
		c.log.Warnf("thing.disable failed, %+v", err)
	}

	err = c.Response(uri.ReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
	if err != nil {
		c.log.Warnf("thing.disable.response failed, %+v", err)
	}
	return c.gwCb.ThingDisable(c, pk, dn)
}

// ProcThingEnable 启用子设备
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/enable
// response:  /sys/{productKey}/{deviceName}/thing/enable_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/enable
func ProcThingEnable(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 5 {
		return ErrInvalidURI
	}
	pk, dn := uris[1], uris[2]
	c.log.Debugf("thing.enable -- %s.%s", pk, dn)

	req := &Request{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		return err
	}

	if err = c.SetDeviceAvail(pk, dn, true); err != nil {
		c.log.Warnf("thing.enable failed, %+v", err)
	}

	_uri := uri.ReplyWithRequestURI(rawURI)
	if err = c.Response(_uri, req.ID, infra.CodeSuccess, "{}"); err != nil {
		c.log.Warnf("thing.enable.response failed, %+v", err)
	}
	return c.gwCb.ThingEnable(c, pk, dn)
}

// ProcThingDelete 子设备删除,网关类型设备
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/delete
// response:  /sys/{productKey}/{deviceName}/thing/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/delete
func ProcThingDelete(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 5 {
		return ErrInvalidURI
	}
	pk, dn := uris[1], uris[2]
	c.log.Debugf("thing.delete -- %s.%s", pk, dn)

	req := &Request{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}

	c.Delete(pk, dn)
	_uri := uri.ReplyWithRequestURI(rawURI)
	err := c.Response(_uri, req.ID, infra.CodeSuccess, "{}")
	if err != nil {
		c.log.Warnf("thing.delete.response failed, %+v", err)
	}
	return c.gwCb.ThingDelete(c, pk, dn)
}
