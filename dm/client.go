package dm

import (
	"encoding/json"
	"sync/atomic"

	"github.com/thinkgos/aliyun-iot/uri"
)

// nextRequestID 获得下一个requestID,协程安全
func (sf *Client) nextRequestID() uint {
	return uint(atomic.AddUint32(&sf.requestID, 1))
}

// Request 发送请求,API内部已实现json序列化
// _uri 唯一定位服务器或(topic)
// requestID: 请求ID
// method: 方法
// params: 消息体Request的params
func (sf *Client) Request(_uri string, requestID uint, method string, params interface{}) error {
	out, err := json.Marshal(&Request{requestID, Version, params, method})
	if err != nil {
		return err
	}
	return sf.Publish(_uri, 1, out)
}

// SendRequest 发送请求,API内部已实现json序列化,requestID内部生成
// _uri 唯一定位服务器或(topic)
// method: 方法
// params: 消息体Request的params
func (sf *Client) SendRequest(_uri, method string, params interface{}) (*Token, error) {
	id := sf.nextRequestID()
	sf.log.Debugf("%s @%d", method, id)
	if err := sf.Request(_uri, id, method, params); err != nil {
		return nil, err
	}
	return sf.putPending(id), nil
}

// Response 发送回复
// _uri 唯一定位服务器或(topic)
// responseID: 回复ID
// code: 回复code
// Data: 数据域
// API内部已实现json序列化
func (sf *Client) Response(_uri string, responseID uint, code int, data interface{}) error {
	out, err := json.Marshal(&Response{ID: responseID, Code: code, Data: data})
	if err != nil {
		return err
	}
	return sf.Publish(_uri, 1, out)
}

