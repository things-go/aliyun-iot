package dm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
	"github.com/thinkgos/aliyun-iot/uri"
)

/*
 - 子设备上下线、批量上下线消息，只支持QoS=0，不支持QoS=1。
 - 一个网关下，同时在线的子设备数量不能超过1500。在线子设备数量达到1500个后，新的子设备上线请求将被拒绝。
 - 调用子设备批量上下线接口，单个批次上下线的子设备数量不超过5个。
 - 设备批量上下线接口为原子接口，调用结果为全部成功或全部失败，失败时返回的data中会包含具体的失败信息。
*/

// GwCombineLoginParams 子设备上线参数域
type GwCombineLoginParams struct {
	ProductKey string `json:"productKey"`       // 子设备的产品key
	DeviceName string `json:"deviceName"`       // 子设备的设备名称
	ClientID   string `json:"clientId"`         // 设备端标识
	Timestamp  int64  `json:"timestamp,string"` // 时间戳,单位ms
	SignMethod string `json:"signMethod"`       // 签名方法，支持hmacSha1、hmacSha256、hmacMd5和Sha256。
	Sign       string `json:"sign"`
	// 如果取值是true，则清理所有子设备离线时的消息，即所有未接收的QoS1消息将被清除。
	// 如果取值是false，则不清理子设备离线时的消息
	CleanSession bool `json:"cleanSession,string"`
}

// GwCombineLoginRequest 子设备上线请求
type GwCombineLoginRequest struct {
	ID     uint                 `json:"id,string"`
	Params GwCombineLoginParams `json:"params"`
}

// GwCombineBatchLoginRequest 子设备上线请求
type GwCombineBatchLoginRequest struct {
	ID     uint `json:"id,string"`
	Params struct {
		DeviceList []GwCombineLoginParams `json:"deviceList"`
	} `json:"params"`
}

// GwCombineLoginData 子设备上线应答数据域
type GwCombineLoginData struct {
	ProductKey string `json:"productKey"` // 子设备的产品key
	DeviceName string `json:"deviceName"` // 子设备的设备名称
}

// GwCombineLoginResponse 子设备上线回复
type GwCombineLoginResponse struct {
	ID      uint               `json:"id,string"`
	Code    int                `json:"code"`
	Data    GwCombineLoginData `json:"data"`
	Message string             `json:"message,omitempty"`
}

// GwCombineBatchLoginResponse 子设备批量上线回复
type GwCombineBatchLoginResponse struct {
	ID      uint                 `json:"id,string"`
	Code    int                  `json:"code"`
	Data    []GwCombineLoginData `json:"data"`
	Message string               `json:"message,omitempty"`
}

// ExtCombineLogin 子设备上线
// NOTE: topic 应使用网关的productKey和deviceName,且只支持qos = 0
// request： /ext/session/${productKey}/${deviceName}/combine/login
// response：/ext/session/${productKey}/${deviceName}/combine/login_reply
func (sf *Client) ExtCombineLogin(devID int) (*Token, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	clientID := fmt.Sprintf("%s.%s|_v=%s|", node.ProductKey(), node.DeviceName(), sign.SDKVersion)
	timestamp := int64(time.Now().Nanosecond()) / 1000000
	signs, err := generateSign(node.ProductKey(), node.DeviceName(), node.DeviceSecret(), clientID, timestamp)
	if err != nil {
		return nil, err
	}
	id := sf.RequestID()
	req, err := json.Marshal(&GwCombineLoginRequest{
		id,
		GwCombineLoginParams{
			node.ProductKey(),
			node.DeviceName(),
			clientID,
			timestamp,
			"hmacsha1",
			signs,
			true,
		},
	})
	if err != nil {
		return nil, err
	}
	err = sf.Publish(sf.GatewayURI(uri.ExtSessionPrefix, uri.CombineLogin), 0, req)
	if err != nil {
		return nil, err
	}

	sf.log.Debugf("upstream Ext GW <sub>: login @%d", id)
	return sf.putPending(id), nil
}

// ExtCombineBatchLogin 子设备批量上线
// NOTE: topic 应使用网关的productKey和deviceName,且只支持qos = 0
// request： /ext/session/${productKey}/${deviceName}/combine/batch_login
// response：/ext/session/${productKey}/${deviceName}/combine/batch_login_reply
func (sf *Client) ExtCombineBatchLogin(devID ...int) (*Token, error) {
	// TODO:
	return nil, nil
}

