package model

import (
	"encoding/json"
	"strings"
)

func ProcThingModelDownRaw(m *Manager, uri string, payload []byte) error {
	uris := strings.Split(uri, SEP)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}
	return ThingModelDownRaw(uris[m.uriOffset+1], uris[m.uriOffset+2], payload)
}

func ProcThingModelUpRawReply(m *Manager, uri string, payload []byte) error {
	uris := strings.Split(uri, SEP)
	if len(uris) != (m.uriOffset + 6) {
		return ErrInvalidURI
	}
	return ThingModelUpRawReply(uris[m.uriOffset+1], uris[m.uriOffset+2], payload)
}

// ProcThingEventPost 处理ThingEventXXX的应答
func ProcThingEventPostReply(m *Manager, uri string, payload []byte) error {
	uris := strings.Split(uri, SEP)
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

func ProcThingDeviceInfoUpdateReply(m *Manager, uri string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ThingDeviceInfoUpdateReply(&rsp)
}

func ProcThingDeviceInfoDeleteReply(m *Manager, uri string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ThingDeviceInfoDeleteReply(&rsp)
}
func ProcThingDsltemplateGetReply(m *Manager, uri string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return ThingDsltemplateGetReply(&rsp)
}
func ProcThingDynamictslGetReply(m *Manager, uri string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	return ThingDynamictslGetReply(&rsp)
}

func ProcExtNtpResponse(m *Manager, uri string, payload []byte) error {
	return ExtNtpResponse(payload)
}

func ProcExtErrorResponse(m *Manager, uri string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	return ExtErrorResponse(&rsp)
}
