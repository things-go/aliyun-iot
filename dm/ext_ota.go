package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	uri "github.com/thinkgos/aliyun-iot/uri"
)

// OtaRequest OTA请求体
type OtaRequest struct {
	ID     uint        `json:"id,string"`
	Params interface{} `json:"params"`
}

// OtaInformParams OTA固件参数域
type OtaInformParams struct {
	Version string `json:"version"`
	// 上报默认（default）模块的固件版本号时,可以不上报module参数.
	// 设备的默认（default）模块的固件版本号等同于整个设备的固件版本号.
	Module string `json:"module"`
}

// OtaInform 上报固件版本
// request：/ota/device/inform/${YourProductKey}/${YourDeviceName}。
func (sf *Client) OtaInform(pk, dn string, params OtaInformParams) error {
	if !sf.hasOTA {
		return ErrNotSupportFeature
	}
	id := sf.nextRequestID()
	req, err := json.Marshal(OtaRequest{id, params})
	if err != nil {
		return err
	}
	_uri := uri.URI(uri.OtaDeviceInformPrefix, "", pk, dn)
	sf.log.Debugf("ota.device.inform @%d", id)
	return sf.Publish(_uri, 1, req)
}

// OTA下载进度比
const (
	OtaProgressStepUpgradeFailed  = -1
	OtaProgressStepDownloadFailed = -2
	OtaProgressStepVerifyFailed   = -3
	OtaProgressStepProgramFailed  = -4
)

// OtaProgressParams 下载过程上报参数域
type OtaProgressParams struct {
	Step   int    `json:"step,string"` // 固件升级进度信息 [1，100] 之间的数字：表示升级进度百分比,其它见上
	Desc   string `json:"desc"`
	Module string `json:"module"`
}

// OtaProgress 固件升级过程中，设备可以通过这个Topic上报固件升级的进度百分比
// request：/ota/device/progress/${YourProductKey}/${YourDeviceName
func (sf *Client) OtaProgress(pk, dn string, params OtaProgressParams) error {
	if !sf.hasOTA {
		return ErrNotSupportFeature
	}
	id := sf.nextRequestID()
	req, err := json.Marshal(OtaRequest{id, params})
	if err != nil {
		return err
	}
	sf.log.Debugf("ota.device.process @%d", id)
	_uri := uri.URI(uri.OtaDeviceProcessPrefix, "", pk, dn)
	return sf.Publish(_uri, 1, req)
}

// OtaFirmwareParam 请求固件信息参数域
type OtaFirmwareParam struct {
	Module string `json:"module"`
}

// OtaFirmwareData 请求固件信息回复数据域
type OtaFirmwareData struct {
	Size       int64  `json:"size"`
	Sign       string `json:"sign"`
	Version    string `json:"version"`
	IsDiff     int    `json:"isDiff"`
	URL        string `json:"url"`
	SignMethod string `json:"signMethod"`
	MD5        string `json:"md5"`
	Module     string `json:"module"`
}

// OtaFirmwareResponse ota firmware response
type OtaFirmwareResponse struct {
	ID      uint            `json:"id,string"`
	Code    int             `json:"code"`
	Data    OtaFirmwareData `json:"data"`
	Message string          `json:"message"`
}

// ThingOtaFirmwareGet 请求固件信息
// module: 不指定则表示请求默认（default）模块的固件信息
// request： /sys/{productKey}/{deviceName}/thing/ota/firmware/get
// response：/sys/{productKey}/{deviceName}/thing/ota/firmware/get_reply
func (sf *Client) ThingOtaFirmwareGet(pk, dn string, param OtaFirmwareParam) (*Token, error) {
	if !sf.hasOTA {
		return nil, ErrNotSupportFeature
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingOtaFirmwareGet, pk, dn)
	return sf.SendRequest(_uri, infra.MethodOtaFirmwareGet, param)
}

// ProcThingOtaFirmwareGetReply 处理请求固件信息应答
// request：  /sys/{productKey}/{deviceName}/thing/ota/firmware/get
// response： /sys/{productKey}/{deviceName}/thing/ota/firmware/get_reply
// subscribe：/sys/{productKey}/{deviceName}/thing/ota/firmware/get_reply
func ProcThingOtaFirmwareGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	rsp := &OtaFirmwareResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.log.Debugf("thing.ota.firmware.get.reply @%d", rsp.ID)
	c.signalPending(Message{rsp.ID, rsp.Data, err})
	pk, dn := uris[1], uris[2]
	return c.cb.ThingOtaFirmwareGetReply(c, pk, dn, rsp.Data)
}

// ProcOtaUpgrade 处理物联网平台推送固件信息
// request：  /ota/device/upgrade/${YourProductKey}/${YourDeviceName}
// subscribe：/ota/device/upgrade/${YourProductKey}/${YourDeviceName}
func ProcOtaUpgrade(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 5 {
		return ErrInvalidURI
	}
	rsp := &OtaFirmwareResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	c.log.Debugf("thing.device.upgrade")
	pk, dn := uris[3], uris[4]
	return c.cb.OtaUpgrade(c, pk, dn, rsp)
}
