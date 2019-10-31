package dm

import (
	"encoding/json"
)

// NopEvt 实现DevUserProc接口的空实现
type NopEvt struct{}

// EvtThingModelUpRawReply see interface EventProc
func (NopEvt) EvtThingModelUpRawReply(c *Client, productKey, deviceName string, payload []byte) error {
	return nil
}

// EvtThingEventPropertyPostReply see interface EventProc
func (NopEvt) EvtThingEventPropertyPostReply(c *Client, err error, productKey, deviceName string) error {
	return nil
}

// EvtThingEventPostReply see interface EventProc
func (NopEvt) EvtThingEventPostReply(c *Client, err error, eventID, productKey, deviceName string) error {
	return nil
}

// EvtThingDeviceInfoUpdateReply see interface EventProc
func (NopEvt) EvtThingDeviceInfoUpdateReply(c *Client, err error, productKey, deviceName string) error {
	return nil
}

// EvtThingDeviceInfoDeleteReply see interface EventProc
func (NopEvt) EvtThingDeviceInfoDeleteReply(c *Client, err error, productKey, deviceName string) error {
	return nil
}

// EvtThingDesiredPropertyGetReply see interface EventProc
func (NopEvt) EvtThingDesiredPropertyGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error {
	return nil
}

// EvtThingDesiredPropertyDeleteReply see interface EventProc
func (NopEvt) EvtThingDesiredPropertyDeleteReply(c *Client, err error, productKey, deviceName string) error {
	return nil
}

// EvtThingDsltemplateGetReply see interface EventProc
func (NopEvt) EvtThingDsltemplateGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error {
	return nil
}

// EvtThingDynamictslGetReply see interface EventProc
func (NopEvt) EvtThingDynamictslGetReply(c *Client, err error, productKey, deviceName string, data json.RawMessage) error {
	return nil
}

// EvtExtNtpResponse see interface EventProc
func (NopEvt) EvtExtNtpResponse(c *Client, productKey, deviceName string, rsp NtpResponsePayload) error {
	return nil
}

// EvtThingConfigGetReply see interface EventProc
func (NopEvt) EvtThingConfigGetReply(c *Client, err error, productKey, deviceName string, data ConfigParamsAndData) error {
	return nil
}

// DownstreamExtErrorResponse see interface EventProc
func (NopEvt) DownstreamExtErrorResponse(c *Client, rsp *Response) error {
	return nil
}

// DownstreamThingModelDownRaw see interface EventProc
func (NopEvt) DownstreamThingModelDownRaw(c *Client, productKey, deviceName string, payload []byte) error {
	return nil
}

// DownstreamThingConfigPush see interface EventProc
func (NopEvt) DownstreamThingConfigPush(c *Client, rsp *ConfigPushRequest) error {
	return nil
}

// DownstreamThingServicePropertySet see interface EventProc
func (NopEvt) DownstreamThingServicePropertySet(c *Client, topic string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}
	return c.SendResponse(URIServiceReplyWithRequestURI(topic), rsp.ID, CodeSuccess, "{}")
}

// DownstreamThingServiceRequest see interface EventProc
func (NopEvt) DownstreamThingServiceRequest(c *Client, productKey, deviceName, srvID string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}

	return c.SendResponse(c.URIService(URISysPrefix, URIThingServiceRequest, productKey, deviceName, srvID),
		rsp.ID, CodeSuccess, "{}")
}

// DownStreamRRPCRequest see interface EventProc
func (NopEvt) DownStreamRRPCRequest(c *Client, productKey, deviceName, messageID string, payload []byte) error {
	return c.Publish(c.URIService(URISysPrefix, URIRRPCResponse, productKey, deviceName, messageID),
		0, `{"note":"default system RRPC implementation"}`)
}

// DownStreamExtRRPCRequest see interface EventProc
func (NopEvt) DownStreamExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	return c.Publish(rawURI, 0, "default ext RRPC implementation")
}
