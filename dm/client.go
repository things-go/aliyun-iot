package dm

// SubscribeAllTopic 对某个设备类型订阅相关所有主题
func (sf *Client) SubscribeAllTopic(devType DevType, productKey, deviceName string) error {
	var err error

	if sf.workOnWho == WorkOnHTTP {
		return nil
	}

	// model raw订阅
	if sf.hasRawModel {
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingModelUpRawReply, productKey, deviceName),
			ProcThingModelUpRawReply); err != nil {
			sf.warnf(err.Error())
		}
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingModelDownRaw, productKey, deviceName),
			ProcThingModelDownRaw); err != nil {
			sf.warnf(err.Error())
		}
	} else {
		// event 主题订阅
		err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingEventPostReplySingleWildcard, productKey, deviceName),
			ProcThingEventPostReply)
		if err != nil {
			sf.warnf(err.Error())
		}
	}

	// desired 期望属性订阅
	if sf.hasDesired {
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDesiredPropertyGetReply, productKey, deviceName),
			ProcThingDesiredPropertyGetReply); err != nil {
			sf.warnf(err.Error())
		}
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDesiredPropertyDelete, productKey, deviceName),
			ProcThingDesiredPropertyDeleteReply); err != nil {
			sf.warnf(err.Error())
		}
	}
	// deviceInfo 主题订阅
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDeviceInfoUpdateReply, productKey, deviceName),
		ProcThingDeviceInfoUpdateReply); err != nil {
		sf.warnf(err.Error())
	}
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDeviceInfoDeleteReply, productKey, deviceName),
		ProcThingDeviceInfoDeleteReply); err != nil {
		sf.warnf(err.Error())
	}

	// 服务调用
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingServicePropertySet, productKey, deviceName),
		ProcThingServicePropertySet); err != nil {
		sf.warnf(err.Error())
	}
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingServiceRequestSingleWildcard, productKey, deviceName),
		ProcThingServiceRequest); err != nil {
		sf.warnf(err.Error())
	}

	// dsltemplate 订阅
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDslTemplateGetReply, productKey, deviceName),
		ProcThingDsltemplateGetReply); err != nil {
		sf.warnf(err.Error())
	}

	// TODO: 不使用??
	// dynamictsl
	// if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDynamicTslGetReply, productKey, deviceName),
	//	ProcThingDynamictslGetReply); err != nil {
	//	sf.warnf(err.Error())
	// }

	// RRPC
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIRRPCRequestSingleWildcard, productKey, deviceName),
		ProcRRPCRequest); err != nil {
		sf.warnf(err.Error())
	}

	// ntp订阅, 只有网关和独立设备支持ntp
	if sf.hasNTP && devType != DevTypeSubDev {
		if err = sf.Subscribe(sf.URIService(URIExtNtpPrefix, URINtpResponse, productKey, deviceName),
			ProcExtNtpResponse); err != nil {
			sf.warnf(err.Error())
		}
	}

	// config 主题订阅
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingConfigGetReply, productKey, deviceName),
		ProcThingConfigGetReply); err != nil {
		sf.warnf(err.Error())
	}
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingConfigPush, productKey, deviceName),
		ProcThingConfigPush); err != nil {
		sf.warnf(err.Error())
	}

	// error 订阅
	if err = sf.Subscribe(sf.URIService(URIExtErrorPrefix, "", productKey, deviceName),
		ProcExtErrorResponse); err != nil {
		sf.warnf(err.Error())
	}

	if sf.isGateway {
		if devType == DevTypeGateway {
			// 网关批量上报数据
			if err = sf.Subscribe(sf.URIServiceSelf(URISysPrefix, URIThingEventPropertyPackPostReply),
				ProcThingEventPropertyPackPostReply); err != nil {
				sf.warnf(err.Error())
			}

			// 添加该网关和子设备的拓扑关系
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingTopoAddReply, productKey, deviceName),
				ProcThingGwSubTopoAddReply); err != nil {
				sf.warnf(err.Error())
			}

			// 删除该网关和子设备的拓扑关系
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingTopoDeleteReply, productKey, deviceName),
				ProcThingGwSubTopoDeleteReply); err != nil {
				sf.warnf(err.Error())
			}

			// 获取该网关和子设备的拓扑关系
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingTopoGetReply, productKey, deviceName),
				ProcThingGwSubTopoGetReply); err != nil {
				sf.warnf(err.Error())
			}

			// 发现设备列表上报
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingListFoundReply, productKey, deviceName),
				ProcThingGwSubListFoundReply); err != nil {
				sf.warnf(err.Error())
			}

			// 添加设备拓扑关系通知,topic需要用网关的productKey,deviceName
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingTopoAddNotify, productKey, deviceName),
				ProcThingGwTopoAddNotify); err != nil {
				sf.warnf(err.Error())
			}

			// 网关网络拓扑关系变化通知,topic需要用网关的productKey,deviceName
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingTopoChange, productKey, deviceName),
				ProcThingGwTopoChange); err != nil {
				sf.warnf(err.Error())
			}

			// 子设备动态注册,topic需要用网关的productKey,deviceName
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingSubDevRegisterReply, productKey, deviceName),
				ProcThingGwSubRegisterReply); err != nil {
				sf.warnf(err.Error())
			}
			// 子设备上线,下线,topic需要用网关的productKey,deviceName,
			// 使用的是网关的通道,所以子设备不注册相关主题
			if err = sf.Subscribe(sf.URIService(URIExtSessionPrefix, URICombineLoginReply, productKey, deviceName),
				ProcExtCombineLoginReply); err != nil {
				sf.warnf(err.Error())
			}
			if err = sf.Subscribe(sf.URIService(URIExtSessionPrefix, URICombineLogoutReply, productKey, deviceName),
				ProcExtCombineLogoutReply); err != nil {
				sf.warnf(err.Error())
			}
		}
		if devType == DevTypeSubDev {
			// 子设备禁用,启用,删除
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDisable, productKey, deviceName),
				ProcThingDisable); err != nil {
				sf.warnf(err.Error())
			}
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingEnable, productKey, deviceName),
				ProcThingEnable); err != nil {
				sf.warnf(err.Error())
			}
			if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDelete, productKey, deviceName),
				ProcThingDelete); err != nil {
				sf.warnf(err.Error())
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
			sf.URIService(URISysPrefix, URIThingModelUpRawReply, productKey, deviceName),
			sf.URIService(URISysPrefix, URIThingModelDownRawReply, productKey, deviceName))
	} else {
		// event 取消订阅
		topicList = append(topicList,
			sf.URIService(URISysPrefix, URIThingEventPostReplySingleWildcard, productKey, deviceName))
	}

	// desired 期望属性取消订阅
	if sf.hasDesired {
		topicList = append(topicList,
			sf.URIService(URISysPrefix, URIThingDesiredPropertyGetReply, productKey, deviceName),
			sf.URIService(URISysPrefix, URIThingDesiredPropertyDelete, productKey, deviceName))
	}
	topicList = append(topicList,
		// deviceInfo
		sf.URIService(URISysPrefix, URIThingDeviceInfoUpdateReply, productKey, deviceName),
		sf.URIService(URISysPrefix, URIThingDeviceInfoDeleteReply, productKey, deviceName),
		// service
		sf.URIService(URISysPrefix, URIThingServicePropertySet, productKey, deviceName),
		sf.URIService(URISysPrefix, URIThingServiceRequestSingleWildcard, productKey, deviceName),
		// dystemplate
		sf.URIService(URISysPrefix, URIThingDslTemplateGetReply, productKey, deviceName),
		// dynamictsl 不使用??
		// sf.URIService(URISysPrefix, URIThingDynamicTslGetReply, productKey, deviceName),
		// RRPC
		sf.URIService(URISysPrefix, URIRRPCRequestSingleWildcard, productKey, deviceName),
		// config
		sf.URIService(URISysPrefix, URIThingConfigGetReply, productKey, deviceName),
		sf.URIService(URISysPrefix, URIThingConfigPush, productKey, deviceName),
		// error
		sf.URIService(URIExtErrorPrefix, "", productKey, deviceName))
	return sf.UnSubscribe(topicList...)
}
