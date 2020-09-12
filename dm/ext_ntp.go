package dm

import (
	"encoding/json"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
)

// NtpRequest ntp请求payload
type NtpRequest struct {
	DeviceSendTime int64 `json:"deviceSendTime,string"`
}

// NtpResponse ntp回复payload
type NtpResponse struct {
	DeviceSendTime int64 `json:"deviceSendTime,string"` // 设备发送时间,单位ms
	ServerRecvTime int64 `json:"serverRecvTime,string"` // 平台接收时间,单位ms
	ServerSendTime int64 `json:"serverSendTime,string"` // 平台发送时间,单位ms
}

// ExtNtpRequest ntp请求
// 发送一条Qos = 0的消息,并带上设备当前的时间戳,平台将回复 设备的发送时间,平台的接收时间, 平台的发送时间.
// 设备计算当前精确时间 = (平台接收时间 + 平台发送时间 + 设备接收时间 - 设备发送时间) / 2
// 请求Topic：/ext/ntp/${YourProductKey}/${YourDeviceName}/request
// 响应Topic：/ext/ntp/${YourProductKey}/${YourDeviceName}/response
func (sf *Client) ExtNtpRequest() error {
	if !sf.hasNTP || sf.hasRawModel {
		return ErrNotSupportFeature
	}
	err := sf.Publish(sf.URIServiceSelf(infra.URIExtNtpPrefix, infra.URINtpRequest), 0,
		NtpRequest{int64(time.Now().Nanosecond()) / 1000000})
	if err != nil {
		return err
	}
	sf.log.Debugf("upstream ext <ntp>: request")
	return nil
}

// ProcExtNtpResponse 处理ntp请求的应答
// 上行
// request: /ext/ntp/${YourProductKey}/${YourDeviceName}/request
// response: /ext/ntp/${YourProductKey}/${YourDeviceName}/response
// subscribe: /ext/ntp/${YourProductKey}/${YourDeviceName}/response
func ProcExtNtpResponse(c *Client, rawURI string, payload []byte) error {
	uris := infra.SpiltURI(rawURI)
	if len(uris) < (c.uriOffset + 5) {
		return ErrInvalidURI
	}
	rsp := NtpResponse{}
	if err := json.Unmarshal(payload, &rsp); err != nil {
		return err
	}
	c.log.Debugf("downstream extend <ntp>: response - %+v", rsp)
	return c.cb.ExtNtpResponse(c, uris[c.uriOffset+2], uris[c.uriOffset+3], rsp)
}
