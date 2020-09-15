package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// GwSubRegisterData 子设备注册应答数据域
type GwSubRegisterData struct {
	IotID        int64  `json:"iotId,string"`
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
}

// GwSubRegisterResponse 子设备注册应答
type GwSubRegisterResponse struct {
	ID      uint                `json:"id,string"`
	Code    int                 `json:"code"`
	Data    []GwSubRegisterData `json:"data"`
	Message string              `json:"message,omitempty"`
}

// ThingGwSubRegister 子设备动态注册
// 以通过上行请求为子设备发起动态注册，返回成功注册的子设备的设备证书
// request:   /sys/{productKey}/{deviceName}/thing/sub/register
// response:  /sys/{productKey}/{deviceName}/thing/sub/register_reply
func (sf *Client) ThingGwSubRegister(pk, dn string) (*Token, error) {
	_uri := sf.GatewayURI(uri.SysPrefix, uri.ThingSubRegister)
	return sf.SendRequest(_uri, infra.MethodSubDevRegister, []infra.MetaPair{
		{ProductKey: pk, DeviceName: dn},
	})
}

// ProcThingGwSubRegisterReply 子设备动态注册处理
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/sub/register
// response:  /sys/{productKey}/{deviceName}/thing/sub/register_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/sub/register_replyc.uriOffset+
func ProcThingGwSubRegisterReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := &GwSubRegisterResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else {
		for _, v := range rsp.Data {
			_ = c.SetDeviceSecret(v.ProductKey, v.DeviceName, v.DeviceSecret)
			_ = c.SetDeviceStatus(v.ProductKey, v.DeviceName, DevStatusRegistered)
		}
	}
	c.signalPending(Message{rsp.ID, nil, err})
	c.log.Debugf("downstream GW thing <sub>: register reply @%d", rsp.ID)
	return nil
}
