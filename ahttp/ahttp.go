// Package ahttp 实现http client 上传数据. 授权方式为自动调用授权,可手动调用,也可以直接调用发送数据接口
package ahttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/thinkgos/go-core-package/lib/algo"
	"golang.org/x/sync/singleflight"

	"github.com/thinkgos/aliyun-iot/clog"
	"github.com/thinkgos/aliyun-iot/infra"
)

// 错误码
const (
	CodeSuccess              = 0
	CodeUnknown              = 10000
	CodeParamException       = 10001
	CodeAuthFailed           = 20000
	CodeTokenExpired         = 20001 // 需重新调用auth进行鉴权，获取token
	CodeTokenIsNull          = 20002 // 需重新调用auth进行鉴权，获取token
	CodeTokenCheckFailed     = 20003 // 根据token获取identify信息失败。需重新调用auth进行鉴权，获取token
	CodeUpdateSessionFailed  = 20004
	CodePublishMessageFailed = 30001
	CodeRequestTooMany       = 40000
)

// Sign method
const (
	hmacsha1 = "hmacsha1"
	hmacmd5  = "hmacmd5"
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
	meta infra.MetaInfo

	host       string
	version    string
	signMethod string

	token atomic.Value
	group singleflight.Group

	httpc *http.Client
	log   clog.LogProvider
}

// New 新建alink http client
// 默认加签算法: hmacmd5
// 默认host: https://iot-as-http.cn-shanghai.aliyuncs.com
// 默认使用 http.DefaultClient
func New(meta infra.MetaInfo, opts ...Option) *Client {
	c := &Client{
		meta:       meta,
		host:       "https://iot-as-http.cn-shanghai.aliyuncs.com",
		version:    "default",
		signMethod: hmacmd5,
		httpc:      http.DefaultClient,
		log:        clog.NewDiscard(),
	}
	c.token.Store("")
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// 鉴权
func (sf *Client) getToken() (string, error) {
	if sf.meta.ProductKey == "" || sf.meta.DeviceName == "" || sf.meta.DeviceSecret == "" {
		return "", errors.New("invalid device meta info")
	}

	if token := sf.token.Load().(string); token != "" {
		return token, nil
	}

	tk, err, _ := sf.group.Do("auth", func() (interface{}, error) {
		// 生成body加签
		method := algo.MethodMD5
		if sf.signMethod == hmacsha1 {
			method = algo.MethodSha1
		}
		clientID, tm := sf.meta.ProductKey+"."+sf.meta.DeviceName, time.Now().Unix()*1000
		signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
			clientID, sf.meta.DeviceName, sf.meta.ProductKey, tm)
		authReq := AuthRequest{
			sf.version,
			clientID,
			sf.signMethod,
			algo.Hmac(method, signSource, sf.meta.DeviceSecret),
			sf.meta.ProductKey,
			sf.meta.DeviceName,
			tm,
		}

		b, err := json.Marshal(&authReq)
		if err != nil {
			return "", err
		}

		request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, sf.host+"/auth", bytes.NewBuffer(b))
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
		if err := json.NewDecoder(response.Body).Decode(&authRsp); err != nil {
			return "", err
		}

		if authRsp.Code != CodeSuccess {
			return "", infra.NewCodeError(authRsp.Code, authRsp.Message)
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

// Publish push message
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
			return errors.New("unknown payload type, must be string or []byte")
		}

		request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, sf.host+uri, buf)
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

		if err := json.NewDecoder(response.Body).Decode(&py); err != nil {
			return err
		}
		sf.log.Debugf("publish response, %+v", py)
		if py.Code == 0 {
			return nil
		}

		if !(py.Code == CodeTokenExpired ||
			py.Code == CodeTokenCheckFailed ||
			py.Code == CodeTokenIsNull) {
			sf.token.Store("")
			return infra.NewCodeError(py.Code, py.Message)
		}
	}
	return infra.NewCodeError(py.Code, py.Message)
}
