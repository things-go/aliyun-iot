// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package dynamic 实现动态注册,只限直连设备动态注册,阿里云目前限制激活过的设备不可再注册
package dynamic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/thinkgos/go-core-package/extrand"
	"github.com/thinkgos/go-core-package/extstr"
	"github.com/thinkgos/go-core-package/lib/algo"

	"github.com/thinkgos/aliyun-iot/infra"
)

// sign method 动态注册只支持以下签名方法
const (
	hmacSHA256 = "hmacsha256"
	hmacSHA1   = "hmacsha1"
	hmacMD5    = "hmacmd5"
)

// Option option
type Option func(*Client)

// WithHTTPClient with custom http.Client
func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpc = c
	}
}

// Client dynamic client
type Client struct {
	httpc *http.Client
}

// New new a dynamic client
func New(opts ...Option) *Client {
	c := &Client{
		http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Response 应答
type Response struct {
	Code int `json:"code"`
	Data struct {
		ProductKey   string `json:"productKey"`
		DeviceName   string `json:"deviceName"`
		DeviceSecret string `json:"deviceSecret"`
	} `json:"data"`
	Message string `json:"message"`
}

// RegisterCloud 一型一密动态注册,传入三元组,根据ProductKey,ProductSecret和deviceName获得DeviceSecret,
// meta: 成功将直接修改meta的DeviceSecret
// crd: 指定注册的云端,地址: [https://, http://]host:port/auth/register/device
// signMethods: 可选指定签名算法hmacmd5,hmacsha1,hmacsha256(默认)
// NOTE: 设备联网前，需要在物联网平台预注册设备DeviceName，建议采用设备的MAC地址、IMEI、SN码等作为DeviceName
// @see https://help.aliyun.com/document_detail/89298.html?spm=a2c4g.11186623.6.703.53265bc9vmD6q8
func (sf *Client) RegisterCloud(meta *infra.MetaTetrad, crd infra.CloudRegionDomain, signMethods ...string) error {
	var domain string

	if meta == nil || meta.ProductKey == "" || meta.ProductSecret == "" || meta.DeviceName == "" {
		return errors.New("invalid parameter")
	}

	if crd.Region == infra.CloudRegionCustom {
		if crd.CustomDomain == "" {
			return errors.New("invalid custom domain")
		}
		if !strings.Contains(crd.CustomDomain, "://") {
			domain = "http://" + crd.CustomDomain
		}
	} else {
		domain = "https://" + infra.HTTPCloudDomain[crd.Region]
	}

	requestBody := requestBody(meta, signMethods...)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		domain+"/auth/register/device", bytes.NewBufferString(requestBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "text/xml,text/javascript,text/html,application/json")

	response, err := sf.httpc.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responsePy := &Response{}
	if err := json.NewDecoder(response.Body).Decode(responsePy); err != nil {
		return err
	}

	if responsePy.Code != infra.CodeSuccess {
		return infra.NewCodeError(responsePy.Code, responsePy.Message)
	}
	meta.DeviceSecret = responsePy.Data.DeviceSecret
	return nil
}

func requestBody(meta *infra.MetaTetrad, signMethods ...string) string {
	signMd := hmacSHA256
	if len(signMethods) > 0 {
		signMd = signMethods[0]
	}
	if !extstr.Contains([]string{hmacMD5, hmacSHA1, hmacSHA256}, signMd) {
		signMd = hmacSHA256 // 非法签名使用默认签名方法sha256
	}
	//  "8Ygb7ULYh53B6OA"
	random := extrand.RandString(16)
	// deviceName{deviceName}productKey{productKey}random{random}
	source := "deviceName" + meta.DeviceName + "productKey" + meta.ProductKey + "random" + random
	// 计算签名 Signature
	sign := algo.Hmac(signMd, meta.ProductSecret, source)

	return fmt.Sprintf("productKey=%s&deviceName=%s&random=%s&sign=%s&signMethod=%s",
		meta.ProductKey, meta.DeviceName, random, sign, signMd)
}
