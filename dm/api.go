package dm

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/thinkgos/cache-go"
)

// MsgType 消息类型
type MsgType byte

// 消息类型定义
const (
	MsgTypeModelUpRaw            MsgType = iota //!< post raw data to cloud
	MsgTypeEventPropertyPost                    //!< post property value to cloud
	MsgTypeEventPost                            //!< post event identifies value to cloud
	MsgTypeEventPropertyPackPost                //!<
	MsgTypeDesiredPropertyGet                   //!< get a device's desired property
	MsgTypeDesiredPropertyDelete                //!< delete a device's desired property
	MsgTypeDeviceInfoUpdate                     //!< post device info update message to cloud
	MsgTypeDeviceInfoDelete                     //!< post device info delete message to cloud
	MsgTypeDsltemplateGet                       //<! get a device's dsltemplate
	MsgTypeDynamictslGet                        //!< ??
	MsgTypeExtNtpRequest                        //!< query ntp time from cloud
	MsgTypeConfigGet                            //!< 获取配置

	MsgTypeTopoAdd                     //!< 网关,添加设备拓扑关系
	MsgTypeTopoDelete                  //!< 网关,删除设备拓扑关系
	MsgTypeTopoGet                     //!< 网关,查询设备拓扑关系
	MsgTypeDevListFound                //!< 网关,设备发现链表上报
	MsgTypeSubDevRegister              //!< 子设备,动态注册
	MsgTypeSubDevLogin                 //!< only for slave device, send login request to cloud
	MsgTypeSubDevLogout                //!< only for slave device, send logout request to cloud
	MsgTypeSubDevDeleteTopo            //!< only for slave device, send delete topo request to cloud
	MsgTypeQueryFOTAData               //!< only for master device, qurey firmware ota data
	MsgTypeQueryCOTAData               //!< only for master device, qurey config ota data
	MsgTypeRequestCOTA                 //!< only for master device, request config ota data from cloud
	MsgTypeRequestFOTAImage            //!< only for master device, request fota image from cloud
	MsgTypeReportSubDevFirmwareVersion //!< report subdev's firmware version
)

// Meta meta 信息
type Meta struct {
	ProductKey    string
	ProductSecret string
	DeviceName    string
	DeviceSecret  string
}

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

// Client 客户端
type Client struct {
	requestID int32

	cfg Config

	*DevMgr
	syncHub  *SyncHub
	msgCache *cache.Cache
	pool     *pool
	Conn
	ipc       chan *ipcMessage
	eventProc EventProc
}

// New 创建一个物管理客户端
func New(cfg *Config) *Client {
	sf := &Client{
		cfg:       *cfg,
		DevMgr:    NewDevMgr(),
		syncHub:   NewSyncHub(),
		ipc:       make(chan *ipcMessage, 1024),
		eventProc: NopEvt{},
	}
	if cfg.hasCache {
		sf.pool = newPool()
		sf.msgCache = cache.New(time.Second*10, time.Second*30)
	}
	sf.cacheInit()
	err := sf.insert(DevNodeLocal, DevTypeSingle, cfg.productKey, cfg.deviceName, cfg.deviceSecret)
	if err != nil {
		panic(fmt.Sprintf("device local duplicate,cause: %+v", err))
	}
	go sf.ipcRunMessage()
	return sf
}

func (sf *Client) NewSubDevice(devType int, meta Meta) (int, error) {
	if !sf.cfg.hasGateway {
		return 0, ErrNotSupportFeature
	}
	return sf.Create(DevTypeSubDev, meta.ProductKey, meta.DeviceName, meta.DeviceSecret)
}

// SetConn 设置连接接口
func (sf *Client) SetConn(conn Conn) *Client {
	sf.Conn = conn
	return sf
}

