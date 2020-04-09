package dm

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/thinkgos/aiot/infra"
)

// ipc 事件类型
type ipcEvtType byte

// ipc 事件类型定义
const (
	// 上行应答
	ipcEvtUpRawReply ipcEvtType = iota
	ipcEvtEventPropertyPostReply
	ipcEvtEventPostReply
	ipcEvtPropertyPackPostReply
	ipcEvtDeviceInfoUpdateReply
	ipcEvtDeviceInfoDeleteReply
	ipcEvtDesiredPropertyGetReply
	ipcEvtDesiredPropertyDeleteReply
	ipcEvtDsltemplateGetReply
	ipcEvtDynamictslGetReply
	ipcEvtExtNtpResponse
	ipcEvtConfigGetReply
	ipcEvtErrorResponse

	// 下行
	ipcEvtDownRaw
	ipcEvtConfigPush
	ipcEvtServicePropertySet
	ipcEvtServiceRequest
	ipcEvtRRPCRequest
	ipcEvtExtRRPCRequest

	// gateway
	// 上行应答
	ipcEvtTopoGetReply
	ipcEvtListFoundReply
	// 下行
	ipcEvtTopoAddNotify
	ipcTopoChange
	ipcThingDisable
	ipcThingEnable
	ipcThingDelete

	// 内部,请求超时
	ipcEvtRequestTimeout
)

type ipcMessage struct {
	evt        ipcEvtType
	err        error
	productKey string
	deviceName string
	extend     string
	payload    interface{}
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
	// 下行应答
	case ipcEvtUpRawReply:
		return sf.eventProc.EvtThingModelUpRawReply(sf, msg.productKey, msg.deviceName, msg.payload.([]byte))
	case ipcEvtEventPropertyPostReply:
		return sf.eventProc.EvtThingEventPropertyPostReply(sf, msg.err, msg.productKey, msg.deviceName)
	case ipcEvtEventPostReply:
		return sf.eventProc.EvtThingEventPostReply(sf, msg.err, msg.extend, msg.productKey, msg.deviceName)
	case ipcEvtPropertyPackPostReply:
		return sf.eventProc.EvtThingEventPropertyPackPostReply(sf, msg.err, msg.productKey, msg.deviceName)
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
	case ipcEvtErrorResponse:
		err := msg.err.(*infra.CodeError)
		data := msg.payload.(ExtErrorData)
		sf.debug("ext evt error response, %+v", err)

		code := err.Code()
		if code == infra.CodeSubDevSessionError {
			node, err := sf.SearchNodeByPkDn(data.ProductKey, data.DeviceName)
			if err != nil {
				return err
			}
			_, _ = sf.upstreamExtGwSubDevCombineLogin(node.ID())
		}
		return sf.eventGwProc.EvtExtErrorResponse(sf, err, data.ProductKey, data.DeviceName)

		// 下行
	case ipcEvtDownRaw:
		return sf.eventProc.EvtThingModelDownRaw(sf, msg.productKey, msg.deviceName, msg.payload.([]byte))
	case ipcEvtConfigPush:
		return sf.eventProc.EvtThingConfigPush(sf, msg.productKey, msg.deviceName, msg.payload.(ConfigParamsAndData))
	case ipcEvtServicePropertySet:
		return sf.eventProc.EvtThingServicePropertySet(sf, msg.productKey, msg.deviceName, msg.payload.([]byte))
	case ipcEvtServiceRequest:
		return sf.eventProc.EvtThingServiceRequest(sf, msg.extend, msg.productKey, msg.deviceName, msg.payload.([]byte))
	case ipcEvtRRPCRequest:
		return sf.eventProc.EvtRRPCRequest(sf, msg.extend, msg.productKey, msg.deviceName, msg.payload.([]byte))
	case ipcEvtExtRRPCRequest:
		ext := strings.SplitN(msg.extend, SEP, 2)
		return sf.eventProc.EvtExtRRPCRequest(sf, ext[0], ext[1], msg.payload.([]byte))

		// 上行应答
	case ipcEvtTopoGetReply:
		return sf.eventGwProc.EvtThingTopoGetReply(sf, msg.err, msg.payload.([]GwTopoGetData))
	case ipcEvtListFoundReply:
		return sf.eventGwProc.EvtThingListFoundReply(sf, msg.err)
		// 下行
	case ipcEvtTopoAddNotify:
		return sf.eventGwProc.EvtThingTopoAddNotify(sf, msg.payload.([]GwTopoAddNotifyParams))
	case ipcTopoChange:
		return sf.eventGwProc.EvtThingTopoChange(sf, msg.payload.(GwTopoChangeParams))
	case ipcThingDisable:
		return sf.eventGwProc.EvtThingDisable(sf, msg.productKey, msg.deviceName)
	case ipcThingEnable:
		return sf.eventGwProc.EvtThingEnable(sf, msg.productKey, msg.deviceName)
	case ipcThingDelete:
		return sf.eventGwProc.EvtThingDelete(sf, msg.productKey, msg.deviceName)

	case ipcEvtRequestTimeout:
		devID, _ := strconv.Atoi(msg.extend)
		return sf.eventProc.EvtRequestWaitResponseTimeout(sf, msg.payload.(MsgType), devID)
	}

	return errors.New("not support ipc event type")
}
