package infra

const (
	CodeSuccess            = 200
	CodeRequestError       = 400
	CodeRequestAuthError   = 401 // 签名校验失败
	CodeRequestParamsError = 460 // 请求参数错误
	CodeRequestTooMany     = 429 //  请求过多
	CodeNoActiveSession    = 520
	CodeTimeout            = 100000
)

const (
	CodeSubDevOnlineRateLimit       = 429  // 单个设备认证过于频繁被限流
	CodeTooManySubDevUnderGateway   = 428  // 网关下同时在线子设备过多
	CodeTopoRelationNotExist        = 6401 // 网关和子设备没有拓扑关系
	CodeTopoRelationCannotAddBySelf = 6402 // 设备不能将自己添加为子设备
	CodeDeviceNotFound              = 6100 // 子设备不存在
	CodeDeviceDeleted               = 521  // 子设备已被禁用
	CodeDeviceForbidden             = 522  // 子设备已被禁用
	CodeInvalidSign                 = 6287 // 子设备密码或签名错误

	// 直接设备一型一密动态注册错误码
	CodeDeviceDynamicRegisterNotEnable = 6288 // 设备动态注册未打开
	CodeDeviceSignMethodNotSupport     = 6292 // 设备签名方法不支持
	CodeDeviceSignVerifyFailed         = 6600 // 签名校验失败
	CodeDeviceAlreadyActive            = 6289 // 设备已激活
)
