package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliIOT/infra"
)

type DevNopUserProc struct{}

func (DevNopUserProc) DownstreamThingModelUpRawReply(m *Client, productKey, deviceName string, payload []byte) error {
	return nil
}

func (DevNopUserProc) DownstreamThingEventPropertyPostReply(m *Client, rsp *Response) error {
	return nil
}

func (DevNopUserProc) DownstreamThingEventPostReply(m *Client, eventID string, rsp *Response) error {
	return nil
}

func (DevNopUserProc) DownstreamThingDeviceInfoUpdateReply(m *Client, rsp *Response) error {
	return nil
}
func (DevNopUserProc) DownstreamThingDeviceInfoDeleteReply(m *Client, rsp *Response) error {
	return nil
}

func (DevNopUserProc) DownstreamThingDesiredPropertyGetReply(m *Client, rsp *Response) error {
	return nil
}

func (DevNopUserProc) DownstreamThingDesiredPropertyDeleteReply(m *Client, rsp *Response) error {
	return nil
}
func (DevNopUserProc) DownstreamThingDsltemplateGetReply(m *Client, rsp *Response) error {
	return nil
}

func (DevNopUserProc) DownstreamThingDynamictslGetReply(m *Client, rsp *Response) error {
	return nil
}
func (DevNopUserProc) DownstreamExtNtpResponse(m *Client, rsp *NtpResponsePayload) error {
	return nil
}

func (DevNopUserProc) DownstreamThingConfigGetReply(m *Client, rsp *ConfigGetResponse) error {
	return nil
}

// TODO: deprecated
func (DevNopUserProc) DownstreamExtErrorResponse(m *Client, rsp *Response) error {
	return nil
}

func (DevNopUserProc) DownstreamThingModelDownRaw(m *Client, productKey, deviceName string, payload []byte) error {
	return nil
}

func (DevNopUserProc) DownstreamThingConfigPush(m *Client, rsp *ConfigPushRequest) error {
	return nil
}

// DownstreamThingServicePropertySet 设置设备属性
func (DevNopUserProc) DownstreamThingServicePropertySet(m *Client, topic string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}
	return m.SendResponse(URIServiceReplyWithRequestURI(topic), rsp.ID, infra.CodeSuccess, "{}")
}

// DownstreamThingServiceRequest 设备服务调用请求
func (DevNopUserProc) DownstreamThingServiceRequest(m *Client, productKey, deviceName, srvID string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}

	return m.SendResponse(m.URIService(URISysPrefix, URIThingServiceRequest, productKey, deviceName, srvID),
		rsp.ID, infra.CodeSuccess, "{}")
}

func (DevNopUserProc) DownStreamRRPCRequest(m *Client, productKey, deviceName, messageID string, payload []byte) error {
	return m.Publish(m.URIService(URISysPrefix, URIRRPCResponse, productKey, deviceName, messageID),
		0, `{"note":"default system RRPC implementation"}`)
}

func (DevNopUserProc) DownStreamExtRRPCRequest(m *Client, rawURI string, payload []byte) error {
	return m.Publish(rawURI, 0, "default ext RRPC implementation")
}
