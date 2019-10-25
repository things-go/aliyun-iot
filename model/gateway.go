package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
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

// SubDevRegisterParams 子设备注册参数域
type SubDevRegisterParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// SubDevRegisterDataReply 子设备注册应答数据域
type SubDevRegisterData struct {
	IotId        int64  `json:"iotId,string"`
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
}

type SubDevRegisterResponse struct {
	Response
	Data []SubDevRegisterData `json:"data"`
}

// UpstreamSubDevRegister 子设备动态注册
func (sf *Manager) UpstreamSubDevRegister(meta ...*MetaInfo) error {
	sublist := make([]SubDevRegisterParams, 0, len(meta))

	for _, v := range meta {
		sublist = append(sublist, SubDevRegisterParams{v.ProductKey, v.DeviceName})
	}

	return sf.SendRequest(sf.URIService(URISysPrefix, URIThingSubDevRegister), sf.RequestID(), methodSubDevRegister, sublist)
}

// SubDevCombineLoginParams 子设备上线参数域
type SubDevCombineLoginParams struct {
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	ClientId     string `json:"clientId"`
	Timestamp    int64  `json:"timestamp,string"`
	SignMethod   string `json:"signMethod"`
	Sign         string `json:"sign"`
	CleanSession bool   `json:"cleanSession,string"`
}

// SubDevCombineLoginRequest 子设备上线请求
type SubDevCombineLoginRequest struct {
	ID     int                      `json:"id,string"`
	Params SubDevCombineLoginParams `json:"params"`
}

// 可以采用Sha256,hmacMd5,hmacSha1,hmacSha256
func generateSign(productKey, deviceName, deviceSecret, clientID string, timestamp int64) (string, error) {
	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
		clientID, deviceName, productKey, timestamp)

	/* setup Password */
	h := hmac.New(sha1.New, []byte(deviceSecret))
	if _, err := h.Write([]byte(signSource)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// UpstreamExtSubDevCombineLogin 子设备上线
// 子设备上下线只支持Qos = 0.
func (sf *Manager) UpstreamExtSubDevCombineLogin(devID int) error {
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
	req, err := json.Marshal(&SubDevCombineLoginRequest{
		ID: sf.RequestID(),
		Params: SubDevCombineLoginParams{
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

// SubDevCombineLogoutParams 子设备下线参数域
type SubDevCombineLogoutParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// SubDevCombineLogoutRequest 子设备下线请求
type SubDevCombineLogoutRequest struct {
	ID     int                       `json:"id,string"`
	Params SubDevCombineLogoutParams `json:"params"`
}

// UpstreamExtSubDevCombineLogout 子设备下线
// 子设备上下线只支持Qos = 0.
func (sf *Manager) UpstreamExtSubDevCombineLogout(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	req, err := json.Marshal(&SubDevCombineLogoutRequest{
		sf.RequestID(),
		SubDevCombineLogoutParams{
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

type SubDevTopoAddParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	ClientId   string `json:"clientId"`
	Timestamp  int64  `json:"timestamp,string"`
	SignMethod string `json:"signMethod"`
	Sign       string `json:"sign"`
}

func (sf *Manager) UpstreamThingTopoAdd(metas ...*MetaInfo) error {
	var clientID string
	var sign string
	var err error

	sublist := make([]SubDevTopoAddParams, 0, len(metas))

	timestamp := time.Now().Unix()
	for _, v := range metas {
		clientID = fmt.Sprintf("%s.%s|_v=%s|", v.ProductKey, v.DeviceName, infra.IOTSDKVersion)
		sign, err = generateSign(v.ProductKey, v.DeviceName, v.DeviceSecret, clientID, timestamp)
		if err != nil {
			return err
		}
		sublist = append(sublist, SubDevTopoAddParams{
			v.ProductKey,
			v.DeviceSecret,
			clientID,
			timestamp,
			infra.SignMethodHMACSHA1,
			sign,
		})
	}

	return sf.SendRequest(sf.URIService(URISysPrefix, URIThingTopoAdd),
		sf.RequestID(), methodTopoAdd, sublist)
}

func (sf *Manager) UpstreamThingTopoDelete(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	_, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	// TODO
	return nil
}

type gwUserProc struct{}

func (gwUserProc) DownstreamExtSubDevRegisterReply(m *Manager, rsp *SubDevRegisterResponse) error {
	return nil
}

func (gwUserProc) DownstreamExtSubDevCombineLoginReply(m *Manager, rsp *Response) error {
	return nil
}

// DownstreamExtSubDevCombineLogoutReply
func (gwUserProc) DownstreamExtSubDevCombineLogoutReply(m *Manager, rsp *Response) error {
	return nil
}

func (gwUserProc) DownstreamThingTopoAddReply(m *Manager, rsp *Response) error {
	return nil
}

func (gwUserProc) DownstreamThingTopoDeleteReply(m *Manager, rsp *Response) error {
	return nil
}
