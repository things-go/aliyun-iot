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

//SHA256       "Sha256"
//HMACMD5      "hmacMd5"
//HMACSHA1     "hmacSha1"
//HMACSHA256   "hmacSha256"

// CombineSUbDevLoginParams 子设备上线参数域
type CombineSUbDevLoginParams struct {
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	ClientId     string `json:"clientId"`
	Timestamp    int64  `json:"timestamp,string"`
	SignMethod   string `json:"signMethod"`
	Sign         string `json:"sign"`
	CleanSession bool   `json:"cleanSession,string"`
}

// CombineSubDevLoginRequest 子设备上线请求
type CombineSubDevLoginRequest struct {
	ID     int                      `json:"id,string"`
	Params CombineSUbDevLoginParams `json:"params"`
}

func generateLoginSign(productKey, deviceName, deviceSecret, clientID string, timestamp int64) (string, error) {
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
	sign, err := generateLoginSign(
		node.ProductKey, node.DeviceName, node.DeviceSecret,
		clientID, timestamp)
	if err != nil {
		return err
	}
	req, err := json.Marshal(&CombineSubDevLoginRequest{
		ID: sf.RequestID(),
		Params: CombineSUbDevLoginParams{
			ProductKey:   node.ProductKey,
			DeviceName:   node.DeviceName,
			ClientId:     clientID,
			Timestamp:    timestamp,
			SignMethod:   "hmacSha1",
			Sign:         sign,
			CleanSession: true,
		},
	})
	if err != nil {
		return err
	}
	// NOTE: 子设备登陆,要用网关的productKey和deviceName
	return sf.Publish(sf.URIService(URIExtSessionPrefix, CombineSubDevLogin), 0, req)
}

// CombineSUbDevLogoutParams 子设备下线参数域
type CombineSUbDevLogoutParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// CombineSubDevLogoutRequest 子设备下线请求
type CombineSubDevLogoutRequest struct {
	ID     int                       `json:"id,string"`
	Params CombineSUbDevLogoutParams `json:"params"`
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

	req, err := json.Marshal(&CombineSubDevLogoutRequest{
		sf.RequestID(),
		CombineSUbDevLogoutParams{
			ProductKey: node.ProductKey,
			DeviceName: node.DeviceName,
		},
	})
	if err != nil {
		return err
	}

	// NOTE: 子设备下线,要用网关的productKey和deviceName
	return sf.Publish(sf.URIService(URIExtSessionPrefix, CombineSubDevLogout), 0, req)
}
