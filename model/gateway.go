package model

type CombineSubDevRequest struct {
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	ClientId     string `json:"clientId"`
	Timestamp    string `json:"timestamp"`
	SignMethod   string `json:"signMethod"`
	Sign         string `json:"sign"`
	CleanSession string `json:"cleanSession"`
}

func (sf *Manager) ExtSubDevCombineLogin(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	_, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}
	// 子设备登陆,要用网关的productKey和deviceName
	_ := sf.URIService(URIExtSessionPrefix, CombineSubDevLogin)

	return nil
}

func (sf *Manager) ExtSubDevCombineLogout(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	_, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}
	// 子设备登陆,要用网关的productKey和deviceName
	_ := sf.URIService(URIExtSessionPrefix, CombineSubDevLogout)

	return nil
}
