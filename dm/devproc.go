package dm

import (
	"encoding/json"
)

// ProcDownStreamFunc 处理下行数据
type ProcDownStreamFunc func(c *Client, rawURI string, payload []byte) error

// ProcThingModelUpRawReply 处理透传上行的应答
func ProcThingModelUpRawReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.debug("downstream thing <model>: up raw reply")
	return c.ipcSendMessage(&ipcMessage{
		err:        nil,
		evt:        ipcEvtRawReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    payload,
	})
}

// ProcThingEventPostReply 处理ThingEvent XXX上行的应答
func ProcThingEventPostReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}

	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	c.CacheRemove(rsp.ID)
	eventID := uris[c.cfg.uriOffset+5]
	c.debug("downstream thing <event>: %s post reply,@%d", eventID, rsp.ID)

	var err error
	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}
	if eventID == "property" {
		return c.ipcSendMessage(&ipcMessage{
			err:        err,
			evt:        ipcEvtEventPropertyPostReply,
			productKey: uris[c.cfg.uriOffset+1],
			deviceName: uris[c.cfg.uriOffset+2],
			payload:    nil,
		})
	}

	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtEventPostReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    nil,
		ext:        eventID,
	})
}

// ProcThingDeviceInfoUpdateReply 处理设备信息更新应答
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
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <deviceInfo>: update reply,@%d", rsp.ID)

	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}

	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDeviceInfoUpdateReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
	})
}

// ProcThingDeviceInfoDeleteReply 处理设备信息删除的应答
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
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <deviceInfo>: delete reply,@%d", rsp.ID)
	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDeviceInfoUpdateReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
	})
}

// ProcThingDesiredPropertyGetReply 处理期望属性获取的应答
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
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <desired>: property get reply,@%d", rsp.ID)
	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDesiredPropertyGetReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    rsp.Data,
	})
}

// ProcThingDesiredPropertyDeleteReply 处理期望属性删除的应答
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
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <desired>: property delete reply,@%d", rsp.ID)
	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDesiredPropertyDeleteReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
	})
}

// ProcThingDsltemplateGetReply 处理dsltemplate获取的应答
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
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <dsl template>: get reply,@%d - %s", rsp.ID, string(rsp.Data))
	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDsltemplateGetReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    rsp.Data,
	})
}

// TODO: 不使用??
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
	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtDynamictslGetReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    rsp.Data,
	})
}

// ProcExtNtpResponse 处理ntp请求的应答
func ProcExtNtpResponse(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 5) {
		return ErrInvalidURI
	}
	rsp := NtpResponsePayload{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.debug("downstream ext <ntp>: response - %+v", rsp)
	return c.ipcSendMessage(&ipcMessage{
		err:        nil,
		evt:        ipcEvtExtNtpResponse,
		productKey: uris[c.cfg.uriOffset+2],
		deviceName: uris[c.cfg.uriOffset+3],
		payload:    rsp,
	})
}

// ProcThingConfigGetReply 处理获取配置的应答
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
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <config>: get reply,@%d,payload@%+v", rsp.ID, rsp)
	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}
	return c.ipcSendMessage(&ipcMessage{
		err:        err,
		evt:        ipcEvtConfigGetReply,
		productKey: uris[c.cfg.uriOffset+1],
		deviceName: uris[c.cfg.uriOffset+2],
		payload:    rsp.Data,
	})
}

// ProcExtErrorResponse 处理错误的回复
func ProcExtErrorResponse(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	// TODO: 处理这个ERROR
	c.CacheRemove(rsp.ID)
	c.debug("downstream ext <Error>: response,@%d", rsp.ID)
	return c.eventProc.DownstreamExtErrorResponse(c, &rsp)
}

// ProcThingModelDownRaw 处理透传下行数据
func ProcThingModelDownRaw(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.debug("downstream thing <model>: down raw request")
	return c.eventProc.DownstreamThingModelDownRaw(c, uris[c.cfg.uriOffset+1], uris[c.cfg.uriOffset+2], payload)
}

// ProcThingConfigPush 处理配置推送
func ProcThingConfigPush(c *Client, rawURI string, payload []byte) error {
	req := ConfigPushRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	c.debug("downstream thing <config>: push request")
	if err := c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}"); err != nil {
		return err
	}
	return c.eventProc.DownstreamThingConfigPush(c, &req)
}

// ProcThingServicePropertySet 处理属性设置
// 处理 thing/service/property/set
func ProcThingServicePropertySet(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}
	c.debug("downstream thing <service>: property set requst")
	return c.eventProc.DownstreamThingServicePropertySet(c, rawURI, payload)
}

// ProcThingServiceRequest 处理服务调用
// 处理 thing/service/{tsl.event.identifier}
func ProcThingServiceRequest(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	serviceID := uris[c.cfg.uriOffset+5]
	c.debug("downstream thing <service>: %s set requst", serviceID)
	return c.eventProc.DownstreamThingServiceRequest(c, uris[c.cfg.uriOffset+1], uris[c.cfg.uriOffset+2], serviceID, payload)
}

// ProcRRPCRequest 处理RRPC请求
func ProcRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	messageID := uris[c.cfg.uriOffset+5]
	c.debug("downstream sys <RRPC>: request - messageID: %s", messageID)
	return c.eventProc.DownStreamRRPCRequest(c,
		uris[c.cfg.uriOffset+1], uris[c.cfg.uriOffset+2], messageID,
		payload)
}

// ProcExtRRPCRequest 处理扩展RRPC请求
func ProcExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	c.debug("downstream ext <RRPC>: Request - URI: ", rawURI)
	return c.eventProc.DownStreamExtRRPCRequest(c, rawURI, payload)
}
