package dm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
)

// GwSubTopoAddParams 添加设备拓扑关系参数域
type GwSubTopoAddParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	ClientID   string `json:"clientId"`
	Timestamp  int64  `json:"timestamp,string"` // 时间戳
	SignMethod string `json:"signMethod"`       // 支持hmacSha1、hmacSha256、hmacMd5、Sha256。
	Sign       string `json:"sign"`
}

// upstreamThingGwSubTopoAdd 添加设备拓扑关系
// 子设备身份注册后,需网关上报与子设备的关系,然后才进行子设备上线
// request:   /sys/{productKey}/{deviceName}/thing/topo/add
// response:  /sys/{productKey}/{deviceName}/thing/topo/add_reply
func (sf *Client) upstreamThingGwSubTopoAdd(devID int) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}

	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	timestamp := int64(time.Now().Nanosecond()) / 1000000
	clientID := fmt.Sprintf("%s.%s|_v=%s|", node.ProductKey(), node.DeviceName(), sign.AlinkSDKVersion)
	signs, err := generateSign(node.ProductKey(), node.DeviceName(), node.deviceSecret, clientID, timestamp)
	if err != nil {
		return nil, err
	}
	id := sf.RequestID()
	uri := sf.URIServiceSelf(URISysPrefix, URIThingTopoAdd)
	err = sf.SendRequest(uri, id, MethodTopoAdd, []GwSubTopoAddParams{
		{
			node.ProductKey(),
			node.DeviceName(),
			clientID,
			timestamp,
			"hmacsha1",
			signs,
		},
	})
	if err != nil {
		return nil, err
	}

	sf.debugf("upstream GW thing <topo>: add @%d", id)
	return sf.Insert(id), nil
}

// GwSubTopoDeleteParams 删除网关与子设备的拓扑关系参数域
type GwSubTopoDeleteParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// upstreamThingGwSubTopoDelete 删除网关与子设备的拓扑关系
func (sf *Client) upstreamThingGwSubTopoDelete(devID int) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}
	id := sf.RequestID()
	uri := sf.URIServiceSelf(URISysPrefix, URIThingTopoDelete)
	err = sf.SendRequest(uri, id, MethodTopoDelete,
		[]GwSubTopoDeleteParams{
			{
				node.ProductKey(),
				node.DeviceName(),
			},
		})
	if err != nil {
		return nil, err
	}
	sf.debugf("upstream GW thing <topo>: delete @%d", id)
	return sf.Insert(id), nil
}

// GwSubTopoGetData 获取网关和子设备的拓扑关系应答的数据域
type GwSubTopoGetData struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwSubTopoGetResponse 获取网关和子设备的拓扑关系应答
type GwSubTopoGetResponse struct {
	ResponseRawData
	Data []GwSubTopoGetData `json:"data"`
}

// UpstreamThingGwSubTopoGet 获取该网关和子设备的拓扑关系
// request:   /sys/{productKey}/{deviceName}/thing/topo/get
// response:  /sys/{productKey}/{deviceName}/thing/topo/get_reply
func (sf *Client) UpstreamThingGwSubTopoGet() (*Entry, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	id := sf.RequestID()
	uri := sf.URIServiceSelf(URISysPrefix, URIThingTopoGet)
	if err := sf.SendRequest(uri, id, MethodTopoGet, "{}"); err != nil {
		return nil, err
	}
	sf.debugf("upstream GW thing <topo>: Get @%d", id)
	return sf.Insert(id), nil
}

// GwSubListFoundParams 发现设备列表上报参数域
type GwSubListFoundParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// UpstreamThingGwListFound 发现设备列表上报
// 场景,网关可以发现新接入的子设备,发现后,需将新接入的子设备的信息上报云端,
// 然后转到第三方应用,选择哪些子设备可以接入该网关
func (sf *Client) UpstreamThingGwListFound(devID int) (*Entry, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}
	id := sf.RequestID()
	uri := sf.URIServiceSelf(URISysPrefix, URIThingListFound)
	err = sf.SendRequest(uri, id, MethodListFound,
		[]GwSubListFoundParams{
			{
				node.ProductKey(),
				node.DeviceName(),
			},
		})
	if err != nil {
		return nil, err
	}
	sf.Insert(id)
	sf.debugf("upstream GW thing <list>: found @%d", id)
	return sf.Insert(id), nil
}

