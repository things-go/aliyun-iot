package model

type gwUserProc struct{}

func (gwUserProc) DownstreamGwExtSubDevRegisterReply(m *Manager, rsp *GwSubDevRegisterResponse) error {
	return nil
}

func (gwUserProc) DownstreamGwExtSubDevCombineLoginReply(m *Manager, rsp *Response) error {
	return nil
}

// DownstreamGwExtSubDevCombineLogoutReply
func (gwUserProc) DownstreamGwExtSubDevCombineLogoutReply(m *Manager, rsp *Response) error {
	return nil
}

func (gwUserProc) DownstreamGwSubDevThingDisable(m *Manager, productKey, deviceName string) error {
	return nil
}
func (gwUserProc) DownstreamGwSubDevThingEnable(m *Manager, productKey, deviceName string) error {
	return nil
}
func (gwUserProc) DownstreamGwSubDevThingDelete(m *Manager, productKey, deviceName string) error {
	return nil
}
func (gwUserProc) DownstreamGwThingTopoAddReply(m *Manager, rsp *Response) error {
	return nil
}

func (gwUserProc) DownstreamGwThingTopoDeleteReply(m *Manager, rsp *Response) error {
	return nil
}

func (gwUserProc) DownstreamGwThingTopoGetReply(m *Manager, rsp *GwTopoGetResponse) error {
	return nil
}
