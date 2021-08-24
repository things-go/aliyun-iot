// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package ahttp 实现http client 上传数据. 授权方式为自动调用授权,可手动调用,也可以直接调用发送数据接口
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"

	aiot "github.com/things-go/aliyun-iot"
	"github.com/things-go/aliyun-iot/infra"
	"github.com/things-go/aliyun-iot/logger"
	"github.com/things-go/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/58034.html?spm=a2c4g.11186623.6.609.54316764YJj5MQ

// 错误码
const (
	CodeSuccess              = 0
	CodeUnknown              = 10000 // 未知错误
	CodeParamException       = 10001 // 请求的参数异常
	CodeAuthFailed           = 20000 // 设备鉴权失败
	CodeTokenExpired         = 20001 // 需重新调用auth进行鉴权,获取token
	CodeTokenIsNull          = 20002 // 需重新调用auth进行鉴权,获取token
	CodeTokenCheckFailed     = 20003 // 根据token获取identify信息失败.需重新调用auth进行鉴权,获取token
	CodeUpdateSessionFailed  = 20004 // 更新失败
	CodePublishMessageFailed = 30001 //
	CodeRequestTooMany       = 40000 // 请求次数过多，流控限制.
)

// Sign method
const (
	hmacsha1 = "hmacsha1"
	hmacmd5  = "hmacmd5"
)

// AuthRequest 鉴权请求
type AuthRequest struct {
	Version    string `json:"version"`    // 版本号, 不传入此参数,则默认 default.
	ClientID   string `json:"clientId"`   // 长度为64字符内,建议以MAC地址或SN码作为clientId. 目前productKey.deviceName
	SignMethod string `json:"signmethod"` // 签名方法: 支持hmacsha1, hmacmd5(默认)
	Sign       string `json:"sign"`       // 签名: hmacmd5(DeviceSecret,content)
	ProductKey string `json:"productKey"` // 设备所属的产品
	DeviceName string `json:"deviceName"` // 设备名称
	// 校验时间戳15分钟内的请求有效.时间戳格式为数值,单位: ms
	// 值为自GMT 1970年1月1日0时0分到当前时间点所经过的毫秒数.
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
	triad infra.MetaTriad

	endpoint   string
	version    string
	signMethod string

	token atomic.Value
	group singleflight.Group

	httpc *http.Client
	log   logger.Logger
}

var _ aiot.Conn = (*Client)(nil)

// New 新建alink http client
// 默认加签算法: hmacmd5(目前支持 hmacsha1, hmacmd5(默认))
// 默认host: https://iot-as-http.cn-shanghai.aliyuncs.com
// 默认使用 http.DefaultClient
func New(meta infra.MetaTriad, opts ...Option) *Client {
	c := &Client{
		triad:      meta,
		endpoint:   "https://iot-as-http.cn-shanghai.aliyuncs.com",
		version:    "default",
		signMethod: hmacmd5,
		httpc:      http.DefaultClient,
		log:        logger.NewDiscard(),
	}
	c.token.Store("")
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// 鉴权
func (sf *Client) getToken() (string, error) {
	if token := sf.token.Load().(string); token != "" {
		return token, nil
	}
	return sf.refreshToken()
}

func (sf *Client) refreshToken() (string, error) {
	if sf.triad.ProductKey == "" || sf.triad.DeviceName == "" || sf.triad.DeviceSecret == "" {
		return "", errors.New("invalid device meta triad")
	}
	tk, err, _ := sf.group.Do("auth", func() (interface{}, error) {
		sf.token.Store("") // 置token无效, 让新进来的直接进入刷新token

		// 生成body加签
		signMethod := sf.signMethod
		switch signMethod {
		case hmacmd5, hmacsha1:
		default:
			signMethod = hmacmd5
		}
		timestamp := infra.Millisecond(time.Now())
		clientID, sign := infra.CalcSign(signMethod, sf.triad, timestamp)

		b, err := json.Marshal(&AuthRequest{
			sf.version,
			clientID,
			signMethod,
			sign,
			sf.triad.ProductKey,
			sf.triad.DeviceName,
			timestamp,
		})
		if err != nil {
			return "", err
		}

		request, err := http.NewRequestWithContext(context.Background(),
			http.MethodPost, sf.endpoint+"/auth", bytes.NewBuffer(b))
		if err != nil {
			return "", err
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := sf.httpc.Do(request)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()

		authRsp := &AuthResponse{}
		err = json.NewDecoder(response.Body).Decode(authRsp)
		if err != nil {
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

// Publish push message, payload support []byte and string
func (sf *Client) Publish(_uri string, _ byte, payload interface{}) error {
	py := &DataResponse{}
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

		request, err := http.NewRequestWithContext(context.Background(),
			http.MethodPost, sf.endpoint+uri.TopicPrefix+_uri, buf)
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

		if err := json.NewDecoder(response.Body).Decode(py); err != nil {
			return err
		}
		sf.log.Debugf("publish response, %+v", py)
		if py.Code == 0 {
			return nil
		}

		if !(py.Code == CodeTokenExpired ||
			py.Code == CodeTokenCheckFailed ||
			py.Code == CodeTokenIsNull) {
			return infra.NewCodeError(py.Code, py.Message)
		}
		sf.refreshToken() // nolint: errcheck
	}
	return infra.NewCodeError(py.Code, py.Message)
}

// Subscribe 实现dm.Conn接口
func (*Client) Subscribe(string, aiot.ProcDownStream) error { return nil }

// UnSubscribe 实现dm.Conn接口
func (*Client) UnSubscribe(...string) error { return nil }

// Close 实现dm.Conn接口
func (sf *Client) Close() error { return nil }
