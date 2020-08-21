// Package ahttp 实现http client 上传数据. 授权方式为自动调用授权,可手动调用,也可以直接调用发送数据接口
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
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/thinkgos/aliyun-iot/clog"
)

// Sign method
const (
	HMACSHA1 = "hmacsha1"
	HMACMD5  = "hmacmd5"
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

	token atomic.Value
	group singleflight.Group

	httpc *http.Client
	*clog.Clog
}

// Option client option
type Option func(c *Client)

// WithHTTPClient with custom http.Client
func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpc = c
	}
}

// WithHost 设置远程主机,
func WithHost(h string) Option {
	return func(c *Client) {
		if !strings.Contains(h, "://") {
			h = "http://" + h
		}
		if h != "" {
			c.host = h
		}
	}
}

// WithSignMethod 设置签名方法,目前支持hmacMD5和hmacSHA1
func WithSignMethod(method string) Option {
	return func(c *Client) {
		if method == HMACSHA1 {
			c.signMethod = HMACSHA1
		} else {
			c.signMethod = HMACMD5
		}
	}
}

// New 新建alink http client
// 默认hmacmd5加签算法
// 默认上海host
// 请求超时2秒
func New(opts ...Option) *Client {
	c := &Client{
		host:       "https://iot-as-http.cn-shanghai.aliyuncs.com",
		version:    "default",
		signMethod: HMACMD5,
		httpc:      http.DefaultClient,
		Clog:       clog.New(clog.WithLogger(clog.NewLogger(log.New(os.Stderr, "alink http --> ", log.LstdFlags)))),
	}
	c.token.Store("")
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// SetDeviceMetaInfo 设置设备三元组信息
func (sf *Client) SetDeviceMetaInfo(productKey, deviceName, deviceSecret string) *Client {
	sf.productKey = productKey
	sf.deviceName = deviceName
	sf.deviceSecret = deviceSecret
	return sf
}

// 鉴权
func (sf *Client) getToken() (string, error) {
	if sf.productKey == "" || sf.deviceName == "" || sf.deviceSecret == "" {
		return "", errors.New("invalid device meta info")
	}

	if token := sf.token.Load().(string); token != "" {
		return token, nil
	}

	tk, err, _ := sf.group.Do("auth", func() (interface{}, error) {
		authReq := AuthRequest{
			Version:    sf.version,
			ClientID:   sf.productKey + "." + sf.deviceName,
			SignMethod: sf.signMethod,
			ProductKey: sf.productKey,
			DeviceName: sf.deviceName,
			Timestamp:  time.Now().Unix() * 1000,
		}

		if err := authReq.generateSign(sf.deviceSecret); err != nil {
			return "", err
		}

		b, err := json.Marshal(&authReq)
		if err != nil {
			return "", err
		}

		request, err := http.NewRequest(http.MethodPost, sf.host+"/auth", bytes.NewBuffer(b))
		if err != nil {
			return "", err
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := sf.httpc.Do(request)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()

		authRsp := AuthResponse{}
		if err = json.NewDecoder(response.Body).Decode(&authRsp); err != nil {
			return "", err
		}

		if authRsp.Code != CodeSuccess {
			return "", NewCodeError(authRsp.Code, authRsp.Message)
		}
		sf.token.Store(authRsp.Info.Token)
		return authRsp.Info.Token, nil
	})
	return tk.(string), err
}

// DataResponse 上报数据回复
type DataResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Info    struct {
		MessageID int64 `json:"messageID"`
	} `json:"info"`
}

func (sf *Client) Publish(uri string, payload interface{}) error {
	py := DataResponse{}

	for retry := 0; retry < 1; retry++ {
		token, err := sf.getToken()
		if err != nil {
			return err
		}

		var buf *bytes.Buffer
		switch v := payload.(type) {
		case string:
			buf = bytes.NewBufferString(v)
		case []byte:
			buf = bytes.NewBuffer(v)
		default:
			return errors.New("Unknown payload type, must be string or []byte")
		}

		request, err := http.NewRequest(http.MethodPost, sf.host+uri, buf)
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", "application/octet-stream")
		request.Header.Set("password", token)
		response, err := sf.httpc.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if err = json.NewDecoder(response.Body).Decode(&py); err != nil {
			return err
		}
		sf.Debug("publish response, %+v", py)
		if py.Code == 0 {
			return nil
		}

		if !(py.Code == CodeTokenExpired ||
			py.Code == CodeTokenCheckFailed ||
			py.Code == CodeTokenIsNull) {
			sf.token.Store("")
			return NewCodeError(py.Code, py.Message)
		}
	}
	return NewCodeError(py.Code, py.Message)
}

func (sf *AuthRequest) generateSign(deviceSecret string) error {
	hashFunc := md5.New
	if sf.SignMethod == HMACSHA1 {
		hashFunc = sha1.New
	}

	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
		sf.ClientID, sf.DeviceName, sf.ProductKey, sf.Timestamp)
	h := hmac.New(hashFunc, []byte(deviceSecret))
	if _, err := h.Write([]byte(signSource)); err != nil {
		return err
	}

	sf.Sign = hex.EncodeToString(h.Sum(nil))
	return nil
}
