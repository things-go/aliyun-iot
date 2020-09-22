package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// ThingDsltemplateGet 设备可以通过上行请求获取设备的TSL模板（包含属性、服务和事件的定义）
// see https://help.aliyun.com/document_detail/89305.html?spm=a2c4g.11186623.6.672.5d3d70374hpPcx
// request:   /sys/{productKey}/{deviceName}/thing/dsltemplate/get
// response:  /sys/{productKey}/{deviceName}/thing/dsltemplate/get_reply
func (sf *Client) ThingDsltemplateGet(pk, dn string) (*Token, error) {
	if !sf.IsActive(pk, dn) {
		return nil, ErrNotActive
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingDslTemplateGet, pk, dn)
	return sf.SendRequest(_uri, infra.MethodDslTemplateGet, "{}")
}

// ThingDynamictslGet 获取动态tsl
func (sf *Client) ThingDynamictslGet(pk, dn string) (*Token, error) {
	if !sf.IsActive(pk, dn) {
		return nil, ErrNotActive
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingDynamicTslGet, pk, dn)
	return sf.SendRequest(_uri, infra.MethodDynamicTslGet, map[string]interface{}{
		"nodes":      []string{"type", "identifier"},
		"addDefault": false,
	})
}

// ProcThingDsltemplateGetReply 处理dsltemplate获取的应答
// request:   /sys/{productKey}/{deviceName}/thing/dsltemplate/get
// response:  /sys/{productKey}/{deviceName}/thing/dsltemplate/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/dsltemplate/get_reply
func ProcThingDsltemplateGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := &ResponseRawData{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, cloneJSONRawMessage(rsp.Data), err})
	c.Log.Debugf("thing.dsltemplate.get.reply @%d", rsp.ID)
	pk, dn := uris[1], uris[2]
	return c.cb.ThingDsltemplateGetReply(c, err, pk, dn, rsp.Data)
}

// ProcThingDynamictslGetReply 处理获取动态tsl应答
// request: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get
// response: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get_reply
// subscribe: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get_reply
func ProcThingDynamictslGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
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

	c.signalPending(Message{rsp.ID, cloneJSONRawMessage(rsp.Data), err})
	pk, dn := uris[1], uris[2]
	c.Log.Debugf("thing.dynamictsl.get.reply @%d", rsp.ID)
	return c.cb.ThingDynamictslGetReply(c, err, pk, dn, rsp.Data)
}
