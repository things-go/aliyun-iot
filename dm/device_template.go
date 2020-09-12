package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
)

// ThingDsltemplateGet 设备可以通过上行请求获取设备的TSL模板（包含属性、服务和事件的定义）
// see https://help.aliyun.com/document_detail/89305.html?spm=a2c4g.11186623.6.672.5d3d70374hpPcx
// request:   /sys/{productKey}/{deviceName}/thing/dsltemplate/get
// response:  /sys/{productKey}/{deviceName}/thing/dsltemplate/get_reply
func (sf *Client) ThingDsltemplateGet(devID int) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	uri := infra.URI(infra.URISysPrefix, infra.URIThingDslTemplateGet, node.ProductKey(), node.DeviceName())
	err = sf.SendRequest(uri, id, infra.MethodDslTemplateGet, "{}")
	if err != nil {
		return nil, err
	}

	sf.log.Debugf("upstream thing <dsl template>: get,@%d", id)
	return sf.Insert(id), nil
}

// ThingDynamictslGet 获取动态tsl
func (sf *Client) ThingDynamictslGet(devID int) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	uri := infra.URI(infra.URISysPrefix, infra.URIThingDynamicTslGet, node.ProductKey(), node.DeviceName())
	err = sf.SendRequest(uri, id, infra.MethodDynamicTslGet, map[string]interface{}{
		"nodes":      []string{"type", "identifier"},
		"addDefault": false,
	})
	if err != nil {
		return nil, err
	}

	sf.log.Debugf("upstream thing <dynamic tsl>: get,@%d", id)
	return sf.Insert(id), nil
}

// ProcThingDsltemplateGetReply 处理dsltemplate获取的应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/dsltemplate/get
// response:  /sys/{productKey}/{deviceName}/thing/dsltemplate/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/dsltemplate/get_reply
func ProcThingDsltemplateGetReply(c *Client, rawURI string, payload []byte) error {
	uris := infra.URISpilt(rawURI)
	if len(uris) < 6 {
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
	c.log.Debugf("downstream thing <dsl template>: get reply,@%d - %s", rsp.ID, string(rsp.Data))
	return c.cb.ThingDsltemplateGetReply(c, err, pk, dn, rsp.Data)
}

// ProcThingDynamictslGetReply 处理
// 上行
// request: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get
// response: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get_reply
// subscribe: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get_reply
func ProcThingDynamictslGetReply(c *Client, rawURI string, payload []byte) error {
	uris := infra.URISpilt(rawURI)
	if len(uris) < 6 {
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
	c.log.Debugf("downstream thing <dynamic tsl>: get reply,@%d - %+v", rsp.ID, rsp)
	return c.cb.ThingDynamictslGetReply(c, err, pk, dn, rsp.Data)
}
