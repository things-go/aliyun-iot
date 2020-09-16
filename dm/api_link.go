package dm

import (
	"encoding/json"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
)

/**************************************** config *****************************/

// LinkThingConfigGet 获取配置参数,同步
func (sf *Client) LinkThingConfigGet(pk, dn string) (ConfigParamsData, error) {
	token, err := sf.ThingConfigGet(pk, dn)
	if err != nil {
		return ConfigParamsData{}, err
	}
	msg, err := token.Wait(time.Second)
	if err != nil {
		return ConfigParamsData{}, err
	}
	return msg.Data.(ConfigParamsData), nil
}

/**************************************** event *****************************/

// LinkThingEventPropertyPost 设备上报属性数据,同步
func (sf *Client) LinkThingEventPropertyPost(pk, dn string, params interface{}) error {
	token, err := sf.ThingEventPropertyPost(pk, dn, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

// LinkThingEventPost 设备事件上报,同步
func (sf *Client) LinkThingEventPost(pk, dn, eventID string, params interface{}) error {
	token, err := sf.ThingEventPost(pk, dn, eventID, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

// LinkThingEventPropertyPackPost 网关批量上报数据,同步
func (sf *Client) LinkThingEventPropertyPackPost(params interface{}) error {
	token, err := sf.ThingEventPropertyPackPost(params)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

// LinkThingEventPropertyHistoryPost 物模型历史数据上报,同步
func (sf *Client) LinkThingEventPropertyHistoryPost(params interface{}) error {
	token, err := sf.ThingEventPropertyHistoryPost(params)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

/**************************************** desired *****************************/

// LinkThingDesiredPropertyGet 获取期望属性值,同步
func (sf *Client) LinkThingDesiredPropertyGet(pk, dn string, params []string) (json.RawMessage, error) {
	token, err := sf.ThingDesiredPropertyGet(pk, dn, params)
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(time.Second)
	if err != nil {
		return nil, err
	}
	return msg.Data.(json.RawMessage), nil
}

// LinkThingDesiredPropertyDelete 清空期望属性值,同步
func (sf *Client) LinkThingDesiredPropertyDelete(pk, dn string, params interface{}) error {
	token, err := sf.ThingDesiredPropertyDelete(pk, dn, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

/**************************************** label *****************************/

// LinkThingDeviceInfoUpdate 设备信息上传(如厂商,设备型号等,可以保存为设备标签),同步
func (sf *Client) LinkThingDeviceInfoUpdate(pk, dn string, params interface{}) error {
	token, err := sf.ThingDeviceInfoUpdate(pk, dn, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

// LinkThingDeviceInfoDelete 删除标签信息.同步
func (sf *Client) LinkThingDeviceInfoDelete(pk, dn string, params interface{}) error {
	token, err := sf.ThingDeviceInfoDelete(pk, dn, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

/**************************************** template *****************************/

// LinkThingDsltemplateGet 设备可以通过上行请求获取设备的TSL模板(包含属性、服务和事件的定义),同步
func (sf *Client) LinkThingDsltemplateGet(pk, dn string) (json.RawMessage, error) {
	token, err := sf.ThingDsltemplateGet(pk, dn)
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(time.Second)
	if err != nil {
		return nil, err
	}
	return msg.Data.(json.RawMessage), nil
}

// LinkThingDynamictslGet 获取动态tsl,同步
func (sf *Client) LinkThingDynamictslGet(pk, dn string) (json.RawMessage, error) {
	token, err := sf.ThingDsltemplateGet(pk, dn)
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(time.Second)
	if err != nil {
		return nil, err
	}
	return msg.Data.(json.RawMessage), nil
}

/**************************************** reqister *****************************/

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
	for _, v := range msg.Data.([]SubRegisterData) {
		sf.SetDeviceSecret(v.ProductKey, v.DeviceName, v.DeviceSecret)      // nolint: errcheck
		sf.SetDeviceStatus(v.ProductKey, v.DeviceName, DevStatusRegistered) // nolint: errcheck
	}
	return nil
}

/**************************************** network *****************************/

// LinkThingTopoAdd 添加设备拓扑关系,同步
func (sf *Client) LinkThingTopoAdd(pk, dn string) error {
	token, err := sf.ThingTopoAdd(pk, dn)
	if err != nil {
		return err
	}
	msg, err := token.Wait(time.Second)
	if err != nil {
		return err
	}
	for _, pair := range msg.Data.([]infra.MetaPair) {
		sf.SetDeviceStatus(pair.ProductKey, pair.DeviceName, DevStatusAttached) // nolint: errcheck
	}
	return nil
}

// LinkThingTopoDelete 删除网关与子设备的拓扑关系
func (sf *Client) LinkThingTopoDelete(pk, dn string) error {
	token, err := sf.ThingTopoDelete(pk, dn)
	if err != nil {
		return err
	}
	msg, err := token.Wait(time.Second)
	if err != nil {
		return err
	}
	for _, pair := range msg.Data.([]infra.MetaPair) {
		sf.SetDeviceStatus(pair.ProductKey, pair.DeviceName, DevStatusRegistered) // nolint: errcheck
	}
	return nil
}

// LinkThingTopoGet 获取该网关和子设备的拓扑关系,同步
func (sf *Client) LinkThingTopoGet() ([]infra.MetaPair, error) {
	token, err := sf.ThingTopoGet()
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(time.Second)
	if err != nil {
		return nil, err
	}
	return msg.Data.([]infra.MetaPair), err
}

// LinkThingListFound 发现设备列表上报,同步
func (sf *Client) LinkThingListFound(pairs []infra.MetaPair) error {
	token, err := sf.ThingListFound(pairs)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	return err
}

/**************************************** session *****************************/

// LinkExtCombineLogin 子设备上线,同步
func (sf *Client) LinkExtCombineLogin(cp CombinePair) error {
	token, err := sf.ExtCombineLogin(cp)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	if err != nil {
		return err
	}
	sf.SetDeviceStatus(cp.ProductKey, cp.DeviceName, DevStatusLogined) // nolint: errcheck
	return nil
}

// LinkExtCombineBatchLogin 子设备批量上线,同步
func (sf *Client) LinkExtCombineBatchLogin(pairs []CombinePair) error {
	token, err := sf.ExtCombineBatchLogin(pairs)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	if err != nil {
		return err
	}

	for _, cp := range pairs {
		sf.SetDeviceStatus(cp.ProductKey, cp.DeviceName, DevStatusLogined) // nolint: errcheck
	}
	return nil
}

// LinkExtCombineLogout 子设备下线,同步
func (sf *Client) LinkExtCombineLogout(pk, dn string) error {
	token, err := sf.ExtCombineLogout(pk, dn)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	if err != nil {
		return err
	}
	sf.SetDeviceStatus(pk, dn, DevStatusAttached) // nolint: errcheck
	return nil
}

// LinkExtCombineBatchLogout 子设备批量下线,同步
func (sf *Client) LinkExtCombineBatchLogout(pairs []infra.MetaPair) error {
	token, err := sf.ExtCombineBatchLogout(pairs)
	if err != nil {
		return err
	}
	_, err = token.Wait(time.Second)
	if err != nil {
		return err
	}
	for _, cp := range pairs {
		sf.SetDeviceStatus(cp.ProductKey, cp.DeviceName, DevStatusLogined) // nolint: errcheck
	}
	return nil
}

/**************************************** ota *****************************/

// LinkThingOtaFirmwareGet 请求固件信息,同步
func (sf *Client) LinkThingOtaFirmwareGet(pk, dn string, param OtaFirmwareParam) (OtaFirmwareData, error) {
	token, err := sf.ThingOtaFirmwareGet(pk, dn, param)
	if err != nil {
		return OtaFirmwareData{}, err
	}
	msg, err := token.Wait(time.Second)
	if err != nil {
		return OtaFirmwareData{}, err
	}
	return msg.Data.(OtaFirmwareData), nil
}
