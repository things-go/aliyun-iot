package dm

import (
	"encoding/json"
	"fmt"

	"github.com/thinkgos/aliyun-iot/infra"
)

// ThingEventPropertyPost 上传属性数据
// request:  /sys/{productKey}/{deviceName}/thing/event/property/post
// response: /sys/{productKey}/{deviceName}/thing/event/property/post_reply
func (sf *Client) ThingEventPropertyPost(devID int, params interface{}) (*Entry, error) {
	if sf.hasRawModel {
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
	uri := infra.URI(infra.URISysPrefix, infra.URIThingEventPropertyPost, node.ProductKey(), node.DeviceName())
	err = sf.SendRequest(uri, id, infra.MethodEventPropertyPost, params)
	if err != nil {
		return nil, err
	}

	sf.log.Debugf("upstream thing <event>: property post,@%d", id)
	return sf.Insert(id), nil
}

// ThingEventPost 事件上传
// request:  /sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post
// response: /sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post_reply
func (sf *Client) ThingEventPost(devID int, eventID string, params interface{}) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	uri := infra.URI(infra.URISysPrefix, infra.URIThingEventPost, node.ProductKey(), node.DeviceName(), eventID)
	err = sf.SendRequest(uri, id, fmt.Sprintf(infra.MethodEventFormatPost, eventID), params)
	if err != nil {
		return nil, err
	}
	sf.log.Debugf("upstream thing <event>: %s post,@%d", eventID, id)
	return sf.Insert(id), nil
}

// ThingEventPropertyPackPost 网关批量上报数据
// NOTE: 仅网关支持,一次最多200个属性,20个事件,一次最多为20个子设备上报数据
// request:  /sys/{productKey}/{deviceName}/thing/event/property/pack/post
// response: /sys/{productKey}/{deviceName}/thing/event/property/pack/post_reply
func (sf *Client) ThingEventPropertyPackPost(params interface{}) (*Entry, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	id := sf.RequestID()
	err := sf.SendRequest(sf.URIGateway(infra.URISysPrefix, infra.URIThingEventPropertyPackPost),
		id, infra.MethodEventPropertyPackPost, params)
	if err != nil {
		return nil, err
	}

	sf.log.Debugf("upstream thing <deviceInfo>: update,@%d", id)
	return sf.Insert(id), nil
}

// ProcThingEventPostReply 处理ThingEvent XXX上行的应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/event/[{tsl.event.identifier},property]/post
// response:  /sys/{productKey}/{deviceName}/thing/event/[{tsl.event.identifier},property]/post_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/event/+/post_reply
func ProcThingEventPostReply(c *Client, rawURI string, payload []byte) error {
	uris := infra.URISpilt(rawURI)
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
	c.done(rsp.ID, err, nil)

	pk, dn := uris[1], uris[2]
	eventID := uris[5]
	c.log.Debugf("downstream thing <event>: %s post reply,@%d", eventID, rsp.ID)
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
	uris := infra.URISpilt(rawURI)
	if len(uris) < 8 {
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
	pk, dn := uris[1], uris[2]
	c.log.Debugf("downstream thing <event>: property pack post reply,@%d", rsp.ID)
	return c.cb.ThingEventPropertyPackPostReply(c, err, pk, dn)
}

// ProcThingServicePropertySet 处理属性设置
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/service/property/set
// response:  /sys/{productKey}/{deviceName}/thing/service/property/set_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/service/[+,#]
func ProcThingServicePropertySet(c *Client, rawURI string, payload []byte) error {
	uris := infra.URISpilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	c.log.Debugf("downstream thing <service>: property set request")
	pk, dn := uris[1], uris[2]
	return c.cb.ThingServicePropertySet(c, pk, dn, payload)
}

// ProcThingServiceRequest 处理设备服务调用(异步)
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}
// response:  /sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/service/[+,#]
func ProcThingServiceRequest(c *Client, rawURI string, payload []byte) error {
	uris := infra.URISpilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	serviceID := uris[5]
	pk, dn := uris[1], uris[2]
	c.log.Debugf("downstream thing <service>: %s set request", serviceID)
	return c.cb.ThingServiceRequest(c, serviceID, pk, dn, payload)
}
