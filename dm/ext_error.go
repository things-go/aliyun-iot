package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// ExtErrorResponse 子设备错误回复
type ExtErrorResponse struct {
	ID      uint           `json:"id,string"`
	Code    int            `json:"code"`
	Data    infra.MetaPair `json:"data"`
	Message string         `json:"message,omitempty"`
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
	c.Log.Debugf("ext.error.response @%d", rsp.ID)

	pk, dn := rsp.Data.ProductKey, rsp.Data.DeviceName
	// if rsp.Code == infra.CodeSubDevSessionError {
	// _, _ = c.ExtCombineLogin(pk, dn) // TODO: ...
	// }
	return c.gwCb.ExtErrorResponse(c, err, pk, dn)
}
