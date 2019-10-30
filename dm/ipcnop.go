package dm

import (
	"encoding/json"
)

// DevNopUserProc 实现DevUserProc接口的空实现
type DevNopUserProc struct{}

// DownstreamThingModelUpRawReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingModelUpRawReply(c *Client, productKey, deviceName string, payload []byte) error {
	return nil
}

// DownstreamThingEventPropertyPostReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingEventPropertyPostReply(c *Client, rsp *Response) error {
	return nil
}

// DownstreamThingEventPostReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingEventPostReply(c *Client, eventID string, rsp *Response) error {
	return nil
}

// DownstreamThingDeviceInfoUpdateReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingDeviceInfoUpdateReply(c *Client, rsp *Response) error {
	return nil
}

// DownstreamThingDeviceInfoDeleteReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingDeviceInfoDeleteReply(c *Client, rsp *Response) error {
	return nil
}

// DownstreamThingDesiredPropertyGetReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingDesiredPropertyGetReply(c *Client, rsp *Response) error {
	return nil
}

// DownstreamThingDesiredPropertyDeleteReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingDesiredPropertyDeleteReply(c *Client, rsp *Response) error {
	return nil
}

// DownstreamThingDsltemplateGetReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingDsltemplateGetReply(c *Client, rsp *Response) error {
	return nil
}

// DownstreamThingDynamictslGetReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingDynamictslGetReply(c *Client, rsp *Response) error {
	return nil
}

// DownstreamExtNtpResponse see interface DevUserProc
func (DevNopUserProc) DownstreamExtNtpResponse(c *Client, rsp *NtpResponsePayload) error {
	return nil
}

// DownstreamThingConfigGetReply see interface DevUserProc
func (DevNopUserProc) DownstreamThingConfigGetReply(c *Client, rsp *ConfigGetResponse) error {
	return nil
}

// DownstreamExtErrorResponse see interface DevUserProc
func (DevNopUserProc) DownstreamExtErrorResponse(c *Client, rsp *Response) error {
	return nil
}

// DownstreamThingModelDownRaw see interface DevUserProc
func (DevNopUserProc) DownstreamThingModelDownRaw(c *Client, productKey, deviceName string, payload []byte) error {
	return nil
}

// DownstreamThingConfigPush see interface DevUserProc
func (DevNopUserProc) DownstreamThingConfigPush(c *Client, rsp *ConfigPushRequest) error {
	return nil
}

// DownstreamThingServicePropertySet see interface DevUserProc
func (DevNopUserProc) DownstreamThingServicePropertySet(c *Client, topic string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}
	return c.SendResponse(URIServiceReplyWithRequestURI(topic), rsp.ID, CodeSuccess, "{}")
}

// DownstreamThingServiceRequest see interface DevUserProc
func (DevNopUserProc) DownstreamThingServiceRequest(c *Client, productKey, deviceName, srvID string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}

	return c.SendResponse(c.URIService(URISysPrefix, URIThingServiceRequest, productKey, deviceName, srvID),
		rsp.ID, CodeSuccess, "{}")
}

// DownStreamRRPCRequest see interface DevUserProc
func (DevNopUserProc) DownStreamRRPCRequest(c *Client, productKey, deviceName, messageID string, payload []byte) error {
	return c.Publish(c.URIService(URISysPrefix, URIRRPCResponse, productKey, deviceName, messageID),
		0, `{"note":"default system RRPC implementation"}`)
}

// DownStreamExtRRPCRequest see interface DevUserProc
func (DevNopUserProc) DownStreamExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	return c.Publish(rawURI, 0, "default ext RRPC implementation")
}