// GwCombineLogoutParams 子设备下线参数域
type GwCombineLogoutParams struct {
	ProductKey string `json:"productKey"` // 子设备所属产品的ProductKey
	DeviceName string `json:"deviceName"` // 子设备的设备名称
}

// GwCombineLogoutRequest 子设备下线请求
type GwCombineLogoutRequest struct {
	ID     uint                  `json:"id,string"`
	Params GwCombineLogoutParams `json:"params"`
}

// GwCombineBatchLogoutRequest 子设备批量下线请求
type GwCombineBatchLogoutRequest struct {
	ID     uint                    `json:"id,string"`
	Params []GwCombineLogoutParams `json:"params"`
}

// ExtCombineLogout 子设备下线
// NOTE: topic 应使用网关的productKey和deviceName,且只支持qos = 0
// request:   /ext/session/{productKey}/{deviceName}/combine/logout
// response:  /ext/session/{productKey}/{deviceName}/combine/logout_reply
func (sf *Client) ExtCombineLogout(devID int) (*Token, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	id := sf.RequestID()
	req, err := json.Marshal(&GwCombineLogoutRequest{
		id,
		GwCombineLogoutParams{
			node.ProductKey(),
			node.DeviceName(),
		},
	})
	if err != nil {
		return nil, err
	}

	err = sf.Publish(sf.GatewayURI(uri.ExtSessionPrefix, uri.CombineLogout), 0, req)
	if err != nil {
		return nil, err
	}
	sf.log.Debugf("upstream Ext GW <sub>: logout @%d", id)
	return sf.putPending(id), nil
}

// ExtCombineBatchLogout 子设备批量下线
// NOTE: topic 应使用网关的productKey和deviceName,且只支持qos = 0
// request:   /ext/session/{productKey}/{deviceName}/combine/batch_logout
// response:  /ext/session/{productKey}/{deviceName}/combine/batch_logout_reply
func (sf *Client) ExtCombineBatchLogout(devID ...int) (*Token, error) {
	// TODO
	return nil, nil
}

// ProcExtCombineLoginReply 子设备上线应答处理
// 上行
// request:   /ext/session/{productKey}/{deviceName}/combine/login
// response:  /ext/session/{productKey}/{deviceName}/combine/login_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/login_reply
func ProcExtCombineLoginReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := &GwCombineLoginResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	pk, dn := uris[1], uris[2]
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else {
		c.SetDevStatusByPkDn(pk, dn, DevStatusLogined) // nolint: errcheck
	}
	c.signalPending(Message{rsp.ID, nil, err})
	c.log.Debugf("downstream Ext GW <sub>: login reply @%d", rsp.ID)
	return nil
}

// ProcExtCombineBatchLoginReply 子设备批量上线应答处理
// 上行
// request:   /ext/session/{productKey}/{deviceName}/combine/batch_login
// response:  /ext/session/{productKey}/{deviceName}/combine/batch_login_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/batch_login_reply
func ProcExtCombineBatchLoginReply(c *Client, rawURI string, payload []byte) error {
	// TODO:
	return nil
}

// ProcExtCombineLoginoutReply 子设备下线应答处理
// 上行
// request:   /ext/session/{productKey}/{deviceName}/combine/logout
// response:  /ext/session/{productKey}/{deviceName}/combine/logout_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/logout_reply
func ProcExtCombineLoginoutReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := &ResponseRawData{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	pk, dn := uris[1], uris[2]
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else {
		c.SetDevStatusByPkDn(pk, dn, DevStatusAttached) // nolint: errcheck
	}
	c.signalPending(Message{rsp.ID, nil, err})
	c.log.Debugf("downstream Ext GW <sub>: logout reply @%d", rsp.ID)
	return nil
}

// ProcExtCombineBatchLogoutReply 子设备批量下线应答处理
// 上行
// request:   /ext/session/{productKey}/{deviceName}/combine/batch_logout
// response:  /ext/session/{productKey}/{deviceName}/combine/batch_logout_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/batch_logout_reply
func ProcExtCombineBatchLogoutReply(c *Client, rawURI string, payload []byte) error {
	// TODO:
	return nil
}
