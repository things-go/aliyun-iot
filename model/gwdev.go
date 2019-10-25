package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thinkgos/aliIOT/infra"
)

// MetaInfo 产品与设备三元组
type MetaInfo struct {
	ProductKey    string
	ProductSecret string
	DeviceName    string
	DeviceSecret  string
}

// GwSubDevRegisterParams 子设备注册参数域
type GwSubDevRegisterParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// SubDevRegisterDataReply 子设备注册应答数据域
type GwSubDevRegisterData struct {
	IotId        int64  `json:"iotId,string"`
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
}

type GwSubDevRegisterResponse struct {
	Response
	Data []GwSubDevRegisterData `json:"data"`
}

// UpstreamGwSubDevRegister 子设备动态注册
func (sf *Manager) UpstreamGwSubDevRegister(meta ...*MetaInfo) error {
	sublist := make([]GwSubDevRegisterParams, 0, len(meta))

	for _, v := range meta {
		sublist = append(sublist, GwSubDevRegisterParams{v.ProductKey, v.DeviceName})
	}

	return sf.SendRequest(sf.URIService(URISysPrefix, URIThingSubDevRegister), sf.RequestID(), methodSubDevRegister, sublist)
}

// GwSubDevCombineLoginParams 子设备上线参数域
type GwSubDevCombineLoginParams struct {
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	ClientId     string `json:"clientId"`
	Timestamp    int64  `json:"timestamp,string"`
	SignMethod   string `json:"signMethod"`
	Sign         string `json:"sign"`
	CleanSession bool   `json:"cleanSession,string"`
}

// GwSubDevCombineLoginRequest 子设备上线请求
type GwSubDevCombineLoginRequest struct {
	ID     int                        `json:"id,string"`
	Params GwSubDevCombineLoginParams `json:"params"`
}

// UpstreamGwExtSubDevCombineLogin 子设备上线
// 子设备上下线只支持Qos = 0.
func (sf *Manager) UpstreamGwExtSubDevCombineLogin(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	clientID := fmt.Sprintf("%s.%s|_v=%s|", node.ProductKey, node.DeviceName, infra.IOTSDKVersion)
	timestamp := time.Now().Unix()
	sign, err := generateSign(
		node.ProductKey, node.DeviceName, node.DeviceSecret,
		clientID, timestamp)
	if err != nil {
		return err
	}
	req, err := json.Marshal(&GwSubDevCombineLoginRequest{
		ID: sf.RequestID(),
		Params: GwSubDevCombineLoginParams{
			ProductKey:   node.ProductKey,
			DeviceName:   node.DeviceName,
			ClientId:     clientID,
			Timestamp:    timestamp,
			SignMethod:   infra.SignMethodHMACSHA1,
			Sign:         sign,
			CleanSession: true,
		},
	})
	if err != nil {
		return err
	}
	// NOTE: 子设备登陆,要用网关的productKey和deviceName
	return sf.Publish(sf.URIService(URIExtSessionPrefix, URISubDevCombineLogin), 0, req)
}

// GwSubDevCombineLogoutParams 子设备下线参数域
type GwSubDevCombineLogoutParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// GwSubDevCombineLogoutRequest 子设备下线请求
type GwSubDevCombineLogoutRequest struct {
	ID     int                         `json:"id,string"`
	Params GwSubDevCombineLogoutParams `json:"params"`
}

// UpstreamGwExtSubDevCombineLogout 子设备下线
// 子设备上下线只支持Qos = 0.
func (sf *Manager) UpstreamGwExtSubDevCombineLogout(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	req, err := json.Marshal(&GwSubDevCombineLogoutRequest{
		sf.RequestID(),
		GwSubDevCombineLogoutParams{
			ProductKey: node.ProductKey,
			DeviceName: node.DeviceName,
		},
	})
	if err != nil {
		return err
	}

	// NOTE: 子设备下线,要用网关的productKey和deviceName
	return sf.Publish(sf.URIService(URIExtSessionPrefix, URISubDevCombineLogout), 0, req)
}
