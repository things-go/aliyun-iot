// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import (
	"encoding/json"

	"github.com/things-go/aliyun-iot/infra"
	"github.com/things-go/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/120329.html?spm=5176.11065259.1996646101.searchclickresult.3703243801uSYu

// ExtErrorResponse 子设备错误回复
type ExtErrorResponse struct {
	ID      uint           `json:"id,string"`
	Code    int            `json:"code"`
	Data    infra.MetaPair `json:"data"`
	Message string         `json:"message,omitempty"`
}

// ProcExtErrorResponse 处理错误的回复,仅与子设备
// response:  ext/error/{productKey}/{deviceName}
// subscribe: ext/error/{productKey}/{deviceName}
func ProcExtErrorResponse(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 4 {
		return ErrInvalidURI
	}

	rsp := &ExtErrorResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.signalPending(Message{rsp.ID, nil, err})
	c.Log.Debugf("ext.error.response @%d", rsp.ID)

	pk, dn := rsp.Data.ProductKey, rsp.Data.DeviceName
	return c.gwCb.ExtErrorResponse(c, err, pk, dn)
}
