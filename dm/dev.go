// Package dm 实现阿里去物模型
package dm

import (
	"fmt"
	"time"
)

// UpstreamThingModelUpRaw 上传透传数据
func (sf *Client) UpstreamThingModelUpRaw(devID int, payload interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	err = sf.Publish(sf.URIService(URISysPrefix, URIThingModelUpRaw, node.productKey, node.deviceName), 1, payload)
	if err != nil {
		return err
	}
	sf.debug("upstream thing <model>: up raw")
	return nil
}

// UpstreamThingEventPropertyPost 上传属性数据
func (sf *Client) UpstreamThingEventPropertyPost(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingEventPropertyPost, node.productKey, node.deviceName),
		id, methodEventPropertyPost, params)
	if err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypeEventPropertyPost, "property")
	sf.debug("upstream thing <event>: property post,@%d", id)
	return nil
}

// UpstreamThingEventPost 事件上传
func (sf *Client) UpstreamThingEventPost(devID int, eventID string, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}
	id := sf.RequestID()
	method := fmt.Sprintf(methodEventFormatPost, eventID)
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingEventPost, node.productKey, node.deviceName, eventID),
		id, method, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeEventPost, method)
	sf.debug("upstream thing <event>: %s post,@%d", eventID, id)
	return nil
}

// UpstreamThingDeviceInfoUpdate 设备信息上传
func (sf *Client) UpstreamThingDeviceInfoUpdate(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDeviceInfoUpdate, node.productKey, node.deviceName),
		id, methodDeviceInfoUpdate, params)
	if err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypeDeviceInfoUpdate, methodDeviceInfoUpdate)
	sf.debug("upstream thing <deviceInfo>: update,@%d", id)
	return nil
}

// UpstreamThingDeviceInfoDelete 设备信息删除
func (sf *Client) UpstreamThingDeviceInfoDelete(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDeviceInfoDelete, node.productKey, node.deviceName),
		id, methodDeviceInfoDelete, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDeviceInfoDelete, methodDeviceInfoDelete)
	sf.debug("upstream thing <deviceInfo>: delete,@%d", id)
	return nil
}

// UpstreamThingDesiredPropertyGet 获取期望值
func (sf *Client) UpstreamThingDesiredPropertyGet(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDesiredPropertyGet, node.productKey, node.deviceName),
		id, methodDesiredPropertyGet, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDesiredPropertyGet, methodDesiredPropertyGet)
	sf.debug("upstream thing <desired>: property get,@%d", id)
	return nil
}

// UpstreamThingDesiredPropertyDelete 清空期望值
func (sf *Client) UpstreamThingDesiredPropertyDelete(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDesiredPropertyDelete, node.productKey, node.deviceName),
		id, methodDesiredPropertyDelete, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDesiredPropertyDelete, methodDesiredPropertyDelete)
	sf.debug("upstream thing <desired>: property delete,@%d", id)
	return nil
}

// UpstreamThingDsltemplateGet 设备可以通过上行请求获取设备的TSL模板（包含属性、服务和事件的定义）
// see https://help.aliyun.com/document_detail/89305.html?spm=a2c4g.11186623.6.672.5d3d70374hpPcx
func (sf *Client) UpstreamThingDsltemplateGet(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDslTemplateGet, node.productKey, node.deviceName),
		id, methodDslTemplateGet, "{}"); err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypeDsltemplateGet, methodDslTemplateGet)
	sf.debug("upstream thing <dsl template>: get,@%d", id)
	return nil
}

// UpstreamThingDynamictslGet 获取
// TODO: 不使用??
func (sf *Client) UpstreamThingDynamictslGet() error {
	id := sf.RequestID()
	err := sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingDynamicTslGet), id,
		methodDynamicTslGet, `{"nodes":["type","identifier"],"addDefault":false}`)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, DevLocal, MsgTypeDynamictslGet, methodDynamicTslGet)
	sf.debug("upstream thing <dynamic tsl>: get,@%d", id)
	return nil
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
func (sf *Client) UpstreamExtNtpRequest() error {
	err := sf.Publish(sf.URIServiceSelf(URIExtNtpPrefix, URINtpRequest),
		0, fmt.Sprintf(`{"deviceSendTime":"%d"}`, time.Now().Unix()))
	if err != nil {
		return err
	}
	sf.debug("upstream ext <ntp>: request")
	return nil
}

// ConfigGetParams 配置参数
type ConfigGetParams struct {
	ConfigScope string `json:"configScope"`
	GetType     string `json:"getType"`
}

// ConfigParamsAndData 配置获取参数域,或推送数据域
type ConfigParamsAndData struct {
	ConfigID   string `json:"configId"`
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
func (sf *Client) UpstreamThingConfigGet(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingConfigGet, node.productKey, node.deviceName),
		id, methodConfigGet, `{"configScope":"product","getType":"file"}`); err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeConfigGet, methodConfigGet)
	sf.debug("upstream thing <config>: get,@%d", id)
	return nil
}
