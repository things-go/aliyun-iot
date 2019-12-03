// Package dynamic 实现动态注册,只限直连设备动态注册,阿里云目前限制激活过的设备不可再注册
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
	"strings"
	"time"

	"github.com/thinkgos/aliIOT/infra"
)

// sign method 动态注册只支持以下签名方法
const (
	signMethodHMACSHA256 = "hmacsha256"
	signMethodHMACSHA1   = "hmacsha1"
	signMethodHMACMD5    = "hmacmd5"
)

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
// 指定签名算法,默认hmacsha256加签算法(支持hmacmd5,hmacsha1,hmacsha256)
func Register2Cloud(meta *infra.MetaInfo, crd infra.CloudRegionDomain, signMethod ...string) error {
	if meta == nil || meta.ProductKey == "" || meta.ProductSecret == "" || meta.DeviceName == "" {
		return errors.New("invalid params")
	}

	signMd := append(signMethod, signMethodHMACSHA256)[0]
	if !(signMd == signMethodHMACMD5 || signMd == signMethodHMACSHA1 || (signMd == signMethodHMACSHA256)) {
		signMd = signMethodHMACSHA256 // 非法签名使用默认签名方法
	}

	ms := MetaSign{
		ProductKey:    meta.ProductKey,
		ProductSecret: meta.ProductSecret,
		DeviceName:    meta.DeviceName,
		Random:        "8Ygb7ULYh53B6OA",
		SignMethod:    signMd,
	}
	// 计算签名 Signature
	sign, err := calcSign(&ms)
	if err != nil {
		return err
	}

	var domain string
	if crd.Region == infra.CloudRegionCustom {
		if crd.CustomDomain == "" {
			return errors.New("custom domain invalid")
		}
		domain = crd.CustomDomain
	} else {
		domain = infra.HTTPCloudDomain[crd.Region]
	}

	if !strings.Contains(domain, "://") {
		domain = "http://" + domain
	}

	requestBody := fmt.Sprintf("productKey=%s&deviceName=%s&random=%s&sign=%s&signMethod=%s",
		meta.ProductKey, meta.DeviceName, ms.Random, sign, signMd)

	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/auth/register/device", domain),
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

	responsePy := Response{}
	if err = json.NewDecoder(response.Body).Decode(&responsePy); err != nil {
		return err
	}

	if responsePy.Code != infra.CodeSuccess {
		return infra.NewCodeError(responsePy.Code, responsePy.Message)
	}
	meta.DeviceSecret = responsePy.Data.DeviceSecret
	return nil
}

// MetaSign 签名
type MetaSign struct {
	ProductKey    string
	ProductSecret string
	DeviceName    string
	Random        string
	SignMethod    string
}

// calcSign 计算动态签名,以productKey为key
func calcSign(info *MetaSign) (string, error) {
	var h hash.Hash

	/* setup password */
	switch info.SignMethod {
	case signMethodHMACSHA1:
		h = hmac.New(sha1.New, []byte(info.ProductSecret))
	case signMethodHMACMD5:
		h = hmac.New(md5.New, []byte(info.ProductSecret))
	case "hmacsha256":
		fallthrough
	case "":
		h = hmac.New(sha256.New, []byte(info.ProductSecret))
	default:
		return "", errors.New("sign method not support")
	}

	signSource := fmt.Sprintf("deviceName%sproductKey%srandom%s", info.DeviceName, info.ProductKey, info.Random)

	if _, err := h.Write([]byte(signSource)); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
