package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
)

// ProcThingTopoAddReply 处理网络拓扑添加
// 上行
// request:  /sys/{productKey}/{deviceName}/thing/topo/add
// response:  /sys/{productKey}/{deviceName}/thing/topo/add_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/add_reply
func ProcThingTopoAddReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else if devID, ok := c.CacheGet(rsp.ID); ok {
		_ = c.SetDevStatusByID(devID, DevStatusAttached)
	}

	c.CacheDone(rsp.ID, err)
	c.debugf("downstream GW thing <topo>: add reply @%d", rsp.ID)
	return nil
}

// ProcThingTopoDeleteReply 处理删除网络拓扑
// 上行
// request:  /sys/{productKey}/{deviceName}/thing/topo/delete
// response:  /sys/{productKey}/{deviceName}/thing/topo/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/delete_reply
func ProcThingTopoDeleteReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else if devID, ok := c.CacheGet(rsp.ID); ok {
		_ = c.SetDevStatusByID(devID, DevStatusRegistered)
	}

	c.CacheDone(rsp.ID, err)
	c.debugf("downstream GW thing <topo>: delete reply @%d", rsp.ID)
	return nil
}

// ProcThingTopoGetReply 处理获取该网关和子设备的拓扑关系
// 上行
// request:  /sys/{productKey}/{deviceName}/thing/topo/get
// response:  /sys/{productKey}/{deviceName}/thing/topo/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/get_reply
func ProcThingTopoGetReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
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

	c.CacheDone(rsp.ID, err)
	c.debugf("downstream GW thing <topo>: get reply @%d", rsp.ID)
	return c.ipcSendMessage(&ipcMessage{
		err:     err,
		evt:     ipcEvtTopoGetReply,
		payload: rsp.Data,
	})
}

// ProcThingListFoundReply 处理发现设备列表上报应答
// 上行
// request:  /sys/{productKey}/{deviceName}/thing/list/found
// response:  /sys/{productKey}/{deviceName}/thing/list/found_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/list/found_reply
func ProcThingListFoundReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.CacheDone(rsp.ID, err)
	c.debugf("downstream GW thing <list>: found reply @%d", rsp.ID)
	return c.ipcSendMessage(&ipcMessage{
		err: err,
		evt: ipcEvtListFoundReply,
	})
}

// GwTopoAddNotifyParams 添加设备拓扑关系通知参数域
type GwTopoAddNotifyParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwTopoAddNotifyRequest 添加设备拓扑关系通知请求
type GwTopoAddNotifyRequest struct {
	Request
	Params []GwTopoAddNotifyParams `json:"params"`
}

// ProcThingTopoAddNotify 通知网关添加设备拓扑关系
// 下行
// request:  /sys/{productKey}/{deviceName}/thing/topo/add/notify
// response:  /sys/{productKey}/{deviceName}/thing/topo/add/notify_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/topo/add/notify
func ProcThingTopoAddNotify(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 7) {
		return ErrInvalidURI
	}
	c.debugf("downstream GW thing <topo>: add notify")

	req := GwTopoAddNotifyRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := c.ipcSendMessage(&ipcMessage{
		evt:     ipcEvtTopoAddNotify,
		payload: req.Params,
	}); err != nil {
		c.warnf("ipc send message failed, %+v", err)
	}
	return c.SendResponse(uriServiceReplyWithRequestURI(rawURI),
		req.ID, infra.CodeSuccess, "{}")
}

// GwTopoChangeDev 网络拓扑关系变化请求参数域 设备结构
type GwTopoChangeDev struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwTopoChangeParams 网络拓扑关系变化请求参数域
type GwTopoChangeParams struct {
	Status  int               `json:"status"` // 0: 创建 1:删除 2: 启用 8: 禁用
	SubList []GwTopoChangeDev `json:"subList"`
}

// GwTopoChangeRequest 网络拓扑关系变化请求
type GwTopoChangeRequest struct {
	Request
	Params GwTopoChangeParams `json:"params"`
}

// ProcThingTopoChange 通知网关拓扑关系变化
// 下行
// request:  /sys/{productKey}/{deviceName}/thing/topo/change
// response:  /sys/{productKey}/{deviceName}/thing/topo/change_reply
// subscribe:  /sys/{productKey}/{deviceName}/thing/topo/change
func ProcThingTopoChange(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	c.debugf("downstream GW thing <topo>: change")

	req := GwTopoChangeRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	if err := c.ipcSendMessage(&ipcMessage{
		evt:     ipcTopoChange,
		payload: req.Params,
	}); err != nil {
		c.warnf("ipc send message failed, %+v", err)
	}
	return c.SendResponse(uriServiceReplyWithRequestURI(rawURI),
		req.ID, infra.CodeSuccess, "{}")
}

/*************************************** 子设备相关处理 *************************************************************/

