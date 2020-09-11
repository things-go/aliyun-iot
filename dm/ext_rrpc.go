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
	c.debugf("downstream sys <RRPC>: request - messageID: %s", messageID)

	return c.ipcSendMessage(&ipcMessage{
		evt:        ipcEvtRRPCRequest,
		productKey: uris[c.uriOffset+1],
		deviceName: uris[c.uriOffset+2],
		payload:    payload,
		extend:     messageID,
	})
}

// ProcExtRRPCRequest 处理扩展RRPC请求
// 下行
// ${topic} 不为空,设备建立要求clientID传ext = 1
// request: /ext/rrpc/${messageId}/${topic}
// response: /ext/rrpc/${messageId}/${topic}
// subscribe: /ext/rrpc/+/${topic}
func ProcExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := strings.SplitN(strings.TrimLeft(rawURI, SEP), SEP, c.uriOffset+3)
	if len(uris) < (c.uriOffset + 3) {
		return ErrInvalidParameter
	}

	c.debugf("downstream extend <RRPC>: Request - URI: ", rawURI)
	return c.ipcSendMessage(&ipcMessage{
		evt:     ipcEvtExtRRPCRequest,
		extend:  uris[c.uriOffset+2], // ${messageId}/${topic}
		payload: payload,
	})
}
