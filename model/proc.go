package model

import (
	"encoding/json"

	"github.com/thinkgos/aliIOT/infra"
)

type ProcDownStreamFunc func(m *Manager, rawURI string, payload []byte) error

// ProcThingModelUpRawReply 处理透传上行的应答
func ProcThingModelUpRawReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}

	return m.devUserProc.DownstreamThingModelUpRawReply(m, uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], payload)
}

// ProcThingModelDownRaw 处理透传下行
func ProcThingModelDownRaw(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}
	return m.devUserProc.DownstreamThingModelDownRaw(m, uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], payload)
}

// ProcThingEventPost 处理ThingEvent XXX的应答
func ProcThingEventPostReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.opt.uriOffset + 7) {
		return ErrInvalidURI
	}
	EventID := uris[m.opt.uriOffset+5]

	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	m.CacheRemove(rsp.ID)
	if EventID == "property" {
		return m.devUserProc.DownstreamThingEventPropertyPostReply(&rsp)
	}
	return m.devUserProc.DownstreamThingEventPostReply(EventID, &rsp)
}

func ProcThingDeviceInfoUpdateReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	return m.devUserProc.DownstreamThingDeviceInfoUpdateReply(&rsp)
}

func ProcThingDeviceInfoDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	return m.devUserProc.DownstreamThingDeviceInfoDeleteReply(&rsp)
}

func ProcThingDesiredPropertyGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	return m.devUserProc.DownstreamThingDesiredPropertyGetReply(&rsp)
}

func ProcThingDesiredPropertyDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	return m.devUserProc.DownstreamThingDesiredPropertyDeleteReply(&rsp)
}

func ProcThingDsltemplateGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	return m.devUserProc.DownstreamThingDsltemplateGetReply(&rsp)
}

// TODO: 需确认
func ProcThingDynamictslGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return m.devUserProc.DownstreamThingDynamictslGetReply(&rsp)
}

func ProcExtNtpResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := NtpResponsePayload{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.devUserProc.DownstreamExtNtpResponse(&rsp)
}

func ProcThingConfigGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := ConfigGetResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	m.CacheRemove(rsp.ID)
	return m.devUserProc.DownstreamThingConfigGetReply(&rsp)
}

func ProcThingConfigPush(m *Manager, rawURI string, payload []byte) error {
	req := ConfigPushRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if err := m.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, infra.CodeSuccess, "{}"); err != nil {
		return err
	}
	return m.devUserProc.DownstreamThingConfigPush(&req)
}

// deprecated
func ProcExtErrorResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.devUserProc.DownstreamExtErrorResponse(&rsp)
}

// ProcThingServiceRequest 服务调用
// 主要处理 thing/service/property/set, thing/service/{tsl.event.identifier}
func ProcThingServiceRequest(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}

	serviceID := uris[m.opt.uriOffset+5]
	if serviceID == "property" &&
		len(uris) >= (m.opt.uriOffset+7) &&
		uris[m.opt.uriOffset+6] == "set" {
		return m.devUserProc.DownstreamThingServicePropertySet(m, rawURI, payload)
	}

	return m.devUserProc.DownstreamThingServiceRequest(m, uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], serviceID, payload)
}

func ProcRRPCRequest(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}

	return m.devUserProc.DownStreamRRPCRequest(m,
		uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], uris[m.opt.uriOffset+5],
		payload)
}

func ProcExtRRPCRequest(m *Manager, rawURI string, payload []byte) error {
	return m.devUserProc.DownStreamExtRRPCRequest(m, rawURI, payload)
}

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
