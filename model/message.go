// Package model 实现阿里去物模型
package model

import (
	"fmt"
	"log"
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

// UpstreamThingDsltemplateGet 获取
func (sf *Manager) UpstreamThingDsltemplateGet() error {
	uri := URIService(URISysPrefix, URIThingDslTemplateGet, sf.ProductKey, sf.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDslTemplateGet, "{}")
}

// UpstreamThingDynamictslGet 获取
func (sf *Manager) UpstreamThingDynamictslGet() error {
	uri := URIService(URISysPrefix, URIThingDynamicTslGet, sf.ProductKey, sf.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDynamicTslGet, `{"nodes\":["type","identifier"],"addDefault":false}`)
}

// UpstreamExtNtpRequest ntp请求
func (sf *Manager) UpstreamExtNtpRequest() error {
	uri := URIService(URIExtNtpPrefix, URINtpRequest, sf.ProductKey, sf.DeviceName)
	return sf.Publish(uri, 1, `{"deviceSendTime":"1234"}`)
}

func ThingModelDownRaw(productKey, deviceName string, payload []byte) error {
	// hex 2 string
	return nil
}
func ThingServicePropertySet(payload []byte) error {
	return nil
}

// deprecated
func ThingServicePropertyGet(productKey, deviceName string, payload []byte) error {
	return nil
}

func ThingServiceRequest(productKey, deviceName, srvID string, payload []byte) error {
	return nil
}

func ThingModelUpRawReply(productKey, deviceName string, payload []byte) error {
	// hex 2 string
	return nil
}

func ThingEventPropertyPostReply(rsp *Response) error {
	log.Println("ThingEventPropertyPostReply")
	return nil
}

func ThingEventPostReply(eventID string, rsp *Response) error {
	log.Println("ThingEventPostReply")
	return nil
}

func ThingDeviceInfoUpdateReply(rsp *Response) error {
	log.Println("ThingDeviceInfoUpdateReply")
	return nil
}
func ThingDeviceInfoDeleteReply(rsp *Response) error {
	log.Println("ThingDeviceInfoDeleteReply")
	return nil
}

func ThingDsltemplateGetReply(rsp *Response) error {
	log.Println("ThingDsltemplateGetReply")
	return nil
}

func ThingDynamictslGetReply(rsp *Response) error {
	log.Println("ThingDynamictslGetReply")
	return nil
}
func ExtNtpResponse(payload []byte) error {
	return nil
}

func ExtErrorResponse(rsp *Response) error {
	return nil
}

func ExtExtSubDevCombineLoginReply(rsp *Response) error {
	return nil
}

func ExtExtSubDevCombineLogoutReply(rsp *Response) error {
	return nil
}
func ExtExtSubDevRegisterReply(rsp *Response) error {
	return nil
}
