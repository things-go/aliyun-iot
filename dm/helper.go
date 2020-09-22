package dm

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"

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

// 可以采用hmacmd5,hmacsha1,hmacsha256
func generateSign(method, productKey, deviceName, deviceSecret, clientID string, timestamp int64) (string, error) {
	var f func() hash.Hash

	switch method {
	case "hmacmd5":
		f = md5.New
	case "hmacsha1":
		f = sha1.New
	case "hmacsha256":
		f = sha256.New
	default:
		return "", errors.New("invalid method")
	}

	// setup Password
	h := hmac.New(f, []byte(deviceSecret))
	source := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
		clientID, deviceName, productKey, timestamp)
	if _, err := h.Write([]byte(source)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func cloneJSONRawMessage(d json.RawMessage) json.RawMessage {
	v := make(json.RawMessage, len(d))
	copy(v, d)
	return v
}
