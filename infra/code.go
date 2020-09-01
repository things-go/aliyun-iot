package infra

import (
	"fmt"
)

// 通用公共错误
// see https://help.aliyun.com/document_detail/120329.html?spm=a2c4g.11186623.6.677.efd1a684yZutHX
const (
	CodeSuccess                       = 200  // 请求成功
	CodeRequestError                  = 400  // 请求错误
	CodeRequestParamsError            = 460  // 请求参数错误,主要设备上报的数据为空，或参数格式错误、参数的数量超过限制等原因
	CodeRequestTooMany                = 429  // 请求过于频繁,触发系统限流
	CodeSystemUnknownException        = 500  // 系统发生未知异常。
	CodeQueryProductInfoFailed        = 5005 // 查询产品信息失败。
	CodeQueryProductInfoFailed1       = 6250 // 查询产品信息失败。
	CodeQueryLoRaWANProductInfoFailed = 5244 // 查询LoRaWAN类型产品的元信息失败。
	CodeDeviceNotFound                = 6100 // 查询设备信息时,未查询到指定设备信息. 子设备不存在
	CodeParseTopicFailed              = 6203 // 解析Topic时失败。
	CodeDeviceDisabled                = 6204 // 设备已被禁用，不能对设备进行操作。
	CodeRawModelMissMethod            = 6450 // 自定义/透传格式数据经过脚本解析为Alink标准格式数据后，无method。
	CodeSystemException               = 6760 // 系统异常。
	CodeTimeout                       = 100000
)

// 数据解析公共错误码
const (
	CodeDpScriptEmpty          = 26001 // 执行数据解析时，获取的脚本内容为空
	CodeDpScriptIssue          = 26002 // 脚本执行正常，但脚本内容有问题，如脚本中语法错误
	CodeDpScriptFailed         = 26006 // 脚本执行正常，脚本内容有误.脚本中，要求有protocolToRawData和rawDataToProtocol这两个服务。如果缺失，会报错。
	CodeDpScriptFailed1        = 26007 // 脚本执行正常，但返回结果不符合格式要求
	CodeDpScriptRequestTooMuch = 26010 // 请求过于频繁，导致被限流。
)

// 设备身份注册相关错误码
const (
	// 子设备身份注册错误码
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/sub/register
	// 错误码 460、5005、5244、500、6288、6100、6619、6292、6203
	// 以下为特有错误码
	CodeDevDynamicRegisterNotEnable = 6288 // 设备动态注册未打开
	CodeDevHasBindGateway           = 6619

	// 直接设备一型一密动态注册错误码
	// 错误码：460、6250、6288、6600、6289、500、6292
	// 以下为特有错误码

	CodeDevSignMethodNotSupport = 6292 // 校验签名时，发现传入的签名方法不支持。
	CodeDevSignVerifyFailed     = 6600 // 签名校验失败
	CodeDevAlreadyActive        = 6289 // 一型一密动态注册直连设备时，发现设备已激活。
)

// 设备拓扑关系相关错误码
const (
	// 添加设备拓扑关系
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/topo/add
	// 错误码：460、429、6402、6100、401、6204、6400、6203
	// 以下为特有错误码
	CodeTopoRequestAuthError        = 401  // 添加拓扑关系时，校验签名信息失败。
	CodeTopoRelationCannotAddBySelf = 6402 // 网关与子设备是同一个设备。添加拓扑关系时，不能把当前网关作为子设备添加到当前网关下。
	CodeTopoRelationSubDevOverLimit = 6400 // 为网关添加的子设备数量超过限制。

	// 删除拓扑关系
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/topo/delete
	// 错误码：460、429、6100、6401、6203
	CodeTopoRelationNotExist = 6401 // 网关和子设备没有拓扑关系
	// 获取拓扑关系
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/topo/get
	// 错误码： 460、429、500、6203

	// 网关上报发现子设备
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/list/found
	// 错误码：460、500、6250、6280、6203
	CodeTopoSubDevNameNotStandard = 6280 // 网关上报的子设备名称不符合规范。
)

// 子设备上下线相关错误码
const (
	// 子设备上线
	// 请求Topic：/ext/session/${productKey}/${deviceName}/combine/login
	// 错误码： 460、429、6100、6204、6287、6401、500

	// 子设备主动下线异常
	// 接收消息的网关Topic：/ext/session/{productKey}/{deviceName}/combine/logout_reply
	// 错误码：460、520、500

	// 子设备被踢下线
	// 接收消息的网关Topic：/ext/error/{productKey}/{deviceName}
	// 错误码：427、521、522、6401

	// 子设备发送消息失败
	// 接收消息的网关Topic：/ext/error/{productKey}/{deviceName}
	// 错误码：520

	// 以下子设备上、下线特有错误码
	CodeSubDevLoginDump           = 527 // 设备重复登录。有使用相同设备证书信息的设备连接物联网平台，导致当前连接被断开。
	CodeSubDevTooManyUnderGateway = 428 // 网关下同时在线子设备过多
	// 子设备会话错误。 子设备会话不存在，可能子设备没有上线，也可能已经被下线。
	// 子设备会话在线，但是并不是通过当前网关会话上线的。
	CodeSubDevSessionError = 520
	CodeSubDevDeleted      = 521  // 子设备已被删除
	CodeSubDevDisabled     = 522  // 子设备已被禁用
	CodeSubDevSignInvalid  = 6287 // 子设备密码或签名错误
)

