package dm

import (
	"encoding/json"
	"fmt"
	"time"
)

// NtpResponse ntp回复payload
type NtpResponse struct {
	DeviceSendTime int64 `json:"deviceSendTime,string"` // 设备发送时间,单位ms
	ServerRecvTime int64 `json:"serverRecvTime,string"` // 平台接收时间,单位ms
	ServerSendTime int64 `json:"serverSendTime,string"` // 平台发送时间,单位ms
}

// upstreamExtNtpRequest ntp请求
// 发送一条Qos = 0的消息,并带上设备当前的时间戳,平台将回复 设备的发送时间,平台的接收时间, 平台的发送时间.
// 设备计算当前精确时间 = (平台接收时间 + 平台发送时间 + 设备接收时间 - 设备发送时间) / 2
// 请求Topic：/ext/ntp/${YourProductKey}/${YourDeviceName}/request
// 响应Topic：/ext/ntp/${YourProductKey}/${YourDeviceName}/response
func (sf *Client) upstreamExtNtpRequest() error {
	err := sf.Publish(sf.URIServiceSelf(URIExtNtpPrefix, URINtpRequest),
		0, fmt.Sprintf(`{"deviceSendTime":"%d"}`, time.Now().Nanosecond()/1000000))
	if err != nil {
		return err
	}
	sf.debugf("upstream ext <ntp>: request")
	return nil
}

// ProcExtNtpResponse 处理ntp请求的应答
// 上行
// request: /ext/ntp/${YourProductKey}/${YourDeviceName}/request
// response: /ext/ntp/${YourProductKey}/${YourDeviceName}/response
// subscribe: /ext/ntp/${YourProductKey}/${YourDeviceName}/response
func ProcExtNtpResponse(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 5) {
		return ErrInvalidURI
	}
	rsp := NtpResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.debugf("downstream extend <ntp>: response - %+v", rsp)
	return c.eventProc.EvtExtNtpResponse(c, uris[c.uriOffset+2], uris[c.uriOffset+3], rsp)
}
