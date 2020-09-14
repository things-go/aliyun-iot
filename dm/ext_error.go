package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// ExtErrorData 子设备错误回复数据域
type ExtErrorData struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

// ExtErrorResponse 子设备错误回复
type ExtErrorResponse struct {
	ID      uint         `json:"id,string"`
	Code    int          `json:"code"`
	Data    ExtErrorData `json:"data"`
	Message string       `json:"message,omitempty"`
}

// ProcExtErrorResponse 处理错误的回复
// response:  ext/error/{productKey}/{deviceName}
// subscribe: ext/error/{productKey}/{deviceName}
func ProcExtErrorResponse(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 4 {
		return ErrInvalidURI
	}

	rsp := ExtErrorResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, nil, err})
	c.log.Debugf("downstream extend <Error>: response,@%d", rsp.ID)

	pk, dn := rsp.Data.ProductKey, rsp.Data.DeviceName
	if rsp.Code == infra.CodeSubDevSessionError {
		node, err := c.SearchNodeByPkDn(pk, dn)
		if err != nil {
			return err
		}
		_, _ = c.ExtCombineLogin(node.ID())
	}
	return c.gwCb.ExtErrorResponse(c, err, pk, dn)
}
