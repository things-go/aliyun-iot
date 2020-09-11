package dm

import (
	"encoding/json"
	"log"
)

// NopEvt 实现EventProc接口的空实现
type NopEvt struct{}

// 确保 NopEvt 实现 Callback 接口
var _ Callback = (*NopEvt)(nil)

// ThingModelUpRawReply see interface Callback
func (NopEvt) ThingModelUpRawReply(*Client, string, string, []byte) error {
	return nil
}

// ThingEventPropertyPostReply see interface Callback
func (NopEvt) ThingEventPropertyPostReply(*Client, error, string, string) error {
	return nil
}

// ThingEventPostReply see interface Callback
func (NopEvt) ThingEventPostReply(*Client, error, string, string, string) error {
	return nil
}

// ThingEventPropertyPackPostReply see interface Callback
func (NopEvt) ThingEventPropertyPackPostReply(*Client, error, string, string) error {
	return nil
}

// ThingDeviceInfoUpdateReply see interface Callback
func (NopEvt) ThingDeviceInfoUpdateReply(*Client, error, string, string) error {
	return nil
}

// ThingDeviceInfoDeleteReply see interface Callback
func (NopEvt) ThingDeviceInfoDeleteReply(*Client, error, string, string) error {
	return nil
}

// ThingDesiredPropertyGetReply see interface Callback
func (NopEvt) ThingDesiredPropertyGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// ThingDesiredPropertyDeleteReply see interface Callback
func (NopEvt) ThingDesiredPropertyDeleteReply(*Client, error, string, string) error {
	return nil
}

// ThingDsltemplateGetReply see interface Callback
func (NopEvt) ThingDsltemplateGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// ThingDynamictslGetReply see interface Callback
func (NopEvt) ThingDynamictslGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// ExtNtpResponse see interface Callback
func (NopEvt) ExtNtpResponse(*Client, string, string, NtpResponse) error {
	return nil
}

// ThingConfigGetReply see interface Callback
func (NopEvt) ThingConfigGetReply(*Client, error, string, string, ConfigParamsData) error {
	return nil
}

// ThingModelDownRaw see interface Callback
func (NopEvt) ThingModelDownRaw(*Client, string, string, []byte) error {
	return nil
}

// ThingConfigPush see interface Callback
func (NopEvt) ThingConfigPush(*Client, string, string, ConfigParamsData) error {
	return nil
}

// ThingServicePropertySet see interface Callback
func (NopEvt) ThingServicePropertySet(*Client, string, string, []byte) error {
	log.Println("ThingServicePropertySet is not implementation")
	return nil
}

// ThingServiceRequest see interface Callback
func (NopEvt) ThingServiceRequest(*Client, string, string, string, []byte) error {
	log.Println("ThingServiceRequest is not implementation")
	return nil
}

// RRPCRequest see interface Callback
func (NopEvt) RRPCRequest(*Client, string, string, string, []byte) error {
	log.Println("RRPCRequest is not implementation")
	return nil
}

// ExtRRPCRequest see interface Callback
func (NopEvt) ExtRRPCRequest(*Client, string, string, []byte) error {
	log.Println("RRPCRequest is not implementation")
	return nil
}

// EvtRequestWaitResponseTimeout  see interface Callback
func (NopEvt) EvtRequestWaitResponseTimeout(_ *Client, devID int) error {
	log.Printf("%d request wait response timeout", devID)
	return nil
}

/******************************** event gateway proc ************************************************************/

// NopGwEvt 实现EventGwProc接口的空实现
type NopGwEvt struct{}

// 确保 NopGwEvt 实现 GwCallback 接口
var _ GwCallback = (*NopGwEvt)(nil)

// ExtErrorResponse see interface GwCallback
func (NopGwEvt) ExtErrorResponse(*Client, error, string, string) error {
	return nil
}

// ThingGwTopoGetReply see interface GwCallback
func (NopGwEvt) ThingGwTopoGetReply(*Client, error, []GwTopoGetData) error {
	return nil
}

// ThingGwListFoundReply see interface GwCallback
func (NopGwEvt) ThingGwListFoundReply(*Client, error) error {
	return nil
}

// ThingGwTopoAddNotify see interface GwCallback
func (NopGwEvt) ThingGwTopoAddNotify(*Client, []GwTopoAddNotifyParams) error {
	return nil
}

// ThingGwTopoChange see interface GwCallback
func (NopGwEvt) ThingGwTopoChange(*Client, GwTopoChangeParams) error {
	return nil
}

// ThingGwDisable see interface GwCallback
func (NopGwEvt) ThingGwDisable(*Client, string, string) error {
	return nil
}

// ThingGwEnable see interface GwCallback
func (NopGwEvt) ThingGwEnable(*Client, string, string) error {
	return nil
}

// ThingGwDelete see interface GwCallback
func (NopGwEvt) ThingGwDelete(*Client, string, string) error {
	return nil
}
