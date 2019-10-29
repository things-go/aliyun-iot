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

/******************************** gateway ****************************************************************/

func ProcThingTopoAddReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("upstream GW thing <topo>: add reply @%d", rsp.ID)
	return c.gwUserProc.DownstreamGwThingTopoAddReply(c, &rsp)
}

func ProcThingTopoDeleteReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("upstream GW thing <topo>: delete reply @%d", rsp.ID)
	return c.gwUserProc.DownstreamGwThingTopoDeleteReply(c, &rsp)
}

func ProcThingTopoGetReply(c *Client, rawURI string, payload []byte) error {
	rsp := GwTopoGetResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("upstream GW thing <topo>: get reply @%d", rsp.ID)
	return c.gwUserProc.DownstreamGwThingTopoGetReply(c, &rsp)
}

func ProcThingListFoundReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	c.CacheRemove(rsp.ID)
	c.debug("upstream GW thing <list>: found reply @%d", rsp.ID)
	return nil
}

type GwTopoAddNotifyParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}
type GwTopoAddNotifyRequest struct {
	Request
	Params []GwTopoAddNotifyParams `json:"params"`
}

// 通知网关添加设备拓扑关系
func ProcThingTopoAddNotify(c *Client, rawURI string, payload []byte) error {
	req := GwTopoAddNotifyRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	// TODO: 处理通知的
	c.debug("upstream GW thing <topo>: notify @%d")
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}")
}

type GwTopoChangeDev struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type GwTopoChangeParams struct {
	Status  int               `json:"status"` // 0: 创建 1:删除 2: 解除禁用(启用) 8: 禁用
	SubList []GwTopoChangeDev `json:"subList"`
}

type GwTopoChangeRequest struct {
	Request
	Params GwTopoChangeParams `json:"params"`
}

// 通知网关拓扑关系变化
func ProcThingTopoChange(c *Client, rawURI string, payload []byte) error {
	req := GwTopoChangeRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	// TODO: 处理拓扑关系变更
	c.debug("upstream GW thing <topo>: change @%d")
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}")
}

func ProcThingSubDevRegisterReply(c *Client, rawURI string, payload []byte) error {
	rsp := GwSubDevRegisterResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("upstream GW thing <sub>: register reply @%d", rsp.ID)
	// TODO
	return c.gwUserProc.DownstreamGwExtSubDevRegisterReply(c, &rsp)
}

func ProcExtSubDevCombineLoginReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("upstream Ext GW <sub>: login reply @%d", rsp.ID)
	return c.gwUserProc.DownstreamGwExtSubDevCombineLoginReply(c, &rsp)
}

func ProcExtSubDevCombineLogoutReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.debug("upstream Ext GW <sub>: logout reply @%d", rsp.ID)
	return c.gwUserProc.DownstreamGwExtSubDevCombineLogoutReply(c, &rsp)
}

// ProcThingDisable 禁用子设备
func ProcThingDisable(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}"); err != nil {
		return err
	}
	c.debug("downstream <thing>: disable @%s - %s", uris[1], uris[2])
	return c.gwUserProc.DownstreamGwSubDevThingDisable(c, uris[1], uris[2])
}

// ProcThingEnable 启用子设备
func ProcThingEnable(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}"); err != nil {
		return err
	}
	c.debug("downstream <thing>: enable @%s - %s", uris[1], uris[2])
	return c.gwUserProc.DownstreamGwSubDevThingDisable(c, uris[1], uris[2])
}
func ProcThingDelete(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}"); err != nil {
		return err
	}
	c.debug("downstream <thing>: delete @%s - %s", uris[1], uris[2])
	return c.gwUserProc.DownstreamGwSubDevThingDisable(c, uris[1], uris[2])
}
