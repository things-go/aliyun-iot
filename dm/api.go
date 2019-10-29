package dm

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
	MsgTypeDesiredPropertyGet                   //!< get a device's desired property
	MsgTypeDesiredPropertyDelete                //!< delete a device's desired property
	MsgTypeDeviceInfoUpdate                     //!< post device info update message to cloud
	MsgTypeDeviceInfoDelete                     //!< post device info delete message to cloud
	MsgTypeDsltemplateGet                       //<! get a device's dsltemplate
	MsgTypeDynamictslGet                        //!< ??
	MsgTypeExtNtpRequest                        //!< query ntp time from cloud
	MsgTypeConfigGet                            //!< 获取配置
	MsgTypeExtErrorRequest

	MsgTypeSubDevLogin                 //!< only for slave device, send login request to cloud
	MsgTypeSubDevLogout                //!< only for slave device, send logout request to cloud
	MsgTypeSubDevDeleteTopo            //!< only for slave device, send delete topo request to cloud
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

// Client 管理
type Client struct {
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
func New(cfg *Config) *Client {
	sf := &Client{
		cfg:         *cfg,
		devMgr:      newDevMgr(),
		gwUserProc:  GwNopUserProc{},
		devUserProc: DevNopUserProc{},
	}
	if cfg.hasCache {
		sf.pool = newPool()
		sf.msgCache = cache.New(time.Second*10, time.Second*30)
	}
	sf.CacheInit()
	err := sf.insert(DevLocal, DevTypeSingle, cfg.productKey, cfg.deviceName, cfg.deviceSecret)
	if err != nil {
		panic(fmt.Sprintf("device local duplicate,cause: %+v", err))
	}
	return sf
}

// Connect 将订阅所有相关主题,主题有config配置
func (sf *Client) Connect() error {
	var devType DevType
	if sf.cfg.hasGateway {
		devType = DevTypeGateway
	} else {
		devType = DevTypeSingle
	}
	return sf.SubscribeAllTopic(devType, sf.cfg.productKey, sf.cfg.deviceName)
}

//
//func (sf *Client) NewSubDevice(devType int, info *MetaInfo) (int, error) {
//	if !sf.cfg.hasGateway {
//		return 0, ErrFeatureNotSupport
//	}
//	return sf.Create(DevTypeSubDev, info.ProductKey, info.DeviceName, info.DeviceSecret)
//}
//
//func (sf *Client) SubDeviceConnect(id int) {
//
//}

// SetConn 设置连接接口
func (sf *Client) SetConn(conn Conn) *Client {
	sf.Conn = conn
	return sf
}

// SetDevUserProc 设置设备用户处理回调
func (sf *Client) SetDevUserProc(proc DevUserProc) *Client {
	sf.devUserProc = proc
	return sf
}

// SetGwUserProc 设置网关处理回调
func (sf *Client) SetGwUserProc(proc GatewayUserProc) *Client {
	sf.gwUserProc = proc
	return sf
}

// RequestID 获得下一个requestID,协程安全
func (sf *Client) RequestID() int {
	return int(atomic.AddInt32(&sf.requestID, 1))
}

// SendRequest 发送请求
// uriService 唯一定位服务器或(topic)
// requestID: 请求ID
// method: 方法
// params: 消息体
// API内部已实现json序列化
func (sf *Client) SendRequest(uriService string, requestID int, method string, params interface{}) error {
	out, err := json.Marshal(&Request{requestID, Version, params, method})
	if err != nil {
		return err
	}
	return sf.Publish(uriService, 1, out)
}

// SendResponse
// uriService 唯一定位服务器或(topic)
// responseID: 回复ID
// code: 回复code
// data: 数据域
// API内部已实现json序列化
func (sf *Client) SendResponse(uriService string, responseID int, code int, data interface{}) error {
	out, err := json.Marshal(
		struct {
			*Response
			Data interface{} `json:"data"`
		}{
			&Response{ID: responseID, Code: code},
			data,
		})
	if err != nil {
		return err
	}
	return sf.Publish(uriService, 1, out)
}

// AlinkReport 上报消息
// msgType 消息类型,支持:
//	- MsgTypeModelUpRaw
//  - MsgTypeEventPropertyPost
//  - MsgTypeDesiredPropertyGet
//  - MsgTypeDesiredPropertyDelete
//  - MsgTypeEventPropertyPost
//  - MsgTypeDesiredPropertyGet
//  - MsgTypeDeviceInfoUpdate
//  - MsgTypeDeviceInfoDelete
// devID 设备ID,独立设备或网关发送使用DevLocal
func (sf *Client) AlinkReport(msgType MsgType, devID int, payload interface{}) error {
	switch msgType {
	case MsgTypeModelUpRaw:
		if !sf.cfg.hasRawModel {
			return ErrNotSupportFeature
		}
		return sf.UpstreamThingModelUpRaw(devID, payload)
	case MsgTypeEventPropertyPost:
		if sf.cfg.hasRawModel {
			return ErrNotSupportFeature
		}
		return sf.UpstreamThingEventPropertyPost(devID, payload)
	case MsgTypeDesiredPropertyGet:
		if !sf.cfg.hasDesired {
			return ErrNotSupportFeature
		}
		return sf.UpstreamThingDesiredPropertyGet(devID, payload)
	case MsgTypeDesiredPropertyDelete:
		if !sf.cfg.hasDesired {
			return ErrNotSupportFeature
		}
		return sf.UpstreamThingDesiredPropertyDelete(devID, payload)
	case MsgTypeDeviceInfoUpdate:
		return sf.UpstreamThingDeviceInfoUpdate(devID, payload)
	case MsgTypeDeviceInfoDelete:
		return sf.UpstreamThingDeviceInfoDelete(devID, payload)

	case MsgTypeSubDevLogin:
		// TODO
	case MsgTypeSubDevLogout:
		//TODO
	case MsgTypeSubDevDeleteTopo:
		// todo
	case MsgTypeReportSubDevFirmwareVersion:
		// TODO
	}
	return ErrNotSupportMsgType
}

// AlinkQuery 请求查询
// msgType 消息类型,支持:
//  - MsgTypeExtNtpRequest
//  - MsgTypeDsltemplateGet
//  - MsgTypeConfigGet
func (sf *Client) AlinkQuery(msgType MsgType, devID int, payload interface{}) error {
	switch msgType {
	case MsgTypeDsltemplateGet:
		return sf.UpstreamThingDsltemplateGet(devID)
	case MsgTypeExtNtpRequest:
		if !sf.cfg.hasNTP {
			return ErrNotSupportFeature
		}
		return sf.UpstreamExtNtpRequest()
	case MsgTypeConfigGet:
		return sf.UpstreamThingConfigGet(devID)
	case MsgTypeExtErrorRequest:
		return sf.UpstreamExtErrorRequest()
	case MsgTypeQueryTopoList:
		// TODO
	case MsgTypeQueryCOTAData:
	case MsgTypeQueryFOTAData:
	case MsgTypeRequestCOTA:
	case MsgTypeRequestFOTAImage:
	}
	return ErrNotSupportMsgType
}

// AlinkTriggerEvent
func (sf *Client) AlinkTriggerEvent(devID int, eventID string, payload interface{}) error {
	return sf.UpstreamThingEventPost(devID, eventID, payload)
}
