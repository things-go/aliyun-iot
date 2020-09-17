package uri

import (
	"fmt"
	"strings"
)

// 分隔符定义
const (
	Sep = "/"
)

// URI 前缀定义
const (
	TopicPrefix = "/topic"

	SysPrefix        = "/sys/%s/%s/"
	ExtSessionPrefix = "/ext/session/%s/%s/"
	ExtNtpPrefix     = "/ext/ntp/%s/%s/"
	ExtErrorPrefix   = "/ext/errorf/%s/%s"

	ReplySuffix = "reply"
)

// RRPC URI定义
const (
	//  系统RRPC调用
	RRPCResponsePrefix     = "rrpc/response/%s"
	RRPCRequestWildcardOne = "rrpc/request/+"

	// 自定义RRPC
	ExtRRPCPrefix            = "/ext/rrpc/%s"
	ExtRRPCWildcardOnePrefix = "/ext/rrpc/+"
)

// ota uri定义
const (
	OtaDeviceInformPrefix    = "/ota/device/inform/%s/%s"
	OtaDeviceUpgradePrefix   = "/ota/device/upgrade/%s/%s"
	OtaDeviceProcessPrefix   = "/ota/device/progress/%s/%s"
	ThingOtaFirmwareGet      = "thing/ota/firmware/get"
	ThingOtaFirmwareGetReply = "thing/ota/firmware/get_reply"
)

// 设备URI 定义
const (
	// 透传数据上行,下行云端
	ThingModelDownRaw      = "thing/model/down_raw"
	ThingModelDownRawReply = "thing/model/down_raw_reply"
	ThingModelUpRaw        = "thing/model/up_raw"
	ThingModelUpRawReply   = "thing/model/up_raw_reply"

	// 事件上行,下行云端
	ThingEventPropertyPost             = "thing/event/property/post"
	ThingEventPropertyPostReply        = "thing/event/property/post_reply"
	ThingEventPost                     = "thing/event/%s/post"
	ThingEventPostReply                = "thing/event/%s/post_reply"
	ThingEventPostReplyWildcardOne     = "thing/event/+/post_reply"
	ThingEventPropertyHistoryPost      = "thing/event/property/history/post"
	ThingEventPropertyHistoryPostReply = "thing/event/property/history/post_reply"

	// 设备信息上行,下行云端
	ThingDeviceInfoUpdate      = "thing/deviceinfo/update"
	ThingDeviceInfoUpdateReply = "thing/deviceinfo/update_reply"
	ThingDeviceInfoDelete      = "thing/deviceinfo/delete"
	ThingDeviceInfoDeleteReply = "thing/deviceinfo/delete_reply"

	// 期望属性值上行,下行云端
	ThingDesiredPropertyGet         = "thing/property/desired/get"
	ThingDesiredPropertyGetReply    = "thing/property/desired/get_reply"
	ThingDesiredPropertyDelete      = "thing/property/desired/delete"
	ThingDesiredPropertyDeleteReply = "thing/property/desired/delete_reply"

	// 服务调用上行,下行云端
	ThingServicePropertySet         = "thing/service/property/set"
	ThingServicePropertySetReply    = "thing/service/property/set_reply"
	ThingServiceRequest             = "thing/service/%s"
	ThingServiceResponse            = "thing/service/%s_reply"
	ThingServiceRequestWildcardOne  = "thing/service/+"
	ThingServiceRequestWildcardSome = "thing/service/#"

	// dsl template From Local To Cloud Request And ResponseRawData
	ThingDslTemplateGet      = "thing/dsltemplate/get"
	ThingDslTemplateGetReply = "thing/dsltemplate/get_reply"

	// dynamic tsl From Local To Cloud Request And ResponseRawData
	ThingDynamicTslGet      = "thing/dynamicTsl/get"
	ThingDynamicTslGetReply = "thing/dynamicTsl/get_reply"

	// ntp From Local To Cloud Request And ResponseRawData
	NtpRequest  = "request"
	NtpResponse = "response"

	// config
	ThingConfigGet       = "thing/config/get"
	ThingConfigGetReply  = "thing/config/get_reply"
	ThingConfigPush      = "thing/config/push"
	ThingConfigPushReply = "thing/config/push_reply"

	// log
	ThingConfigLogGet      = "thing/config/log/get"
	ThingConfigLogGetReply = "thing/config/log/get_reply"
	ThingConfigLogPush     = "thing/config/log/push"
	ThingLogPost           = "thing/config/log/post"
	ThingLogPostReply      = "thing/config/log/post_reply"
)

