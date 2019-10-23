package model

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

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

type Request struct {
	ID      int32       `json:"id,string"`
	Version string      `json:"version"`
	Params  interface{} `json:"params"`
	Method  string      `json:"method"`
}

type Response struct {
	ID   string
	Code int
	Data interface{}
}

type ResponseInfo struct {
	ID                     int
	Qos                    int
	UriPrefix, UriName     string
	ProductKey, DeviceName string
}

type Conn interface {
	// Publish will publish a message with the specified QoS and content
	Publish(topic string, payload interface{}) error
}

type Manager struct {
	Conn
	requestID              int32
	reportID               int32
	ProductKey, DeviceName string
}

// RequestID 获得下一个requestID
func (this *Manager) RequestID() int32 {
	return atomic.AddInt32(&this.requestID, 1)
}

// ReportID 获得下一个reportID
func (this *Manager) ReportID() int32 {
	return atomic.AddInt32(&this.reportID, 1)
}

// 发送请求,
// uriService 唯一定位服务器或(topic)
// requestID: 请求ID
// method: 方法
// params: 消息体
// API内部已实现json序列化
func (this *Manager) SendRequest(uriService string, requestID int32, method string, params interface{}) error {
	out, err := json.Marshal(&Request{requestID, Version, params, method})
	if err != nil {
		return err
	}
	return this.Publish(uriService, out)
}

// UpstreamThingModelUpRaw 上传透传数据
func (this *Manager) UpstreamThingModelUpRaw(payload interface{}) error {
	uri := UriService(UriSysPrefix, UriThingModelUpRaw, this.ProductKey, this.DeviceName)

	return this.SendRequest(uri, this.RequestID(), methodUpRaw, payload)
}

// UpstreamThingPropertyPost 上传数属性数据
func (this *Manager) UpstreamThingPropertyPost(params interface{}) error {
	uri := UriService(UriSysPrefix, UriThingEventPropertyPost, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID(), methodPropertyPost, params)
}

// UpstreamThingEventPost 事件上传
func (this *Manager) UpstreamThingEventPost(EventID string, params interface{}) error {
	uri := UriService(UriSysPrefix, fmt.Sprintf(UriThingEventPost, EventID), this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID(), fmt.Sprintf(methodEventPostFormat, EventID), params)
}

// UpstreamThingDeviceInfoUpdate 设备信息上传
func (this *Manager) UpstreamThingDeviceInfoUpdate(params interface{}) error {
	uri := UriService(UriSysPrefix, UriThingDeviceInfoUpdate, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID(), methodDeviceInfoUpdate, params)
}

// UpstreamThingDeviceInfoDelete 设备信息删除
func (this *Manager) UpstreamThingDeviceInfoDelete(params interface{}) error {
	uri := UriService(UriSysPrefix, UriThingDeviceInfoDelete, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID(), methodDeviceInfoDelete, params)
}

// UpstreamThingDsltemplateGet
func (this *Manager) UpstreamThingDsltemplateGet() error {
	uri := UriService(UriSysPrefix, UriThingDslTemplateGet, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID(), methodDslTemplateGet, "{}")
}

// UpstreamThingDynamictslGet
func (this *Manager) UpstreamThingDynamictslGet() error {
	uri := UriService(UriSysPrefix, UriThingDynamicTslGet, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID(), methodDynamicTslGet, `{"nodes\":["type","identifier"],"addDefault":false}`)
}

// UpstreamThingNtpRequest
func (this *Manager) UpstreamThingNtpRequest() error {
	uri := UriService(UriExtNtpPrefix, UriNtpRequest, this.ProductKey, this.DeviceName)
	return this.Publish(uri, `{"deviceSendTime":"1234"}`)
}
