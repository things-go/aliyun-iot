// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package infra

// method 定义
const (
	MethodEventPropertyPost        = "thing.event.property.post"
	MethodEventFormatPost          = "thing.event.%s.post"
	MethodEventPropertyPackPost    = "thing.event.property.pack.post"
	MethodEventPropertyHistoryPost = "thing.event.property.history.post"
	MethodDeviceInfoUpdate         = "thing.deviceinfo.update"
	MethodDeviceInfoDelete         = "thing.deviceinfo.delete"
	MethodDesiredPropertyGet       = "thing.property.desired.get"
	MethodDesiredPropertyDelete    = "thing.property.desired.delete"
	MethodOtaFirmwareGet           = "thing.ota.firmware.get"
	MethodDslTemplateGet           = "thing.dsltemplate.get"
	MethodDynamicTslGet            = "thing.dynamicTsl.get"
	MethodConfigGet                = "thing.config.get"
	MethodConfigLogGet             = "thing.config.log.get"
	MethodLogPost                  = "thing.log.post"
	MethodSubDevRegister           = "thing.sub.register"
	MethodTopoAdd                  = "thing.topo.add"
	MethodTopoDelete               = "thing.topo.delete"
	MethodTopoGet                  = "thing.topo.get"
	MethodListFound                = "thing.list.found"
	MethodCombineLogin             = "combine.login"
	MethodCombineLogout            = "combine.logout"
	MethodCombineBatchLogin        = "combine.batch.login"
	MethodCombineBatchLogout       = "combine.batch.logout"
)
