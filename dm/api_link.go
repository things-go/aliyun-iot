package dm

import (
	"time"
)

// LinkThingSubRegister 同步子设备注册
func (sf *Client) LinkThingSubRegister(pk, dn string) error {
	token, err := sf.ThingSubRegister(pk, dn)
	if err != nil {
		return err
	}
	msg, err := token.Wait(time.Second * 3)
	if err != nil {
		return err
	}
	for _, v := range msg.Data.([]GwSubRegisterData) {
		sf.SetDeviceSecret(v.ProductKey, v.DeviceName, v.DeviceSecret)      // nolint: errcheck
		sf.SetDeviceStatus(v.ProductKey, v.DeviceName, DevStatusRegistered) // nolint: errcheck
	}
	return nil
}

func (sf *Client) LinkThingGwTopoAdd(pk, dn string) error {
	token, err := sf.ThingGwTopoAdd(pk, dn)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

// LinkThingGwTopoDelete 删除网关与子设备的拓扑关系
func (sf *Client) LinkThingGwTopoDelete(pk, dn string) error {
	token, err := sf.ThingGwTopoDelete(pk, dn)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

func (sf *Client) LinkExtCombineLogin(pk, dn string) error {
	token, err := sf.ExtCombineLogin(pk, dn)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

func (sf *Client) LinkExtCombineLogout(pk, dn string) error {
	token, err := sf.ExtCombineLogout(pk, dn)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}