// ProcThingSubDevRegisterReply 子设备动态注册处理
// 上行
// request: /sys/{productKey}/{deviceName}/thing/sub/register
// response: /sys/{productKey}/{deviceName}/thing/sub/register_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/sub/register_reply
func ProcThingSubDevRegisterReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}
	rsp := GwSubDevRegisterResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else {
		for _, v := range rsp.Data {
			node, er := c.SearchNodeByPkDn(v.ProductKey, v.DeviceName)
			if er != nil {
				c.warnf("downstream GW thing <sub>: register reply, %+v <%s - %s - %s>",
					er, v.ProductKey, v.DeviceName, v.DeviceSecret)
				continue
			}
			_ = c.SetDeviceSecretByID(node.ID(), v.DeviceSecret)
			_ = c.SetDevStatusByID(node.ID(), DevStatusRegistered)
		}
	}
	c.CacheDone(rsp.ID, err)
	c.debugf("downstream GW thing <sub>: register reply @%d", rsp.ID)
	return nil
}

// ProcExtSubDevCombineLoginReply 子设备上线应答处理
// 上行
// request: /ext/session/{productKey}/{deviceName}/combine/login
// response: /ext/session/{productKey}/{deviceName}/combine/login_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/login_reply
func ProcExtSubDevCombineLoginReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else if devID, ok := c.CacheGet(rsp.ID); ok {
		_ = c.SetDevStatusByID(devID, DevStatusLogined)
	}
	c.CacheDone(rsp.ID, err)
	c.debugf("downstream Ext GW <sub>: login reply @%d", rsp.ID)
	return nil
}

// ProcExtSubDevCombineLogoutReply 子设备下线应答处理
// 上行
// request: /ext/session/{productKey}/{deviceName}/combine/logout
// response: /ext/session/{productKey}/{deviceName}/combine/logout_reply
// subscribe: /ext/session/{productKey}/{deviceName}/combine/logout_reply
func ProcExtSubDevCombineLogoutReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := Response{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	} else if devID, ok := c.CacheGet(rsp.ID); ok {
		_ = c.SetDevStatusByID(devID, DevStatusAttached)
	}
	c.CacheDone(rsp.ID, err)
	c.debugf("downstream Ext GW <sub>: logout reply @%d", rsp.ID)
	return nil
}

// ProcThingDisable 禁用子设备
// 下行
// request: /sys/{productKey}/{deviceName}/thing/disable
// response: /sys/{productKey}/{deviceName}/thing/disable_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/disable
func ProcThingDisable(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 5) {
		return ErrInvalidURI
	}
	c.debugf("downstream <thing>: disable >> %s - %s", uris[1], uris[2])

	req := Request{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return err
	}
	if err = c.SetDevAvailByPkDN(uris[1], uris[2], false); err != nil {
		c.warnf("<thing> disable failed, %+v", err)
	}
	if err = c.ipcSendMessage(&ipcMessage{
		evt:        ipcThingDisable,
		productKey: uris[1],
		deviceName: uris[2],
	}); err != nil {
		c.warnf("<thing> disable, ipc send message failed, %+v", err)
	}
	return c.SendResponse(uriServiceReplyWithRequestURI(rawURI),
		req.ID, infra.CodeSuccess, "{}")
}

// ProcThingEnable 启用子设备
// 下行
// request: /sys/{productKey}/{deviceName}/thing/enable
// response: /sys/{productKey}/{deviceName}/thing/enable_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/enable
func ProcThingEnable(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 5) {
		return ErrInvalidURI
	}
	c.debugf("downstream <thing>: enable >> %s - %s", uris[1], uris[2])

	req := Request{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return err
	}
	if err = c.SetDevAvailByPkDN(uris[1], uris[2], true); err != nil {
		c.warnf("<thing> enable failed, %+v", err)
	}
	if err = c.ipcSendMessage(&ipcMessage{
		evt:        ipcThingEnable,
		productKey: uris[1],
		deviceName: uris[2],
	}); err != nil {
		c.warnf("<thing> enable, ipc send message failed, %+v", err)
	}
	return c.SendResponse(uriServiceReplyWithRequestURI(rawURI),
		req.ID, infra.CodeSuccess, "{}")
}

// ProcThingDelete 子设备删除,网关类型设备
// 下行
// request: /sys/{productKey}/{deviceName}/thing/delete
// response: /sys/{productKey}/{deviceName}/thing/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/delete
func ProcThingDelete(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.cfg.uriOffset + 5) {
		return ErrInvalidURI
	}
	c.debugf("downstream <thing>: delete >> %s - %s", uris[1], uris[2])

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	c.DeleteByPkDn(uris[1], uris[2])
	if err := c.ipcSendMessage(&ipcMessage{
		evt:        ipcThingDelete,
		productKey: uris[1],
		deviceName: uris[2],
	}); err != nil {
		c.warnf("<thing> delete, ipc send message failed, %+v", err)
	}
	return c.SendResponse(uriServiceReplyWithRequestURI(rawURI),
		req.ID, infra.CodeSuccess, "{}")
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

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.CacheDone(rsp.ID, err)
	c.debugf("downstream extend <Error>: response,@%d", rsp.ID)
	return c.ipcSendMessage(&ipcMessage{
		err:     err,
		evt:     ipcEvtErrorResponse,
		payload: rsp.Data,
	})
}
