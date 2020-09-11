package dm

import (
	"fmt"
	"time"

	sign2 "github.com/thinkgos/aliyun-iot/sign"
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

// upstreamGwThingTopoAdd 添加设备拓扑关系
// 子设备身份注册后,需网关上报与子设备的关系,然后才进行子设备上线
func (sf *Client) upstreamGwThingTopoAdd(devID int) (uint, error) {
	if devID < 0 {
		return 0, ErrInvalidParameter
	}

	node, err := sf.SearchNode(devID)
	if err != nil {
		return 0, err
	}

	timestamp := time.Now().Unix()
	clientID := fmt.Sprintf("%s.%s|_v=%s|", node.ProductKey(), node.DeviceName(), sign2.AlinkSDKVersion)
	sign, err := generateSign(node.ProductKey(), node.DeviceName(), node.deviceSecret, clientID, timestamp)
	if err != nil {
		return 0, err
	}
	id := sf.RequestID()
	err = sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingTopoAdd),
		id, MethodTopoAdd, []SubDevTopoAddParams{
			{
				node.ProductKey(),
				node.DeviceName(),
				clientID,
				timestamp,
				"hmacsha1",
				sign,
			},
		})
	if err != nil {
		return 0, err
	}
	sf.CacheInsert(id, devID, MsgTypeTopoAdd)
	sf.debugf("upstream GW thing <topo>: add @%d", id)
	return id, nil
}

// GwTopoDeleteParams 删除网关与子设备的拓扑关系参数域
type GwTopoDeleteParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// upstreamGwThingTopoDelete 删除网关与子设备的拓扑关系
func (sf *Client) upstreamGwThingTopoDelete(devID int) (uint, error) {
	if devID < 0 {
		return 0, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return 0, err
	}
	id := sf.RequestID()
	err = sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingTopoDelete),
		id, MethodTopoDelete, []GwTopoDeleteParams{
			{
				node.ProductKey(),
				node.DeviceName(),
			},
		})
	if err != nil {
		return 0, err
	}
	sf.CacheInsert(id, devID, MsgTypeTopoDelete)
	sf.debugf("upstream GW thing <topo>: delete @%d", id)
	return id, nil
}

// GwTopoGetData 获取网关和子设备的拓扑关系应答的数据域
type GwTopoGetData struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwTopoGetResponse 获取网关和子设备的拓扑关系应答
type GwTopoGetResponse struct {
	ResponseRawData
	Data []GwTopoGetData `json:"data"`
}

// upstreamGwThingTopoGet 获取该网关和子设备的拓扑关系
func (sf *Client) upstreamGwThingTopoGet() error {
	id := sf.RequestID()
	if err := sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingTopoGet),
		id, MethodTopoGet, "{}"); err != nil {
		return err
	}
	sf.CacheInsert(id, DevNodeLocal, MsgTypeTopoGet)
	sf.debugf("upstream GW thing <topo>: Get @%d", id)
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
	node, err := sf.SearchNode(devID)
	if err != nil {
		return err
	}
	id := sf.RequestID()
	err = sf.SendRequest(sf.URIServiceSelf(URISysPrefix, URIThingListFound),
		id, MethodListFound, []GwDevListFoundParams{
			{
				node.ProductKey(),
				node.DeviceName(),
			},
		})
	if err != nil {
		return err
	}
	sf.CacheInsert(id, DevNodeLocal, MsgTypeDevListFound)
	sf.debugf("upstream GW thing <list>: found @%d", id)
	return nil
}
