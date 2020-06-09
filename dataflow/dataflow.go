package dataflow

import (
	"encoding/json"
)

// DeviceStatus 设备的状态
const (
	DeviceStatusOnline  = "online"
	DeviceStatusOffline = "offline"
)

// DeviceStatus 设备上下线状态
// Topic: /as/mqtt/status/{productKey}/{deviceName}
type DeviceStatus struct {
	Status      string  `json:"status"`      // "online"|"offline"
	ProductKey  string  `json:"productKey"`  // 设备所属的产品productKey
	DeviceName  string  `json:"deviceName"`  // 设各名称
	Time        Time    `json:"time"`        // 发送通知的时间点
	UtcTime     UTCtime `json:"utcTime"`     // 发送通知的UTC时间点
	LastTime    Time    `json:"lastTime"`    // 状态变更前最后一次通信的时间, 根据lastTime来维护最终设备的状态
	UtcLastTime UTCtime `json:"utcLastTime"` // 状态变更前最后一次通信的UTC时间。
	ClientIp    string  `json:"clientIp"`    // 设备公网出口IP
}

// MarshalBinary 实现接口
func (sf *DeviceStatus) MarshalBinary() ([]byte, error) {
	return json.Marshal(sf)
}

// 实现接口 UnmarshalBinary
func (sf *DeviceStatus) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, sf)
}

// DeviceProperty 设备属性上报
// Topic: /{productKey}/{deviceName}/thing/event/property/post
type DeviceProperty struct {
	IotId      string      `json:"iotId"`
	ProductKey string      `json:"productKey"` // 设备所属产品的productKey
	DeviceName string      `json:"deviceName"` // 设各名称
	GmtCreate  int64       `json:"gmtCreate"`  // 数据流转产生的时间,单位ms
	DeviceType string      `json:"deviceType"` // 设备类型
	Items      interface{} `json:"items"`      // 设备属性条目,由用户在平台定义
}

// DeviceEvent 设备上报的事件信息。
// Topic：/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post
type DeviceEvent struct {
	Identifier string      `json:"identifier"` // 事件唯一标识,事件ID
	Name       string      `json:"name"`       // 事件名称
	Type       string      // 事件类型,参见产品的TSL描述
	IotId      string      `json:"iotId"`      // 设备在平台内的唯一标识.
	ProductKey string      `json:"productKey"` // 设备所属productKey
	DeviceName string      `json:"deviceName"` // 设备名称
	GmtCreate  int64       `json:"gmtCreate"`  // 数据流转产生的时间,单位ms
	Values     interface{} `json:"values"`     // 设备事件的参数,由用户在平台根据事件ID定义,
	Time       int64       `json:"time"`       // 事件产生时间，如果设备没有上报默认采用远端时间。单位ms
}

// DeviceLifecycle action 设备生命周期action
const (
	DeviceLifeActionCreate  = "create"  // 创建设备
	DeviceLifeActionDelete  = "delete"  // 删除设备
	DeviceLifeActionEnable  = "enable"  // 启用设备
	DeviceLifeActionDisable = "disable" // 禁用设备
)

// DeviceLifecycle 设备生命周期
// Topic：/{productKey}/{deviceName}/thing/lifecycle
type DeviceLifecycle struct {
	Action            string `json:"action"`            // create|delete|enable|disable
	IotId             string `json:"iotId"`             // 设备在平台内的唯一标识.
	ProductKey        string `json:"productKey"`        // 设备所属productKey
	DeviceName        string `json:"deviceName"`        // 设备名称
	DeviceSecret      string `json:"deviceSecret"`      // 仅在设备create时包含
	MessageCreateTime int64  `json:"messageCreateTime"` // 消息产生时间戳,单位ms.
}

// DeviceTopologyRelation action 设备与网关拓扑关系action
const (
	DeviceTopoActionAdd     = "add"     // 新增
	DeviceTopoActionRemove  = "remove"  // 移除
	DeviceTopoActionEnable  = "enable"  // 启用
	DeviceTopoActionDisable = "disable" // 禁用
)

// DeviceSubInfo 子设备信息
type DeviceSubInfo struct {
	IotId      string `json:"iotId"`      // 设备在平台内的唯一标识.
	ProductKey string `json:"productKey"` // 设备所属productKey
	DeviceName string `json:"deviceName"` // 设各名称
}

