package aiot

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/89298.html?spm=a2c4g.11186623.6.703.31c552ce6TPRuP

// SubRegisterData 子设备注册应答数据域
type SubRegisterData struct {
	IotID        string `json:"iotId"`
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
}

// SubRegisterResponse 子设备注册应答
type SubRegisterResponse struct {
	ID      uint              `json:"id,string"`
	Code    int               `json:"code"`
	Data    []SubRegisterData `json:"data"`
	Message string            `json:"message,omitempty"`
}

// thingSubRegister 子设备动态注册
// 网关类型的设备,通过上行请求为子设备发起动态注册,返回成功注册的子设备的设备证书
// request:   /sys/{productKey}/{deviceName}/thing/sub/register
// response:  /sys/{productKey}/{deviceName}/thing/sub/register_reply
func (sf *Client) thingSubRegister(pk, dn string) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	_uri := sf.URIGateway(uri.SysPrefix, uri.ThingSubRegister)
	return sf.SendRequest(_uri, infra.MethodSubDevRegister, []infra.MetaPair{
		{ProductKey: pk, DeviceName: dn},
	})
}

// ProcThingSubRegisterReply 处理子设备动态注册回复
// request:   /sys/{productKey}/{deviceName}/thing/sub/register
// response:  /sys/{productKey}/{deviceName}/thing/sub/register_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/sub/register_reply
func ProcThingSubRegisterReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := &SubRegisterResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, rsp.Data, err})
	c.Log.Debugf("thing.sub.register.reply @%d", rsp.ID)
	return nil
}
