package model

import (
	"encoding/json"

	"github.com/thinkgos/aliIOT/infra"
)

type ProcDownStreamFunc func(m *Manager, rawURI string, payload []byte) error

// ProcThingModelUpRawReply 处理透传上行的应答
func ProcThingModelUpRawReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	m.debug("downstream thing <model>: up raw reply")
	return m.devUserProc.DownstreamThingModelUpRawReply(m, uris[m.cfg.uriOffset+1], uris[m.cfg.uriOffset+2], payload)
}

// ProcThingEventPost 处理ThingEvent XXX的应答
func ProcThingEventPostReply(m *Manager, rawURI string, payload []byte) error {
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

func ProcThingDeviceInfoUpdateReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <deviceInfo>: update reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDeviceInfoUpdateReply(m, &rsp)
}

func ProcThingDeviceInfoDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <deviceInfo>: delete reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDeviceInfoDeleteReply(m, &rsp)
}

func ProcThingDesiredPropertyGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <desired>: property get reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDesiredPropertyGetReply(m, &rsp)
}

func ProcThingDesiredPropertyDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <desired>: property delete reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDesiredPropertyDeleteReply(m, &rsp)
}

func ProcThingDsltemplateGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <dsl template>: get reply,@%d - %s", rsp.ID, string(rsp.Data))
	return m.devUserProc.DownstreamThingDsltemplateGetReply(m, &rsp)
}

// TODO: 需确认
func ProcThingDynamictslGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.debug("downstream thing <dynamic tsl>: get reply,@%d", rsp.ID)
	return m.devUserProc.DownstreamThingDynamictslGetReply(m, &rsp)
}

func ProcExtNtpResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := NtpResponsePayload{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.debug("downstream ext <ntp>: response - %+v", rsp)
	return m.devUserProc.DownstreamExtNtpResponse(m, &rsp)
}

func ProcThingConfigGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := ConfigGetResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	m.debug("downstream thing <config>: get reply,@%d,payload@%+v", rsp.ID, rsp)
	return m.devUserProc.DownstreamThingConfigGetReply(m, &rsp)
}

// TODO: deprecated
func ProcExtErrorResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.debug("downstream ext <Error>: response,@%d", rsp.ID)
	return m.devUserProc.DownstreamExtErrorResponse(m, &rsp)
}

// ProcThingModelDownRaw 处理透传下行
func ProcThingModelDownRaw(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	m.debug("downstream thing <model>: down raw request")
	return m.devUserProc.DownstreamThingModelDownRaw(m, uris[m.cfg.uriOffset+1], uris[m.cfg.uriOffset+2], payload)
}

func ProcThingConfigPush(m *Manager, rawURI string, payload []byte) error {
	req := ConfigPushRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	m.debug("downstream thing <config>: push request")
	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, infra.CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.devUserProc.DownstreamThingConfigPush(m, &req)
}

// ProcThingServiceRequest 服务调用
// 主要处理 thing/service/property/set, thing/service/{tsl.event.identifier}
func ProcThingServiceRequest(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}

	serviceID := uris[m.cfg.uriOffset+5]
	m.debug("downstream thing <service>: %s set requst", serviceID)
	if serviceID == "property" &&
		len(uris) >= (m.cfg.uriOffset+7) &&
		uris[m.cfg.uriOffset+6] == "set" {
		return m.devUserProc.DownstreamThingServicePropertySet(m, rawURI, payload)
	}

	return m.devUserProc.DownstreamThingServiceRequest(m, uris[m.cfg.uriOffset+1], uris[m.cfg.uriOffset+2], serviceID, payload)
}

func ProcRRPCRequest(m *Manager, rawURI string, payload []byte) error {
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

func ProcExtRRPCRequest(m *Manager, rawURI string, payload []byte) error {
	m.debug("downstream ext <RRPC>: Request - URI: ", rawURI)
	return m.devUserProc.DownStreamExtRRPCRequest(m, rawURI, payload)
}

/******************************** gateway ****************************************************************/

func ProcThingSubDevRegisterReply(m *Manager, rawURI string, payload []byte) error {
	rsp := GwSubDevRegisterResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwExtSubDevRegisterReply(m, &rsp)
}

func ProcExtSubDevCombineLoginReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwExtSubDevCombineLoginReply(m, &rsp)
}

func ProcExtSubDevCombineLogoutReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwExtSubDevCombineLogoutReply(m, &rsp)
}

func ProcThingTopoAddReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwThingTopoAddReply(m, &rsp)
}

func ProcThingTopoDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwThingTopoDeleteReply(m, &rsp)
}

func ProcThingTopoGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := GwTopoGetResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwThingTopoGetReply(m, &rsp)
}

func ProcThingDisable(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwSubDevThingDisable(m, uris[1], uris[2])
}
func ProcThingEnable(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwSubDevThingDisable(m, "", "")
}
func ProcThingDelete(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamGwSubDevThingDisable(m, "", "")
}
