// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import (
	"encoding/json"
	"fmt"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/89301.html?spm=a2c4g.11186623.6.706.78b524baCoL1Gf

// ThingEventPropertyPost 设备上报属性数据
// request:  /sys/{productKey}/{deviceName}/thing/event/property/post
// response: /sys/{productKey}/{deviceName}/thing/event/property/post_reply
func (sf *Client) ThingEventPropertyPost(pk, dn string, params interface{}) (*Token, error) {
	if sf.hasRawModel {
		return nil, ErrNotSupportFeature
	}
	if !sf.IsActive(pk, dn) {
		return nil, ErrNotActive
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingEventPropertyPost, pk, dn)
	return sf.SendRequest(_uri, infra.MethodEventPropertyPost, params)
}

// ThingEventPost 设备事件上报
// request:  /sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post
// response: /sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post_reply
func (sf *Client) ThingEventPost(pk, dn, eventID string, params interface{}) (*Token, error) {
	if !sf.IsActive(pk, dn) {
		return nil, ErrNotActive
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingEventPost, pk, dn, eventID)
	method := fmt.Sprintf(infra.MethodEventFormatPost, eventID)
	return sf.SendRequest(_uri, method, params)
}

// ThingEventPropertyPackPost 网关批量上报数据
// NOTE: 仅网关支持,一次最多200个属性,20个事件,一次最多为20个子设备上报数据
// request:  /sys/{productKey}/{deviceName}/thing/event/property/pack/post
// response: /sys/{productKey}/{deviceName}/thing/event/property/pack/post_reply
func (sf *Client) ThingEventPropertyPackPost(params interface{}) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	if !sf.IsActive(sf.tetrad.ProductKey, sf.tetrad.DeviceName) {
		return nil, ErrNotActive
	}
	_uri := sf.URIGateway(uri.SysPrefix, uri.ThingEventPropertyPackPost)
	return sf.SendRequest(_uri, infra.MethodEventPropertyPackPost, params)
}

// ThingEventPropertyHistoryPost  物模型历史数据上报
// 直连设备仅能上报自己的物模型历史数据,网关设备可以上报其子设备的物模型历史数据
// request： /sys/{productKey}/{deviceName}/thing/event/property/history/post
// response：/sys/{productKey}/{deviceName}/thing/event/property/history/post_reply
func (sf *Client) ThingEventPropertyHistoryPost(params interface{}) (*Token, error) {
	if !sf.IsActive(sf.tetrad.ProductKey, sf.tetrad.DeviceName) {
		return nil, ErrNotActive
	}
	_uri := sf.URIGateway(uri.SysPrefix, uri.ThingEventPropertyHistoryPost)
	return sf.SendRequest(_uri, infra.MethodEventPropertyHistoryPost, params)
}

// ProcThingEventPostReply 处理ThingEvent XXX上行的应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/event/[{tsl.event.identifier},property]/post
// response:  /sys/{productKey}/{deviceName}/thing/event/[{tsl.event.identifier},property]/post_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/event/+/post_reply
func ProcThingEventPostReply(c *Client, rawURI string, payload []byte) error {
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

	pk, dn := uris[1], uris[2]
	eventID := uris[5]
	c.Log.Debugf("thing.event.%s.post.reply @%d", eventID, rsp.ID)
	if eventID == property {
		return c.cb.ThingEventPropertyPostReply(c, err, pk, dn)
	}
	return c.cb.ThingEventPostReply(c, err, eventID, pk, dn)
}

// ProcThingEventPropertyPackPostReply 网关批量上报数据
// 上行,仅网关支持
// request:   /sys/{productKey}/{deviceName}/thing/event/property/pack/post
// response:  /sys/{productKey}/{deviceName}/thing/event/property/pack/post_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/event/property/pack/post_reply
func ProcThingEventPropertyPackPostReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 8 {
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
	pk, dn := uris[1], uris[2]
	c.Log.Debugf("thing.event.property.pack.post.reply @%d", rsp.ID)
	return c.cb.ThingEventPropertyPackPostReply(c, err, pk, dn)
}

// ProcThingEventPropertyHistoryPostReply 物模型历史数据上报应答
// request：  /sys/{productKey}/{deviceName}/thing/event/property/history/post
// response： /sys/{productKey}/{deviceName}/thing/event/property/history/post_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/event/property/history/post_reply
func ProcThingEventPropertyHistoryPostReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 8 {
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
	pk, dn := uris[1], uris[2]
	c.Log.Debugf("thing.event.property.history.post.reply @%d", rsp.ID)
	return c.cb.ThingEventPropertyHistoryPostReply(c, err, pk, dn)
}
