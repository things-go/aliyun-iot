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

// URIService 生成URI, inName的作用是有一些需要格式化到name里.
func (sf *Client) URIService(prefix, name, productKey, deviceName string, inName ...string) string {
	if sf.cfg.uriOffset == 1 {
		prefix = URICOAPHTTPPrePrefix + prefix
	}
	return URIService(prefix, name, productKey, deviceName, inName...)
}

func (sf *Client) _URIExtRRPCService(prefix, messageID, topic string) string {
	if sf.cfg.uriOffset == 1 {
		return URICOAPHTTPPrePrefix + fmt.Sprintf(prefix, messageID) + topic
	}
	return fmt.Sprintf(prefix, messageID) + topic
}

// URIExtRRPCService 生成Ext RRPC URI
func (sf *Client) URIExtRRPCService(messageID, topic string) string {
	return sf._URIExtRRPCService(URIExtRRPCPrefix, messageID, topic)
}

// URIExtRRPCWildcardService 生成Ext RRPC 订阅 URI
func (sf *Client) URIExtRRPCWildcardService(topic string) string {
	return sf._URIExtRRPCService(URIExtRRPCPrefix, "+", topic)
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
