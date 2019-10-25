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
	return DownstreamThingModelDownRaw(uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], payload)
}

func ProcThingServicePropertySet(m *Manager, rawURI string, payload []byte) error {
	return DownstreamThingServicePropertySet(payload)
}

// deprecated
func ProcThingServicePropertyGet(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}
	return DownstreamThingServicePropertyGet(uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], payload)
}

func ProcThingServiceRequest(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}

	return DownstreamThingServiceRequest(uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], uris[m.opt.uriOffset+5], payload)
}

// ProcThingModelUpRawReply 处理透传上行的应答
func ProcThingModelUpRawReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.opt.uriOffset + 6) {
		return ErrInvalidURI
	}

	return DownstreamThingModelUpRawReply(uris[m.opt.uriOffset+1], uris[m.opt.uriOffset+2], payload)
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

	if EventID == "property" {
		_ = DownstreamThingEventPropertyPostReply(&rsp)
	} else {
		_ = DownstreamThingEventPostReply(EventID, &rsp)
	}

	return nil
}

func ProcThingDeviceInfoUpdateReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return DownstreamThingDeviceInfoUpdateReply(&rsp)
}

func ProcThingDeviceInfoDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return DownstreamThingDeviceInfoDeleteReply(&rsp)
}
func ProcThingDsltemplateGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return DownstreamThingDsltemplateGetReply(&rsp)
}
func ProcThingDynamictslGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return DownstreamThingDynamictslGetReply(&rsp)
}

func ProcExtNtpResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := NtpResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return DownstreamExtNtpResponse(&rsp)
}

func ProcExtErrorResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return DownstreamExtErrorResponse(&rsp)
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
