package aliIOT

import (
	"github.com/thinkgos/aliIOT/model"
)

type dm_client_uri_map_t struct {
	devType   int
	uriPrefix string
	uriName   string
	proc      model.ProcDownStreamFunc
}

//func (sf *Client) subscribeAll(devType model.DevType, productKey, deviceName string) error {
//
//}
