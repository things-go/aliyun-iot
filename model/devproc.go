package model

import (
	"encoding/json"

	"github.com/thinkgos/aliIOT/infra"
)

type DevNopUserProc struct{}

func (DevNopUserProc) DownstreamThingModelUpRawReply(m *Manager, productKey, deviceName string, payload []byte) error {
	m.debug("downstream thing <model>: up raw reply")
	return nil
}

func (DevNopUserProc) DownstreamThingEventPropertyPostReply(m *Manager, rsp *Response) error {
	m.debug("downstream thing <event>: property post reply")
	return nil
}

func (DevNopUserProc) DownstreamThingEventPostReply(m *Manager, eventID string, rsp *Response) error {
	m.debug("downstream thing <event>: %s post reply", eventID)
	return nil
}

func (DevNopUserProc) DownstreamThingDeviceInfoUpdateReply(m *Manager, rsp *Response) error {
	m.debug("downstream thing <deviceInfo>: update reply")
	return nil
}
func (DevNopUserProc) DownstreamThingDeviceInfoDeleteReply(m *Manager, rsp *Response) error {
	m.debug("downstream thing <deviceInfo>: delete reply")
	return nil
}

func (DevNopUserProc) DownstreamThingDesiredPropertyGetReply(m *Manager, rsp *Response) error {
	m.debug("downstream thing <desired>: property get reply")
	return nil
}

func (DevNopUserProc) DownstreamThingDesiredPropertyDeleteReply(m *Manager, rsp *Response) error {
	m.debug("downstream thing <desired>: property delete reply")
	return nil
}
func (DevNopUserProc) DownstreamThingDsltemplateGetReply(m *Manager, rsp *Response) error {
	m.debug("downstream thing <dsl template>: get reply")
	return nil
}

func (DevNopUserProc) DownstreamThingDynamictslGetReply(m *Manager, rsp *Response) error {
	m.debug("downstream thing <dynamic tsl>: get reply")
	return nil
}
func (DevNopUserProc) DownstreamExtNtpResponse(m *Manager, rsp *NtpResponsePayload) error {
	m.debug("downstream ext <ntp>: response")
	return nil
}

func (DevNopUserProc) DownstreamThingConfigGetReply(m *Manager, rsp *ConfigGetResponse) error {
	m.debug("downstream thing <config>: get reply")
	return nil
}

// TODO: deprecated
func (DevNopUserProc) DownstreamExtErrorResponse(m *Manager, rsp *Response) error {
	m.debug("downstream ext <Error>: response")
	return nil
}

func (DevNopUserProc) DownstreamThingModelDownRaw(m *Manager, productKey, deviceName string, payload []byte) error {
	m.debug("downstream thing <model>: down raw request")
	return nil
}

func (DevNopUserProc) DownstreamThingConfigPush(m *Manager, rsp *ConfigPushRequest) error {
	m.debug("downstream thing <config>: push request")
	return nil
}

// DownstreamThingServicePropertySet 设置设备属性
func (DevNopUserProc) DownstreamThingServicePropertySet(m *Manager, topic string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}
	m.debug("downstream thing <service>: property set requst")
	return m.SendResponse(URIServiceReplyWithRequestURI(topic), rsp.ID, infra.CodeSuccess, "{}")
}

// DownstreamThingServiceRequest 设备服务调用请求
func (DevNopUserProc) DownstreamThingServiceRequest(m *Manager, productKey, deviceName, srvID string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}

	m.debug("downstream thing <service>: %s request", srvID)
	return m.SendResponse(m.URIService(URISysPrefix, URIThingServiceRequest, productKey, deviceName, srvID),
		rsp.ID, infra.CodeSuccess, "{}")
}

func (DevNopUserProc) DownStreamRRPCRequest(m *Manager, productKey, deviceName, messageID string, payload []byte) error {
	m.debug("downstream sys <RRPC>: request - messageID: ", messageID)
	return m.Publish(m.URIService(URISysPrefix, URIRRPCResponse, productKey, deviceName, messageID),
		0, "default system RRPC implementation")
}

func (DevNopUserProc) DownStreamExtRRPCRequest(m *Manager, rawURI string, payload []byte) error {
	m.debug("downstream ext <RRPC>: Request - URI: ", rawURI)
	return m.Publish(rawURI, 0, "default ext RRPC implementation")
}
