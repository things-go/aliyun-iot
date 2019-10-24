package model

import (
	"fmt"
	"strings"
)

const (
	SEP = "/"
)

// URI 前缀定义
const (
	URISysPrefix             = "/sys/%s/%s/"
	URIExtSessionPrefix      = "/ext/session/%s/%s/"
	URIExtNtpPrefix          = "/ext/ntp/%s/%s/"
	URIExtErrorPrefix        = "/ext/error/%s/%s"
	URIOtaDeviceInformPrefix = "/ota/device/inform/%s/%s"

	URIReplySuffix = "reply"
)

// URI thing定义
const (
	/* Model Raw From Cloud To Local Request And Response*/
	URIThingModelDownRaw      = "thing/model/down_raw"
	URIThingModelDownRawReply = "thing/model/down_raw_reply"
	URIThingModelUpRaw        = "thing/model/up_raw"
	URIThingModelUpRawReply   = "thing/model/up_raw_reply"

	URIRRPCRequestWildcard = "rrpc/request/+"

	/* service From Cloud To Local Request And Response*/
	URIThingServicePropertySet      = "thing/service/property/set"
	URIThingServicePropertySetReply = "thing/service/property/set_reply"
	URIThingServicePropertyGet      = "thing/service/property/get"
	URIThingServicePropertyGetReply = "thing/service/property/get_reply"
	URIThingServiceRequestWildcard  = "thing/service/+"
	URIThingServiceRequestWildcard2 = "thing/service/#"
	URIThingServiceRequest          = "thing/service/%s"
	URIThingServiceResponse         = "thing/service/%s_reply"

	/* event From Local To Cloud Request And Response*/
	URIThingEventPropertyPost      = "thing/event/property/post"
	URIThingEventPropertyPostReply = "thing/event/property/post_reply"
	URIThingEventPost              = "thing/event/%s/post"
	URIThingEventPostReply         = "thing/event/%s/post_reply"
	URIThingEventPostReplyWildcard = "thing/event/+/post_reply"

	/* device info From Local To Cloud Request And Response */
	URIThingDeviceInfoUpdate      = "thing/deviceinfo/update"
	URIThingDeviceInfoUpdateReply = "thing/deviceinfo/update_reply"
	URIThingDeviceInfoDelete      = "thing/deviceinfo/delete"
	URIThingDeviceInfoDeleteReply = "thing/deviceinfo/delete_reply"

	/* dsl template From Local To Cloud Request And Response */
	URIThingDslTemplateGet      = "thing/dsltemplate/get"
	URIThingDslTemplateGetReply = "thing/dsltemplate/get_reply"

	/* dynamic tsl From Local To Cloud Request And Response */
	URIThingDynamicTslGet      = "thing/dynamicTsl/get"
	URIThingDynamicTslGetReply = "thing/dynamicTsl/get_reply"

	/* ntp From Local To Cloud Request And Response */
	URINtpRequest  = "request"
	URINtpResponse = "response"

	//! config
	URIThingConfigGet       = "thing/config/get"
	URIThingConfigGetReply  = "thing/config/get_reply"
	URIThingConfigPush      = "thing/config/push"
	URIThingConfigPushReply = "thing/config/push_reply"
)

// URIService 生成URI定义符
func URIService(prefix, name, productKey, deviceName string) string {
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