// ProcThingGwSubTopoAddReply 处理网络拓扑添加
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/topo/add
// response:  /sys/{productKey}/{deviceName}/thing/topo/add_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/add_reply
func ProcThingGwSubTopoAddReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := &ResponseRawData{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else {
		_ = c.SetDevStatusByPkDn(pk, dn, DevStatusAttached)
	}

	c.done(rsp.ID, err, nil)
	c.debugf("downstream GW thing <topo>: add reply @%d", rsp.ID)
	return nil
}

// ProcThingGwSubTopoDeleteReply 处理删除网络拓扑
// 上行
// request:  /sys/{productKey}/{deviceName}/thing/topo/delete
// response:  /sys/{productKey}/{deviceName}/thing/topo/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/delete_reply
func ProcThingGwSubTopoDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}
	rsp := ResponseRawData{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else {
		c.SetDevStatusByPkDn(pk, dn, DevStatusRegistered) // nolint: errcheck
	}

	c.done(rsp.ID, err, nil)
	c.debugf("downstream GW thing <topo>: delete reply @%d", rsp.ID)
	return nil
}

// ProcThingGwSubTopoGetReply 处理获取该网关和子设备的拓扑关系
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/topo/get
// response:  /sys/{productKey}/{deviceName}/thing/topo/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/get_reply
func ProcThingGwSubTopoGetReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}
	rsp := GwSubTopoGetResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.done(rsp.ID, err, nil)
	c.debugf("downstream GW thing <topo>: get reply @%d", rsp.ID)
	return c.eventGwProc.EvtThingGwSubTopoGetReply(c, err, rsp.Data)
}

// ProcThingGwSubListFoundReply 处理发现设备列表上报应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/list/found
// response:  /sys/{productKey}/{deviceName}/thing/list/found_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/list/found_reply
func ProcThingGwSubListFoundReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
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

	c.done(rsp.ID, err, nil)
	c.debugf("downstream GW thing <list>: found reply @%d", rsp.ID)
	return c.eventGwProc.EvtThingListFoundReply(c, err)
}

// GwSubTopoAddNotifyParams 添加设备拓扑关系通知参数域
type GwSubTopoAddNotifyParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwSubTopoAddNotifyRequest 添加设备拓扑关系通知请求
type GwSubTopoAddNotifyRequest struct {
	ID      uint                       `json:"id,string"`
	Version string                     `json:"version"`
	Params  []GwSubTopoAddNotifyParams `json:"params"`
	Method  string                     `json:"method"`
}

// ProcThingGwTopoAddNotify 通知网关添加设备拓扑关系
// 下行
// request:  /sys/{productKey}/{deviceName}/thing/topo/add/notify
// response:  /sys/{productKey}/{deviceName}/thing/topo/add/notify_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/add/notify
func ProcThingGwTopoAddNotify(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 7) {
		return ErrInvalidURI
	}
	c.debugf("downstream GW thing <topo>: add notify")

	req := &GwSubTopoAddNotifyRequest{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}

	if err := c.eventGwProc.EvtThingTopoAddNotify(c, req.Params); err != nil {
		c.warnf("ipc send message failed, %+v", err)
	}
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
}

// GwSubTopoChange 网络拓扑关系变化请求参数域 设备结构
type GwSubTopoChange struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwTopoChangeParams 网络拓扑关系变化请求参数域
type GwTopoChangeParams struct {
	Status  int               `json:"status"` // 0: 创建 1:删除 2: 启用 8: 禁用
	SubList []GwSubTopoChange `json:"subList"`
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
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.debugf("downstream GW thing <topo>: change")

	req := &GwTopoChangeRequest{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}

	if err := c.eventGwProc.EvtThingTopoChange(c, req.Params); err != nil {
		c.warnf("ipc send message failed, %+v", err)
	}
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
}
