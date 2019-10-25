package model

import (
	"encoding/json"
)

type ProcDownStreamFunc func(m *Manager, rawURI string, payload []byte) error

// ProcThingModelDownRaw 处理透传
func ProcThingModelDownRaw(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}
	return DownstreamThingModelDownRaw(uris[m.uriOffset+1], uris[m.uriOffset+2], payload)
}

func ProcThingServicePropertySet(m *Manager, rawURI string, payload []byte) error {
	return DownstreamThingServicePropertySet(payload)
}

// deprecated
func ProcThingServicePropertyGet(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}
	return DownstreamThingServicePropertyGet(uris[m.uriOffset+1], uris[m.uriOffset+2], payload)
}

func ProcThingServiceRequest(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}

	return DownstreamThingServiceRequest(uris[m.uriOffset+1], uris[m.uriOffset+2], uris[m.uriOffset+5], payload)
}

// ProcThingModelUpRawReply 处理透传上行的应答
func ProcThingModelUpRawReply(m *Manager, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}

	return DownstreamThingModelUpRawReply(uris[m.uriOffset+1], uris[m.uriOffset+2], payload)
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
	rsp := SubDevRegisterResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamExtSubDevRegisterReply(m, &rsp)
}

func ProcExtSubDevCombineLoginReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamExtSubDevCombineLoginReply(m, &rsp)
}

func ProcExtSubDevCombineLogoutReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamExtSubDevCombineLogoutReply(m, &rsp)
}

func ProcThingTopoAddReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamThingTopoAddReply(m, &rsp)
}

func ProcThingTopoDeleteReply(m *Manager, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return m.gwUserProc.DownstreamThingTopoDeleteReply(m, &rsp)
}
