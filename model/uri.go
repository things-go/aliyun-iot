package model

import (
	"fmt"
	"strings"
)

// uri 前缀定义
const (
	UriSysPrefix             = "/sys/%s/%s/"
	UriExtSessionPrefix      = "/ext/session/%s/%s/"
	UriExtNtpPrefix          = "/ext/ntp/%s/%s/"
	UriExtErrorPrefix        = "/ext/error/%s/%s"
	UriOtaDeviceInformPrefix = "/ota/device/inform/%s/%s"

	UriReplySuffix = "reply"
)

const (
	/* Model Raw From Cloud To Local Request And Response*/
	UriThingModelDownRaw      = "thing/model/down_raw"
	UriThingModelDownRawReply = "thing/model/down_raw_reply"
	UriThingModelUpRaw        = "thing/model/up_raw"
	UriThingModelUpRawReply   = "thing/model/up_raw_reply"

	UriRRPCRequestWildcard = "rrpc/request/+"

	/* service From Cloud To Local Request And Response*/
	UriThingServicePropertySet      = "thing/service/property/set"
	UriThingServicePropertySetReply = "thing/service/property/set_reply"
	UriThingServicePropertyGet      = "thing/service/property/get"
	UriThingServicePropertyGetReply = "thing/service/property/get_reply"
	UriThingServiceRequestWildcard  = "thing/service/+"
	UriThingServiceRequestWildcard2 = "thing/service/#"
	UriThingServiceRequest          = "thing/service/%s"
	UriThingServiceResponse         = "thing/service/%s_reply"

	/* event From Local To Cloud Request And Response*/
	UriThingEventPropertyPost      = "thing/event/property/post"
	UriThingEventPropertyPostReply = "thing/event/property/post_reply"
	UriThingEventPost              = "thing/event/%s/post"
	UriThingEventPostReply         = "thing/event/%s/post_reply"
	UriThingEventPostReplyWildcard = "thing/event/+/post_reply"

	/* device info From Local To Cloud Request And Response */
	UriThingDeviceInfoUpdate      = "thing/deviceinfo/update"
	UriThingDeviceInfoUpdateReply = "thing/deviceinfo/update_reply"
	UriThingDeviceInfoDelete      = "thing/deviceinfo/delete"
	UriThingDeviceInfoDeleteReply = "thing/deviceinfo/delete_reply"

	/* dsl template From Local To Cloud Request And Response */
	UriThingDslTemplateGet      = "thing/dsltemplate/get"
	UriThingDslTemplateGetReply = "thing/dsltemplate/get_reply"

	/* dynamic tsl From Local To Cloud Request And Response */
	UriThingDynamicTslGet      = "thing/dynamicTsl/get"
	UriThingDynamicTslGetReply = "thing/dynamicTsl/get_reply"

	/* ntp From Local To Cloud Request And Response */
	UriNtpRequest  = "request"
	UriNtpResponse = "response"

	//! config
	UriThingConfigGet       = "thing/config/get"
	UriThingConfigGetReply  = "thing/config/get_reply"
	UriThingConfigPush      = "thing/config/push"
	UriThingConfigPushReply = "thing/config/push_reply"
)

// UriService 生成uri定义符
func UriService(prefix, name, productKey, deviceName string) string {
	str := strings.Builder{}
	str.Grow(len(prefix) + len(name) + len(productKey) + len(deviceName))
	if prefix != "" {
		str.WriteString(fmt.Sprintf(prefix, productKey, deviceName))
	}
	if name != "" {
		str.WriteString(name)
	}
	return str.String()
}
