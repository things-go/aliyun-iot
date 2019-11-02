package dm

// 平台通信版本
const (
	Version = "1.0"
)

// method 定义
const (
	methodEventPropertyPost     = "thing.event.property.post"
	methodEventFormatPost       = "thing.event.%s.post"
	methodEventPropertyPackPost = "thing.event.property.pack.post"
	methodDeviceInfoUpdate      = "thing.deviceinfo.update"
	methodDeviceInfoDelete      = "thing.deviceinfo.delete"
	methodDesiredPropertyGet    = "thing.property.desired.get"
	methodDesiredPropertyDelete = "thing.property.desired.delete"
	methodDslTemplateGet        = "thing.dsltemplate.get"
	methodDynamicTslGet         = "thing.dynamicTsl.get"
	methodConfigGet             = "thing.config.get"
	methodSubDevRegister        = "thing.sub.register"
	methodTopoAdd               = "thing.topo.add"
	methodTopoDelete            = "thing.topo.delete"
	methodTopoGet               = "thing.topo.get"
	methodListFound             = "thing.list.found"
)
