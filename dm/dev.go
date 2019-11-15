// Package dm 实现阿里去物模型
package dm

import (
	"encoding/json"
	"fmt"
	"time"
)

// upstreamThingModelUpRaw 上传透传数据
func (sf *Client) upstreamThingModelUpRaw(devID int, payload interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	err = sf.Publish(sf.URIService(URISysPrefix, URIThingModelUpRaw, node.ProductKey(), node.DeviceName()), 1, payload)
	if err != nil {
		return err
	}
	sf.debug("upstream thing <model>: up raw")
	return nil
}

// upstreamThingEventPropertyPost 上传属性数据
func (sf *Client) upstreamThingEventPropertyPost(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingEventPropertyPost, node.ProductKey(), node.DeviceName()),
		id, methodEventPropertyPost, params)
	if err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypeEventPropertyPost)
	sf.debug("upstream thing <event>: property post,@%d", id)
	return nil
}

// upstreamThingEventPost 事件上传
func (sf *Client) upstreamThingEventPost(devID int, eventID string, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}
	id := sf.RequestID()
	method := fmt.Sprintf(methodEventFormatPost, eventID)
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingEventPost, node.ProductKey(), node.DeviceName(), eventID),
		id, method, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeEventPost)
	sf.debug("upstream thing <event>: %s post,@%d", eventID, id)
	return nil
}

func (sf *Client) upstreamThingEventPropertyPackPost(params interface{}) error {
	id := sf.RequestID()
	err := sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingEventPropertyPackPost),
		id, methodEventPropertyPackPost, params)
	if err != nil {
		return err
	}

	sf.CacheInsert(id, DevNodeLocal, MsgTypeEventPropertyPackPost)
	sf.debug("upstream thing <deviceInfo>: update,@%d", id)
	return nil
}

// upstreamThingDeviceInfoUpdate 设备信息上传
func (sf *Client) upstreamThingDeviceInfoUpdate(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDeviceInfoUpdate, node.ProductKey(), node.DeviceName()),
		id, methodDeviceInfoUpdate, params)
	if err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypeDeviceInfoUpdate)
	sf.debug("upstream thing <deviceInfo>: update,@%d", id)
	return nil
}

// upstreamThingDeviceInfoDelete 设备信息删除
func (sf *Client) upstreamThingDeviceInfoDelete(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDeviceInfoDelete, node.ProductKey(), node.DeviceName()),
		id, methodDeviceInfoDelete, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDeviceInfoDelete)
	sf.debug("upstream thing <deviceInfo>: delete,@%d", id)
	return nil
}

// upstreamThingDesiredPropertyGet 获取期望值
func (sf *Client) upstreamThingDesiredPropertyGet(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDesiredPropertyGet, node.ProductKey(), node.DeviceName()),
		id, methodDesiredPropertyGet, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDesiredPropertyGet)
	sf.debug("upstream thing <desired>: property get,@%d", id)
	return nil
}

// upstreamThingDesiredPropertyDelete 清空期望值
func (sf *Client) upstreamThingDesiredPropertyDelete(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDesiredPropertyDelete, node.ProductKey(), node.DeviceName()),
		id, methodDesiredPropertyDelete, params)
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeDesiredPropertyDelete)
	sf.debug("upstream thing <desired>: property delete,@%d", id)
	return nil
}

// upstreamThingDsltemplateGet 设备可以通过上行请求获取设备的TSL模板（包含属性、服务和事件的定义）
// see https://help.aliyun.com/document_detail/89305.html?spm=a2c4g.11186623.6.672.5d3d70374hpPcx
func (sf *Client) upstreamThingDsltemplateGet(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDslTemplateGet, node.ProductKey(), node.DeviceName()),
		id, methodDslTemplateGet, "{}"); err != nil {
		return err
	}

	sf.CacheInsert(id, devID, MsgTypeDsltemplateGet)
	sf.debug("upstream thing <dsl template>: get,@%d", id)
	return nil
}

