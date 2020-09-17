package dm

import (
	"encoding/json"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/89299.html?spm=a2c4g.11186623.6.704.6237597f3s8Q9t

// TopoAddParams 添加设备拓扑关系参数域
type TopoAddParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	ClientID   string `json:"clientId"`
	Timestamp  int64  `json:"timestamp,string"` // 时间戳
	SignMethod string `json:"signMethod"`       // 支持hmacSha1、hmacSha256、hmacMd5、Sha256。
	Sign       string `json:"sign"`
}

// ThingTopoAdd 添加设备拓扑关系
// 子设备身份注册后,需网关上报与子设备的关系,然后才进行子设备上线
// request:   /sys/{productKey}/{deviceName}/thing/topo/add
// response:  /sys/{productKey}/{deviceName}/thing/topo/add_reply
func (sf *Client) ThingTopoAdd(pk, dn string) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	ds, err := sf.DeviceSecret(pk, dn)
	if err != nil {
		return nil, err
	}

	timestamp := infra.Millisecond(time.Now())
	clientID := ClientID(pk, dn)
	signs, err := generateSign(pk, dn, ds, clientID, timestamp)
	if err != nil {
		return nil, err
	}
	_uri := sf.URIGateway(uri.SysPrefix, uri.ThingTopoAdd)
	return sf.SendRequest(_uri, infra.MethodTopoAdd, []TopoAddParams{
		{
			pk,
			dn,
			clientID,
			timestamp,
			"hmacsha1",
			signs,
		},
	})
}

// ThingTopoDelete 删除网关与子设备的拓扑关系
// request： /sys/{productKey}/{deviceName}/thing/topo/delete
// response：/sys/{productKey}/{deviceName}/thing/topo/delete_reply
func (sf *Client) ThingTopoDelete(pk, dn string) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	_uri := sf.URIGateway(uri.SysPrefix, uri.ThingTopoDelete)
	return sf.SendRequest(_uri, infra.MethodTopoDelete, []infra.MetaPair{
		{ProductKey: pk, DeviceName: dn},
	})
}

// ThingTopoGet 获取该网关和子设备的拓扑关系
// request:   /sys/{productKey}/{deviceName}/thing/topo/get
// response:  /sys/{productKey}/{deviceName}/thing/topo/get_reply
func (sf *Client) ThingTopoGet() (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	_uri := sf.URIGateway(uri.SysPrefix, uri.ThingTopoGet)
	return sf.SendRequest(_uri, infra.MethodTopoGet, "{}")
}

// ThingListFound 发现设备列表上报
// 场景,网关可以发现新接入的子设备,发现后,需将新接入的子设备的信息上报云端,
// 然后转到第三方应用,选择哪些子设备可以接入该网关
// request： /sys/{productKey}/{deviceName}/thing/list/found
// response：/sys/{productKey}/{deviceName}/thing/list/found_reply
func (sf *Client) ThingListFound(pairs []infra.MetaPair) (*Token, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	if len(pairs) == 0 {
		return nil, ErrInvalidParameter
	}
	_uri := sf.URIGateway(uri.SysPrefix, uri.ThingListFound)
	return sf.SendRequest(_uri, infra.MethodListFound, pairs)
}

// TopoAddResponse 添加网络拓扑应答
type TopoAddResponse struct {
	ID      uint             `json:"id,string"`
	Code    int              `json:"code"`
	Data    []infra.MetaPair `json:"Data"`
	Message string           `json:"message,omitempty"`
}

// ProcThingTopoAddReply 处理网络拓扑添加
// request:   /sys/{productKey}/{deviceName}/thing/topo/add
// response:  /sys/{productKey}/{deviceName}/thing/topo/add_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/add_reply
func ProcThingTopoAddReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	rsp := &TopoAddResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, rsp.Data, err})
	c.Log.Debugf("thing.topo.add.reply @%d", rsp.ID)
	return nil
}

// TopoDeleteResponse 删除网络拓扑应答
type TopoDeleteResponse struct {
	ID      uint             `json:"id,string"`
	Code    int              `json:"code"`
	Data    []infra.MetaPair `json:"Data"`
	Message string           `json:"message,omitempty"`
}

