package dm

import (
	"encoding/json"
)

// NopCb 实现Callback接口的空实现
type NopCb struct{}

// 确保 NopCb 实现 Callback 接口
var _ Callback = (*NopCb)(nil)

// ThingModelUpRawReply see interface Callback
func (NopCb) ThingModelUpRawReply(*Client, string, string, []byte) error { return nil }

// ThingEventPropertyPostReply see interface Callback
func (NopCb) ThingEventPropertyPostReply(*Client, error, string, string) error { return nil }

// ThingEventPostReply see interface Callback
func (NopCb) ThingEventPostReply(*Client, error, string, string, string) error { return nil }

// ThingEventPropertyPackPostReply see interface Callback
func (NopCb) ThingEventPropertyPackPostReply(*Client, error, string, string) error { return nil }

// ThingEventPropertyHistoryPostReply see interface Callback
func (NopCb) ThingEventPropertyHistoryPostReply(*Client, error, string, string) error { return nil }

// ThingDeviceInfoUpdateReply see interface Callback
func (NopCb) ThingDeviceInfoUpdateReply(*Client, error, string, string) error { return nil }

// ThingDeviceInfoDeleteReply see interface Callback
func (NopCb) ThingDeviceInfoDeleteReply(*Client, error, string, string) error { return nil }

// ThingDesiredPropertyGetReply see interface Callback
func (NopCb) ThingDesiredPropertyGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// ThingDesiredPropertyDeleteReply see interface Callback
func (NopCb) ThingDesiredPropertyDeleteReply(*Client, error, string, string) error { return nil }

// ThingDsltemplateGetReply see interface Callback
func (NopCb) ThingDsltemplateGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// ThingDynamictslGetReply see interface Callback
func (NopCb) ThingDynamictslGetReply(*Client, error, string, string, json.RawMessage) error {
	return nil
}

// ExtNtpResponse see interface Callback
func (NopCb) ExtNtpResponse(*Client, string, string, NtpResponse) error { return nil }

// ThingConfigGetReply see interface Callback
func (NopCb) ThingConfigGetReply(*Client, error, string, string, ConfigParamsData) error { return nil }

// ThingModelDownRaw see interface Callback
func (NopCb) ThingModelDownRaw(*Client, string, string, []byte) error { return nil }

// ThingConfigPush see interface Callback
func (NopCb) ThingConfigPush(*Client, string, string, ConfigParamsData) error { return nil }

// ThingServicePropertySet see interface Callback
func (NopCb) ThingServicePropertySet(*Client, string, string, []byte) error { return nil }

// ThingServiceRequest see interface Callback
func (NopCb) ThingServiceRequest(*Client, string, string, string, []byte) error { return nil }

// RRPCRequest see interface Callback
func (NopCb) RRPCRequest(*Client, string, string, string, []byte) error { return nil }

// ExtRRPCRequest see interface Callback
func (NopCb) ExtRRPCRequest(*Client, string, string, []byte) error { return nil }

/********************************  gateway callback ************************************************************/

// 确保 NopGwCb 实现 GwCallback 接口
var _ GwCallback = (*NopGwCb)(nil)

// NopGwCb 实现EventGwProc接口的空实现
type NopGwCb struct{}

// ExtErrorResponse see interface GwCallback
func (NopGwCb) ExtErrorResponse(*Client, error, string, string) error { return nil }

// ThingGwTopoGetReply see interface GwCallback
func (NopGwCb) ThingGwTopoGetReply(*Client, error, []GwTopoGetData) error { return nil }

// ThingGwListFoundReply see interface GwCallback
func (NopGwCb) ThingGwListFoundReply(*Client, error) error { return nil }

// ThingGwTopoAddNotify see interface GwCallback
func (NopGwCb) ThingGwTopoAddNotify(*Client, []GwTopoAddNotifyParams) error { return nil }

// ThingGwTopoChange see interface GwCallback
func (NopGwCb) ThingGwTopoChange(*Client, GwTopoChangeParams) error { return nil }

// ThingGwDisable see interface GwCallback
func (NopGwCb) ThingGwDisable(*Client, string, string) error { return nil }

// ThingGwEnable see interface GwCallback
func (NopGwCb) ThingGwEnable(*Client, string, string) error { return nil }

// ThingGwDelete see interface GwCallback
func (NopGwCb) ThingGwDelete(*Client, string, string) error { return nil }
