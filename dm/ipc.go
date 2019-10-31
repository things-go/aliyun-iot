package dm

import (
	"encoding/json"
	"errors"
)

type ipcEvtType byte

const (
	ipcEvtRawReply ipcEvtType = iota
	ipcEvtEventPropertyPostReply
	ipcEvtEventPostReply
	ipcEvtDeviceInfoUpdateReply
	ipcEvtDeviceInfoDeleteReply
	ipcEvtDesiredPropertyGetReply
	ipcEvtDesiredPropertyDeleteReply
	ipcEvtDsltemplateGetReply
	ipcEvtDynamictslGetReply
	ipcEvtExtNtpResponse
	ipcEvtConfigGetReply
)

type ipcMessage struct {
	err        error
	evt        ipcEvtType
	productKey string
	deviceName string
	payload    interface{}
	ext        string
}

func (sf *Client) ipcSendMessage(msg *ipcMessage) error {
	select {
	case sf.ipc <- msg:
	default:
		return ErrIPCMessageBuffFull
	}
	return nil
}

func (sf *Client) ipcRunMessage() {
	var err error

	for msg := range sf.ipc {
		err = sf.ipcEventProc(msg)
		if err != nil {
			sf.error("ipc event process failed, %+v", err)
		}
	}
	//for {
	//	select {
	//	case msg := <-sf.ipc:
	//		err = sf.ipcEventProc(msg)
	//		if err != nil {
	//			sf.error("ipc event process failed, %+v", err)
	//		}
	//	}
	//}
}

func (sf *Client) ipcEventProc(msg *ipcMessage) error {
	defer func() {
		if err := recover(); err != nil {
			sf.critical("panic happen, %+v", err)
		}
	}()

	switch msg.evt {
	case ipcEvtRawReply:
		return sf.eventProc.EvtThingModelUpRawReply(sf, msg.productKey, msg.deviceName, msg.payload.([]byte))
	case ipcEvtEventPropertyPostReply:
		return sf.eventProc.EvtThingEventPropertyPostReply(sf, msg.err, msg.productKey, msg.deviceName)
	case ipcEvtEventPostReply:
		return sf.eventProc.EvtThingEventPostReply(sf, msg.err, msg.ext, msg.productKey, msg.deviceName)
	case ipcEvtDeviceInfoUpdateReply:
		return sf.eventProc.EvtThingDeviceInfoUpdateReply(sf, msg.err, msg.productKey, msg.deviceName)
	case ipcEvtDeviceInfoDeleteReply:
		return sf.eventProc.EvtThingDeviceInfoDeleteReply(sf, msg.err, msg.productKey, msg.deviceName)
	case ipcEvtDesiredPropertyGetReply:
		return sf.eventProc.EvtThingDesiredPropertyGetReply(sf, msg.err, msg.productKey, msg.deviceName, msg.payload.(json.RawMessage))
	case ipcEvtDesiredPropertyDeleteReply:
		return sf.eventProc.EvtThingDesiredPropertyDeleteReply(sf, msg.err, msg.productKey, msg.deviceName)
	case ipcEvtDsltemplateGetReply:
		return sf.eventProc.EvtThingDsltemplateGetReply(sf, msg.err, msg.productKey, msg.deviceName, msg.payload.(json.RawMessage))
	case ipcEvtDynamictslGetReply:
		return sf.eventProc.EvtThingDynamictslGetReply(sf, msg.err, msg.productKey, msg.deviceName, msg.payload.(json.RawMessage))
	case ipcEvtExtNtpResponse:
		return sf.eventProc.EvtExtNtpResponse(sf, msg.productKey, msg.deviceName, msg.payload.(NtpResponsePayload))
	case ipcEvtConfigGetReply:
		return sf.eventProc.EvtThingConfigGetReply(sf, msg.err, msg.productKey, msg.deviceName, msg.payload.(ConfigParamsAndData))
	}

	return errors.New("not support ipc event type")
}
