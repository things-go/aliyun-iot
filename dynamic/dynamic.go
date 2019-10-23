package dynamic

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/thinkgos/aliIOT/infra"
)

// MetaInfo 产品与设备三元组
type MetaInfo struct {
	ProductKey    string
	ProductSecret string
	DeviceName    string
	DeviceSecret  string
	CustomDomain  string // 如果使用CloudRegionCustom,需要定义此字段
}

type Response struct {
	Code int `json:"code"`
	Data struct {
		ProductKey   string `json:"productKey"`
		DeviceName   string `json:"deviceName"`
		DeviceSecret string `json:"deviceSecret"`
	} `json:"data"`
	Message string `json:"message"`
}

// DynamicRegister 动态注册,传入三元组,获得DeviceSecret
func DynamicRegister(meta *MetaInfo, region infra.CloudRegion) error {
	var domain string

	if meta.ProductKey == "" || meta.ProductSecret == "" ||
		meta.DeviceName == "" {
		return errors.New("invalid params")
	}

	// Calcute Signature
	random, sign, err := calcDynregSign(meta)
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

	reqPayload := fmt.Sprintf("productKey=%s&deviceName=%s&random=%s&sign=%s&signMethod=%s",
		meta.ProductKey, meta.DeviceName, random, sign, "hmacsha256")

	request, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("http://%s/auth/register/device", domain),
		bytes.NewBufferString(reqPayload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "text/xml,text/javascript,text/html,application/json")
	client := http.Client{Timeout: time.Millisecond * 1000}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	rspPayload := &Response{}
	if err = json.NewDecoder(response.Body).Decode(rspPayload); err != nil {
		return err
	}
	meta.DeviceSecret = rspPayload.Data.DeviceSecret
	return nil
}

// 计算动态签名
func calcDynregSign(info *MetaInfo) (random, sign string, err error) {
	random = "8Ygb7ULYh53B6OA"
	signSource := fmt.Sprintf("deviceName%sproductKey%srandom%s", info.DeviceName, info.ProductKey, random)
	/* setup password */
	h := hmac.New(sha256.New, []byte(info.ProductSecret))
	if _, err = h.Write([]byte(signSource)); err != nil {
		return
	}
	sign = hex.EncodeToString(h.Sum(nil))
	return
}
