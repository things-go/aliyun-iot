package dm

// GwNopUserProc 实现GatewayUserProc的接口的空实现
type GwNopUserProc struct{}

// DownstreamGwExtSubDevRegisterReply see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwExtSubDevRegisterReply(m *Client, rsp *GwSubDevRegisterResponse) error {
	return nil
}

// DownstreamGwExtSubDevCombineLoginReply see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwExtSubDevCombineLoginReply(m *Client, rsp *Response) error {
	return nil
}

// DownstreamGwExtSubDevCombineLogoutReply see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwExtSubDevCombineLogoutReply(m *Client, rsp *Response) error {
	return nil
}

// DownstreamGwSubDevThingDisable see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwSubDevThingDisable(m *Client, productKey, deviceName string) error {
	return nil
}

// DownstreamGwSubDevThingEnable see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwSubDevThingEnable(m *Client, productKey, deviceName string) error {
	return nil
}

// DownstreamGwSubDevThingDelete see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwSubDevThingDelete(m *Client, productKey, deviceName string) error {
	return nil
}

// DownstreamGwThingTopoAddReply see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwThingTopoAddReply(m *Client, rsp *Response) error {
	return nil
}

// DownstreamGwThingTopoDeleteReply see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwThingTopoDeleteReply(m *Client, rsp *Response) error {
	return nil
}

// DownstreamGwThingTopoGetReply see interface GatewayUserProc
func (GwNopUserProc) DownstreamGwThingTopoGetReply(m *Client, rsp *GwTopoGetResponse) error {
	return nil
}
