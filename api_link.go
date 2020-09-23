// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import (
	"encoding/json"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
)

/**************************************** config *****************************/

// LinkThingConfigGet 获取配置参数,同步
func (sf *Client) LinkThingConfigGet(pk, dn string, timeout time.Duration) (ConfigParamsData, error) {
	token, err := sf.ThingConfigGet(pk, dn)
	if err != nil {
		return ConfigParamsData{}, err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return ConfigParamsData{}, err
	}
	return msg.Data.(ConfigParamsData), nil
}

/**************************************** event *****************************/

// LinkThingEventPropertyPost 设备上报属性数据,同步
func (sf *Client) LinkThingEventPropertyPost(pk, dn string, params interface{}, timeout time.Duration) error {
	token, err := sf.ThingEventPropertyPost(pk, dn, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

// LinkThingEventPost 设备事件上报,同步
func (sf *Client) LinkThingEventPost(pk, dn, eventID string, params interface{}, timeout time.Duration) error {
	token, err := sf.ThingEventPost(pk, dn, eventID, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

// LinkThingEventPropertyPackPost 网关批量上报数据,同步
func (sf *Client) LinkThingEventPropertyPackPost(params interface{}, timeout time.Duration) error {
	token, err := sf.ThingEventPropertyPackPost(params)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

// LinkThingEventPropertyHistoryPost 物模型历史数据上报,同步
func (sf *Client) LinkThingEventPropertyHistoryPost(params interface{}, timeout time.Duration) error {
	token, err := sf.ThingEventPropertyHistoryPost(params)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

/**************************************** desired *****************************/

// LinkThingDesiredPropertyGet 获取期望属性值,同步
func (sf *Client) LinkThingDesiredPropertyGet(pk, dn string,
	params []string, timeout time.Duration) (json.RawMessage, error) {
	token, err := sf.ThingDesiredPropertyGet(pk, dn, params)
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data.(json.RawMessage), nil
}

// LinkThingDesiredPropertyDelete 清空期望属性值,同步
func (sf *Client) LinkThingDesiredPropertyDelete(pk, dn string, params interface{}, timeout time.Duration) error {
	token, err := sf.ThingDesiredPropertyDelete(pk, dn, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

/**************************************** label *****************************/

// LinkThingDeviceInfoUpdate 设备信息上传(如厂商,设备型号等,可以保存为设备标签),同步
func (sf *Client) LinkThingDeviceInfoUpdate(pk, dn string, params interface{}, timeout time.Duration) error {
	token, err := sf.ThingDeviceInfoUpdate(pk, dn, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

// LinkThingDeviceInfoDelete 删除标签信息.同步
func (sf *Client) LinkThingDeviceInfoDelete(pk, dn string, params interface{}, timeout time.Duration) error {
	token, err := sf.ThingDeviceInfoDelete(pk, dn, params)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

/**************************************** template *****************************/

// LinkThingDsltemplateGet 设备可以通过上行请求获取设备的TSL模板(包含属性、服务和事件的定义),同步
func (sf *Client) LinkThingDsltemplateGet(pk, dn string, timeout time.Duration) (json.RawMessage, error) {
	token, err := sf.ThingDsltemplateGet(pk, dn)
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data.(json.RawMessage), nil
}

// LinkThingDynamictslGet 获取动态tsl,同步
func (sf *Client) LinkThingDynamictslGet(pk, dn string, timeout time.Duration) (json.RawMessage, error) {
	token, err := sf.ThingDynamictslGet(pk, dn)
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data.(json.RawMessage), nil
}

// LinkThingConfigLogGet 获取日志配置,同步
func (sf *Client) LinkThingConfigLogGet(pk, dn string,
	clp ConfigLogParam, timeout time.Duration) (ConfigLogParamData, error) {
	token, err := sf.ThingConfigLogGet(pk, dn, clp)
	if err != nil {
		return ConfigLogParamData{}, err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return ConfigLogParamData{}, err
	}
	return msg.Data.(ConfigLogParamData), nil
}

// LinkThingLogPost 设备上报日志内容,同步
func (sf *Client) LinkThingLogPost(pk, dn string, lp []LogParam, timeout time.Duration) error {
	token, err := sf.ThingLogPost(pk, dn, lp)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

/**************************************** reqister *****************************/

// LinkThingSubRegister 同步子设备注册,
func (sf *Client) LinkThingSubRegister(pk, dn string, timeout time.Duration) ([]SubRegisterData, error) {
	token, err := sf.thingSubRegister(pk, dn)
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return nil, err
	}
	data := msg.Data.([]SubRegisterData)
	for _, v := range msg.Data.([]SubRegisterData) {
		sf.SetDeviceSecret(v.ProductKey, v.DeviceName, v.DeviceSecret)      // nolint: errcheck
		sf.SetDeviceStatus(v.ProductKey, v.DeviceName, DevStatusRegistered) // nolint: errcheck
	}
	return data, nil
}

/**************************************** network *****************************/

// LinkThingTopoAdd 添加设备拓扑关系,同步
func (sf *Client) LinkThingTopoAdd(pk, dn string, timeout time.Duration) error {
	token, err := sf.thingTopoAdd(pk, dn)
	if err != nil {
		return err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return err
	}
	for _, pair := range msg.Data.([]infra.MetaPair) {
		sf.SetDeviceStatus(pair.ProductKey, pair.DeviceName, DevStatusAttached) // nolint: errcheck
	}
	return nil
}

// LinkThingTopoDelete 删除网关与子设备的拓扑关系
func (sf *Client) LinkThingTopoDelete(pk, dn string, timeout time.Duration) error {
	token, err := sf.thingTopoDelete(pk, dn)
	if err != nil {
		return err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return err
	}
	for _, pair := range msg.Data.([]infra.MetaPair) {
		sf.SetDeviceStatus(pair.ProductKey, pair.DeviceName, DevStatusRegistered) // nolint: errcheck
	}
	return nil
}

// LinkThingTopoGet 获取该网关和子设备的拓扑关系,同步
func (sf *Client) LinkThingTopoGet(timeout time.Duration) ([]infra.MetaPair, error) {
	token, err := sf.ThingTopoGet()
	if err != nil {
		return nil, err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data.([]infra.MetaPair), err
}

// LinkThingListFound 发现设备列表上报,同步
func (sf *Client) LinkThingListFound(pairs []infra.MetaPair, timeout time.Duration) error {
	token, err := sf.ThingListFound(pairs)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

/**************************************** session *****************************/

// LinkExtCombineLogin 子设备上线,同步
func (sf *Client) LinkExtCombineLogin(cp CombinePair, timeout time.Duration) error {
	token, err := sf.extCombineLogin(cp)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	if err != nil {
		return err
	}
	sf.SetDeviceStatus(cp.ProductKey, cp.DeviceName, DevStatusLogined) // nolint: errcheck
	return nil
}

// LinkExtCombineBatchLogin 子设备批量上线,同步
func (sf *Client) LinkExtCombineBatchLogin(pairs []CombinePair, timeout time.Duration) error {
	token, err := sf.extCombineBatchLogin(pairs)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	if err != nil {
		return err
	}

	for _, cp := range pairs {
		sf.SetDeviceStatus(cp.ProductKey, cp.DeviceName, DevStatusLogined) // nolint: errcheck
	}
	return nil
}

// LinkExtCombineLogout 子设备下线,同步
func (sf *Client) LinkExtCombineLogout(pk, dn string, timeout time.Duration) error {
	token, err := sf.extCombineLogout(pk, dn)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	if err != nil {
		return err
	}
	sf.SetDeviceStatus(pk, dn, DevStatusAttached) // nolint: errcheck
	return nil
}

// LinkExtCombineBatchLogout 子设备批量下线,同步
func (sf *Client) LinkExtCombineBatchLogout(pairs []infra.MetaPair, timeout time.Duration) error {
	token, err := sf.extCombineBatchLogout(pairs)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	if err != nil {
		return err
	}
	for _, cp := range pairs {
		sf.SetDeviceStatus(cp.ProductKey, cp.DeviceName, DevStatusAttached) // nolint: errcheck
	}
	return nil
}

/**************************************** ota *****************************/

// LinkThingOtaFirmwareGet 请求固件信息,同步
func (sf *Client) LinkThingOtaFirmwareGet(pk, dn string,
	param OtaFirmwareParam, timeout time.Duration) (OtaFirmwareData, error) {
	token, err := sf.ThingOtaFirmwareGet(pk, dn, param)
	if err != nil {
		return OtaFirmwareData{}, err
	}
	msg, err := token.Wait(timeout)
	if err != nil {
		return OtaFirmwareData{}, err
	}
	return msg.Data.(OtaFirmwareData), nil
}

/**************************************** diag *****************************/

// LinkThingDiagPost 设备主动上报当前网络状态,同步
func (sf *Client) LinkThingDiagPost(pk, dn string, p P, timeout time.Duration) error {
	token, err := sf.ThingDiagPost(pk, dn, p)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}

// LinkThingDiagHistoryPost 设备主动上报历史网络状态,同步
func (sf *Client) LinkThingDiagHistoryPost(pk, dn string, p []P, timeout time.Duration) error {
	token, err := sf.ThingDiagHistoryPost(pk, dn, p)
	if err != nil {
		return err
	}
	_, err = token.Wait(timeout)
	return err
}
