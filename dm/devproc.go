package dm

import (
	"encoding/json"
	"strings"

	"github.com/thinkgos/aliIOT/infra"
)

// ProcDownStreamFunc 处理下行数据
type ProcDownStreamFunc func(c *Client, rawURI string, payload []byte) error

// ProcThingModelUpRawReply 处理透传上行的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/model/up_raw
// response: /sys/{productKey}/{deviceName}/thing/model/up_raw_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/model/up_raw_reply
func ProcThingModelUpRawReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.debug("downstream thing <model>: up raw reply")
	return c.ipcSendMessage(&ipcMessage{
		err:        nil,
		evt:        ipcEvtUpRawReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    payload,
	})
}

// ProcThingEventPostReply 处理ThingEvent XXX上行的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/event/[{tsl.event.identifier},property]/post
// response: /sys/{productKey}/{deviceName}/thing/event/[{tsl.event.identifier},property]/post_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/event/+/post_reply
func ProcThingEventPostReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	eventID := uris[c.cfg.uriOffset+5]
	c.debug("downstream thing <event>: %s post reply,@%d", eventID, rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	if eventID == "property" {
		return c.ipcSendMessage(&ipcMessage{
			err:        err,
			evt:        ipcEvtEventPropertyPostReply,
			productKey: uris[c.cfg.uriOffset+1],
			deviceName: uris[c.cfg.uriOffset+2],
		})
	}

	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtEventPostReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		extend:     eventID,
	})
}

// ProcThingEventPropertyPackPostReply 网关批量上报数据
// 上行,仅网关支持
// request: /sys/{productKey}/{deviceName}/thing/event/property/pack/post
// response: /sys/{productKey}/{deviceName}/thing/event/property/pack/post_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/event/property/pack/post_reply
func ProcThingEventPropertyPackPostReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 8) {
		return ErrInvalidURI
	}
	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debug("downstream thing <event>: property pack post reply,@%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtEventPropertyPostReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2]})
}

// ProcThingDeviceInfoUpdateReply 处理设备信息更新应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/deviceinfo/update
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/deviceinfo/update_reply
func ProcThingDeviceInfoUpdateReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debug("downstream thing <deviceInfo>: update reply,@%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)

	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDeviceInfoUpdateReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
	})
}

// ProcThingDeviceInfoDeleteReply 处理设备信息删除的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete
// response: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/deviceinfo/delete_reply
func ProcThingDeviceInfoDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debug("downstream thing <deviceInfo>: delete reply,@%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDeviceInfoUpdateReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
	})
}

// ProcThingDesiredPropertyGetReply 处理期望属性获取的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/property/desired/get
// response: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/property/desired/get_reply
func ProcThingDesiredPropertyGetReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}
	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debug("downstream thing <desired>: property get reply,@%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDesiredPropertyGetReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    rsp.Data,
	})
}

// ProcThingDesiredPropertyDeleteReply 处理期望属性删除的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/property/desired/delete
// response: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/property/desired/delete_reply
func ProcThingDesiredPropertyDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}
	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debug("downstream thing <desired>: property delete reply,@%d", rsp.ID)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDesiredPropertyDeleteReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
	})
}

// ProcThingDsltemplateGetReply 处理dsltemplate获取的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/dsltemplate/get
// response: /sys/{productKey}/{deviceName}/thing/dsltemplate/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/dsltemplate/get_reply
func ProcThingDsltemplateGetReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debug("downstream thing <dsl template>: get reply,@%d - %s", rsp.ID, string(rsp.Data))
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDsltemplateGetReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    rsp.Data,
	})
}

// ProcThingDynamictslGetReply 处理
// 上行
// request: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get
// response: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get_reply
// subscribe: /sys/${YourProductKey}/${YourDeviceName}/thing/dynamicTsl/get_reply
func ProcThingDynamictslGetReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debug("downstream thing <dynamic tsl>: get reply,@%d - %+v", rsp.ID, rsp)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDynamictslGetReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    rsp.Data,
	})
}

// ProcExtNtpResponse 处理ntp请求的应答
// 上行
// request: /ext/ntp/${YourProductKey}/${YourDeviceName}/request
// response: /ext/ntp/${YourProductKey}/${YourDeviceName}/response
// subscribe: /ext/ntp/${YourProductKey}/${YourDeviceName}/response
func ProcExtNtpResponse(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 5) {
		return ErrInvalidURI
	}
	rsp := NtpResponsePayload{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.debug("downstream extend <ntp>: response - %+v", rsp)
	return c.ipcSendMessage(&ipcMessage{
		err:        nil,
		evt:        ipcEvtExtNtpResponse,
		productKey: uris[c.cfg.uriOffset+2],
		deviceName: uris[c.cfg.uriOffset+3],
		payload:    rsp,
	})
}

