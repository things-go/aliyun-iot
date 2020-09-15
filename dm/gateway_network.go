package dm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
	"github.com/thinkgos/aliyun-iot/uri"
)

// GwTopoAddParams 添加设备拓扑关系参数域
type GwTopoAddParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	ClientID   string `json:"clientId"`
	Timestamp  int64  `json:"timestamp,string"` // 时间戳
	SignMethod string `json:"signMethod"`       // 支持hmacSha1、hmacSha256、hmacMd5、Sha256。
	Sign       string `json:"sign"`
}

// ThingGwTopoAdd 添加设备拓扑关系
// 子设备身份注册后,需网关上报与子设备的关系,然后才进行子设备上线
// request:   /sys/{productKey}/{deviceName}/thing/topo/add
// response:  /sys/{productKey}/{deviceName}/thing/topo/add_reply
func (sf *Client) ThingGwTopoAdd(pk, dn string) (*Token, error) {
	ds, err := sf.DeviceSecret(pk, dn)
	if err != nil {
		return nil, err
	}

	timestamp := int64(time.Now().Nanosecond()) / 1000000
	clientID := fmt.Sprintf("%s.%s|_v=%s|", pk, dn, sign.SDKVersion)
	signs, err := generateSign(pk, dn, ds, clientID, timestamp)
	if err != nil {
		return nil, err
	}
	id := sf.RequestID()
	_uri := sf.GatewayURI(uri.SysPrefix, uri.ThingTopoAdd)
	err = sf.SendRequest(_uri, id, infra.MethodTopoAdd, []GwTopoAddParams{
		{
			pk,
			dn,
			clientID,
			timestamp,
			"hmacsha1",
			signs,
		},
	})
	if err != nil {
		return nil, err
	}

	sf.log.Debugf("upstream GW thing <topo>: add @%d", id)
	return sf.putPending(id), nil
}

// GwTopoDeleteParams 删除网关与子设备的拓扑关系参数域
type GwTopoDeleteParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// ThingGwTopoDelete 删除网关与子设备的拓扑关系
func (sf *Client) ThingGwTopoDelete(pk, dn string) (*Token, error) {
	id := sf.RequestID()
	_uri := sf.GatewayURI(uri.SysPrefix, uri.ThingTopoDelete)
	err := sf.SendRequest(_uri, id, infra.MethodTopoDelete,
		[]GwTopoDeleteParams{
			{pk, dn},
		})
	if err != nil {
		return nil, err
	}
	sf.log.Debugf("upstream GW thing <topo>: delete @%d", id)
	return sf.putPending(id), nil
}

// GwTopoGetData 获取网关和子设备的拓扑关系应答的数据域
type GwTopoGetData struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwTopoGetResponse 获取网关和子设备的拓扑关系应答
type GwTopoGetResponse struct {
	ResponseRawData
	Data []GwTopoGetData `json:"Data"`
}

// ThingGwTopoGet 获取该网关和子设备的拓扑关系
// request:   /sys/{productKey}/{deviceName}/thing/topo/get
// response:  /sys/{productKey}/{deviceName}/thing/topo/get_reply
func (sf *Client) ThingGwTopoGet() (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	id := sf.RequestID()
	_uri := sf.GatewayURI(uri.SysPrefix, uri.ThingTopoGet)
	if err := sf.SendRequest(_uri, id, infra.MethodTopoGet, "{}"); err != nil {
		return nil, err
	}
	sf.log.Debugf("upstream GW thing <topo>: Get @%d", id)
	return sf.putPending(id), nil
}

// GwListFoundParams 发现设备列表上报参数域
type GwListFoundParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// ThingGwListFound 发现设备列表上报
// 场景,网关可以发现新接入的子设备,发现后,需将新接入的子设备的信息上报云端,
// 然后转到第三方应用,选择哪些子设备可以接入该网关
func (sf *Client) ThingGwListFound(pk, dn string) (*Token, error) {
	id := sf.RequestID()
	_uri := sf.GatewayURI(uri.SysPrefix, uri.ThingListFound)
	err := sf.SendRequest(_uri, id, infra.MethodListFound,
		[]GwListFoundParams{
			{pk, dn},
		})
	if err != nil {
		return nil, err
	}
	sf.putPending(id)
	sf.log.Debugf("upstream GW thing <list>: found @%d", id)
	return sf.putPending(id), nil
}

