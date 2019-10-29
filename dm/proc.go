package dm

import (
	"encoding/json"
)

// ProcDownStreamFunc 处理下行数据
type ProcDownStreamFunc func(m *Client, rawURI string, payload []byte) error

// ProcThingModelUpRawReply 处理透传上行的应答
func ProcThingModelUpRawReply(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	m.debug("downstream thing <model>: up raw reply")
	return m.devUserProc.DownstreamThingModelUpRawReply(m, uris[m.cfg.uriOffset+1], uris[m.cfg.uriOffset+2], payload)
}

// ProcThingEventPostReply 处理ThingEvent XXX上行的应答
func ProcThingEventPostReply(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}

	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	m.CacheRemove(rsp.ID)
	eventID := uris[m.cfg.uriOffset+5]
	m.debug("downstream thing <event>: %s post reply,@%d", eventID, rsp.ID)
	if eventID == "property" {
		return m.devUserProc.DownstreamThingEventPropertyPostReply(m, &rsp)
	}
	return m.devUserProc.DownstreamThingEventPostReply(m, eventID, &rsp)
}

// ProcThingDeviceInfoUpdateReply 处理设备信息更新应答
func ProcThingDeviceInfoUpdateReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <deviceInfo>: update reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDeviceInfoUpdateReply(m, &rsp)
}

// ProcThingDeviceInfoDeleteReply 处理设备信息删除的应答
func ProcThingDeviceInfoDeleteReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <deviceInfo>: delete reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDeviceInfoDeleteReply(m, &rsp)
}

// ProcThingDesiredPropertyGetReply 处理期望属性获取的应答
func ProcThingDesiredPropertyGetReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <desired>: property get reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDesiredPropertyGetReply(m, &rsp)
}

// ProcThingDesiredPropertyDeleteReply 处理期望属性删除的应答
func ProcThingDesiredPropertyDeleteReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <desired>: property delete reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDesiredPropertyDeleteReply(m, &rsp)
}

// ProcThingDsltemplateGetReply 处理dsltemplate获取的应答
func ProcThingDsltemplateGetReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <dsl template>: get reply,@%d - %s", rsp.ID, string(rsp.Data))
	return m.devUserProc.DownstreamThingDsltemplateGetReply(m, &rsp)
}

// TODO: 不使用??
func ProcThingDynamictslGetReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.debug("downstream thing <dynamic tsl>: get reply,@%d - %+v", rsp.ID, rsp)
	return m.devUserProc.DownstreamThingDynamictslGetReply(m, &rsp)
}

// ProcExtNtpResponse 处理ntp请求的应答
func ProcExtNtpResponse(m *Client, rawURI string, payload []byte) error {
	rsp := NtpResponsePayload{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.debug("downstream ext <ntp>: response - %+v", rsp)
	return m.devUserProc.DownstreamExtNtpResponse(m, &rsp)
}

// ProcThingConfigGetReply 处理获取配置的应答
func ProcThingConfigGetReply(m *Client, rawURI string, payload []byte) error {
	rsp := ConfigGetResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <config>: get reply,@%d,payload@%+v", rsp.ID, rsp)
	return m.devUserProc.DownstreamThingConfigGetReply(m, &rsp)
}

// ProcExtErrorResponse 处理错误的回复
func ProcExtErrorResponse(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	// TODO: 处理这个ERROR
	m.CacheRemove(rsp.ID)
	m.debug("downstream ext <Error>: response,@%d", rsp.ID)
	return m.devUserProc.DownstreamExtErrorResponse(m, &rsp)
}

// ProcThingModelDownRaw 处理透传下行数据
func ProcThingModelDownRaw(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	m.debug("downstream thing <model>: down raw request")
	return m.devUserProc.DownstreamThingModelDownRaw(m, uris[m.cfg.uriOffset+1], uris[m.cfg.uriOffset+2], payload)
}

// ProcThingConfigPush 处理配置推送
func ProcThingConfigPush(m *Client, rawURI string, payload []byte) error {
	req := ConfigPushRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	m.debug("downstream thing <config>: push request")
	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.devUserProc.DownstreamThingConfigPush(m, &req)
}

// ProcThingServicePropertySet 处理属性设置
// 处理 thing/service/property/set
func ProcThingServicePropertySet(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}
	m.debug("downstream thing <service>: property set requst")
	return m.devUserProc.DownstreamThingServicePropertySet(m, rawURI, payload)
}

// ProcThingServiceRequest 处理服务调用
// 处理 thing/service/{tsl.event.identifier}
func ProcThingServiceRequest(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	serviceID := uris[m.cfg.uriOffset+5]
	m.debug("downstream thing <service>: %s set requst", serviceID)
	return m.devUserProc.DownstreamThingServiceRequest(m, uris[m.cfg.uriOffset+1], uris[m.cfg.uriOffset+2], serviceID, payload)
}

// ProcRRPCRequest 处理RRPC请求
func ProcRRPCRequest(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	messageID := uris[m.cfg.uriOffset+5]
	m.debug("downstream sys <RRPC>: request - messageID: %s", messageID)
	return m.devUserProc.DownStreamRRPCRequest(m,
		uris[m.cfg.uriOffset+1], uris[m.cfg.uriOffset+2], messageID,
		payload)
}

// ProcExtRRPCRequest 处理扩展RRPC请求
func ProcExtRRPCRequest(m *Client, rawURI string, payload []byte) error {
	m.debug("downstream ext <RRPC>: Request - URI: ", rawURI)
	return m.devUserProc.DownStreamExtRRPCRequest(m, rawURI, payload)
}

/******************************** gateway ****************************************************************/

func ProcThingSubDevRegisterReply(m *Client, rawURI string, payload []byte) error {
	rsp := GwSubDevRegisterResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwExtSubDevRegisterReply(m, &rsp)
}

func ProcExtSubDevCombineLoginReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwExtSubDevCombineLoginReply(m, &rsp)
}

func ProcExtSubDevCombineLogoutReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwExtSubDevCombineLogoutReply(m, &rsp)
}

func ProcThingTopoAddReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwThingTopoAddReply(m, &rsp)
}

func ProcThingTopoDeleteReply(m *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwThingTopoDeleteReply(m, &rsp)
}

func ProcThingTopoGetReply(m *Client, rawURI string, payload []byte) error {
	rsp := GwTopoGetResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwThingTopoGetReply(m, &rsp)
}

func ProcThingDisable(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwSubDevThingDisable(m, uris[1], uris[2])
}
func ProcThingEnable(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwSubDevThingDisable(m, "", "")
}
func ProcThingDelete(m *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwSubDevThingDisable(m, "", "")
}
