package ahttp

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"time"
)

const (
	signMethodSHA1 = "hmacsha1"
	signMethodMD5  = "hmacmd5"
	defaultTimeout = time.Second * 2
)

// AuthRequest 鉴权请求
type AuthRequest struct {
	Version    string `json:"version"`
	ClientID   string `json:"clientId"`
	SignMethod string `json:"signmethod"`
	Sign       string `json:"sign"`
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	// 校验时间戳15分钟内的请求有效。时间戳格式为数值，
	// 值为自GMT 1970年1月1日0时0分到当前时间点所经过的毫秒数。
	Timestamp int64 `json:"timestamp"`
}

// AuthResponse 鉴权回复
type AuthResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Info    struct {
		Token string `json:"token"`
	} `json:"info"`
}

// Client 客户端
type Client struct {
	productKey   string
	deviceName   string
	deviceSecret string
	host         string
	version      string
	signMethod   string

	token string

	c *http.Client
}

// 默认
func New() *Client {
	return &Client{
		host:       "https://iot-as-http.cn-shanghai.aliyuncs.com",
		version:    "default",
		signMethod: signMethodMD5,
		c: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// SetHost 设置主机
func (sf *Client) SetHost(h string) *Client {
	if h != "" {
		sf.host = h
	}

	return sf
}

// SetRequestTimeout 设置请求超时时间
func (sf *Client) SetRequestTimeout(t time.Duration) *Client {
	sf.c.Timeout = t
	return sf
}

// SetDeviceMetaInfo 设置设备三元组信息
func (sf *Client) SetDeviceMetaInfo(productKey, deviceName, deviceSecret string) *Client {
	sf.productKey = productKey
	sf.deviceName = deviceName
	sf.deviceSecret = deviceSecret
	return sf
}

// SetSignMethod 设置签名方法
func (sf *Client) SetSignMethod(method string) *Client {
	if method == signMethodMD5 || method == signMethodSHA1 {
		sf.signMethod = method
	} else {
		sf.signMethod = signMethodMD5
	}
	return sf
}

func (sf *AuthRequest) generateSign(deviceSecret string) error {
	var f func() hash.Hash

	if sf.SignMethod == signMethodSHA1 {
		f = sha1.New
	} else {
		f = md5.New
		sf.SignMethod = signMethodMD5
	}
	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
		sf.ClientID, sf.DeviceName, sf.ProductKey, sf.Timestamp)
	h := hmac.New(f, []byte(deviceSecret))
	if _, err := h.Write([]byte(signSource)); err != nil {
		return err
	}

	sf.Sign = hex.EncodeToString(h.Sum(nil))
	return nil
}

// SendAuth 鉴权
func (sf *Client) SendAuth() error {
	if sf.productKey == "" ||
		sf.deviceName == "" ||
		sf.deviceSecret == "" {
		return errors.New("invalid meta info")
	}

	authPy := AuthRequest{
		Version:    sf.version,
		ClientID:   sf.productKey + "." + sf.deviceName,
		SignMethod: sf.signMethod,
		ProductKey: sf.productKey,
		DeviceName: sf.deviceName,
		Timestamp:  time.Now().Unix() * 1000,
	}

	if err := authPy.generateSign(sf.deviceSecret); err != nil {
		return err
	}

	b, err := json.Marshal(&authPy)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, sf.host+"/auth", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := sf.c.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	rspPy := AuthResponse{}
	if err = json.NewDecoder(response.Body).Decode(&rspPy); err != nil {
		return err
	}

	if rspPy.Code != 0 {
		switch rspPy.Code {
		case 10000:
			err = ErrUnknown
		case 10001:
			err = ErrParamException
		case 20000:
			err = ErrAuthFailed
		case 20004:
			err = ErrUpdateSessionFailed
		case 40000:
			err = ErrRequestTooMany
		default:
			err = ErrUnknown
		}
		return err
	}
	sf.token = rspPy.Info.Token
	return nil
}
