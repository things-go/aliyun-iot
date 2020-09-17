package dm

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/thinkgos/aliyun-iot/uri"
)

// URIGateway 获得本设备网关URI
func (sf *Client) URIGateway(prefix, name string) string {
	return uri.URI(prefix, name, sf.tetrad.ProductKey, sf.tetrad.DeviceName)
}

// ClientID to client id like {pk}.{dn}
func ClientID(pk, dn string) string {
	return pk + "." + dn
}

// 可以采用Sha256,hmacMd5,hmacSha1,hmacSha256
func generateSign(productKey, deviceName, deviceSecret, clientID string, timestamp int64) (string, error) {
	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
		clientID, deviceName, productKey, timestamp)

	// setup Password
	h := hmac.New(sha1.New, []byte(deviceSecret))
	if _, err := h.Write([]byte(signSource)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func cloneJSONRawMessage(d json.RawMessage) json.RawMessage {
	v := make(json.RawMessage, len(d))
	copy(v, d)
	return v
}
