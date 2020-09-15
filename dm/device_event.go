package dm

import (
	"encoding/json"
	"fmt"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// ThingEventPropertyPost 上传属性数据
// request:  /sys/{productKey}/{deviceName}/thing/event/property/post
// response: /sys/{productKey}/{deviceName}/thing/event/property/post_reply
func (sf *Client) ThingEventPropertyPost(pk, dn string, params interface{}) (*Token, error) {
	if sf.hasRawModel {
		return nil, ErrNotSupportFeature
	}

	id := sf.RequestID()
	_uri := uri.URI(uri.SysPrefix, uri.ThingEventPropertyPost, pk, dn)
	err := sf.SendRequest(_uri, id, infra.MethodEventPropertyPost, params)
	if err != nil {
		return nil, err
	}
	sf.log.Debugf("thing <event>: property post, @%d", id)
	return sf.putPending(id), nil
}

// ThingEventPost 事件上传
// request:  /sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post
// response: /sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post_reply
func (sf *Client) ThingEventPost(pk, dn, eventID string, params interface{}) (*Token, error) {
	id := sf.RequestID()
	_uri := uri.URI(uri.SysPrefix, uri.ThingEventPost, pk, dn, eventID)
	err := sf.SendRequest(_uri, id, fmt.Sprintf(infra.MethodEventFormatPost, eventID), params)
	if err != nil {
		return nil, err
	}
	sf.log.Debugf("thing <event>: %s post, @%d", eventID, id)
	return sf.putPending(id), nil
}

// ThingEventPropertyPackPost 网关批量上报数据
// NOTE: 仅网关支持,一次最多200个属性,20个事件,一次最多为20个子设备上报数据
// request:  /sys/{productKey}/{deviceName}/thing/event/property/pack/post
// response: /sys/{productKey}/{deviceName}/thing/event/property/pack/post_reply
func (sf *Client) ThingEventPropertyPackPost(params interface{}) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	id := sf.RequestID()
	_uri := sf.GatewayURI(uri.SysPrefix, uri.ThingEventPropertyPackPost)
	err := sf.SendRequest(_uri, id, infra.MethodEventPropertyPackPost, params)
	if err != nil {
		return nil, err
	}
	sf.log.Debugf("thing <event>: property pack post, @%d", id)
	return sf.putPending(id), nil
}

// ThingEventPropertyHistoryPost 直连设备仅能上报自己的物模型历史数据,网关设备可以上报其子设备的物模型历史数据
// request： /sys/{productKey}/{deviceName}/thing/event/property/history/post
// response：/sys/{productKey}/{deviceName}/thing/event/property/history/post_reply
func (sf *Client) ThingEventPropertyHistoryPost(params interface{}) (*Token, error) {
	id := sf.RequestID()
	_uri := sf.GatewayURI(uri.SysPrefix, uri.ThingEventPropertyHistoryPost)
	err := sf.SendRequest(_uri, id, infra.MethodEventPropertyHistoryPost, params)
	if err != nil {
		return nil, err
	}
	sf.log.Debugf("thing <event>: property history post, @%d", id)
	return sf.putPending(id), nil
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
	c.log.Debugf("thing <event>: %s post reply, @%d", eventID, rsp.ID)
	if eventID == "property" {
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
	c.log.Debugf("thing <event>: property pack post reply, @%d", rsp.ID)
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
	c.log.Debugf("thing <event>: property pack post reply, @%d", rsp.ID)
	return c.cb.ThingEventPropertyHistoryPostReply(c, err, pk, dn)
}
