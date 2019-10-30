package dm

import (
	"fmt"
	"time"

	"github.com/thinkgos/aliIOT/infra"
)

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
	clientID := fmt.Sprintf("%s.%s|_v=%s|", node.productKey, node.deviceName, infra.IOTSDKVersion)
	sign, err := generateSign(node.productKey, node.deviceName, node.deviceSecret, clientID, timestamp)
	if err != nil {
		return err
	}
	id := sf.RequestID()
	err = sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingTopoAdd),
		id, methodTopoAdd, []SubDevTopoAddParams{
			{
				node.productKey,
				node.deviceName,
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

type GwTopoDeleteParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// UpstreamGwThingTopoDelete 删除与子设备的关系
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
				node.productKey,
				node.deviceName,
			},
		}); err != nil {
		return err
	}
	sf.CacheInsert(id, devID, MsgTypeTopoDelete, methodTopoDelete)
	sf.debug("upstream GW thing <topo>: delete @%d", id)
	return nil
}

type GwTopoGetData struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

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

type GwDevListFoundParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

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
		id, methodListFound, []GwDevListFoundParams{{
			node.productKey,
			node.deviceName}}); err != nil {
		return err
	}
	sf.CacheInsert(id, DevNodeLocal, MsgTypeDevListFound, methodListFound)
	sf.debug("upstream GW thing <list>: found @%d", id)
	return nil
}
