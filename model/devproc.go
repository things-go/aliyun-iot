package model

import "log"

type DevNopUserProc struct{}

func (DevNopUserProc) DownstreamThingModelUpRawReply(productKey, deviceName string, payload []byte) error {
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

func (DevNopUserProc) DownstreamThingDesiredPropertyGetReply(rsp *Response) error {
	log.Println("DownstreamThingDesiredPropertyGetReply")
	return nil
}

func (DevNopUserProc) DownstreamThingDesiredPropertyDeleteReply(rsp *Response) error {
	log.Println("DownstreamThingDesiredPropertyDeleteReply")
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

func (DevNopUserProc) DownstreamThingDsltemplateGetReply(rsp *Response) error {
	log.Println("DownstreamThingDsltemplateGetReply")
	return nil
}

func (DevNopUserProc) DownstreamThingDynamictslGetReply(rsp *Response) error {
	log.Println("DownstreamThingDynamictslGetReply")
	return nil
}
func (DevNopUserProc) DownstreamExtNtpResponse(rsp *NtpResponse) error {
	return nil
}

func (DevNopUserProc) DownstreamExtErrorResponse(rsp *Response) error {
	return nil
}

func (DevNopUserProc) DownstreamThingModelDownRaw(productKey, deviceName string, payload []byte) error {
	// hex 2 string
	return nil
}

// deprecated
func (DevNopUserProc) DownstreamThingServicePropertyGet(productKey, deviceName string, payload []byte) error {
	return nil
}

func (DevNopUserProc) DownstreamThingServiceRequest(productKey, deviceName, srvID string, payload []byte) error {
	return nil
}

func (DevNopUserProc) DownstreamThingServicePropertySet(payload []byte) error {
	return nil
}
