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
	return c.devUserProc.DownstreamThingModelUpRawReply(c, uris[c.cfg.uriOffset+1], uris[c.cfg.uriOffset+2], payload)
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
	if eventID == "property" {
		return c.devUserProc.DownstreamThingEventPropertyPostReply(c, &rsp)
	}
	return c.devUserProc.DownstreamThingEventPostReply(c, eventID, &rsp)
}

// ProcThingDeviceInfoUpdateReply 处理设备信息更新应答
func ProcThingDeviceInfoUpdateReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <deviceInfo>: update reply,@%d", rsp.ID)
	return c.devUserProc.DownstreamThingDeviceInfoUpdateReply(c, &rsp)
}

// ProcThingDeviceInfoDeleteReply 处理设备信息删除的应答
func ProcThingDeviceInfoDeleteReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <deviceInfo>: delete reply,@%d", rsp.ID)
	return c.devUserProc.DownstreamThingDeviceInfoDeleteReply(c, &rsp)
}

// ProcThingDesiredPropertyGetReply 处理期望属性获取的应答
func ProcThingDesiredPropertyGetReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <desired>: property get reply,@%d", rsp.ID)
	return c.devUserProc.DownstreamThingDesiredPropertyGetReply(c, &rsp)
}

// ProcThingDesiredPropertyDeleteReply 处理期望属性删除的应答
func ProcThingDesiredPropertyDeleteReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <desired>: property delete reply,@%d", rsp.ID)
	return c.devUserProc.DownstreamThingDesiredPropertyDeleteReply(c, &rsp)
}

// ProcThingDsltemplateGetReply 处理dsltemplate获取的应答
func ProcThingDsltemplateGetReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <dsl template>: get reply,@%d - %s", rsp.ID, string(rsp.Data))
	return c.devUserProc.DownstreamThingDsltemplateGetReply(c, &rsp)
}

// TODO: 不使用??
func ProcThingDynamictslGetReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.debug("downstream thing <dynamic tsl>: get reply,@%d - %+v", rsp.ID, rsp)
	return c.devUserProc.DownstreamThingDynamictslGetReply(c, &rsp)
}

// ProcExtNtpResponse 处理ntp请求的应答
func ProcExtNtpResponse(c *Client, rawURI string, payload []byte) error {
	rsp := NtpResponsePayload{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.debug("downstream ext <ntp>: response - %+v", rsp)
	return c.devUserProc.DownstreamExtNtpResponse(c, &rsp)
}

// ProcThingConfigGetReply 处理获取配置的应答
func ProcThingConfigGetReply(c *Client, rawURI string, payload []byte) error {
	rsp := ConfigGetResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("downstream thing <config>: get reply,@%d,payload@%+v", rsp.ID, rsp)
	return c.devUserProc.DownstreamThingConfigGetReply(c, &rsp)
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
	return c.devUserProc.DownstreamExtErrorResponse(c, &rsp)
}

// ProcThingModelDownRaw 处理透传下行数据
func ProcThingModelDownRaw(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.debug("downstream thing <model>: down raw request")
	return c.devUserProc.DownstreamThingModelDownRaw(c, uris[c.cfg.uriOffset+1], uris[c.cfg.uriOffset+2], payload)
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
	return c.devUserProc.DownstreamThingConfigPush(c, &req)
}

// ProcThingServicePropertySet 处理属性设置
// 处理 thing/service/property/set
func ProcThingServicePropertySet(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}
	c.debug("downstream thing <service>: property set requst")
	return c.devUserProc.DownstreamThingServicePropertySet(c, rawURI, payload)
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
	return c.devUserProc.DownstreamThingServiceRequest(c, uris[c.cfg.uriOffset+1], uris[c.cfg.uriOffset+2], serviceID, payload)
}

// ProcRRPCRequest 处理RRPC请求
func ProcRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	messageID := uris[c.cfg.uriOffset+5]
	c.debug("downstream sys <RRPC>: request - messageID: %s", messageID)
	return c.devUserProc.DownStreamRRPCRequest(c,
		uris[c.cfg.uriOffset+1], uris[c.cfg.uriOffset+2], messageID,
		payload)
}

// ProcExtRRPCRequest 处理扩展RRPC请求
func ProcExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	c.debug("downstream ext <RRPC>: Request - URI: ", rawURI)
	return c.devUserProc.DownStreamExtRRPCRequest(c, rawURI, payload)
}
