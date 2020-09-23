package dm

import (
	"encoding/json"
	"fmt"

	"github.com/thinkgos/go-core-package/lib/algo"

	"github.com/thinkgos/aliyun-iot/infra"
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

// 可以采用hmacmd5,hmacsha1,hmacsha256,返回clientID和加签后的值
func generateSign(method string, meta infra.MetaTriad, timestamp int64) (string, string) {
	clientID := ClientID(meta.ProductKey, meta.DeviceName)
	source := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
		clientID, meta.DeviceName, meta.ProductKey, timestamp)
	return clientID, algo.Hmac(method, meta.DeviceSecret, source)
}

func cloneJSONRawMessage(d json.RawMessage) json.RawMessage {
	v := make(json.RawMessage, len(d))
	copy(v, d)
	return v
}
