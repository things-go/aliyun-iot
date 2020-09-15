package dm

import (
	"strings"

	"github.com/thinkgos/aliyun-iot/uri"
)

// ProcRRPCRequest 处理RRPC请求
// 下行
// request: /sys/${YourProductKey}/${YourDeviceName}/rrpc/request/${messageId}
// response: /sys/${YourProductKey}/${YourDeviceName}/rrpc/response/${messageId}
// subscribe: /sys/${YourProductKey}/${YourDeviceName}/rrpc/request/+
func ProcRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}

	pk, dn := uris[1], uris[2]
	messageID := uris[5]
	c.log.Debugf("rrpc.request.%s", messageID)
	return c.cb.RRPCRequest(c, messageID, pk, dn, payload)
}

// ProcExtRRPCRequest 处理扩展RRPC请求
// 下行
// ${topic} 不为空,设备建立要求clientID传ext = 1
// request:   /ext/rrpc/${messageId}/${topic}
// response:  /ext/rrpc/${messageId}/${topic}
// subscribe: /ext/rrpc/+/${topic}
func ProcExtRRPCRequest(c *Client, rawURI string, payload []byte) error {
	uris := strings.SplitN(strings.TrimLeft(rawURI, uri.Sep), uri.Sep, 4)
	if len(uris) < 3 {
		return ErrInvalidParameter
	}
	messageID, topic := uris[2], uris[3]
	c.log.Debugf("ext.rrpc.%s -- topic: %s", messageID, topic)
	return c.cb.ExtRRPCRequest(c, messageID, topic, payload)
}
