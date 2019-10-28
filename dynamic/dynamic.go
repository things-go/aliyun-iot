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
)

// sign method 动态注册只支持以下签名方法
const (
	signMethodSHA256 = "hmacsha256"
	signMethodSHA1   = "hmacsha1"
	signMethodMD5    = "hmacmd5"
)

// HTTPCloudDomain http 域名
var HTTPCloudDomain = []string{
	"iot-auth.cn-shanghai.aliyuncs.com",    /* Shanghai */
	"iot-auth.ap-southeast-1.aliyuncs.com", /* Singapore */
	"iot-auth.ap-northeast-1.aliyuncs.com", /* Japan */
	"iot-auth.us-west-1.aliyuncs.com",      /* America */
	"iot-auth.eu-central-1.aliyuncs.com",   /* Germany */
}

// CloudRegion HTPP云端地域
type CloudRegion byte

// 云平台地域定义CloudRegionRegion
const (
	CloudRegionShangHai CloudRegion = iota
	CloudRegionSingapore
	CloudRegionJapan
	CloudRegionAmerica
	CloudRegionGermany
	CloudRegionCustom
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

// Register2Cloud 动态注册,传入三元组,获得DeviceSecret,直接修改meta,
// 指定签名算法,支持hmacmd5,hmacsha1,hmacsha256将采用默认hmacsha256加签算法
func Register2Cloud(meta *MetaInfo, region CloudRegion, signMethod ...string) error {
	if meta == nil || meta.ProductKey == "" ||
		meta.ProductSecret == "" || meta.DeviceName == "" {
		return errors.New("invalid params")
	}

	signMd := append(signMethod, signMethodSHA256)[0]
	if !(signMd == signMethodMD5 || signMd == signMethodSHA1 || (signMd == signMethodSHA256)) {
		signMd = signMethodSHA256
	}
	// 计算签名 Signature
	random, sign, err := calcDynregSign(meta, signMd)
	if err != nil {
		return err
	}

	var domain string
	if region == CloudRegionCustom {
		if meta.CustomDomain == "" {
			return errors.New("custom domain invalid")
		}
		domain = meta.CustomDomain
	} else {
		domain = HTTPCloudDomain[region]
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
	// TODO: 根据不同的code返回不同的错误
	if responsePayload.Code != 200 {
		return fmt.Errorf("got response but payload failed, %#v", responsePayload)
	}
	meta.DeviceSecret = responsePayload.Data.DeviceSecret
	return nil
}

// calcDynregSign 计算动态签名,以productKey为key
func calcDynregSign(info *MetaInfo, signMethod string) (random, sign string, err error) {
	var h hash.Hash

	/* setup password */
	switch signMethod {
	case signMethodSHA1:
		h = hmac.New(sha1.New, []byte(info.ProductSecret))
	case signMethodMD5:
		h = hmac.New(md5.New, []byte(info.ProductSecret))
	default: // signMethodSHA256
		h = hmac.New(sha256.New, []byte(info.ProductSecret))
	}

	random = "8Ygb7ULYh53B6OA"
	signSource := fmt.Sprintf("deviceName%sproductKey%srandom%s", info.DeviceName, info.ProductKey, random)

	if _, err = h.Write([]byte(signSource)); err != nil {
		return
	}
	sign = hex.EncodeToString(h.Sum(nil))
	return
}
