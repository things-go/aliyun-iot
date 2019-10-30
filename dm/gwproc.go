package dm

import (
	"encoding/json"
)

/******************************** gateway ****************************************************************/

func ProcThingTopoAddReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	if rsp.Code != CodeSuccess {
		c.syncHub.Done(rsp.ID, NewCodeError(rsp.Code, rsp.Message))
	} else {
		c.syncHub.Done(rsp.ID, nil)
	}
	c.debug("downstream GW thing <topo>: add reply @%d", rsp.ID)
	return nil
}

func ProcThingTopoDeleteReply(c *Client, rawURI string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	if rsp.Code != CodeSuccess {
		c.syncHub.Done(rsp.ID, NewCodeError(rsp.Code, rsp.Message))
	} else {
		c.syncHub.Done(rsp.ID, nil)
	}
	c.debug("downstream GW thing <topo>: delete reply @%d", rsp.ID)
	return nil
}

// ProcThingTopoGetReply 处理获取该网关和子设备的拓扑关系
func ProcThingTopoGetReply(c *Client, rawURI string, payload []byte) error {
	rsp := GwTopoGetResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("downstream GW thing <topo>: get reply @%d", rsp.ID)
	// TODO: 处理网关与子设备的拓扑关系
	return nil
}

// ProcThingListFoundReply 处理发现设备列表上报
func ProcThingListFoundReply(c *Client, _ string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}

	c.CacheRemove(rsp.ID)
	c.debug("downstream GW thing <list>: found reply @%d", rsp.ID)
	return nil
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

// 通知网关添加设备拓扑关系
func ProcThingTopoAddNotify(c *Client, rawURI string, payload []byte) error {
	req := GwTopoAddNotifyRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	// TODO: 处理添加设备拓扑关系通知请求
	c.debug("downstream GW thing <topo>: add notify")
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}")
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
func ProcThingTopoChange(c *Client, rawURI string, payload []byte) error {
	req := GwTopoChangeRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	// TODO: 处理通知网关拓扑关系变化
	c.debug("downstream GW thing <topo>: change")
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}")
}

/*************************************** 子设备相关处理 *************************************************************/

// ProcThingSubDevRegisterReply 子设备动态注册处理
func ProcThingSubDevRegisterReply(c *Client, _ string, payload []byte) error {
	rsp := GwSubDevRegisterResponse{}

	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	c.debug("downstream GW thing <sub>: register reply @%d", rsp.ID)

	if rsp.Code != CodeSuccess {
		c.syncHub.Done(rsp.ID, NewCodeError(rsp.Code, rsp.Message))
		return nil
	}

	for _, v := range rsp.Data {
		if err := c.SetDeviceSecretByPkDn(v.ProductKey, v.DeviceName, v.DeviceSecret); err != nil {
			c.warn("downstream GW thing <sub>: register reply, %+v <%s - %s - %s>",
				err, v.ProductKey, v.DeviceName, v.DeviceSecret)
		}
	}
	c.syncHub.Done(rsp.ID, nil)
	return nil
}

// ProcExtSubDevCombineLoginReply 子设备上线应答处理
func ProcExtSubDevCombineLoginReply(c *Client, _ string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	if rsp.Code != CodeSuccess {
		c.syncHub.Done(rsp.ID, NewCodeError(rsp.Code, rsp.Message))
	} else {
		c.syncHub.Done(rsp.ID, nil)
	}
	c.debug("downstream Ext GW <sub>: login reply @%d", rsp.ID)
	return nil
}

// ProcExtSubDevCombineLogoutReply 子设备下线应答处理
func ProcExtSubDevCombineLogoutReply(c *Client, _ string, payload []byte) error {
	rsp := Response{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.CacheRemove(rsp.ID)
	if rsp.Code != CodeSuccess {
		c.syncHub.Done(rsp.ID, NewCodeError(rsp.Code, rsp.Message))
	} else {
		c.syncHub.Done(rsp.ID, nil)
	}
	c.debug("downstream Ext GW <sub>: logout reply @%d", rsp.ID)
	return nil
}

// ProcThingDisable 禁用子设备
func ProcThingDisable(c *Client, rawURI string, payload []byte) error {
	var err error

	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}
	c.debug("downstream <thing>: disable >> %s - %s", uris[1], uris[2])

	req := Request{}
	if err = json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if err = c.SetDevAvailByPkDN(uris[1], uris[2], false); err != nil {
		c.warn("<thing> disable failed, %+v", err)
	}

	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}")
}

// ProcThingEnable 启用子设备
func ProcThingEnable(c *Client, rawURI string, payload []byte) error {
	var err error
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}
	c.debug("downstream <thing>: enable >> %s - %s", uris[1], uris[2])

	req := Request{}
	if err = json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if err = c.SetDevAvailByPkDN(uris[1], uris[2], true); err != nil {
		c.warn("<thing> enable failed, %+v", err)
	}

	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}")
}

// ProcThingDelete 子设备删除
func ProcThingDelete(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) != 5 {
		return ErrInvalidURI
	}
	c.debug("downstream <thing>: delete >> %s - %s", uris[1], uris[2])

	req := Request{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	c.DeleteByPkDn(uris[1], uris[2])

	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI),
		req.ID, CodeSuccess, "{}")
}
