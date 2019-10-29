package dm

// SubscribeAllTopic 对某个设备类型订阅相关所有主题
func (sf *Client) SubscribeAllTopic(devType DevType, productKey, deviceName string) error {
	var err error

	if sf.cfg.workOnWho == workOnHTTP {
		return nil
	}

	// model raw订阅
	if sf.cfg.hasRawModel {
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingModelUpRawReply, productKey, deviceName),
			procThingModelUpRawReply); err != nil {
			sf.warn(err.Error())
		}
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingModelDownRaw, productKey, deviceName),
			procThingModelDownRaw); err != nil {
			sf.warn(err.Error())
		}
	} else {
		// event 主题订阅
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingEventPostReplySingleWildcard, productKey, deviceName),
			procThingEventPostReply); err != nil {
			sf.warn(err.Error())
		}
	}

	// desired 期望属性订阅
	if sf.cfg.hasDesired {
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDesiredPropertyGetReply, productKey, deviceName),
			procThingDesiredPropertyGetReply); err != nil {
			sf.warn(err.Error())
		}
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDesiredPropertyDelete, productKey, deviceName),
			procThingDesiredPropertyDeleteReply); err != nil {
			sf.warn(err.Error())
		}
	}
	// deviceInfo 主题订阅
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDeviceInfoUpdateReply, productKey, deviceName),
		procThingDeviceInfoUpdateReply); err != nil {
		sf.warn(err.Error())
	}
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDeviceInfoDeleteReply, productKey, deviceName),
		procThingDeviceInfoDeleteReply); err != nil {
		sf.warn(err.Error())
	}

	// 服务调用
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingServicePropertySet, productKey, deviceName),
		procThingServicePropertySet); err != nil {
		sf.warn(err.Error())
	}
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingServiceRequestSingleWildcard, productKey, deviceName),
		procThingServiceRequest); err != nil {
		sf.warn(err.Error())
	}

	// dsltemplate 订阅
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDslTemplateGetReply, productKey, deviceName),
		procThingDsltemplateGetReply); err != nil {
		sf.warn(err.Error())
	}

	//// TODO: 不使用??
	// dynamictsl
	//if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDynamicTslGetReply, productKey, deviceName),
	//	procThingDynamictslGetReply); err != nil {
	//	sf.warn(err.Error())
	//}

	// RRPC
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIRRPCRequestSingleWildcard, productKey, deviceName),
		procRRPCRequest); err != nil {
		sf.warn(err.Error())
	}

	// ntp订阅, 只有网关和独立设备支持ntp
	if sf.cfg.hasNTP && !(devType == DevTypeSubDev) {
		if err = sf.Subscribe(sf.URIService(URIExtNtpPrefix, URINtpResponse, productKey, deviceName),
			procExtNtpResponse); err != nil {
			sf.warn(err.Error())
		}
	}

	// config 主题订阅
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingConfigGetReply, productKey, deviceName),
		procThingConfigGetReply); err != nil {
		sf.warn(err.Error())
	}
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingConfigPush, productKey, deviceName),
		procThingConfigPush); err != nil {
		sf.warn(err.Error())
	}

	// error 订阅
	if err = sf.Subscribe(sf.URIService(URIExtErrorPrefix, "", productKey, deviceName),
		procExtErrorResponse); err != nil {
		sf.warn(err.Error())
	}

	//if sf.cfg.hasGateway {
	//TODO
	//}

	return nil
}

// UnSubscribeSubDevAllTopic 取消子设备相关所有主题
func (sf *Client) UnSubscribeSubDevAllTopic(productKey, deviceName string) error {
	var topicList []string

	if !sf.cfg.hasGateway || sf.cfg.workOnWho == workOnHTTP {
		return nil
	}

	// model raw 取消订阅
	if sf.cfg.hasRawModel {
		topicList = append(topicList,
			sf.URIService(URISysPrefix, URIThingModelUpRawReply, productKey, deviceName),
			sf.URIService(URISysPrefix, URIThingModelDownRawReply, productKey, deviceName))
	} else {
		// event 取消订阅
		topicList = append(topicList,
			sf.URIService(URISysPrefix, URIThingEventPostReplySingleWildcard, productKey, deviceName))
	}

	// desired 期望属性取消订阅
	if sf.cfg.hasDesired {
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
		//sf.URIService(URISysPrefix, URIThingDynamicTslGetReply, productKey, deviceName),
		// RRPC
		sf.URIService(URISysPrefix, URIRRPCRequestSingleWildcard, productKey, deviceName),
		// config
		sf.URIService(URISysPrefix, URIThingConfigGetReply, productKey, deviceName),
		sf.URIService(URISysPrefix, URIThingConfigPush, productKey, deviceName),
		// error
		sf.URIService(URIExtErrorPrefix, "", productKey, deviceName))
	return sf.UnSubscribe(topicList...)
}