// 设备属性、事件、服务相关错误码
const (
	// 设备上报属性
	// 请求Topic（透传数据格式）： /sys/{productKey}/{deviceName}/thing/model/up_raw
	// 请求Topic（Alink数据格式）：/sys/{productKey}/{deviceName}/thing/event/property/post
	// 错误码：460、500、6250、6203、6207、6313、6300、6320、6321、6326、6301、6302、6317、6323、6316、6306、
	//  6307、6322、6308、6309、6310、6311、6312、6324、6328、6325、6200、6201、26001、26002、26006、26007

	// 设备上报事件
	// 请求Topic（透传格式数据）： /sys/{productKey}/{deviceName}/thing/model/up_raw
	// 请求Topic（Alink格式数据）：/sys/{productKey}/{deviceName}/thing/event/{tsl.identifier}/post
	// 错误码：460、500、6250、6203、6207、6313、6300、6320、6321、6326、6301、6302、6317、6323、6316、6306、
	//   6307、6322、6308、6309、6310、6311、6312、6324、6328、6325、6200、6201、26001、26002、26006、26007

	// 网关批量上报子设备数据
	// 请求Topic（透传格式数据）：/sys/{productKey}/{deviceName}/thing/model/up_raw
	// 请求Topic（Alink格式数据）：/sys/{productKey}/{deviceName}/thing/event/property/pack/post
	// 错误码：460、6401、6106、6357、6356、6100、6207、6313、6300、6320、6321、6326、6301、6302、6317、6323、6316、6306、
	//	 6307、6322、6308、6309、6310、6311、6312、6324、6328、6325、6200、6201、26001、26002、26006、26007
	CodeReportTooMuch            = 6106 // 上报的属性数据过多。设备一次上报的有效属性个数不能超过200个。
	CodeSubDevReportTooMuch      = 6357 // 子设备数据过多。网关代替子设备上报数据，一次上报最多可包含20个子设备的数据
	CodeSubDevReportEventTooMuch = 6356 // 上报的事件数据过多。网关代替子设备上报数据，一次上报的事件个数不可超过200
)

// 设备期望属性值相关错误码
const (
	// 设备获取期望属性值
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/property/desired/get
	// 错误码： 460、6104、6661、500
	CodePropertyTooMuch            = 6104 // 请求中包含的属性个数过多。一次请求可包含的属性个数不能超过200个。
	CodeQueryDesiredPropertyFailed = 6661 // 查询期望属性失败。系统异常。

	// 设备清空期望属性值
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/property/desired/delete
	// 错误码：460、6104、6661、500、6207、6313、6300、6320、6321、6326、6301、6302、6317、6323、
	// 6316、6306、6307、6322、6308、6309、6310、6311、6312、6324、6328、6325
)

// 设备标签相关错误码
// 获取TSL模板失败错误码
// 设备请求固件信息失败
// 设备请求配置信息失败
const (
	// 设备上报标签信息
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/deviceinfo/update
	// 错误码：460、6100

	// 设备删除标签信息
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/deviceinfo/delete
	// 错误码： 460、500

	// 获取TSL模板失败错误码
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/dsltemplate/get
	// 错误码：460、5159、5160、5161

	// 设备请求固件信息失败
	// 请求Topic：/ota/device/request/${YourProductKey}/${YourDeviceName}
	// 错误码：429、9112、500
	CodeNotFindDevInfo = 9112 // 未查询到指定的设备信息。

	// 设备请求配置信息失败
	// 请求Topic：/sys/{productKey}/{deviceName}/thing/config/get
	// 错误码：460、500、6713、6710
	CodeRemoteConfigServiceDisabled = 6713 // 远程配置服务不可用。该产品的远程配置开关未打开。
	CodeNotFindRemoteConfig         = 6710 // 未查询到远程配置信息
)

// CodeError code 错误
type CodeError struct {
	code    int
	message string
}

// NewCodeError 生成code错误
func NewCodeError(code int, message string) error {
	return &CodeError{code, message}
}

// Error 实现error接口
func (sf *CodeError) Error() string {
	return fmt.Sprintf("code: %d - message: %s", sf.code, sf.message)
}

// Code unwrap then got code
func (sf *CodeError) Code() int {
	return sf.code
}