// ProcThingTopoDeleteReply 处理删除网络拓扑
// request:   /sys/{productKey}/{deviceName}/thing/topo/delete
// response:  /sys/{productKey}/{deviceName}/thing/topo/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/delete_reply
func ProcThingTopoDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := TopoDeleteResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, rsp.Data, err})
	c.Log.Debugf("thing.topo.delete.reply @%d", rsp.ID)
	return nil
}

// TopoGetResponse 获取网关和子设备的拓扑关系应答
type TopoGetResponse struct {
	ID      uint             `json:"id,string"`
	Code    int              `json:"code"`
	Data    []infra.MetaPair `json:"Data"`
	Message string           `json:"message,omitempty"`
}

// ProcThingTopoGetReply 处理获取该网关和子设备的拓扑关系
// request:   /sys/{productKey}/{deviceName}/thing/topo/get
// response:  /sys/{productKey}/{deviceName}/thing/topo/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/get_reply
func ProcThingTopoGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := TopoGetResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, rsp.Data, err})
	c.Log.Debugf("thing.topo.get.reply @%d", rsp.ID)
	return c.gwCb.ThingTopoGetReply(c, err, rsp.Data)
}

// ProcThingListFoundReply 处理发现设备列表上报应答
// request:   /sys/{productKey}/{deviceName}/thing/list/found
// response:  /sys/{productKey}/{deviceName}/thing/list/found_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/list/found_reply
func ProcThingListFoundReply(c *Client, rawURI string, payload []byte) error {
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
	c.Log.Debugf("thing.list.found.reply @%d", rsp.ID)
	return c.gwCb.ThingListFoundReply(c, err)
}

// TopoAddNotifyRequest 添加设备拓扑关系通知请求
type TopoAddNotifyRequest struct {
	ID      uint             `json:"id,string"`
	Version string           `json:"version"`
	Params  []infra.MetaPair `json:"params"`
	Method  string           `json:"method"`
}

// ProcThingTopoAddNotify 通知网关添加设备拓扑关系
// request:   /sys/{productKey}/{deviceName}/thing/topo/add/notify
// response:  /sys/{productKey}/{deviceName}/thing/topo/add/notify_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/add/notify
func ProcThingTopoAddNotify(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	c.Log.Debugf("thing.topo.add.notify")

	req := &TopoAddNotifyRequest{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}

	_uri := uri.ReplyWithRequestURI(rawURI)
	err := c.Response(_uri, Response{ID: req.ID, Code: infra.CodeSuccess, Data: "{}"})
	if err != nil {
		c.Log.Warnf("thing.topo.add.notify.response, %+v", err)
	}
	return c.gwCb.ThingTopoAddNotify(c, req.Params)
}

// TopoChangeParams 网络拓扑关系变化请求参数域
type TopoChangeParams struct {
	Status  int              `json:"status"` // 0: 创建 1:删除 2: 启用 8: 禁用
	SubList []infra.MetaPair `json:"subList"`
}

// TopoChangeRequest 网络拓扑关系变化请求
type TopoChangeRequest struct {
	ID      uint             `json:"id,string"`
	Version string           `json:"version"`
	Params  TopoChangeParams `json:"params"`
	Method  string           `json:"method"`
}

// ProcThingTopoChange 通知网关拓扑关系变化
// request:    /sys/{productKey}/{deviceName}/thing/topo/change
// response:   /sys/{productKey}/{deviceName}/thing/topo/change_reply
// subscribe:  /sys/{productKey}/{deviceName}/thing/topo/change
func ProcThingTopoChange(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	c.Log.Debugf("thing.topo.change")

	req := &TopoChangeRequest{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}

	_uri := uri.ReplyWithRequestURI(rawURI)
	err := c.Response(_uri, Response{ID: req.ID, Code: infra.CodeSuccess, Data: "{}"})
	if err != nil {
		c.Log.Warnf("thing.topo.change.response, %+v", err)
	}
	return c.gwCb.ThingTopoChange(c, req.Params)
}
