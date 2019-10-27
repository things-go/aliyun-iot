package model

import "log"

type DevNopUserProc struct{}

func (DevNopUserProc) DownstreamThingModelUpRawReply(productKey, deviceName string, payload []byte) error {
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

func (DevNopUserProc) DownstreamThingModelDownRaw(productKey, deviceName string, payload []byte) error {
	log.Println("DownstreamThingModelDownRaw")
	return nil
}

// deprecated
func (DevNopUserProc) DownstreamThingServicePropertyGet(productKey, deviceName string, payload []byte) error {
	log.Println("DownstreamThingServicePropertyGet")
	return nil
}

// DownstreamThingServicePropertySet 设置设备属性
func (DevNopUserProc) DownstreamThingServicePropertySet(payload []byte) error {
	log.Println("DownstreamThingServicePropertySet")
	return nil
}

// DownstreamThingServiceRequest 设备服务调用请求
func (DevNopUserProc) DownstreamThingServiceRequest(productKey, deviceName, srvID string, payload []byte) error {
	log.Println("DownstreamThingServiceRequest")
	return nil
}
