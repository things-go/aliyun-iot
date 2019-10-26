package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// URIServiceItself 获得本设备URI
func (sf *Manager) URIServiceItself(prefix, name string) string {
	return URIService(prefix, name, sf.opt.productKey, sf.opt.deviceName)
}

func (sf *Manager) URIService(prefix, name, productKey, deviceName string) string {
	if sf.opt.uriOffset == 1 {
		prefix = URICOAPHTTPPrePrefix + prefix
	}
	return URIService(prefix, name, productKey, deviceName)
}

// 可以采用Sha256,hmacMd5,hmacSha1,hmacSha256
func generateSign(productKey, deviceName, deviceSecret, clientID string, timestamp int64) (string, error) {
	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
		clientID, deviceName, productKey, timestamp)

	/* setup Password */
	h := hmac.New(sha1.New, []byte(deviceSecret))
	if _, err := h.Write([]byte(signSource)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
