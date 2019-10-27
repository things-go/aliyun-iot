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

	sf.CacheInsert(id, devID, MsgTypeEventPropertyPost, "property")
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
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingEventPost, node.ProductKey, node.DeviceName, EventID),
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

// UpstreamThingDesiredPropertyGet 获取期望值
func (sf *Manager) UpstreamThingDesiredPropertyGet(devID int, params interface{}) error {
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
	sf.CacheInsert(id, devID, MsgTypeDesiredPropertyGet, methodDesiredPropertyGet)
	return nil
}

// UpstreamThingPropertyDesiredDelete 清空期望值
func (sf *Manager) UpstreamThingDesiredPropertyDelete(devID int, params interface{}) error {
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
	sf.CacheInsert(id, devID, MsgTypeDesiredPropertyDelete, methodDesiredPropertyDelete)
	return nil
}

// UpstreamThingDsltemplateGet 设备可以通过上行请求获取设备的TSL模板（包含属性、服务和事件的定义）
// see https://help.aliyun.com/document_detail/89305.html?spm=a2c4g.11186623.6.672.5d3d70374hpPcx
func (sf *Manager) UpstreamThingDsltemplateGet(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDslTemplateGet, node.ProductKey, node.DeviceName),
		id, methodDslTemplateGet, "{}"); err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypeDsltemplateGet, methodDslTemplateGet)
	return nil
}

// UpstreamThingDynamictslGet 获取
func (sf *Manager) UpstreamThingDynamictslGet() error {
	// TODO: 需要确定.未来审核
	uri := sf.URIServiceSelf(URISysPrefix, URIThingDynamicTslGet)
	return sf.SendRequest(uri, sf.RequestID(), methodDynamicTslGet, `{"nodes\":["type","identifier"],"addDefault":false}`)
}

// NtpResponsePayload ntp回复payload
type NtpResponsePayload struct {
	DeviceSendTime int64 `json:"deviceSendTime,string"`
	ServerRecvTime int64 `json:"serverRecvTime,string"`
	ServerSendTime int64 `json:"serverSendTime,string"`
}

// UpstreamExtNtpRequest ntp请求
// 发送一条Qos = 0的消息,并带上设备当前的时间戳,平台将回复 设备的发送时间,平台的接收时间, 平台的发送时间.
// 设备计算当前精确时间 = (平台接收时间 + 平台发送时间 + 设备接收时间 - 设备发送时间) / 2
func (sf *Manager) UpstreamExtNtpRequest() error {
	return sf.Publish(sf.URIServiceSelf(URIExtNtpPrefix, URINtpRequest),
		0, fmt.Sprintf(`{"deviceSendTime":"%d"}`, time.Now().Unix()))
}

// ConfigGetParams 配置参数
type ConfigGetParams struct {
	ConfigScope string `json:"configScope"`
	GetType     string `json:"getType"`
}

// ConfigParamsAndData 配置获取参数域,或推送数据域
type ConfigParamsAndData struct {
	ConfigId   string `json:"configId"`
	ConfigSize int64  `json:"configSize"`
	Sign       string `json:"sign"`
	SignMethod string `json:"signMethod"`
	URL        string `json:"url"`
	GetType    string `json:"getType"`
}

// ConfigGetResponse 配置获取的回复
type ConfigGetResponse struct {
	Response
	Data ConfigParamsAndData `json:"data"`
}

// ConfigPushRequest 配置推送的请求
type ConfigPushRequest struct {
	Request
	Params ConfigParamsAndData `json:"params"`
}

// UpstreamThingConfigGet 获取配置参数
func (sf *Manager) UpstreamThingConfigGet(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingConfigGet, node.ProductKey, node.DeviceName),
		id, methodConfigGet, `{"configScope":"product","getType":"file"}`); err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeConfigGet, methodConfigGet)
	return nil
}
