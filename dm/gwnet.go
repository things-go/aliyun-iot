package dm

import (
	"fmt"
	"time"

	"github.com/thinkgos/aliIOT/infra"
)

// SubDevTopoAddParams 添加设备拓扑关系参数域
type SubDevTopoAddParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	ClientID   string `json:"clientId"`
	Timestamp  int64  `json:"timestamp,string"`
	SignMethod string `json:"signMethod"`
	Sign       string `json:"sign"`
}

// UpstreamGwThingTopoAdd 添加设备拓扑关系
// 子设备身份注册后,需网关上报与子设备的关系,然后才进行子设备上线
func (sf *Client) UpstreamGwThingTopoAdd(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	clientID := fmt.Sprintf("%s.%s|_v=%s|", node.ProductKey(), node.DeviceName(), infra.IOTSDKVersion)
	sign, err := generateSign(node.ProductKey(), node.DeviceName(), node.deviceSecret, clientID, timestamp)
	if err != nil {
		return err
	}
	id := sf.RequestID()
	err = sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingTopoAdd),
		id, methodTopoAdd, []SubDevTopoAddParams{
			{
				node.ProductKey(),
				node.DeviceName(),
				clientID,
				timestamp,
				infra.SignMethodHMACSHA1,
				sign,
			},
		})
	if err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeTopoAdd, methodTopoAdd)
	sf.debug("upstream GW thing <topo>: add @%d", id)
	return nil
}

// GwTopoDeleteParams
type GwTopoDeleteParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// UpstreamGwThingTopoDelete 删除网关与子设备的拓扑关系
func (sf *Client) UpstreamGwThingTopoDelete(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}
	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingTopoDelete),
		id, methodTopoDelete, []GwTopoDeleteParams{
			{
				node.ProductKey(),
				node.DeviceName(),
			},
		}); err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeTopoDelete, methodTopoDelete)
	sf.debug("upstream GW thing <topo>: delete @%d", id)
	return nil
}

// GwTopoGetData 获取网关和子设备的拓扑关系应答的数据域
type GwTopoGetData struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwTopoGetResponse 获取网关和子设备的拓扑关系应答
type GwTopoGetResponse struct {
	Response
	Data []GwTopoGetData `json:"data"`
}

// UpstreamGwThingTopoGet 获取该网关和子设备的拓扑关系
func (sf *Client) UpstreamGwThingTopoGet() error {
	id := sf.RequestID()
	if err := sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingTopoGet),
		id, methodTopoGet, "{}"); err != nil {
		return err
	}
	sf.CacheInsert(id, DevNodeLocal, MsgTypeTopoGet, methodTopoGet)
	sf.debug("upstream GW thing <topo>: Get @%d", id)
	return nil
}

// GwDevListFoundParams 发现设备列表上报参数域
type GwDevListFoundParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// UpstreamGwThingListFound 发现设备列表上报
// 场景,网关可以发现新接入的子设备,发现后,需将新接入的子设备的信息上报云端,
// 然后转到第三方应用,选择哪些子设备可以接入该网关
func (sf *Client) UpstreamGwThingListFound(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}
	id := sf.RequestID()
	if err = sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingListFound),
		id, methodListFound, []GwDevListFoundParams{
			{
				node.ProductKey(),
				node.DeviceName(),
			},
		}); err != nil {
		return err
	}
	sf.CacheInsert(id, DevNodeLocal, MsgTypeDevListFound, methodListFound)
	sf.debug("upstream GW thing <list>: found @%d", id)
	return nil
}
