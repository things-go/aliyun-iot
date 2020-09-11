package dm

import (
	"encoding/json"
	"log"
)

// NopEvt 实现EventProc接口的空实现
type NopEvt struct{}

// 确保 NopEvt 实现 EventProc 接口
var _ EventProc = (*NopEvt)(nil)

// EvtThingModelUpRawReply see interface EventProc
func (NopEvt) EvtThingModelUpRawReply(*Client, string, string, []byte) error {
	return nil
}

// EvtThingEventPropertyPostReply see interface EventProc
func (NopEvt) EvtThingEventPropertyPostReply(*Client, error, string, string) error {
	return nil
}

// EvtThingEventPostReply see interface EventProc
func (NopEvt) EvtThingEventPostReply(*Client, error, string, string, string) error {
	return nil
}

// EvtThingEventPropertyPackPostReply see interface EventProc
func (NopEvt) EvtThingEventPropertyPackPostReply(*Client, error, string, string) error {
	return nil
}

// EvtThingDeviceInfoUpdateReply see interface EventProc
func (NopEvt) EvtThingDeviceInfoUpdateReply(*Client, error, string, string) error {
	return nil
}

// EvtThingDeviceInfoDeleteReply see interface EventProc
func (NopEvt) EvtThingDeviceInfoDeleteReply(*Client, error, string, string) error {
	return nil
}

// EvtThingDesiredPropertyGetReply see interface EventProc
func (NopEvt) EvtThingDesiredPropertyGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// EvtThingDesiredPropertyDeleteReply see interface EventProc
func (NopEvt) EvtThingDesiredPropertyDeleteReply(*Client, error, string, string) error {
	return nil
}

// EvtThingDsltemplateGetReply see interface EventProc
func (NopEvt) EvtThingDsltemplateGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// EvtThingDynamictslGetReply see interface EventProc
func (NopEvt) EvtThingDynamictslGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// EvtExtNtpResponse see interface EventProc
func (NopEvt) EvtExtNtpResponse(*Client, string, string, NtpResponse) error {
	return nil
}

// EvtThingConfigGetReply see interface EventProc
func (NopEvt) EvtThingConfigGetReply(*Client, error, string, string, ConfigParamsData) error {
	return nil
}

// EvtThingModelDownRaw see interface EventProc
func (NopEvt) EvtThingModelDownRaw(*Client, string, string, []byte) error {
	return nil
}

// EvtThingConfigPush see interface EventProc
func (NopEvt) EvtThingConfigPush(*Client, string, string, ConfigParamsData) error {
	return nil
}

// EvtThingServicePropertySet see interface EventProc
func (NopEvt) EvtThingServicePropertySet(*Client, string, string, []byte) error {
	log.Println("EvtThingServicePropertySet is not implementation")
	return nil
}

// EvtThingServiceRequest see interface EventProc
func (NopEvt) EvtThingServiceRequest(*Client, string, string, string, []byte) error {
	log.Println("EvtThingServiceRequest is not implementation")
	return nil
}

// EvtRRPCRequest see interface EventProc
func (NopEvt) EvtRRPCRequest(*Client, string, string, string, []byte) error {
	log.Println("EvtRRPCRequest is not implementation")
	return nil
}

// EvtExtRRPCRequest see interface EventProc
func (NopEvt) EvtExtRRPCRequest(*Client, string, string, []byte) error {
	log.Println("EvtRRPCRequest is not implementation")
	return nil
}

// EvtRequestWaitResponseTimeout  see interface EventProc
func (NopEvt) EvtRequestWaitResponseTimeout(_ *Client, msgType MsgType, devID int) error {
	log.Printf("%d - %s request wait response timeout", devID, msgType.String())
	return nil
}

/******************************** event gateway proc ************************************************************/

// NopGwEvt 实现EventGwProc接口的空实现
type NopGwEvt struct{}

// 确保 NopGwEvt 实现 EventGwProc 接口
var _ EventGwProc = (*NopGwEvt)(nil)

// EvtExtErrorResponse see interface EventGwProc
func (NopGwEvt) EvtExtErrorResponse(*Client, error, string, string) error {
	return nil
}

// EvtThingTopoGetReply see interface EventGwProc
func (NopGwEvt) EvtThingTopoGetReply(*Client, error, []GwTopoGetData) error {
	return nil
}

// EvtThingListFoundReply see interface EventGwProc
func (NopGwEvt) EvtThingListFoundReply(*Client, error) error {
	return nil
}

// EvtThingTopoAddNotify see interface EventGwProc
func (NopGwEvt) EvtThingTopoAddNotify(*Client, []GwTopoAddNotifyParams) error {
	return nil
}

// EvtThingTopoChange see interface EventGwProc
func (NopGwEvt) EvtThingTopoChange(*Client, GwTopoChangeParams) error {
	return nil
}

// EvtThingDisable see interface EventGwProc
func (NopGwEvt) EvtThingDisable(*Client, string, string) error {
	return nil
}

// EvtThingEnable see interface EventGwProc
func (NopGwEvt) EvtThingEnable(*Client, string, string) error {
	return nil
}

// EvtThingDelete see interface EventGwProc
func (NopGwEvt) EvtThingDelete(*Client, string, string) error {
	return nil
}
