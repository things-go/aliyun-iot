package aiot

import (
	"encoding/json"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

/*
 @see https://help.aliyun.com/document_detail/89300.html?spm=a2c4g.11186623.6.705.290045333dGHQV
 - 子设备上下线、批量上下线消息，只支持QoS=0，不支持QoS=1。
 - 一个网关下，同时在线的子设备数量不能超过1500。在线子设备数量达到1500个后，新的子设备上线请求将被拒绝。
 - 调用子设备批量上下线接口，单个批次上下线的子设备数量不超过5个。
 - 设备批量上下线接口为原子接口，调用结果为全部成功或全部失败，失败时返回的data中会包含具体的失败信息。
*/

// CombinePair combine pair
type CombinePair struct {
	ProductKey   string
	DeviceName   string
	CleanSession bool
}

// CombineLoginParams 子设备上线参数域
type CombineLoginParams struct {
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

// CombineLoginRequest 子设备上线请求
type CombineLoginRequest struct {
	ID     uint               `json:"id,string"`
	Params CombineLoginParams `json:"params"`
}

// CombineLoginResponse 子设备上线回复
type CombineLoginResponse struct {
	ID      uint           `json:"id,string"`
	Code    int            `json:"code"`
	Data    infra.MetaPair `json:"data"`
	Message string         `json:"message,omitempty"`
}

// extCombineLogin 子设备上线
// NOTE: topic 应使用网关的productKey和deviceName,且只支持qos = 0
// cleanSession
// 	如果取值是true，则清理所有子设备离线时的消息，即所有未接收的QoS1消息将被清除。
// 	如果取值是false，则不清理子设备离线时的消息
// request： /ext/session/${productKey}/${deviceName}/combine/login
// response：/ext/session/${productKey}/${deviceName}/combine/login_reply
func (sf *Client) extCombineLogin(cp CombinePair) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	ds, err := sf.DeviceSecret(cp.ProductKey, cp.DeviceName)
	if err != nil {
		return nil, err
	}

	timestamp := infra.Millisecond(time.Now())
	clientID, signs := generateSign("hmacsha256",
		infra.MetaTriad{
			ProductKey:   cp.ProductKey,
			DeviceName:   cp.DeviceName,
			DeviceSecret: ds,
		},
		timestamp)
	id := sf.nextRequestID()
	req, err := json.Marshal(&CombineLoginRequest{
		id,
		CombineLoginParams{
			cp.ProductKey,
			cp.DeviceName,
			clientID,
			timestamp,
			"hmacsha256",
			signs,
			cp.CleanSession,
		},
	})
	if err != nil {
		return nil, err
	}
	_uri := sf.URIGateway(uri.ExtSessionPrefix, uri.CombineLogin)
	err = sf.Publish(_uri, 0, req)
	if err != nil {
		return nil, err
	}
	sf.Log.Debugf("ext.session.combine.login @%d", id)
	return sf.putPending(id), nil
}

// CombineBatchLoginRequest 子设备上线请求
type CombineBatchLoginRequest struct {
	ID     uint `json:"id,string"`
	Params struct {
		DeviceList []CombineLoginParams `json:"deviceList"`
	} `json:"params"`
}

// CombineBatchLoginResponse 子设备批量上线回复
type CombineBatchLoginResponse struct {
	ID      uint             `json:"id,string"`
	Code    int              `json:"code"`
	Data    []infra.MetaPair `json:"data"`
	Message string           `json:"message,omitempty"`
}

// extCombineBatchLogin 子设备批量上线
// NOTE: topic 应使用网关的productKey和deviceName,且只支持qos = 0
// request： /ext/session/${productKey}/${deviceName}/combine/batch_login
// response：/ext/session/${productKey}/${deviceName}/combine/batch_login_reply
func (sf *Client) extCombineBatchLogin(pairs []CombinePair) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	if len(pairs) == 0 {
		return nil, ErrInvalidParameter
	}

	timestamp := infra.Millisecond(time.Now())
	clps := make([]CombineLoginParams, 0, len(pairs))
	for _, cp := range pairs {
		ds, err := sf.DeviceSecret(cp.ProductKey, cp.DeviceName)
		if err != nil {
			return nil, err
		}
		clientID, signs := generateSign("hmacsha256", infra.MetaTriad{
			ProductKey:   cp.ProductKey,
			DeviceName:   cp.DeviceName,
			DeviceSecret: ds,
		}, timestamp)
		clps = append(clps, CombineLoginParams{
			cp.ProductKey,
			cp.DeviceName,
			clientID,
			timestamp,
			"hmacsha256",
			signs,
			cp.CleanSession,
		})
	}

	id := sf.nextRequestID()
	req, err := json.Marshal(&CombineBatchLoginRequest{
		id,
		struct {
			DeviceList []CombineLoginParams `json:"deviceList"`
		}{
			clps,
		},
	})
	if err != nil {
		return nil, err
	}
	_uri := sf.URIGateway(uri.ExtSessionPrefix, uri.CombineBatchLogin)
	err = sf.Publish(_uri, 0, req)
	if err != nil {
		return nil, err
	}
	sf.Log.Debugf("ext.session.combine.batch.login @%d", id)
	return sf.putPending(id), nil
}

