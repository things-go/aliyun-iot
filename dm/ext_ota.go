package dm

import (
	"encoding/json"
)

// OTARequest OTA请求体
type OTARequest struct {
	ID     uint        `json:"id,string"`
	Params interface{} `json:"params"`
}

// OTAFirmwareVersionParams OTA固件参数域
type OTAFirmwareVersionParams struct {
	Version string `json:"version"`
}

// UpstreamOATFirmwareVersion 上报固件版本
func (sf *Client) UpstreamOATFirmwareVersion(devID int, params interface{}) error {
	if !sf.hasOTA {
		return ErrNotSupportFeature
	}
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNode(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	req, err := json.Marshal(OTARequest{id, params})
	if err != nil {
		return err
	}
	uri := sf.URIService(URIOtaDeviceInformPrefix, "", node.ProductKey(), node.DeviceName())
	err = sf.Publish(uri, 1, req)
	if err != nil {
		return err
	}

	// sf.Insert(id, devID, MsgTypeReportFirmwareVersion)
	sf.log.Debugf("upstream version <OTA>: inform,@%d", id)
	return nil
}

// OTA下载进度比
const (
	OTAProgressStepUpgradeFailed  = -1
	OTAProgressStepDownloadFailed = -2
	OTAProgressStepVerifyFailed   = -3
	OTAProgressStepProgramFailed  = -4
)

// OTAProgressParams 下载过程上报参数域
type OTAProgressParams struct {
	Step int    `json:"step,string"`
	Desc string `json:"desc"`
}

func (sf *Client) upstreamOTAProgress(devID int, params interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNode(devID)
	if err != nil {
		return err
	}

	id := sf.RequestID()
	req, err := json.Marshal(OTARequest{id, params})
	if err != nil {
		return err
	}
	err = sf.Publish(sf.URIService(URIOtaDeviceProcessPrefix, "", node.ProductKey(), node.DeviceName()),
		1, req)
	if err != nil {
		return err
	}

	// sf.Insert(id, devID, MsgTypeReportFirmwareVersion)
	sf.log.Debugf("upstream step <OTA>: progress,@%d", id)
	return nil
}

// OTAUpgradeData OTA upgrade 数据域
type OTAUpgradeData struct {
	Size    int    `json:"size"`
	Version string `json:"version"`
	URL     string `json:"url"`
	MD5     string `json:"md5"`
}

// OTAUpgradeRequest OTA upgrade 请求
type OTAUpgradeRequest struct {
	Code    int            `json:"code,string"`
	Data    OTAUpgradeData `json:"data"`
	ID      int            `json:"id"`
	Message string         `json:"message"`
}