// upstreamThingDynamictslGet 获取动态tsl
func (sf *Client) upstreamThingDynamictslGet(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingDynamicTslGet, node.ProductKey(), node.DeviceName()), id,
		methodDynamicTslGet, `{"nodes":["type","identifier"],"addDefault":false}`); err != nil {
		return err
	}
	sf.CacheInsert(id, DevNodeLocal, MsgTypeDynamictslGet)
	sf.debug("upstream thing <dynamic tsl>: get,@%d", id)
	return nil
}

// NtpResponsePayload ntp回复payload
type NtpResponsePayload struct {
	DeviceSendTime int64 `json:"deviceSendTime,string"`
	ServerRecvTime int64 `json:"serverRecvTime,string"`
	ServerSendTime int64 `json:"serverSendTime,string"`
}

// upstreamExtNtpRequest ntp请求
// 发送一条Qos = 0的消息,并带上设备当前的时间戳,平台将回复 设备的发送时间,平台的接收时间, 平台的发送时间.
// 设备计算当前精确时间 = (平台接收时间 + 平台发送时间 + 设备接收时间 - 设备发送时间) / 2
func (sf *Client) upstreamExtNtpRequest() error {
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

// upstreamThingConfigGet 获取配置参数
func (sf *Client) upstreamThingConfigGet(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIService(URISysPrefix, URIThingConfigGet, node.ProductKey(), node.DeviceName()),
		id, methodConfigGet, `{"configScope":"product","getType":"file"}`); err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeConfigGet)
	sf.debug("upstream thing <config>: get,@%d", id)
	return nil
}

// OTARequest OTA请求体
type OTARequest struct {
	ID     int         `json:"id,string"`
	Params interface{} `json:"params"`
}

// OTAFirmwareVersionParams OTA固件参数域
type OTAFirmwareVersionParams struct {
	Version string `json:"version"`
}

// upstreamOATFirmwareVersion 上报固件版本
func (sf *Client) upstreamOATFirmwareVersion(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	req, err := json.Marshal(OTARequest{id, params})
	if err != nil {
		return err
	}

	if err = sf.Publish(sf.URIService(URIOtaDeviceInformPrefix, "", node.ProductKey(), node.DeviceName()),
		1, req); err != nil {
		return err
	}

	//sf.CacheInsert(id, devID, MsgTypeReportFirmwareVersion)
	sf.debug("upstream version <OTA>: inform,@%d", id)
	return nil
}

// OTA下载进度比
const (
	OTAProgressStepUpgradeFailed  = -1
	OTAProgressStepDownloadFailed = -2
	OTAProgressStepVerifyFailed   = -3
	OTAProgressStepProgramFailed  = -4
)

// OTAProgressParams 下载过程上报参数域
type OTAProgressParams struct {
	Step int    `json:"step,string"`
	Desc string `json:"desc"`
}

func (sf *Client) upstreamOTAProgress(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	req, err := json.Marshal(OTARequest{id, params})
	if err != nil {
		return err
	}

	if err = sf.Publish(sf.URIService(URIOtaDeviceProcessPrefix, "", node.ProductKey(), node.DeviceName()),
		1, req); err != nil {
		return err
	}

	//sf.CacheInsert(id, devID, MsgTypeReportFirmwareVersion)
	sf.debug("upstream step <OTA>: progress,@%d", id)
	return nil
}

// OTAUpgradeData OTA upgrade 数据域
type OTAUpgradeData struct {
	Size    int    `json:"size"`
	Version string `json:"version"`
	URL     string `json:"url"`
	MD5     string `json:"md5"`
}

// OTAUpgradeRequest OTA upgrade 请求
type OTAUpgradeRequest struct {
	Code    int            `json:"code,string"`
	Data    OTAUpgradeData `json:"data"`
	ID      int            `json:"id"`
	Message string         `json:"message"`
}
