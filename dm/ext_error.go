package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
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
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 4) {
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

	c.done(rsp.ID, err, nil)
	c.debugf("downstream extend <Error>: response,@%d", rsp.ID)

	pk, dn := rsp.Data.ProductKey, rsp.Data.DeviceName
	if rsp.Code == infra.CodeSubDevSessionError {
		node, err := c.SearchNodeByPkDn(pk, dn)
		if err != nil {
			return err
		}
		_, _ = c.upstreamExtGwCombineLogin(node.ID())
	}
	return c.eventGwProc.EvtExtErrorResponse(c, err, pk, dn)
}
