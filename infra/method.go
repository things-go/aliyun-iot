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
	MethodSubDevRegister           = "thing.sub.register"
	MethodTopoAdd                  = "thing.topo.add"
	MethodTopoDelete               = "thing.topo.delete"
	MethodTopoGet                  = "thing.topo.get"
	MethodListFound                = "thing.list.found"
)
