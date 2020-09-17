package dm

import (
	"encoding/json"
	"time"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/102509.html?spm=a2c4g.11186623.6.689.7f67741ai8OqLc

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
	sf.Log.Debugf("ext.ntp.request")
	_uri := sf.URIGateway(uri.ExtNtpPrefix, uri.NtpRequest)
	py, err := json.Marshal(NtpRequest{infra.Millisecond(time.Now())})
	if err != nil {
		return err
	}
	return sf.Publish(_uri, 0, py)
}

// ProcExtNtpResponse 处理ntp请求的应答
// 上行
// request:   /ext/ntp/${YourProductKey}/${YourDeviceName}/request
// response:  /ext/ntp/${YourProductKey}/${YourDeviceName}/response
// subscribe: /ext/ntp/${YourProductKey}/${YourDeviceName}/response
func ProcExtNtpResponse(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 5 {
		return ErrInvalidURI
	}
	rsp := &NtpResponse{}
	if err := json.Unmarshal(payload, rsp); err != nil {
		return err
	}
	tm := (rsp.ServerRecvTime + rsp.ServerSendTime + infra.Millisecond(time.Now()) - rsp.DeviceSendTime) / 2
	exact := infra.Time(tm)
	c.Log.Debugf("ext.ntp.response -- %+v", exact)
	pk, dn := uris[2], uris[3]
	return c.cb.ExtNtpResponse(c, pk, dn, exact)
}