// ProcThingGwTopoAddReply 处理网络拓扑添加
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/topo/add
// response:  /sys/{productKey}/{deviceName}/thing/topo/add_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/add_reply
func ProcThingGwTopoAddReply(c *Client, rawURI string, payload []byte) error {
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
		_ = c.SetDevStatus(pk, dn, DevStatusAttached)
	}

	c.signalPending(Message{rsp.ID, nil, err})
	c.log.Debugf("downstream GW thing <topo>: add reply @%d", rsp.ID)
	return nil
}

// ProcThingGwTopoDeleteReply 处理删除网络拓扑
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/topo/delete
// response:  /sys/{productKey}/{deviceName}/thing/topo/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/delete_reply
func ProcThingGwTopoDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := ResponseRawData{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	pk, dn := uris[1], uris[2]
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else {
		c.SetDevStatus(pk, dn, DevStatusRegistered) // nolint: errcheck
	}

	c.signalPending(Message{rsp.ID, nil, err})
	c.log.Debugf("downstream GW thing <topo>: delete reply @%d", rsp.ID)
	return nil
}

// ProcThingGwTopoGetReply 处理获取该网关和子设备的拓扑关系
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/topo/get
// response:  /sys/{productKey}/{deviceName}/thing/topo/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/get_reply
func ProcThingGwTopoGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := GwTopoGetResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, nil, err})
	c.log.Debugf("downstream GW thing <topo>: get reply @%d", rsp.ID)
	return c.gwCb.ThingGwTopoGetReply(c, err, rsp.Data)
}

// ProcThingGwListFoundReply 处理发现设备列表上报应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/list/found
// response:  /sys/{productKey}/{deviceName}/thing/list/found_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/list/found_reply
func ProcThingGwListFoundReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := ResponseRawData{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, nil, err})
	c.log.Debugf("downstream GW thing <list>: found reply @%d", rsp.ID)
	return c.gwCb.ThingGwListFoundReply(c, err)
}

// GwTopoAddNotifyParams 添加设备拓扑关系通知参数域
type GwTopoAddNotifyParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwTopoAddNotifyRequest 添加设备拓扑关系通知请求
type GwTopoAddNotifyRequest struct {
	ID      uint                    `json:"id,string"`
	Version string                  `json:"version"`
	Params  []GwTopoAddNotifyParams `json:"params"`
	Method  string                  `json:"method"`
}

// ProcThingGwTopoAddNotify 通知网关添加设备拓扑关系
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/topo/add/notify
// response:  /sys/{productKey}/{deviceName}/thing/topo/add/notify_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/add/notify
func ProcThingGwTopoAddNotify(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	c.log.Debugf("downstream GW thing <topo>: add notify")

	req := &GwTopoAddNotifyRequest{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}

	if err := c.gwCb.ThingGwTopoAddNotify(c, req.Params); err != nil {
		c.log.Warnf("ipc send Message failed, %+v", err)
	}
	return c.SendResponse(uri.ReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
}

// GwTopoChange 网络拓扑关系变化请求参数域的设备结构
type GwTopoChange struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwTopoChangeParams 网络拓扑关系变化请求参数域
type GwTopoChangeParams struct {
	Status  int            `json:"status"` // 0: 创建 1:删除 2: 启用 8: 禁用
	SubList []GwTopoChange `json:"subList"`
}

// GwTopoChangeRequest 网络拓扑关系变化请求
type GwTopoChangeRequest struct {
	ID      uint               `json:"id,string"`
	Version string             `json:"version"`
	Params  GwTopoChangeParams `json:"params"`
	Method  string             `json:"method"`
}

// ProcThingGwTopoChange 通知网关拓扑关系变化
// 下行
// request:  /sys/{productKey}/{deviceName}/thing/topo/change
// response:  /sys/{productKey}/{deviceName}/thing/topo/change_reply
// subscribe:  /sys/{productKey}/{deviceName}/thing/topo/change
func ProcThingGwTopoChange(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	c.log.Debugf("downstream GW thing <topo>: change")

	req := &GwTopoChangeRequest{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}

	if err := c.gwCb.ThingGwTopoChange(c, req.Params); err != nil {
		c.log.Warnf("ipc send Message failed, %+v", err)
	}
	return c.SendResponse(uri.ReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
}
