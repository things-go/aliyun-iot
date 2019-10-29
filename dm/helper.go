package dm

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// URIServiceSelf 获得本设备URI
func (sf *Client) URIServiceSelf(prefix, name string) string {
	return URIService(prefix, name, sf.cfg.productKey, sf.cfg.deviceName)
}

func (sf *Client) URIService(prefix, name, productKey, deviceName string, inName ...string) string {
	if sf.cfg.uriOffset == 1 {
		prefix = URICOAPHTTPPrePrefix + prefix
	}
	return URIService(prefix, name, productKey, deviceName, inName...)
}

func (sf *Client) URIExtRRPCService(prefix, messageID, topic string) string {
	if sf.cfg.uriOffset == 1 {
		return URICOAPHTTPPrePrefix + prefix + messageID + topic
	}
	return prefix + messageID + topic
}

func (sf *Client) URIExtRRPCWildcardService(topic string) string {
	return sf.URIExtRRPCService(URIExtRRPCPrefix, "+", topic)
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