// GwDeviceTopologyRelation 设备拓扑关系变更
// Topic：/{productKey}/{deviceName}/thing/topo/lifecycle
type GwDeviceTopologyRelation struct {
	Action            string          `json:"action"`            // add|remove|enable|disable
	GwIotId           string          `json:"gwIotId"`           // 网关设备在平台内的唯一标识.
	GwProductKey      string          `json:"gwProductKey"`      // 网关设备所属productKey
	GwDeviceName      string          `json:"gwDeviceName"`      // 网关设各名称
	Devices           []DeviceSubInfo `json:"devices"`           // 变更的子设备列表
	MessageCreateTime int64           `json:"messageCreateTime"` // 消息产生的时间戳,单秒ms
}

// GwDeviceFound 网关发现子设备
// Topic：/{productKey}/{deviceName}/thing/list/found
type GwDeviceFound struct {
	GwIotId      string          `json:"gwIotId"`      // 网关设备在平台内的唯一标识.
	GwProductKey string          `json:"gwProductKey"` // 网关设备所属productKey
	GwDeviceName string          `json:"gwDeviceName"` // 网关设各名称
	Devices      []DeviceSubInfo `json:"devices"`      // 发现的子设备列表
}

const (
	CodeSuccess             = 200  // 成功
	CodeErrRequest          = 400  // 内部服务错误
	CodeErrRequestPara      = 460  // 请求参数错误
	CodeErrTooManyRequest   = 429  // 请求过于频繁
	CodeErrDeviceNotActive  = 9200 // 设备没有激活
	CodeErrDeviceOffline    = 9201 // 设备不在线
	CodeErrRequestForbidden = 403  // 请求被禁止,由于欠费导致
)

// DeviceDownlinkResult 设备下行指令结果
// Topic：/{productKey}/{deviceName}/thing/downlink/reply/message
type DeviceDownlinkResult struct {
	GmtCreate  int64       `json:"gmtCreate"`  // UTC时间戳
	IotId      string      `json:"iotId"`      // 设备在平台内的唯一标识.
	ProductKey string      `json:"productKey"` // 设备所属productKey
	DeviceName string      `json:"deviceName"` // 设备名称
	RequestId  int64       `json:"requestId"`  // 阿里云产生和设备通信的信息ID
	Code       int         `json:"code"`       // 调用的结果code
	Message    string      `json:"message"`    // 结果code的消息
	Topic      string      `json:"topic"`      // 主题
	Data       interface{} `json:"data"`       // 数据,Alink格式数据返回设备处理的结果
}

// DeviceHistoryProperty 历史属性上报
// Topic：/sys/{productKey}/{deviceName}/thing/event/property/history/post
type DeviceHistoryProperty struct {
	IotId      string      `json:"iotId"`      // 设备在平台内的唯一标识.
	ProductKey string      `json:"productKey"` // 设备所属productKey
	DeviceName string      `json:"deviceName"` // 设备名称
	GmtCreate  int64       `json:"gmtCreate"`  // UTC时间戳
	DeviceType string      `json:"deviceType"` // 物模型类型，详情参见产品的TSL描述。
	Item       interface{} `json:"item"`       // 数据
}

// DeviceHistoryEvent 历史事件上报
// Topic：/sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/history/post
type DeviceHistoryEvent struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Type       string
	IotId      string      `json:"iotId"`      // 设备在平台内的唯一标识.
	ProductKey string      `json:"productKey"` // 设备所属productKey
	DeviceName string      `json:"deviceName"` // 设备名称
	GmtCreate  int64       `json:"gmtCreate"`  // UTC时间戳
	Value      interface{} `json:"value"`
	Time       int64       `json:"time"` // 事件产生时间，如果设备没有上报默认采用远端时间。单位ms
}

// DeviceOtaUpgrade 固件升级状态通知
// Topic：/sys/${productKey}/${deviceName}/ota/upgrade
type DeviceOtaUpgrade struct {
	IotId             string `json:"iotId"`             // 设备在平台内的唯一标识.
	ProductKey        string `json:"productKey"`        // 设备所属productKey
	DeviceName        string `json:"deviceName"`        // 设备名称
	Status            int64  `json:"status"`            // 升级状态 SUCCEEDED|FAILED
	MessageCreateTime int64  `json:"messageCreateTime"` // 消息产生时间戳，单位ms
	SrcVersion        string `json:"srcVersion"`        // 升级前的原固件版本。
	DestVersion       string `json:"destVersion"`       // 升级目标固件版本。
	Desc              string `json:"desc"`              // 升级状态描述信息。
	JobId             string `json:"jobId"`             // 升级批次ID，升级批次的唯一标识符。
	TaskId            string `json:"taskId"`            // 设备升级记录的唯一标识符。
}
