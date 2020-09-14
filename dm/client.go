package dm

import (
	uri2 "github.com/thinkgos/aliyun-iot/uri"
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
		uri = uri2.URI(uri2.SysPrefix, uri2.ThingModelUpRawReply, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingModelUpRawReply); err != nil {
			sf.log.Warnf(err.Error())
		}

		uri = uri2.URI(uri2.SysPrefix, uri2.ThingModelDownRaw, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingModelDownRaw); err != nil {
			sf.log.Warnf(err.Error())
		}
	} else {
		// event 主题订阅
		uri = uri2.URI(uri2.SysPrefix, uri2.ThingEventPostReplyWildcardOne, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingEventPostReply); err != nil {
			sf.log.Warnf(err.Error())
		}
	}

	// desired 期望属性订阅
	if sf.hasDesired {
		uri = uri2.URI(uri2.SysPrefix, uri2.ThingDesiredPropertyGetReply, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingDesiredPropertyGetReply); err != nil {
			sf.log.Warnf(err.Error())
		}
		uri = uri2.URI(uri2.SysPrefix, uri2.ThingDesiredPropertyDelete, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcThingDesiredPropertyDeleteReply); err != nil {
			sf.log.Warnf(err.Error())
		}
	}
	// deviceInfo 主题订阅
	uri = uri2.URI(uri2.SysPrefix, uri2.ThingDeviceInfoUpdateReply, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingDeviceInfoUpdateReply); err != nil {
		sf.log.Warnf(err.Error())
	}
	uri = uri2.URI(uri2.SysPrefix, uri2.ThingDeviceInfoDeleteReply, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingDeviceInfoDeleteReply); err != nil {
		sf.log.Warnf(err.Error())
	}

	// 服务调用
	uri = uri2.URI(uri2.SysPrefix, uri2.ThingServicePropertySet, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingServicePropertySet); err != nil {
		sf.log.Warnf(err.Error())
	}
	uri = uri2.URI(uri2.SysPrefix, uri2.ThingServiceRequestWildcardOne, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingServiceRequest); err != nil {
		sf.log.Warnf(err.Error())
	}

	// dsltemplate 订阅
	uri = uri2.URI(uri2.SysPrefix, uri2.ThingDslTemplateGetReply, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingDsltemplateGetReply); err != nil {
		sf.log.Warnf(err.Error())
	}

	// TODO: 不使用??
	// dynamictsl
	// uri = infra.URI(SysPrefix, URIThingDynamicTslGetReply, productKey, deviceName)
	// if err = sf.Subscribe(uri, ProcThingDynamictslGetReply); err != nil {
	//	sf.log.Warnf(err.Error())
	// }

	// RRPC
	uri = uri2.URI(uri2.SysPrefix, uri2.RRPCRequestWildcardOne, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcRRPCRequest); err != nil {
		sf.log.Warnf(err.Error())
	}

	// ntp订阅, 只有网关和独立设备支持ntp
	if sf.hasNTP && devType != DevTypeSubDev {
		uri = uri2.URI(uri2.ExtNtpPrefix, uri2.NtpResponse, productKey, deviceName)
		if err = sf.Subscribe(uri, ProcExtNtpResponse); err != nil {
			sf.log.Warnf(err.Error())
		}
	}

	// config 主题订阅
	uri = uri2.URI(uri2.SysPrefix, uri2.ThingConfigGetReply, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingConfigGetReply); err != nil {
		sf.log.Warnf(err.Error())
	}
	uri = uri2.URI(uri2.SysPrefix, uri2.ThingConfigPush, productKey, deviceName)
	if err = sf.Subscribe(uri, ProcThingConfigPush); err != nil {
		sf.log.Warnf(err.Error())
	}

	// error 订阅
	uri = uri2.URI(uri2.ExtErrorPrefix, "", productKey, deviceName)
	if err = sf.Subscribe(uri, ProcExtErrorResponse); err != nil {
		sf.log.Warnf(err.Error())
	}

	if sf.isGateway {
		if devType == DevTypeGateway {
			// 网关批量上报数据
			uri = sf.GatewayURI(uri2.SysPrefix, uri2.ThingEventPropertyPackPostReply)
			if err = sf.Subscribe(uri, ProcThingEventPropertyPackPostReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 添加该网关和子设备的拓扑关系
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingTopoAddReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoAddReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 删除该网关和子设备的拓扑关系
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingTopoDeleteReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoDeleteReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 获取该网关和子设备的拓扑关系
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingTopoGetReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoGetReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 发现设备列表上报
			if err = sf.Subscribe(uri2.URI(uri2.SysPrefix, uri2.ThingListFoundReply, productKey, deviceName),
				ProcThingGwListFoundReply); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 添加设备拓扑关系通知,topic需要用网关的productKey,deviceName
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingTopoAddNotify, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoAddNotify); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 网关网络拓扑关系变化通知,topic需要用网关的productKey,deviceName
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingTopoChange, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwTopoChange); err != nil {
				sf.log.Warnf(err.Error())
			}

			// 子设备动态注册,topic需要用网关的productKey,deviceName
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingSubRegisterReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwSubRegisterReply); err != nil {
				sf.log.Warnf(err.Error())
			}
			// 子设备上线,下线,topic需要用网关的productKey,deviceName,
			// 使用的是网关的通道,所以子设备不注册相关主题
			uri = uri2.URI(uri2.ExtSessionPrefix, uri2.CombineLoginReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcExtCombineLoginReply); err != nil {
				sf.log.Warnf(err.Error())
			}
			uri = uri2.URI(uri2.ExtSessionPrefix, uri2.CombineLogoutReply, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcExtCombineLoginoutReply); err != nil {
				sf.log.Warnf(err.Error())
			}
		}
		if devType == DevTypeSubDev {
			// 子设备禁用,启用,删除
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingDisable, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwDisable); err != nil {
				sf.log.Warnf(err.Error())
			}
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingEnable, productKey, deviceName)
			if err = sf.Subscribe(uri, ProcThingGwEnable); err != nil {
				sf.log.Warnf(err.Error())
			}
			uri = uri2.URI(uri2.SysPrefix, uri2.ThingDelete, productKey, deviceName)
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
			uri2.URI(uri2.SysPrefix, uri2.ThingModelUpRawReply, productKey, deviceName),
			uri2.URI(uri2.SysPrefix, uri2.ThingModelDownRawReply, productKey, deviceName),
		)
	} else {
		// event 取消订阅
		topicList = append(topicList,
			uri2.URI(uri2.SysPrefix, uri2.ThingEventPostReplyWildcardOne, productKey, deviceName),
		)
	}

	// desired 期望属性取消订阅
	if sf.hasDesired {
		topicList = append(topicList,
			uri2.URI(uri2.SysPrefix, uri2.ThingDesiredPropertyGetReply, productKey, deviceName),
			uri2.URI(uri2.SysPrefix, uri2.ThingDesiredPropertyDelete, productKey, deviceName),
		)
	}
	topicList = append(topicList,
		// deviceInfo
		uri2.URI(uri2.SysPrefix, uri2.ThingDeviceInfoUpdateReply, productKey, deviceName),
		uri2.URI(uri2.SysPrefix, uri2.ThingDeviceInfoDeleteReply, productKey, deviceName),
		// service
		uri2.URI(uri2.SysPrefix, uri2.ThingServicePropertySet, productKey, deviceName),
		uri2.URI(uri2.SysPrefix, uri2.ThingServiceRequestWildcardOne, productKey, deviceName),
		// dystemplate
		uri2.URI(uri2.SysPrefix, uri2.ThingDslTemplateGetReply, productKey, deviceName),
		// dynamictsl 不使用??
		// infra.URI(SysPrefix, URIThingDynamicTslGetReply, productKey, deviceName),
		// RRPC
		uri2.URI(uri2.SysPrefix, uri2.RRPCRequestWildcardOne, productKey, deviceName),
		// config
		uri2.URI(uri2.SysPrefix, uri2.ThingConfigGetReply, productKey, deviceName),
		uri2.URI(uri2.SysPrefix, uri2.ThingConfigPush, productKey, deviceName),
		// error
		uri2.URI(uri2.ExtErrorPrefix, "", productKey, deviceName),
	)
	return sf.UnSubscribe(topicList...)
}
