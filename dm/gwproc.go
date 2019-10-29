package dm

type GwNopUserProc struct{}

func (GwNopUserProc) DownstreamGwExtSubDevRegisterReply(m *Client, rsp *GwSubDevRegisterResponse) error {
	return nil
}

func (GwNopUserProc) DownstreamGwExtSubDevCombineLoginReply(m *Client, rsp *Response) error {
	return nil
}

// DownstreamGwExtSubDevCombineLogoutReply
func (GwNopUserProc) DownstreamGwExtSubDevCombineLogoutReply(m *Client, rsp *Response) error {
	return nil
}

func (GwNopUserProc) DownstreamGwSubDevThingDisable(m *Client, productKey, deviceName string) error {
	return nil
}
func (GwNopUserProc) DownstreamGwSubDevThingEnable(m *Client, productKey, deviceName string) error {
	return nil
}
func (GwNopUserProc) DownstreamGwSubDevThingDelete(m *Client, productKey, deviceName string) error {
	return nil
}
func (GwNopUserProc) DownstreamGwThingTopoAddReply(m *Client, rsp *Response) error {
	return nil
}

func (GwNopUserProc) DownstreamGwThingTopoDeleteReply(m *Client, rsp *Response) error {
	return nil
}

func (GwNopUserProc) DownstreamGwThingTopoGetReply(m *Client, rsp *GwTopoGetResponse) error {
	return nil
}