// SubscribeAllTopic 对某个设备类型订阅相关所有主题
func (sf *Client) SubscribeAllTopic(productKey, deviceName string, isSub bool) error {
	var err error
	var _uri string

	if sf.workOnWho == WorkOnHTTP {
		return nil
	}

	// model raw订阅
	if sf.hasRawModel {
		_uri = uri.URI(uri.SysPrefix, uri.ThingModelUpRawReply, productKey, deviceName)
		if err = sf.Subscribe(_uri, ProcThingModelUpRawReply); err != nil {
			sf.log.Warnf(err.Error())
		}

		_uri = uri.URI(uri.SysPrefix, uri.ThingModelDownRaw, productKey, deviceName)
		if err = sf.Subscribe(_uri, ProcThingModelDownRaw); err != nil {
			sf.log.Warnf(err.Error())
		}
	} else {
		// event 主题订阅
		_uri = uri.URI(uri.SysPrefix, uri.ThingEventPostReplyWildcardOne, productKey, deviceName)
		if err = sf.Subscribe(_uri, ProcThingEventPostReply); err != nil {
			sf.log.Warnf(err.Error())
		}
	}

	// desired 期望属性订阅
	if sf.hasDesired {
		_uri = uri.URI(uri.SysPrefix, uri.ThingDesiredPropertyGetReply, productKey, deviceName)
		if err = sf.Subscribe(_uri, ProcThingDesiredPropertyGetReply); err != nil {
			sf.log.Warnf(err.Error())
		}
		_uri = uri.URI(uri.SysPrefix, uri.ThingDesiredPropertyDelete, productKey, deviceName)
		if err = sf.Subscribe(_uri, ProcThingDesiredPropertyDeleteReply); err != nil {
			sf.log.Warnf(err.Error())
		}
	}
	// deviceInfo 主题订阅
	_uri = uri.URI(uri.SysPrefix, uri.ThingDeviceInfoUpdateReply, productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcThingDeviceInfoUpdateReply); err != nil {
		sf.log.Warnf(err.Error())
	}
	_uri = uri.URI(uri.SysPrefix, uri.ThingDeviceInfoDeleteReply, productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcThingDeviceInfoDeleteReply); err != nil {
		sf.log.Warnf(err.Error())
	}

	// 服务调用
	_uri = uri.URI(uri.SysPrefix, uri.ThingServicePropertySet, productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcThingServicePropertySet); err != nil {
		sf.log.Warnf(err.Error())
	}
	_uri = uri.URI(uri.SysPrefix, uri.ThingServiceRequestWildcardOne, productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcThingServiceRequest); err != nil {
		sf.log.Warnf(err.Error())
	}

	// dsltemplate 订阅
	_uri = uri.URI(uri.SysPrefix, uri.ThingDslTemplateGetReply, productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcThingDsltemplateGetReply); err != nil {
		sf.log.Warnf(err.Error())
	}

	// TODO: 不使用??
	// dynamictsl
	// _uri = infra.URI(SysPrefix, URIThingDynamicTslGetReply, productKey, deviceName)
	// if err = sf.Subscribe(_uri, ProcThingDynamictslGetReply); err != nil {
	//	sf.log.Warnf(err.Error())
	// }

	// RRPC
	_uri = uri.URI(uri.SysPrefix, uri.RRPCRequestWildcardOne, productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcRRPCRequest); err != nil {
		sf.log.Warnf(err.Error())
	}

	// ntp订阅, 只有网关和独立设备支持ntp
	if sf.hasNTP && !isSub {
		_uri = uri.URI(uri.ExtNtpPrefix, uri.NtpResponse, productKey, deviceName)
		if err = sf.Subscribe(_uri, ProcExtNtpResponse); err != nil {
			sf.log.Warnf(err.Error())
		}
	}

	// config 主题订阅
	_uri = uri.URI(uri.SysPrefix, uri.ThingConfigGetReply, productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcThingConfigGetReply); err != nil {
		sf.log.Warnf(err.Error())
	}
	_uri = uri.URI(uri.SysPrefix, uri.ThingConfigPush, productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcThingConfigPush); err != nil {
		sf.log.Warnf(err.Error())
	}

	// error 订阅
	_uri = uri.URI(uri.ExtErrorPrefix, "", productKey, deviceName)
	if err = sf.Subscribe(_uri, ProcExtErrorResponse); err != nil {
		sf.log.Warnf(err.Error())
	}

	if sf.isGateway {
		if !isSub {
			// 网关批量上报数据
			_uri = sf.URIGateway(uri.SysPrefix, uri.ThingEventPropertyPackPostReply)
			if err = sf.Subscribe(_uri, ProcThingEventPropertyPackPostReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 添加该网关和子设备的拓扑关系
			_uri = uri.URI(uri.SysPrefix, uri.ThingTopoAddReply, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingTopoAddReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 删除该网关和子设备的拓扑关系
			_uri = uri.URI(uri.SysPrefix, uri.ThingTopoDeleteReply, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingTopoDeleteReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 获取该网关和子设备的拓扑关系
			_uri = uri.URI(uri.SysPrefix, uri.ThingTopoGetReply, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingTopoGetReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 发现设备列表上报
			if err = sf.Subscribe(uri.URI(uri.SysPrefix, uri.ThingListFoundReply, productKey, deviceName),
				ProcThingListFoundReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 添加设备拓扑关系通知,topic需要用网关的productKey,deviceName
			_uri = uri.URI(uri.SysPrefix, uri.ThingTopoAddNotify, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingTopoAddNotify); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 网关网络拓扑关系变化通知,topic需要用网关的productKey,deviceName
			_uri = uri.URI(uri.SysPrefix, uri.ThingTopoChange, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingTopoChange); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 子设备动态注册,topic需要用网关的productKey,deviceName
			_uri = uri.URI(uri.SysPrefix, uri.ThingSubRegisterReply, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingSubRegisterReply); err != nil {
				sf.log.Warnf(err.Error())
			}
			// 子设备上线,下线,topic需要用网关的productKey,deviceName,
			// 使用的是网关的通道,所以子设备不注册相关主题
			_uri = uri.URI(uri.ExtSessionPrefix, uri.CombineLoginReply, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcExtCombineLoginReply); err != nil {
				sf.log.Warnf(err.Error())
			}
			_uri = uri.URI(uri.ExtSessionPrefix, uri.CombineLogoutReply, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcExtCombineLoginoutReply); err != nil {
				sf.log.Warnf(err.Error())
			}
		} else {
			// 子设备禁用,启用,删除
			_uri = uri.URI(uri.SysPrefix, uri.ThingDisable, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingGwDisable); err != nil {
				sf.log.Warnf(err.Error())
			}
			_uri = uri.URI(uri.SysPrefix, uri.ThingEnable, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingGwEnable); err != nil {
				sf.log.Warnf(err.Error())
			}
			_uri = uri.URI(uri.SysPrefix, uri.ThingDelete, productKey, deviceName)
			if err = sf.Subscribe(_uri, ProcThingGwDelete); err != nil {
				sf.log.Warnf(err.Error())
			}
		}

		// if sf.hasOTA {
		// TODO
		// }
	}

	return nil
}

// UnSubscribeSubDevAllTopic 取消子设备相关所有主题
func (sf *Client) UnSubscribeSubDevAllTopic(productKey, deviceName string) error {
	var topicList []string

	if !sf.isGateway || sf.workOnWho == WorkOnHTTP {
		return nil
	}

	// model raw 取消订阅
	if sf.hasRawModel {
		topicList = append(topicList,
			uri.URI(uri.SysPrefix, uri.ThingModelUpRawReply, productKey, deviceName),
			uri.URI(uri.SysPrefix, uri.ThingModelDownRawReply, productKey, deviceName),
		)
	} else {
		// event 取消订阅
		topicList = append(topicList,
			uri.URI(uri.SysPrefix, uri.ThingEventPostReplyWildcardOne, productKey, deviceName),
		)
	}

	// desired 期望属性取消订阅
	if sf.hasDesired {
		topicList = append(topicList,
			uri.URI(uri.SysPrefix, uri.ThingDesiredPropertyGetReply, productKey, deviceName),
			uri.URI(uri.SysPrefix, uri.ThingDesiredPropertyDelete, productKey, deviceName),
		)
	}
	topicList = append(topicList,
		// deviceInfo
		uri.URI(uri.SysPrefix, uri.ThingDeviceInfoUpdateReply, productKey, deviceName),
		uri.URI(uri.SysPrefix, uri.ThingDeviceInfoDeleteReply, productKey, deviceName),
		// service
		uri.URI(uri.SysPrefix, uri.ThingServicePropertySet, productKey, deviceName),
		uri.URI(uri.SysPrefix, uri.ThingServiceRequestWildcardOne, productKey, deviceName),
		// dystemplate
		uri.URI(uri.SysPrefix, uri.ThingDslTemplateGetReply, productKey, deviceName),
		// dynamictsl 不使用??
		// infra.URI(SysPrefix, URIThingDynamicTslGetReply, productKey, deviceName),
		// RRPC
		uri.URI(uri.SysPrefix, uri.RRPCRequestWildcardOne, productKey, deviceName),
		// config
		uri.URI(uri.SysPrefix, uri.ThingConfigGetReply, productKey, deviceName),
		uri.URI(uri.SysPrefix, uri.ThingConfigPush, productKey, deviceName),
		// error
		uri.URI(uri.ExtErrorPrefix, "", productKey, deviceName),
	)
	return sf.UnSubscribe(topicList...)
}
