// Package dynamic 实现动态注册,只限直连设备动态注册
package dynamic

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"time"

	"github.com/thinkgos/aliIOT/infra"
)

// sign method 签名方法
const (
	SignMethodSHA256 = "hmacsha256"
	SignMethodSHA1   = "hmacsha1"
	SignMethodMD5    = "hmacmd5"
)

// MetaInfo 产品与设备三元组
type MetaInfo struct {
	ProductKey    string
	ProductSecret string
	DeviceName    string
	DeviceSecret  string
	CustomDomain  string // 如果使用CloudRegionCustom,需要定义此字段
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

// Register2Cloud 动态注册,传入三元组,获得DeviceSecret,直接修改meta,指定签名算法,未设置或错误,将采用默认sha256
func Register2Cloud(meta *MetaInfo, region infra.CloudRegion, signMethod ...string) error {
	var domain string

	if meta.ProductKey == "" || meta.ProductSecret == "" ||
		meta.DeviceName == "" {
		return errors.New("invalid params")
	}
	signMd := append(signMethod, SignMethodSHA256)[0]
	// 计算签名 Signature
	random, sign, err := calcDynregSign(meta, signMd)
	if err != nil {
		return err
	}

	if region == infra.CloudRegionCustom {
		if meta.CustomDomain == "" {
			return errors.New("custom domain invalid")
		}
		domain = meta.CustomDomain
	} else {
		domain = infra.HTTPCloudDomain[region]
	}

	requestBody := fmt.Sprintf("productKey=%s&deviceName=%s&random=%s&sign=%s&signMethod=%s",
		meta.ProductKey, meta.DeviceName, random, sign, signMd)

	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/auth/register/device", domain),
		bytes.NewBufferString(requestBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "text/xml,text/javascript,text/html,application/json")

	response, err := (&http.Client{Timeout: time.Millisecond * 2000}).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responsePayload := Response{}
	if err = json.NewDecoder(response.Body).Decode(&responsePayload); err != nil {
		return err
	}
	if responsePayload.Code != 200 {
		return errors.New("got response but payload failed")
	}
	meta.DeviceSecret = responsePayload.Data.DeviceSecret
	return nil
}

// 计算动态签名,以productKey为key
func calcDynregSign(info *MetaInfo, signMethod string) (random, sign string, err error) {
	var h hash.Hash

	random = "8Ygb7ULYh53B6OA"
	signSource := fmt.Sprintf("deviceName%sproductKey%srandom%s", info.DeviceName, info.ProductKey, random)

	/* setup password */
	switch signMethod {
	case SignMethodSHA1:
		h = hmac.New(sha1.New, []byte(info.ProductSecret))
	case SignMethodMD5:
		h = hmac.New(md5.New, []byte(info.ProductSecret))
	default: // SignMethodSHA256
		h = hmac.New(sha256.New, []byte(info.ProductSecret))
		sign = SignMethodSHA256
	}

	if _, err = h.Write([]byte(signSource)); err != nil {
		return
	}
	sign = hex.EncodeToString(h.Sum(nil))
	return
}
