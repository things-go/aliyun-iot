package dm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thinkgos/aliIOT/infra"
)

// GwSubDevRegisterParams 子设备注册参数域
type GwSubDevRegisterParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwSubDevRegisterData 子设备注册应答数据域
type GwSubDevRegisterData struct {
	IotID        int64  `json:"iotId,string"`
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
}

// GwSubDevRegisterResponse 子设备注册应答
type GwSubDevRegisterResponse struct {
	Response
	Data []GwSubDevRegisterData `json:"data"`
}

// upstreamThingGwSubDevRegister 子设备动态注册
// @return: requestID
// 子设备注册使用网关topic通道
func (sf *Client) upstreamThingGwSubDevRegister(devID int) (int, error) {
	if devID < 0 {
		return 0, ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return 0, err
	}

	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingSubDevRegister),
		id, methodSubDevRegister, []GwSubDevRegisterParams{
			{
				node.ProductKey(),
				node.DeviceName(),
			},
		}); err != nil {
		return 0, err
	}

	sf.CacheInsert(id, devID, MsgTypeSubDevRegister, methodSubDevRegister)
	sf.debug("upstream thing GW <sub>: register @%d", id)
	return id, nil
}

// GwSubDevCombineLoginParams 子设备上线参数域
type GwSubDevCombineLoginParams struct {
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	ClientID     string `json:"clientId"`
	Timestamp    int64  `json:"timestamp,string"`
	SignMethod   string `json:"signMethod"`
	Sign         string `json:"sign"`
	CleanSession bool   `json:"cleanSession,string"`
}

// GwSubDevCombineLoginRequest 子设备上线请求
type GwSubDevCombineLoginRequest struct {
	ID     int                        `json:"id,string"`
	Params GwSubDevCombineLoginParams `json:"params"`
}

// upstreamExtGwSubDevCombineLogin 子设备上线
// 子设备上下线只支持Qos = 0.
func (sf *Client) upstreamExtGwSubDevCombineLogin(devID int) (int, error) {
	if devID < 0 {
		return 0, ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return 0, err
	}

	clientID := fmt.Sprintf("%s.%s|_v=%s|", node.ProductKey(), node.DeviceName(), infra.IOTSDKVersion)
	timestamp := time.Now().Unix()
	sign, err := generateSign(node.ProductKey(), node.DeviceName(), node.DeviceSecret(), clientID, timestamp)
	if err != nil {
		return 0, err
	}
	id := sf.RequestID()
	req, err := json.Marshal(&GwSubDevCombineLoginRequest{
		id,
		GwSubDevCombineLoginParams{
			node.ProductKey(),
			node.DeviceName(),
			clientID,
			timestamp,
			infra.SignMethodHMACSHA1,
			sign,
			true,
		},
	})
	if err != nil {
		return 0, err
	}
	// NOTE: 子设备上线,使用网关的productKey和deviceName,且只支持qos = 0
	if err = sf.Publish(sf.URIServiceSelf(URIExtSessionPrefix, URISubDevCombineLogin),
		0, req); err != nil {
		return 0, err
	}

	sf.CacheInsert(id, devID, MsgTypeSubDevLogin, methodSubDevLogin)
	sf.debug("upstream Ext GW <sub>: login @%d", id)
	return id, nil
}

// GwSubDevCombineLogoutParams 子设备下线参数域
type GwSubDevCombineLogoutParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwSubDevCombineLogoutRequest 子设备下线请求
type GwSubDevCombineLogoutRequest struct {
	ID     int                         `json:"id,string"`
	Params GwSubDevCombineLogoutParams `json:"params"`
}

// upstreamExtGwSubDevCombineLogout 子设备下线
// 子设备上下线只支持Qos = 0.
func (sf *Client) upstreamExtGwSubDevCombineLogout(devID int) (int, error) {
	if devID < 0 {
		return 0, ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return 0, err
	}

	id := sf.RequestID()
	req, err := json.Marshal(&GwSubDevCombineLogoutRequest{
		id,
		GwSubDevCombineLogoutParams{
			node.ProductKey(),
			node.DeviceName(),
		},
	})
	if err != nil {
		return 0, err
	}

	// NOTE: 子设备下线,要用网关的productKey和deviceName
	if err = sf.Publish(sf.URIServiceSelf(URIExtSessionPrefix, URISubDevCombineLogout),
		0, req); err != nil {
		return 0, err
	}
	sf.CacheInsert(id, devID, MsgTypeSubDevLogin, methodSubDevLogout)
	sf.debug("upstream Ext GW <sub>: logout @%d", id)
	return id, nil
}

// ExtErrorData 子设备错误回复数据域
type ExtErrorData struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// ExtErrorResponse 子设备错误回复
type ExtErrorResponse struct {
	Response
	Data ExtErrorData `json:"data"`
}

// ProcExtErrorResponse 处理错误的回复
// response: ext/error/{productKey}/{deviceName}
// subscribe: ext/error/{productKey}/{deviceName}
func ProcExtErrorResponse(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 4) {
		return ErrInvalidURI
	}

	rsp := ExtErrorResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	if rsp.Code != CodeSuccess {
		err = NewCodeError(rsp.Code, rsp.Message)
	}
	c.debug("downstream extend <Error>: response,@%d", rsp.ID)
	return c.ipcSendMessage(&ipcMessage{
		err:     err,
		evt:     ipcEvtErrorResponse,
		payload: rsp.Data,
	})
}
