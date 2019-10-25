// Package model 实现阿里去物模型
package model

import (
	"fmt"
	"log"
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

	uri := URIService(URISysPrefix, URIThingModelUpRaw, node.ProductKey, node.DeviceName)
	return sf.Publish(uri, 1, payload)
}

func DownstreamThingModelUpRawReply(productKey, deviceName string, payload []byte) error {
	// hex 2 string
	return nil
}

// UpstreamThingPropertyPost 上传数属性数据
func (sf *Manager) UpstreamThingEventPropertyPost(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	uri := URIService(URISysPrefix, URIThingEventPropertyPost, node.ProductKey, node.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodPropertyPost, params)
}

func DownstreamThingEventPropertyPostReply(rsp *Response) error {
	log.Println("DownstreamThingEventPropertyPostReply")
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

	uri := URIService(URISysPrefix, fmt.Sprintf(URIThingEventPost, EventID), node.ProductKey, node.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), fmt.Sprintf(methodEventPostFormat, EventID), params)
}

func DownstreamThingEventPostReply(eventID string, rsp *Response) error {
	log.Println("DownstreamThingEventPostReply")
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

	uri := URIService(URISysPrefix, URIThingDeviceInfoUpdate, node.ProductKey, node.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDeviceInfoUpdate, params)
}

func DownstreamThingDeviceInfoUpdateReply(rsp *Response) error {
	log.Println("DownstreamThingDeviceInfoUpdateReply")
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

	uri := URIService(URISysPrefix, URIThingDeviceInfoDelete, node.ProductKey, node.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDeviceInfoDelete, params)
}

func DownstreamThingDeviceInfoDeleteReply(rsp *Response) error {
	log.Println("DownstreamThingDeviceInfoDeleteReply")
	return nil
}

// UpstreamThingDsltemplateGet 获取
func (sf *Manager) UpstreamThingDsltemplateGet() error {
	uri := URIService(URISysPrefix, URIThingDslTemplateGet, sf.ProductKey, sf.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDslTemplateGet, "{}")
}

func DownstreamThingDsltemplateGetReply(rsp *Response) error {
	log.Println("DownstreamThingDsltemplateGetReply")
	return nil
}

// UpstreamThingDynamictslGet 获取
func (sf *Manager) UpstreamThingDynamictslGet() error {
	uri := URIService(URISysPrefix, URIThingDynamicTslGet, sf.ProductKey, sf.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDynamicTslGet, `{"nodes\":["type","identifier"],"addDefault":false}`)
}

func DownstreamThingDynamictslGetReply(rsp *Response) error {
	log.Println("DownstreamThingDynamictslGetReply")
	return nil
}

// UpstreamExtNtpRequest ntp请求
// 发送一条Qos = 0的消息,并带上设备当前的时间戳,平台将回复 设备的发送时间,平台的接收时间, 平台的发送时间.
// 设备计算当前精确时间 = (平台接收时间 + 平台发送时间 + 设备接收时间 - 设备发送时间) / 2
func (sf *Manager) UpstreamExtNtpRequest() error {
	return sf.Publish(URIService(URIExtNtpPrefix, URINtpRequest, sf.ProductKey, sf.DeviceName),
		0, fmt.Sprintf(`{"deviceSendTime":"%d"}`, time.Now().Unix()))
}

// NtpResponse ntp回复payload
type NtpResponse struct {
	DeviceSendTime int `json:"deviceSendTime,string"`
	ServerRecvTime int `json:"serverRecvTime,string"`
	ServerSendTime int `json:"serverSendTime,string"`
}

//
func DownstreamExtNtpResponse(rsp *NtpResponse) error {
	return nil
}

// deprecated
func DownstreamThingServicePropertyGet(productKey, deviceName string, payload []byte) error {
	return nil
}

func DownstreamThingServiceRequest(productKey, deviceName, srvID string, payload []byte) error {
	return nil
}

func DownstreamThingServicePropertySet(payload []byte) error {
	return nil
}

func DownstreamExtErrorResponse(rsp *Response) error {
	return nil
}

func DownstreamThingModelDownRaw(productKey, deviceName string, payload []byte) error {
	// hex 2 string
	return nil
}
