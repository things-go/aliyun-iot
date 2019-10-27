// Package model 实现阿里去物模型
package model

import (
	"fmt"
	"time"
)

// UpstreamThingModelUpRaw 上传透传数据
func (sf *Manager) UpstreamThingModelUpRaw(devID int, payload interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	return sf.Publish(sf.URIService(URISysPrefix, URIThingModelUpRaw, node.ProductKey, node.DeviceName), 1, payload)
}

// UpstreamThingPropertyPost 上传属性数据
func (sf *Manager) UpstreamThingEventPropertyPost(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingEventPropertyPost, node.ProductKey, node.DeviceName),
		id, methodEventPropertyPost, params)
	if err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypePropertyPost, "property")
	return nil
}

// UpstreamThingEventPost 事件上传
func (sf *Manager) UpstreamThingEventPost(devID int, EventID string, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}
	id := sf.RequestID()
	method := fmt.Sprintf(methodEventFormatPost, EventID)
	err = sf.SendRequest(sf.URIService(URISysPrefix, fmt.Sprintf(URIThingEventPost, EventID), node.ProductKey, node.DeviceName),
		id, method, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeEventPost, method)
	return nil

}

// UpstreamThingDeviceInfoUpdate 设备信息上传
func (sf *Manager) UpstreamThingDeviceInfoUpdate(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDeviceInfoUpdate, node.ProductKey, node.DeviceName),
		id, methodDeviceInfoUpdate, params)
	if err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypeDeviceInfoUpdate, methodDeviceInfoUpdate)
	return nil
}

// UpstreamThingDeviceInfoDelete 设备信息删除
func (sf *Manager) UpstreamThingDeviceInfoDelete(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDeviceInfoDelete, node.ProductKey, node.DeviceName),
		id, methodDeviceInfoDelete, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDeviceInfoDelete, methodDeviceInfoDelete)
	return nil
}

// UpstreamThingPropertyDesiredGet 获取期望值
func (sf *Manager) UpstreamThingPropertyDesiredGet(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDesiredPropertyGet, node.ProductKey, node.DeviceName),
		id, methodDesiredPropertyGet, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypePropertyDesiredGet, methodDesiredPropertyGet)
	return nil
}

// UpstreamThingPropertyDesiredDelete 清空期望值
func (sf *Manager) UpstreamThingPropertyDesiredDelete(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDesiredPropertyDelete, node.ProductKey, node.DeviceName),
		id, methodDesiredPropertyDelete, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypePropertyDesiredDelete, methodDesiredPropertyDelete)
	return nil
}

// UpstreamThingDsltemplateGet 获取
func (sf *Manager) UpstreamThingDsltemplateGet() error {
	uri := sf.URIService(URISysPrefix, URIThingDslTemplateGet, sf.opt.productKey, sf.opt.deviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDslTemplateGet, "{}")
}

// UpstreamThingDynamictslGet 获取
func (sf *Manager) UpstreamThingDynamictslGet() error {
	uri := sf.URIService(URISysPrefix, URIThingDynamicTslGet, sf.opt.productKey, sf.opt.deviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDynamicTslGet, `{"nodes\":["type","identifier"],"addDefault":false}`)
}

// UpstreamExtNtpRequest ntp请求
// 发送一条Qos = 0的消息,并带上设备当前的时间戳,平台将回复 设备的发送时间,平台的接收时间, 平台的发送时间.
// 设备计算当前精确时间 = (平台接收时间 + 平台发送时间 + 设备接收时间 - 设备发送时间) / 2
func (sf *Manager) UpstreamExtNtpRequest() error {
	return sf.Publish(sf.URIService(URIExtNtpPrefix, URINtpRequest, sf.opt.productKey, sf.opt.deviceName),
		0, fmt.Sprintf(`{"deviceSendTime":"%d"}`, time.Now().Unix()))
}

// NtpResponse ntp回复payload
type NtpResponse struct {
	DeviceSendTime int `json:"deviceSendTime,string"`
	ServerRecvTime int `json:"serverRecvTime,string"`
	ServerSendTime int `json:"serverSendTime,string"`
}
