package model

import (
	"encoding/json"

	"github.com/thinkgos/aliIOT/infra"
)

type ProcDownStreamFunc func(m *Manager, rawURI string, payload []byte) error

// ProcThingModelDownRaw 处理透传
func ProcThingModelDownRaw(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}
	return ThingModelDownRaw(uris[m.uriOffset+1], uris[m.uriOffset+2], payload)
}

func ProcThingServicePropertySet(m *Manager, rawURI string, payload []byte) error {
	code := infra.CodeSuccess
	if err := ThingServicePropertySet(payload); err != nil {
		code = infra.CodeRequestError
	}
	return m.SendResponse(URIServiceReplyWithRequestURI(rawURI), m.RequestID(), code, "{}")
}

// deprecated
func ProcThingServicePropertyGet(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}
	return ThingServicePropertyGet(uris[m.uriOffset+1], uris[m.uriOffset+2], payload)
}

func ProcThingServiceRequest(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}

	return ThingServiceRequest(uris[m.uriOffset+1], uris[m.uriOffset+2], uris[m.uriOffset+5], payload)
}

// ProcThingModelUpRawReply 处理透传上行的应答
func ProcThingModelUpRawReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}

	return ThingModelUpRawReply(uris[m.uriOffset+1], uris[m.uriOffset+2], payload)
}

// ProcThingEventPost 处理ThingEventXXX的应答
func ProcThingEventPostReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 7) {
		return ErrInvalidURI
	}
	EventID := uris[m.uriOffset+5]

	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	if EventID == "property" {
		_ = ThingEventPropertyPostReply(&rsp)
	} else {
		_ = ThingEventPostReply(EventID, &rsp)
	}

	return nil
}

func ProcThingDeviceInfoUpdateReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ThingDeviceInfoUpdateReply(&rsp)
}

func ProcThingDeviceInfoDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ThingDeviceInfoDeleteReply(&rsp)
}
func ProcThingDsltemplateGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return ThingDsltemplateGetReply(&rsp)
}
func ProcThingDynamictslGetReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return ThingDynamictslGetReply(&rsp)
}

func ProcExtNtpResponse(m *Manager, rawURI string, payload []byte) error {
	return ExtNtpResponse(payload)
}

func ProcExtErrorResponse(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ExtErrorResponse(&rsp)
}

func ProcExtSubDevCombineLoginReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ExtExtSubDevCombineLoginReply(&rsp)
}

func ProcExtSubDevCombineLogoutReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ExtExtSubDevCombineLogoutReply(&rsp)
}

func ProcThingSubDevRegisterReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ExtExtSubDevRegisterReply(&rsp)
}