// ProcThingConfigGetReply 处理获取配置的应答
// 上行
// request: /sys/{productKey}/{deviceName}/thing/config/get
// response: /sys/{productKey}/{deviceName}/thing/config/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/config/get_reply
func ProcThingConfigGetReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	rsp := ConfigGetResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	c.debug("downstream thing <config>: get reply,@%d,payload@%+v", rsp.ID, rsp)
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtConfigGetReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    rsp.Data,
	})
}

// ProcThingModelDownRaw 处理透传下行数据
// 下行
// request: /sys/{productKey}/{deviceName}/thing/model/down_raw
// response: /sys/{productKey}/{deviceName}/thing/model/down_raw_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/model/down_raw
func ProcThingModelDownRaw(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.debug("downstream thing <model>: down raw request")
	return c.ipcSendMessage(&ipcMessage{
		evt:        ipcEvtDownRaw,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    payload,
	})
}

// ProcThingConfigPush 处理配置推送
// 下行
// request: /sys/{productKey}/{deviceName}/thing/config/push
// response: /sys/{productKey}/{deviceName}/thing/config/push_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/config/push
func ProcThingConfigPush(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	req := ConfigPushRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	c.debug("downstream thing <config>: push request")
	if err := c.SendResponse(uriServiceReplyWithRequestURI(rawURI),
		req.ID, infra.CodeSuccess, "{}"); err != nil {
		return err
	}
	return c.ipcSendMessage(&ipcMessage{
		evt:        ipcEvtConfigPush,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    req.Params,
	})
}

// ProcThingServicePropertySet 处理属性设置
// 下行
// request: /sys/{productKey}/{deviceName}/thing/service/property/set
// response: /sys/{productKey}/{deviceName}/thing/service/property/set_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/service/[+,#]
func ProcThingServicePropertySet(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}
	c.debug("downstream thing <service>: property set request")
	return c.ipcSendMessage(&ipcMessage{
		evt:        ipcEvtServicePropertySet,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    payload,
	})
}

// ProcThingServiceRequest 处理服务调用
// 下行
// request: /sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}
// response: /sys/{productKey}/{deviceName}/thing/service/{tsl.service.identifier}_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/service/[+,#]
func ProcThingServiceRequest(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	serviceID := uris[c.cfg.uriOffset+5]
	c.debug("downstream thing <service>: %s set request", serviceID)
	return c.ipcSendMessage(&ipcMessage{
		evt:        ipcEvtServiceRequest,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    payload,
		extend:     serviceID,
	})
}

// ProcRRPCRequest 处理RRPC请求
// 下行
// request: /sys/${YourProductKey}/${YourDeviceName}/rrpc/request/${messageId}
// response: /sys/${YourProductKey}/${YourDeviceName}/rrpc/response/${messageId}
// subscribe: /sys/${YourProductKey}/${YourDeviceName}/rrpc/request/+
func ProcRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	messageID := uris[c.cfg.uriOffset+5]
	c.debug("downstream sys <RRPC>: request - messageID: %s", messageID)

	return c.ipcSendMessage(&ipcMessage{
		evt:        ipcEvtRRPCRequest,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    payload,
		extend:     messageID,
	})
}

// ProcExtRRPCRequest 处理扩展RRPC请求
// 下行
// ${topic} 不为空,设备建立要求clientID传ext = 1
// request: /ext/rrpc/${messageId}/${topic}
// response: /ext/rrpc/${messageId}/${topic}
// subscribe: /ext/rrpc/+/${topic}
func ProcExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := strings.SplitN(strings.TrimLeft(rawURI, SEP), SEP, c.cfg.uriOffset+3)
	if len(uris) < (c.cfg.uriOffset + 3) {
		return ErrInvalidParameter
	}

	c.debug("downstream extend <RRPC>: Request - URI: ", rawURI)
	return c.ipcSendMessage(&ipcMessage{
		evt:     ipcEvtExtRRPCRequest,
		extend:  uris[c.cfg.uriOffset+2], // ${messageId}/${topic}
		payload: payload,
	})
}
