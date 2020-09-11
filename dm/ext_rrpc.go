package dm

import (
	"strings"
)

// ProcRRPCRequest 处理RRPC请求
// 下行
// request: /sys/${YourProductKey}/${YourDeviceName}/rrpc/request/${messageId}
// response: /sys/${YourProductKey}/${YourDeviceName}/rrpc/response/${messageId}
// subscribe: /sys/${YourProductKey}/${YourDeviceName}/rrpc/request/+
func ProcRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}
	messageID := uris[c.uriOffset+5]
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.log.Debugf("downstream sys <RRPC>: request - messageID: %s", messageID)
	return c.cb.RRPCRequest(c, messageID, pk, dn, payload)
}

// ProcExtRRPCRequest 处理扩展RRPC请求
// 下行
// ${topic} 不为空,设备建立要求clientID传ext = 1
// request:   /ext/rrpc/${messageId}/${topic}
// response:  /ext/rrpc/${messageId}/${topic}
// subscribe: /ext/rrpc/+/${topic}
func ProcExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := strings.SplitN(strings.TrimLeft(rawURI, SEP), SEP, c.uriOffset+4)
	if len(uris) < (c.uriOffset + 3) {
		return ErrInvalidParameter
	}
	c.log.Debugf("downstream extend <RRPC>: Request - URI: ", rawURI)
	messageID, topic := uris[c.uriOffset+2], uris[c.uriOffset+3]
	return c.cb.ExtRRPCRequest(c, messageID, topic, payload)
}
