package model

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/thinkgos/aliIOT/infra"
)

type DevNopUserProc struct{}

func (DevNopUserProc) DownstreamThingModelUpRawReply(m *Manager, productKey, deviceName string, payload []byte) error {
	log.Println("DownstreamThingModelUpRawReply")
	return nil
}

func (DevNopUserProc) DownstreamThingEventPropertyPostReply(rsp *Response) error {
	log.Println("DownstreamThingEventPropertyPostReply")
	return nil
}

func (DevNopUserProc) DownstreamThingEventPostReply(eventID string, rsp *Response) error {
	log.Println("DownstreamThingEventPostReply")
	return nil
}

func (DevNopUserProc) DownstreamThingDeviceInfoUpdateReply(rsp *Response) error {
	log.Println("DownstreamThingDeviceInfoUpdateReply")
	return nil
}
func (DevNopUserProc) DownstreamThingDeviceInfoDeleteReply(rsp *Response) error {
	log.Println("DownstreamThingDeviceInfoDeleteReply")
	return nil
}

func (DevNopUserProc) DownstreamThingDesiredPropertyGetReply(rsp *Response) error {
	log.Println("DownstreamThingDesiredPropertyGetReply")
	return nil
}

func (DevNopUserProc) DownstreamThingDesiredPropertyDeleteReply(rsp *Response) error {
	log.Println("DownstreamThingDesiredPropertyDeleteReply")
	return nil
}
func (DevNopUserProc) DownstreamThingDsltemplateGetReply(rsp *Response) error {
	log.Println("DownstreamThingDsltemplateGetReply")
	return nil
}

func (DevNopUserProc) DownstreamThingDynamictslGetReply(rsp *Response) error {
	log.Println("DownstreamThingDynamictslGetReply")
	return nil
}
func (DevNopUserProc) DownstreamExtNtpResponse(rsp *NtpResponsePayload) error {
	log.Println("DownstreamExtNtpResponse")
	return nil
}

func (DevNopUserProc) DownstreamThingConfigGetReply(rsp *ConfigGetResponse) error {
	log.Println("DownstreamThingConfigGetReply")
	return nil
}

func (DevNopUserProc) DownstreamThingConfigPush(rsp *ConfigPushRequest) error {
	log.Println("DownstreamThingConfigPush")
	return nil
}

func (DevNopUserProc) DownstreamExtErrorResponse(rsp *Response) error {
	log.Println("DownstreamExtErrorResponse")
	return nil
}

func (DevNopUserProc) DownstreamThingModelDownRaw(m *Manager, productKey, deviceName string, payload []byte) error {
	log.Println("DownstreamThingModelDownRaw")
	return nil
}

// DownstreamThingServicePropertySet 设置设备属性
func (DevNopUserProc) DownstreamThingServicePropertySet(m *Manager, topic string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}
	return m.SendResponse(URIServiceReplyWithRequestURI(topic), rsp.ID, infra.CodeSuccess, "{}")
}

// DownstreamThingServiceRequest 设备服务调用请求
func (DevNopUserProc) DownstreamThingServiceRequest(m *Manager, productKey, deviceName, srvID string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return nil
	}

	log.Println("DownstreamThingServiceRequest")
	return m.SendResponse(m.URIService(URISysPrefix, fmt.Sprintf(URIThingServiceRequest, srvID), productKey, deviceName),
		rsp.ID, infra.CodeSuccess, "{}")
}

func (DevNopUserProc) DownStreamRRPCRequest(m *Manager, productKey, deviceName, messageID string, payload []byte) error {
	log.Println("DownStreamRRPCRequest")
	return m.Publish(m.URIService(URISysPrefix, URIRRPCResponse, productKey, deviceName, messageID),
		0, "default system RRPC implementation")
}

func (DevNopUserProc) DownStreamExtRRPCRequest(m *Manager, topic string, payload []byte) error {
	log.Println("DownStreamRRPCRequest")
	return m.Publish(topic, 0, "default ext RRPC implementation")
}
