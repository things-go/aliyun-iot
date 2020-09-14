package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// GwSubRegisterParams 子设备注册参数域
type GwSubRegisterParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

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
func (sf *Client) ThingGwSubRegister(devID int) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	uri := sf.GatewayURI(uri.SysPrefix, uri.ThingSubRegister)
	err = sf.SendRequest(uri, id, infra.MethodSubDevRegister,
		[]GwSubRegisterParams{
			{node.ProductKey(), node.DeviceName()},
		})
	if err != nil {
		return nil, err
	}

	sf.log.Debugf("upstream thing GW <sub>: register @%d", id)
	return sf.Insert(id), nil
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
			node, er := c.SearchNodeByPkDn(v.ProductKey, v.DeviceName)
			if er != nil {
				c.log.Warnf("downstream GW thing <sub>: register reply, %+v <%s - %s - %s>",
					er, v.ProductKey, v.DeviceName, v.DeviceSecret)
				continue
			}
			_ = c.SetDeviceSecret(node.ID(), v.DeviceSecret)
			_ = c.SetDevStatus(node.ID(), DevStatusRegistered)
		}
	}
	c.signal(rsp.ID, err, nil)
	c.log.Debugf("downstream GW thing <sub>: register reply @%d", rsp.ID)
	return nil
}
