package model

// 平台通信版本
const (
	Version = "1.0"
)

// method 定义
const (
	methodPropertyPost     = "thing.event.property.post"
	methodDeviceInfoUpdate = "thing.deviceinfo.update"
	methodDeviceInfoDelete = "thing.deviceinfo.delete"
	//methodUpRaw            = "thing.model.up_raw"
	methodEventPostFormat = "thing.event.%s.post"
	methodDslTemplateGet  = "thing.dsltemplate.get"
	methodDynamicTslGet   = "thing.dynamicTsl.get"

	methodSubDevRegister = "thing.sub.register"
	methodSubDevDisable  = "thing.disable"
	methodSubDevEnable   = "thing.enable"
	methodSubDevDelete   = "thing.delete"
	methodTopoAdd        = "thing.topo.add"
	methodTopoDelete     = "thing.topo.delete"
	methodTopoGet        = "thing.topo.get"
)
