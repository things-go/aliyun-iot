package model

import (
	"encoding/json"

	"github.com/thinkgos/aliIOT/infra"
)

type ProcDownStreamFunc func(m *Manager, rawURI string, payload []byte) error

// ProcThingModelDownRaw 处理透传
func ProcThingModelDownRaw(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}
	return m.devUserProc.DownstreamThingModelDownRaw(uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], payload)
}

// ProcThingModelUpRawReply 处理透传上行的应答
func ProcThingModelUpRawReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}

	return m.devUserProc.DownstreamThingModelUpRawReply(uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], payload)
}

// ProcThingEventPost 处理ThingEventXXX的应答
func ProcThingEventPostReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.opt.uriOffset + 7) {
		return ErrInvalidURI
	}
	EventID := uris[m.opt.uriOffset+5]

	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	m.CacheRemove(rsp.ID)
	if EventID == "property" {
		_ = m.devUserProc.DownstreamThingEventPropertyPostReply(&rsp)
	} else {
		_ = m.devUserProc.DownstreamThingEventPostReply(EventID, &rsp)
	}

	return nil
}

func ProcThingDeviceInfoUpdateReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.devUserProc.DownstreamThingDeviceInfoUpdateReply(&rsp)
}

func ProcThingDeviceInfoDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.devUserProc.DownstreamThingDeviceInfoDeleteReply(&rsp)
}
func ProcThingDsltemplateGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return m.devUserProc.DownstreamThingDsltemplateGetReply(&rsp)
}
func ProcThingDynamictslGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return m.devUserProc.DownstreamThingDynamictslGetReply(&rsp)
}

func ProcThingServicePropertySet(m *Manager, rawURI string, payload []byte) error {
	return m.devUserProc.DownstreamThingServicePropertySet(payload)
}

// deprecated
func ProcThingServicePropertyGet(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}
	return m.devUserProc.DownstreamThingServicePropertyGet(uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], payload)
}

func ProcThingServiceRequest(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}

	return m.devUserProc.DownstreamThingServiceRequest(uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], uris[m.opt.uriOffset+5], payload)
}

func ProcExtNtpResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := NtpResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.devUserProc.DownstreamExtNtpResponse(&rsp)
}

func ProcExtErrorResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.devUserProc.DownstreamExtErrorResponse(&rsp)
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
