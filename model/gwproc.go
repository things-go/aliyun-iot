package model

type GwNopUserProc struct{}

func (GwNopUserProc) DownstreamGwExtSubDevRegisterReply(m *Manager, rsp *GwSubDevRegisterResponse) error {
	return nil
}

func (GwNopUserProc) DownstreamGwExtSubDevCombineLoginReply(m *Manager, rsp *Response) error {
	return nil
}

// DownstreamGwExtSubDevCombineLogoutReply
func (GwNopUserProc) DownstreamGwExtSubDevCombineLogoutReply(m *Manager, rsp *Response) error {
	return nil
}

func (GwNopUserProc) DownstreamGwSubDevThingDisable(m *Manager, productKey, deviceName string) error {
	return nil
}
func (GwNopUserProc) DownstreamGwSubDevThingEnable(m *Manager, productKey, deviceName string) error {
	return nil
}
func (GwNopUserProc) DownstreamGwSubDevThingDelete(m *Manager, productKey, deviceName string) error {
	return nil
}
func (GwNopUserProc) DownstreamGwThingTopoAddReply(m *Manager, rsp *Response) error {
	return nil
}

func (GwNopUserProc) DownstreamGwThingTopoDeleteReply(m *Manager, rsp *Response) error {
	return nil
}

func (GwNopUserProc) DownstreamGwThingTopoGetReply(m *Manager, rsp *GwTopoGetResponse) error {
	return nil
}
