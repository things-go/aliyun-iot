// Package model 实现阿里去物模型
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
)

// 平台通信版本
const (
	Version = "1.0"
)

// method 定义
const (
	methodPropertyPost     = "thing.event.property.post"
	methodDeviceInfoUpdate = "thing.deviceinfo.update"
	methodDeviceInfoDelete = "thing.deviceinfo.delete"
	methodUpRaw            = "thing.model.up_raw"
	methodEventPostFormat  = "thing.event.%s.post"
	methodDslTemplateGet   = "thing.dsltemplate.get"
	methodDynamicTslGet    = "thing.dynamicTsl.get"
)

// Request 请求
type Request struct {
	ID      int32       `json:"id,string"`
	Version string      `json:"version"`
	Params  interface{} `json:"params"`
	Method  string      `json:"method"`
}

// Response 应答
type Response struct {
	ID      int32       `json:"id,string"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// ResponseInfo 回复的消息
type ResponseInfo struct {
	ID         int
	Qos        int
	URIPrefix  string
	URIName    string
	ProductKey string
	DeviceName string
}

// Conn conn接口
type Conn interface {
	// Publish will publish a message with the specified QoS and content
	Publish(topic string, payload interface{}) error
}

// Manager 管理
type Manager struct {
	Conn
	requestID  int32
	reportID   int32
	ProductKey string
	DeviceName string
	uriOffset  int
}

// RequestID 获得下一个requestID
func (sf *Manager) RequestID() int32 {
	return atomic.AddInt32(&sf.requestID, 1)
}

// ReportID 获得下一个reportID
func (sf *Manager) ReportID() int32 {
	return atomic.AddInt32(&sf.reportID, 1)
}

// SendRequest 发送请求
// uriService 唯一定位服务器或(topic)
// requestID: 请求ID
// method: 方法
// params: 消息体
// API内部已实现json序列化
func (sf *Manager) SendRequest(uriService string, requestID int32, method string, params interface{}) error {
	out, err := json.Marshal(&Request{requestID, Version, params, method})
	if err != nil {
		return err
	}
	return sf.Publish(uriService, out)
}

func (sf *Manager) SendResponse(uriService string, reportID int32, code int, data interface{}) error {
	out, err := json.Marshal(&Response{reportID, code, data, ""})
	if err != nil {
		return err
	}
	return sf.Publish(uriService, out)
}

// UpstreamThingModelUpRaw 上传透传数据
func (sf *Manager) UpstreamThingModelUpRaw(payload interface{}) error {
	uri := URIService(URISysPrefix, URIThingModelUpRaw, sf.ProductKey, sf.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodUpRaw, payload)
}

// UpstreamThingPropertyPost 上传数属性数据
func (sf *Manager) UpstreamThingEventPropertyPost(params interface{}) error {
	uri := URIService(URISysPrefix, URIThingEventPropertyPost, sf.ProductKey, sf.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodPropertyPost, params)
}

// UpstreamThingEventPost 事件上传
func (sf *Manager) UpstreamThingEventPost(EventID string, params interface{}) error {
	uri := URIService(URISysPrefix, fmt.Sprintf(URIThingEventPost, EventID), sf.ProductKey, sf.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), fmt.Sprintf(methodEventPostFormat, EventID), params)
}

// UpstreamThingDeviceInfoUpdate 设备信息上传
func (sf *Manager) UpstreamThingDeviceInfoUpdate(params interface{}) error {
	uri := URIService(URISysPrefix, URIThingDeviceInfoUpdate, sf.ProductKey, sf.DeviceName)
	return sf.SendRequest(uri, sf.RequestID(), methodDeviceInfoUpdate, params)
}

// UpstreamThingDeviceInfoDelete 设备信息删除
func (sf *Manager) UpstreamThingDeviceInfoDelete(params interface{}) error {
	uri := URIService(URISysPrefix, URIThingDeviceInfoDelete, sf.ProductKey, sf.DeviceName)
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
	return sf.Publish(uri, `{"deviceSendTime":"1234"}`)
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
	return errors.New("got it reply")
}

func ThingEventPostReply(eventID string, rsp *Response) error {
	return nil
}

func ThingDeviceInfoUpdateReply(rsp *Response) error {
	return nil
}
func ThingDeviceInfoDeleteReply(rsp *Response) error {
	return nil
}

func ThingDsltemplateGetReply(rsp *Response) error {
	return nil
}

func ThingDynamictslGetReply(rsp *Response) error {
	return nil
}
func ExtNtpResponse(payload []byte) error {
	return nil
}

func ExtErrorResponse(rsp *Response) error {
	return nil
}