// SetDevUserProc 设置设备用户处理回调
func (sf *Client) SetDevUserProc(proc EventProc) *Client {
	sf.eventProc = proc
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

// SendResponse 发送回复
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

// AlinkConnect 将订阅所有相关主题,主题有config配置
func (sf *Client) AlinkConnect() error {
	var devType DevType

	if sf.cfg.hasGateway {
		devType = DevTypeGateway
	} else {
		devType = DevTypeSingle
	}
	return sf.SubscribeAllTopic(devType, sf.cfg.productKey, sf.cfg.deviceName)
}

// AlinkSubDeviceConnect 子设备连接注册并添加到网关拓扑关系
func (sf *Client) AlinkSubDeviceConnect(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	node, err := sf.SearchNodeByID(devID)
	if err != nil {
		return err
	}
	if node.DeviceSecret() == "" { // 需要注册
		// 子设备注册
		if err := sf.linKitGwSubDevRegister(devID); err != nil {
			return err
		}
	}

	// 子设备添加到拓扑
	return sf.linkKitGwSubDevTopoAdd(devID)
}

// AlinkReport 上报消息
// msgType 消息类型,支持:
//	- MsgTypeModelUpRaw
//  - MsgTypeEventPropertyPost
//  - MsgTypeEventPropertyPackPost
//  - MsgTypeDesiredPropertyGet
//  - MsgTypeDesiredPropertyDelete
//  - MsgTypeEventPropertyPost
//  - MsgTypeDesiredPropertyGet
//  - MsgTypeDeviceInfoUpdate
//  - MsgTypeDeviceInfoDelete
// devID 设备ID,独立设备或网关发送使用DevLocal
func (sf *Client) AlinkReport(msgType MsgType, devID int, payload interface{}) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	switch msgType {
	case MsgTypeModelUpRaw:
		if !sf.cfg.hasRawModel {
			return ErrNotSupportFeature
		}
		return sf.upstreamThingModelUpRaw(devID, payload)
	case MsgTypeEventPropertyPost:
		if sf.cfg.hasRawModel {
			return ErrNotSupportFeature
		}
		return sf.upstreamThingEventPropertyPost(devID, payload)
	case MsgTypeEventPropertyPackPost:
		if !sf.cfg.hasGateway {
			return ErrNotSupportFeature
		}
		return sf.upstreamThingEventPropertyPackPost(payload)
	case MsgTypeDesiredPropertyGet:
		if !sf.cfg.hasDesired {
			return ErrNotSupportFeature
		}
		return sf.upstreamThingDesiredPropertyGet(devID, payload)
	case MsgTypeDesiredPropertyDelete:
		if !sf.cfg.hasDesired {
			return ErrNotSupportFeature
		}
		return sf.upstreamThingDesiredPropertyDelete(devID, payload)
	case MsgTypeDeviceInfoUpdate:
		return sf.upstreamThingDeviceInfoUpdate(devID, payload)
	case MsgTypeDeviceInfoDelete:
		return sf.upstreamThingDeviceInfoDelete(devID, payload)

	case MsgTypeReportSubDevFirmwareVersion:
		// TODO
	}
	return ErrNotSupportMsgType
}

// AlinkRequest 同步请求
// msgType 消息类型,支持:
//  MsgTypeSubDevLogin
//  MsgTypeSubDevLogout
//  MsgTypeSubDevDeleteTopo
func (sf *Client) AlinkRequest(msgType MsgType, devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	switch msgType {
	case MsgTypeSubDevLogin:
		if !sf.cfg.hasGateway {
			return ErrNotSupportFeature
		}
		return sf.linkKitGwSubDevCombineLogin(devID)
	case MsgTypeSubDevLogout:
		if !sf.cfg.hasGateway {
			return ErrNotSupportFeature
		}
		return sf.linkKitGwSubDevCombineLogout(devID)
	case MsgTypeSubDevDeleteTopo:
		if !sf.cfg.hasGateway {
			return ErrNotSupportFeature
		}
		return sf.linkKitGwSubDevTopoDelete(devID)
	}
	return ErrNotSupportMsgType
}

// AlinkQuery 请求查询
// msgType 消息类型,支持:
//  - MsgTypeExtNtpRequest
//  - MsgTypeDsltemplateGet
//  - MsgTypeConfigGet
func (sf *Client) AlinkQuery(msgType MsgType, devID int, payload ...interface{}) error {
	switch msgType {
	case MsgTypeDsltemplateGet:
		return sf.UpstreamThingDsltemplateGet(devID)
	case MsgTypeDynamictslGet:
		// TODO: BUG
		return sf.UpstreamThingDynamictslGet()
	case MsgTypeExtNtpRequest:
		if !sf.cfg.hasNTP || sf.cfg.hasRawModel {
			return ErrNotSupportFeature
		}
		return sf.upstreamExtNtpRequest()
	case MsgTypeConfigGet:
		if !sf.cfg.hasGateway {
			return ErrNotSupportFeature
		}
		return sf.upstreamThingConfigGet(devID)
	case MsgTypeTopoGet:
		if !sf.cfg.hasGateway {
			return ErrNotSupportFeature
		}
		return sf.upstreamGwThingTopoGet()
	case MsgTypeQueryCOTAData:
	case MsgTypeQueryFOTAData:
	case MsgTypeRequestCOTA:
	case MsgTypeRequestFOTAImage:
	}
	return ErrNotSupportMsgType
}

// AlinkTriggerEvent 事件上报
func (sf *Client) AlinkTriggerEvent(devID int, eventID string, payload interface{}) error {
	return sf.UpstreamThingEventPost(devID, eventID, payload)
}
