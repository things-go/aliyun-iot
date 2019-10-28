package model

// 对某个设备类型订阅相关所有主题
func (sf *Manager) SubscribeAllTopic(devType DevType, productKey, deviceName string) error {
	var err error

	if sf.cfg.workOnWho == workOnHTTP {
		return nil
	}

	// model raw订阅
	if sf.cfg.hasRawModel {
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingModelUpRawReply, productKey, deviceName),
			ProcThingModelUpRawReply); err != nil {
			sf.warn(err.Error())
		}
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingModelDownRaw, productKey, deviceName),
			ProcThingModelDownRaw); err != nil {
			sf.warn(err.Error())
		}
	} else {
		// event 主题订阅
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingEventPostReplySingleWildcard, productKey, deviceName),
			ProcThingEventPostReply); err != nil {
			sf.warn(err.Error())
		}
	}

	// desired 期望属性订阅
	if sf.cfg.hasDesired {
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDesiredPropertyGetReply, productKey, deviceName),
			ProcThingDesiredPropertyGetReply); err != nil {
			sf.warn(err.Error())
		}
		if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDesiredPropertyDelete, productKey, deviceName),
			ProcThingDesiredPropertyDeleteReply); err != nil {
			sf.warn(err.Error())
		}
	}
	// deviceInfo 主题订阅
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDeviceInfoUpdateReply, productKey, deviceName),
		ProcThingDeviceInfoUpdateReply); err != nil {
		sf.warn(err.Error())
	}
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDeviceInfoDeleteReply, productKey, deviceName),
		ProcThingDeviceInfoDeleteReply); err != nil {
		sf.warn(err.Error())
	}

	// 服务调用
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingServicePropertySet, productKey, deviceName),
		ProcThingServicePropertySet); err != nil {
		sf.warn(err.Error())
	}
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingServiceRequestSingleWildcard, productKey, deviceName),
		ProcThingServiceRequest); err != nil {
		sf.warn(err.Error())
	}

	// dsltemplate 订阅
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIThingDslTemplateGetReply, productKey, deviceName),
		ProcThingDsltemplateGetReply); err != nil {
		sf.warn(err.Error())
	}

	// RRPC
	if err = sf.Subscribe(sf.URIService(URISysPrefix, URIRRPCRequestSingleWildcard, productKey, deviceName),
		ProcRRPCRequest); err != nil {
		sf.warn(err.Error())
	}

	// ntp订阅, 只有网关和独立设备支持ntp
	if sf.cfg.hasNTP && !(devType == DevTypeSubDev) {
		if err = sf.Subscribe(sf.URIService(URIExtNtpPrefix, URINtpRequest, productKey, deviceName),
			ProcThingDeviceInfoDeleteReply); err != nil {
			sf.warn(err.Error())
		}
	}

	// error 订阅
	if err = sf.Subscribe(sf.URIService(URIExtErrorPrefix, "", productKey, deviceName),
		ProcThingDeviceInfoDeleteReply); err != nil {
		sf.warn(err.Error())
	}

	//if sf.cfg.hasGateway {
	//TODO
	//}

	return nil
}

// UnSubscribeSubDevAllTopic 取消子设备相关所有主题
func (sf *Manager) UnSubscribeSubDevAllTopic(productKey, deviceName string) error {
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

	// deviceInfo 主题取消订阅
	topicList = append(topicList,
		sf.URIService(URISysPrefix, URIThingDeviceInfoUpdateReply, productKey, deviceName),
		sf.URIService(URISysPrefix, URIThingDeviceInfoDeleteReply, productKey, deviceName),
		// 服务调用
		sf.URIService(URISysPrefix, URIThingServicePropertySet, productKey, deviceName),
		sf.URIService(URISysPrefix, URIThingServiceRequestSingleWildcard, productKey, deviceName),
		// dystemplate 订阅
		sf.URIService(URISysPrefix, URIThingDslTemplateGetReply, productKey, deviceName),
		// RRPC
		sf.URIService(URISysPrefix, URIRRPCRequestSingleWildcard, productKey, deviceName),
		// error 订阅
		sf.URIService(URIExtErrorPrefix, "", productKey, deviceName))
	return sf.UnSubscribe(topicList...)
}
