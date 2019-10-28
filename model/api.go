package model

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/thinkgos/cache-go"
)

type MsgType byte

const (
	MsgTypeModelUpRaw            MsgType = iota //!< post raw data to cloud
	MsgTypeEventPropertyPost                    //!< post property value to cloud
	MsgTypeEventPost                            //!< post event identifies value to cloud
	MsgTypeDeviceInfoUpdate                     //!< post device info update message to cloud
	MsgTypeDeviceInfoDelete                     //!< post device info delete message to cloud
	MsgTypeDesiredPropertyGet                   //!< get a device's desired property
	MsgTypeDesiredPropertyDelete                //!< delete a device's desired property
	MsgTypeDsltemplateGet                       //<! get a device's dsltemplate
	MsgTypeDynamictslGet
	MsgTypeConfigGet
	MsgTypeExtErrorRequest

	MsgTypeSubDevLogin                 //!< only for slave device, send login request to cloud
	MsgTypeSubDevLogout                //!< only for slave device, send logout request to cloud
	MsgTypeSubDevDeleteTopo            //!< only for slave device, send delete topo request to cloud
	MsgTypeQueryTimestamp              //!< query ntp time from cloud
	MsgTypeQueryTopoList               //!< only for master device, query topo list
	MsgTypeQueryFOTAData               //!< only for master device, qurey firmware ota data
	MsgTypeQueryCOTAData               //!< only for master device, qurey config ota data
	MsgTypeRequestCOTA                 //!< only for master device, request config ota data from cloud
	MsgTypeRequestFOTAImage            //!< only for master device, request fota image from cloud
	MsgTypeReportSubDevFirmwareVersion //!< report subdev's firmware version

)

// Request 请求
type Request struct {
	ID      int         `json:"id,string"`
	Version string      `json:"version"`
	Params  interface{} `json:"params"`
	Method  string      `json:"method"`
}

// Response 应答
type Response struct {
	ID      int             `json:"id,string"`
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message,omitempty"`
}

// Manager 管理
type Manager struct {
	requestID int32

	cfg Config

	*devMgr
	msgCache *cache.Cache
	pool     *pool
	Conn
	gwUserProc  GatewayUserProc
	devUserProc DevUserProc
}

// New 创建一个物管理
func New(opt *Config) *Manager {
	sf := &Manager{
		cfg:         *opt,
		devMgr:      newDevMgr(),
		gwUserProc:  GwNopUserProc{},
		devUserProc: DevNopUserProc{},
	}
	if opt.enableCache {
		sf.pool = newPool()
		sf.msgCache = cache.New(time.Second*10, time.Second*30)
	}
	sf.CacheInit()
	err := sf.insert(DevLocal, DevTypeSingle, opt.productKey, opt.deviceName, opt.deviceSecret)
	if err != nil {
		panic(fmt.Sprintf("device local duplicate,cause: ", err))
	}
	return sf
}

func (sf *Manager) Connect() error {
	var devType DevType
	if sf.cfg.hasGateway {
		devType = DevTypeGateway
	} else {
		devType = DevTypeSingle
	}
	return sf.SubscribeAllTopic(devType, sf.cfg.productKey, sf.cfg.deviceName)
}

//
//func (sf *Manager) NewSubDevice(devType int, info *MetaInfo) (int, error) {
//	if !sf.cfg.hasGateway {
//		return 0, ErrFeatureNotSupport
//	}
//	return sf.Create(DevTypeSubDev, info.ProductKey, info.DeviceName, info.DeviceSecret)
//}
//
//func (sf *Manager) SubDeviceConnect(id int) {
//
//}

// SetConn 设置连接接口
func (sf *Manager) SetConn(conn Conn) *Manager {
	sf.Conn = conn
	return sf
}

func (sf *Manager) SetGwUserProc(proc GatewayUserProc) *Manager {
	sf.gwUserProc = proc
	return sf
}

func (sf *Manager) SetDevUserProc(proc DevUserProc) *Manager {
	sf.devUserProc = proc
	return sf
}

// RequestID 获得下一个requestID
func (sf *Manager) RequestID() int {
	return int(atomic.AddInt32(&sf.requestID, 1))
}

// SendRequest 发送请求
// uriService 唯一定位服务器或(topic)
// requestID: 请求ID
// method: 方法
// params: 消息体
// API内部已实现json序列化
func (sf *Manager) SendRequest(uriService string, requestID int, method string, params interface{}) error {
	out, err := json.Marshal(&Request{requestID, Version, params, method})
	if err != nil {
		return err
	}
	return sf.Publish(uriService, 1, out)
}

func (sf *Manager) SendResponse(uriService string, id int, code int, data interface{}) error {
	out, err := json.Marshal(struct {
		*Response
		Data interface{} `json:"data"`
	}{
		&Response{
			ID:   id,
			Code: code,
		},
		data,
	})
	if err != nil {
		return err
	}
	return sf.Publish(uriService, 1, out)
}

func (sf *Manager) AlinkReport(msgType MsgType, devID int, payload interface{}) error {
	switch msgType {
	case MsgTypeEventPropertyPost:
		return sf.UpstreamThingEventPropertyPost(devID, payload)
	case MsgTypeDeviceInfoUpdate:
		return sf.UpstreamThingDeviceInfoUpdate(devID, payload)
	case MsgTypeDeviceInfoDelete:
		return sf.UpstreamThingDeviceInfoDelete(devID, payload)
	case MsgTypeModelUpRaw:
		return sf.UpstreamThingModelUpRaw(devID, payload)
	case MsgTypeSubDevLogin:
		// TODO
	case MsgTypeSubDevLogout:
		//TODO
	case MsgTypeSubDevDeleteTopo:
		// todo
	case MsgTypeReportSubDevFirmwareVersion:
		// TODO
	case MsgTypeDesiredPropertyGet:
		// TODO
	case MsgTypeDesiredPropertyDelete:
		// TODO

	}
	return ErrNotSupportMsgType
}

func (sf *Manager) AlinkQuery(msgType MsgType, devID int, payload interface{}) error {
	switch msgType {
	case MsgTypeQueryTimestamp:
		return sf.UpstreamExtNtpRequest()
	case MsgTypeQueryTopoList:
	case MsgTypeQueryCOTAData:
	case MsgTypeQueryFOTAData:
	case MsgTypeRequestCOTA:
	case MsgTypeRequestFOTAImage:
	}
	return ErrNotSupportMsgType
}

func (sf *Manager) AlinkTriggerEvent(devID int, eventID string, payload interface{}) error {
	return sf.UpstreamThingEventPost(devID, eventID, payload)
}