// CombineLogoutRequest 子设备下线请求
type CombineLogoutRequest struct {
	ID     uint           `json:"id,string"`
	Params infra.MetaPair `json:"params"`
}

// CombineLogoutResponse 子设备上线回复
type CombineLogoutResponse struct {
	ID      uint           `json:"id,string"`
	Code    int            `json:"code"`
	Data    infra.MetaPair `json:"data"`
	Message string         `json:"message,omitempty"`
}

// extCombineLogout 子设备下线
// NOTE: topic 应使用网关的productKey和deviceName,且只支持qos = 0
// request:   /ext/session/{productKey}/{deviceName}/combine/logout
// response:  /ext/session/{productKey}/{deviceName}/combine/logout_reply
func (sf *Client) extCombineLogout(pk, dn string) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}

	id := sf.nextRequestID()
	req, err := json.Marshal(&CombineLogoutRequest{
		id,
		infra.MetaPair{ProductKey: pk, DeviceName: dn},
	})
	if err != nil {
		return nil, err
	}

	_uri := sf.URIGateway(uri.ExtSessionPrefix, uri.CombineLogout)
	err = sf.Publish(_uri, 0, req)
	if err != nil {
		return nil, err
	}
	sf.Log.Debugf("ext.session.combine.logout @%d", id)
	return sf.putPending(id), nil
}

// CombineBatchLogoutRequest 子设备批量下线请求
type CombineBatchLogoutRequest struct {
	ID     uint             `json:"id,string"`
	Params []infra.MetaPair `json:"params"`
}

// CombineBatchLogoutResponse 子设备批量下线回复
type CombineBatchLogoutResponse struct {
	ID      uint             `json:"id,string"`
	Code    int              `json:"code"`
	Data    []infra.MetaPair `json:"data"`
	Message string           `json:"message,omitempty"`
}

// extCombineBatchLogout 子设备批量下线
// NOTE: topic 应使用网关的productKey和deviceName,且只支持qos = 0
// request:   /ext/session/{productKey}/{deviceName}/combine/batch_logout
// response:  /ext/session/{productKey}/{deviceName}/combine/batch_logout_reply
func (sf *Client) extCombineBatchLogout(pairs []infra.MetaPair) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	if len(pairs) == 0 {
		return nil, ErrInvalidParameter
	}

	id := sf.nextRequestID()
	req, err := json.Marshal(&CombineBatchLogoutRequest{id, pairs})
	if err != nil {
		return nil, err
	}

	_uri := sf.URIGateway(uri.ExtSessionPrefix, uri.CombineBatchLogout)
	err = sf.Publish(_uri, 0, req)
	if err != nil {
		return nil, err
	}
	sf.Log.Debugf("ext.session.combine.batch.login @%d", id)
	return sf.putPending(id), nil
}

// ProcExtCombineLoginReply 处理子设备上线应答
// request:   /ext/session/{productKey}/{deviceName}/combine/login
// response:  /ext/session/{productKey}/{deviceName}/combine/login_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/login_reply
func ProcExtCombineLoginReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	c.Log.Debugf(string(payload))
	rsp := &CombineLoginResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, nil, err})
	c.Log.Debugf("ext.session.combine.login.reply @%d", rsp.ID)
	return nil
}

// ProcExtCombineBatchLoginReply 子设备批量上线应答处理
// request:   /ext/session/{productKey}/{deviceName}/combine/batch_login
// response:  /ext/session/{productKey}/{deviceName}/combine/batch_login_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/batch_login_reply
func ProcExtCombineBatchLoginReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := &CombineBatchLoginResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, nil, err})
	c.Log.Debugf("ext.session.combine.batch.login.reply @%d", rsp.ID)
	return nil
}

// ProcExtCombineLogoutReply 子设备下线应答处理
// 上行
// request:   /ext/session/{productKey}/{deviceName}/combine/logout
// response:  /ext/session/{productKey}/{deviceName}/combine/logout_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/logout_reply
func ProcExtCombineLogoutReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := &CombineLogoutResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, nil, err})
	c.Log.Debugf("ext.session.combine.logout.reply @%d", rsp.ID)
	return nil
}

// ProcExtCombineBatchLogoutReply 子设备批量下线应答处理
// 上行
// request:   /ext/session/{productKey}/{deviceName}/combine/batch_logout
// response:  /ext/session/{productKey}/{deviceName}/combine/batch_logout_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/batch_logout_reply
func ProcExtCombineBatchLogoutReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := &CombineBatchLogoutResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, nil, err})
	c.Log.Debugf("ext.session.combine.batch.logout.reply @%d", rsp.ID)
	return nil
}
