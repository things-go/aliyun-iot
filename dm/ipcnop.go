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

// EvtExtErrorResponse see interface EventProc
func (NopEvt) EvtExtErrorResponse(c *Client, rsp *Response) error {
	return nil
}

// EvtThingModelDownRaw see interface EventProc
func (NopEvt) EvtThingModelDownRaw(c *Client, productKey, deviceName string, payload []byte) error {
	return nil
}

// EvtThingConfigPush see interface EventProc
func (NopEvt) EvtThingConfigPush(c *Client, productKey, deviceName string, data ConfigParamsAndData) error {
	return nil
}

// EvtThingServicePropertySet see interface EventProc
func (NopEvt) EvtThingServicePropertySet(c *Client, productKey, deviceName string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}
	return c.SendResponse(c.URIService(URISysPrefix, URIThingServicePropertySet, productKey, deviceName),
		rsp.ID, CodeSuccess, "{}")
}

// EvtThingServiceRequest see interface EventProc
func (NopEvt) EvtThingServiceRequest(c *Client, srvID, productKey, deviceName string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}

	return c.SendResponse(c.URIService(URISysPrefix, URIThingServiceRequest, productKey, deviceName, srvID),
		rsp.ID, CodeSuccess, "{}")
}

// EvtRRPCRequest see interface EventProc
func (NopEvt) EvtRRPCRequest(c *Client, messageID, productKey, deviceName string, payload []byte) error {
	return c.Publish(c.URIService(URISysPrefix, URIRRPCResponse, productKey, deviceName, messageID),
		0, `{"note":"default system RRPC implementation"}`)
}

// EvtExtRRPCRequest see interface EventProc
func (NopEvt) EvtExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	return c.Publish(rawURI, 0, "default ext RRPC implementation")
}
