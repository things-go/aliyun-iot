package model

import (
	"fmt"
	"time"

	"github.com/thinkgos/aliIOT/infra"
)

type SubDevTopoAddParams struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	ClientId   string `json:"clientId"`
	Timestamp  int64  `json:"timestamp,string"`
	SignMethod string `json:"signMethod"`
	Sign       string `json:"sign"`
}

func (sf *Manager) UpstreamGwThingTopoAdd(metas ...*MetaInfo) error {
	var clientID string
	var sign string
	var err error

	sublist := make([]SubDevTopoAddParams, 0, len(metas))

	timestamp := time.Now().Unix()
	for _, v := range metas {
		clientID = fmt.Sprintf("%s.%s|_v=%s|", v.ProductKey, v.DeviceName, infra.IOTSDKVersion)
		sign, err = generateSign(v.ProductKey, v.DeviceName, v.DeviceSecret, clientID, timestamp)
		if err != nil {
			return err
		}
		sublist = append(sublist, SubDevTopoAddParams{
			v.ProductKey,
			v.DeviceSecret,
			clientID,
			timestamp,
			infra.SignMethodHMACSHA1,
			sign,
		})
	}

	return sf.SendRequest(sf.URIService(URISysPrefix, URIThingTopoAdd),
		sf.RequestID(), methodTopoAdd, sublist)
}

func (sf *Manager) UpstreamGwThingTopoDelete(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	_, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}

	// TODO
	return nil
}

type GwTopoGetData struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type GwTopoGetResponse struct {
	Response
	Data []GwTopoGetData `json:"data"`
}

func (sf *Manager) UpstreamGwThingTopoGet() error {
	return sf.SendRequest(sf.URIService(URISysPrefix, URIThingTopoGet),
		sf.RequestID(), methodTopoGet, "{}")
}
