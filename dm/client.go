package dm

import (
	"github.com/thinkgos/aliyun-iot/infra"
)

// SubscribeAllTopic 对某个设备类型订阅相关所有主题
func (sf *Client) SubscribeAllTopic(devType DevType, productKey, deviceName string) error {
	var err error
	var uri string

	if sf.workOnWho == WorkOnHTTP {
		return nil
	}

	// model raw订阅
	if sf.hasRawModel {
		uri = sf.uriService(infra.URISysPrefix, infra.URIThingModelUpRawReply, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingModelUpRawReply); err != nil {
			sf.log.Warnf(err.Error())
		}

		uri = sf.uriService(infra.URISysPrefix, infra.URIThingModelDownRaw, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingModelDownRaw); err != nil {
			sf.log.Warnf(err.Error())
		}
	} else {
		// event 主题订阅
		uri = sf.uriService(infra.URISysPrefix, infra.URIThingEventPostReplySingleWildcard, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingEventPostReply); err != nil {
			sf.log.Warnf(err.Error())
		}
	}

	// desired 期望属性订阅
	if sf.hasDesired {
		uri = sf.uriService(infra.URISysPrefix, infra.URIThingDesiredPropertyGetReply, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingDesiredPropertyGetReply); err != nil {
			sf.log.Warnf(err.Error())
		}
		uri = sf.uriService(infra.URISysPrefix, infra.URIThingDesiredPropertyDelete, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingDesiredPropertyDeleteReply); err != nil {
			sf.log.Warnf(err.Error())
		}
	}
	// deviceInfo 主题订阅
	uri = sf.uriService(infra.URISysPrefix, infra.URIThingDeviceInfoUpdateReply, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingDeviceInfoUpdateReply); err != nil {
		sf.log.Warnf(err.Error())
	}
	uri = sf.uriService(infra.URISysPrefix, infra.URIThingDeviceInfoDeleteReply, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingDeviceInfoDeleteReply); err != nil {
		sf.log.Warnf(err.Error())
	}

	// 服务调用
	uri = sf.uriService(infra.URISysPrefix, infra.URIThingServicePropertySet, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingServicePropertySet); err != nil {
		sf.log.Warnf(err.Error())
	}
	uri = sf.uriService(infra.URISysPrefix, infra.URIThingServiceRequestSingleWildcard, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingServiceRequest); err != nil {
		sf.log.Warnf(err.Error())
	}

	// dsltemplate 订阅
	uri = sf.uriService(infra.URISysPrefix, infra.URIThingDslTemplateGetReply, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingDsltemplateGetReply); err != nil {
		sf.log.Warnf(err.Error())
	}

	// TODO: 不使用??
	// dynamictsl
	// uri = sf.uriService(URISysPrefix, URIThingDynamicTslGetReply, productKey, deviceName)
	// if err = sf.Subscribe(uri, ProcThingDynamictslGetReply); err != nil {
	//	sf.log.Warnf(err.Error())
	// }

	// RRPC
	uri = sf.uriService(infra.URISysPrefix, infra.URIRRPCRequestSingleWildcard, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcRRPCRequest); err != nil {
		sf.log.Warnf(err.Error())
	}

	// ntp订阅, 只有网关和独立设备支持ntp
	if sf.hasNTP && devType != DevTypeSubDev {
		uri = sf.uriService(infra.URIExtNtpPrefix, infra.URINtpResponse, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcExtNtpResponse); err != nil {
			sf.log.Warnf(err.Error())
		}
	}

	// config 主题订阅
	uri = sf.uriService(infra.URISysPrefix, infra.URIThingConfigGetReply, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingConfigGetReply); err != nil {
		sf.log.Warnf(err.Error())
	}
	uri = sf.uriService(infra.URISysPrefix, infra.URIThingConfigPush, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingConfigPush); err != nil {
		sf.log.Warnf(err.Error())
	}

	// error 订阅
	uri = sf.uriService(infra.URIExtErrorPrefix, "", productKey, deviceName)
	if err = sf.Subscribe(uri, ProcExtErrorResponse); err != nil {
		sf.log.Warnf(err.Error())
	}

	if sf.isGateway {
		if devType == DevTypeGateway {
			// 网关批量上报数据
			uri = sf.URIServiceSelf(infra.URISysPrefix, infra.URIThingEventPropertyPackPostReply)
			if err = sf.Subscribe(uri, ProcThingEventPropertyPackPostReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 添加该网关和子设备的拓扑关系
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingTopoAddReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoAddReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 删除该网关和子设备的拓扑关系
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingTopoDeleteReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoDeleteReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 获取该网关和子设备的拓扑关系
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingTopoGetReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoGetReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 发现设备列表上报
			if err = sf.Subscribe(sf.uriService(infra.URISysPrefix, infra.URIThingListFoundReply, productKey, deviceName),
				ProcThingGwListFoundReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 添加设备拓扑关系通知,topic需要用网关的productKey,deviceName
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingTopoAddNotify, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoAddNotify); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 网关网络拓扑关系变化通知,topic需要用网关的productKey,deviceName
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingTopoChange, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoChange); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 子设备动态注册,topic需要用网关的productKey,deviceName
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingSubDevRegisterReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwSubRegisterReply); err != nil {
				sf.log.Warnf(err.Error())
			}
			// 子设备上线,下线,topic需要用网关的productKey,deviceName,
			// 使用的是网关的通道,所以子设备不注册相关主题
			uri = sf.uriService(infra.URIExtSessionPrefix, infra.URICombineLoginReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcExtCombineLoginReply); err != nil {
				sf.log.Warnf(err.Error())
			}
			uri = sf.uriService(infra.URIExtSessionPrefix, infra.URICombineLogoutReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcExtCombineLoginoutReply); err != nil {
				sf.log.Warnf(err.Error())
			}
		}
		if devType == DevTypeSubDev {
			// 子设备禁用,启用,删除
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingDisable, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwDisable); err != nil {
				sf.log.Warnf(err.Error())
			}
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingEnable, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwEnable); err != nil {
				sf.log.Warnf(err.Error())
			}
			uri = sf.uriService(infra.URISysPrefix, infra.URIThingDelete, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwDelete); err != nil {
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
			sf.uriService(infra.URISysPrefix, infra.URIThingModelUpRawReply, productKey, deviceName),
			sf.uriService(infra.URISysPrefix, infra.URIThingModelDownRawReply, productKey, deviceName),
		)
	} else {
		// event 取消订阅
		topicList = append(topicList,
			sf.uriService(infra.URISysPrefix, infra.URIThingEventPostReplySingleWildcard, productKey, deviceName),
		)
	}

	// desired 期望属性取消订阅
	if sf.hasDesired {
		topicList = append(topicList,
			sf.uriService(infra.URISysPrefix, infra.URIThingDesiredPropertyGetReply, productKey, deviceName),
			sf.uriService(infra.URISysPrefix, infra.URIThingDesiredPropertyDelete, productKey, deviceName),
		)
	}
	topicList = append(topicList,
		// deviceInfo
		sf.uriService(infra.URISysPrefix, infra.URIThingDeviceInfoUpdateReply, productKey, deviceName),
		sf.uriService(infra.URISysPrefix, infra.URIThingDeviceInfoDeleteReply, productKey, deviceName),
		// service
		sf.uriService(infra.URISysPrefix, infra.URIThingServicePropertySet, productKey, deviceName),
		sf.uriService(infra.URISysPrefix, infra.URIThingServiceRequestSingleWildcard, productKey, deviceName),
		// dystemplate
		sf.uriService(infra.URISysPrefix, infra.URIThingDslTemplateGetReply, productKey, deviceName),
		// dynamictsl 不使用??
		// sf.uriService(URISysPrefix, URIThingDynamicTslGetReply, productKey, deviceName),
		// RRPC
		sf.uriService(infra.URISysPrefix, infra.URIRRPCRequestSingleWildcard, productKey, deviceName),
		// config
		sf.uriService(infra.URISysPrefix, infra.URIThingConfigGetReply, productKey, deviceName),
		sf.uriService(infra.URISysPrefix, infra.URIThingConfigPush, productKey, deviceName),
		// error
		sf.uriService(infra.URIExtErrorPrefix, "", productKey, deviceName),
	)
	return sf.UnSubscribe(topicList...)
}