// 网关URI定义
const (
	ThingEventPropertyPackPost      = "thing/event/property/pack/post"
	ThingEventPropertyPackPostReply = "thing/event/property/pack/post_reply"
	// 子设备动态注册
	ThingSubRegister      = "thing/sub/register"
	ThingSubRegisterReply = "thing/sub/register_reply"

	// 子设备登录
	CombineLogin            = "combine/login"
	CombineLoginReply       = "combine/login_reply"
	CombineBatchLogin       = "combine/batch_login"
	CombineBatchLoginReply  = "combine/batch_login_reply"
	CombineBatchLogout      = "combine/batch_logout"
	CombineBatchLogoutReply = "combine/batch_logout_reply"
	CombineLogout           = "combine/logout"
	CombineLogoutReply      = "combine/logout_reply"

	// 网关网络拓扑
	ThingTopoAdd         = "thing/topo/add"
	ThingTopoAddReply    = "thing/topo/add_reply"
	ThingTopoDelete      = "thing/topo/delete"
	ThingTopoDeleteReply = "thing/topo/delete_reply"
	ThingTopoGet         = "thing/topo/get"
	ThingTopoGetReply    = "thing/topo/get_reply"
	ThingListFound       = "thing/list/found"
	ThingListFoundReply  = "thing/list/found_reply"
	ThingTopoAddNotify   = "thing/topo/add/notify"
	ThingTopoChange      = "thing/topo/change"
	ThingDisable         = "thing/disable"
	ThingDisableReply    = "thing/disable_reply"
	ThingEnable          = "thing/enable"
	ThingEnableReply     = "thing/enable_reply"
	ThingDelete          = "thing/delete"
	ThingDeleteReply     = "thing/delete_reply"

	// 设备网络状态
	ThingDiagPost      = "_thing/diag/post"
	ThingDiagPostReply = "_thing/diag/post_reply"
)

// URI 生成URI
// prefix: 主题前缀
// name: 可为空字符串
// productKey: 产品key
// deviceName: 设备名
// inName的作用是有一些需要格式化到name里.
func URI(prefix, name, productKey, deviceName string, inName ...string) string {
	str := strings.Builder{}
	str.Grow(len(prefix) + len(name) + len(productKey) + len(deviceName))
	if prefix != "" {
		str.WriteString(fmt.Sprintf(prefix, productKey, deviceName))
	}
	if name != "" {
		if len(inName) > 0 {
			str.WriteString(fmt.Sprintf(name, inName[0]))
		} else {
			str.WriteString(name)
		}
	}
	return str.String()
}

// ReplyWithRequestURI 根据请求的URI生成应答的URI
func ReplyWithRequestURI(uri string) string {
	return uri + "_" + ReplySuffix
}

// Spilt 分隔URI,会去除左边的SEP分隔符
func Spilt(uri string) []string {
	return strings.Split(strings.TrimLeft(uri, Sep), Sep)
}

// ExtRRPC 生成Ext RRPC URI
func ExtRRPC(messageID, _uri string) string {
	return fmt.Sprintf(ExtRRPCPrefix, messageID) + _uri
}

// ExtRRPCWildcardOne 生成Ext RRPC 订阅 URI
func ExtRRPCWildcardOne(_uri string) string {
	return ExtRRPCWildcardOnePrefix + _uri
}
