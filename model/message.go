package model

import (
	"encoding/json"
	"fmt"
)

const (
	version = "1.0"
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
	ID      int         `json:"id,string"`
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
	ID                     uint64
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
	RequestID              int
	RePortID               int
	ProductKey, DeviceName string
}

func (this *Manager) SendRequest(uriService string, id int, method string, params interface{}) error {
	req := Request{
		ID:      id,
		Version: version,
		Params:  params,
		Method:  method,
	}
	out, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	return this.Publish(uriService, out)
}

func (this *Manager) UpstreamThingModelUpRaw(payload interface{}) error {
	uri := UriService(UriSysPrefix, UriThingModelUpRaw, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID, methodUpRaw, payload)
}

func (this *Manager) UpstreamThingPropertyPost(params interface{}) error {
	uri := UriService(UriSysPrefix, UriThingEventPropertyPost, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID, methodPropertyPost, params)
}
func (this *Manager) UpstreamThingEventPost(EventID string, params interface{}) error {
	uri := UriService(UriSysPrefix, fmt.Sprintf(UriThingEventPost, EventID), this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID, fmt.Sprintf(methodEventPostFormat, EventID), params)
}
func (this *Manager) UpstreamThingDeviceInfoUpdate(params interface{}) error {
	uri := UriService(UriSysPrefix, UriThingDeviceInfoUpdate, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID, methodDeviceInfoUpdate, params)
}

func (this *Manager) UpstreamThingDeviceInfoDelete(params interface{}) error {
	uri := UriService(UriSysPrefix, UriThingDeviceInfoDelete, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID, methodDeviceInfoDelete, params)
}

func (this *Manager) UpstreamThingDsltemplateGet() error {
	uri := UriService(UriSysPrefix, UriThingDslTemplateGet, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID, methodDslTemplateGet, "{}")
}

func (this *Manager) UpstreamThingDynamictslGet() error {
	uri := UriService(UriSysPrefix, UriThingDynamicTslGet, this.ProductKey, this.DeviceName)
	return this.SendRequest(uri, this.RequestID, methodDynamicTslGet, `{"nodes\":["type","identifier"],"addDefault":false}`)
}

func (this *Manager) UpstreamThingNtpRequest() error {
	uri := UriService(UriExtNtpPrefix, UriNtpRequest, this.ProductKey, this.DeviceName)
	return this.Publish(uri, `{"deviceSendTime":"1234"}`)
}
