package model

const (
	UriReplySuffix = "reply"

	UriSysPrefix        = "/sys/%s/%s/"
	UriExtSessionPrefix = "/ext/session/%s/%s/"
	UriExtNtpPrefix     = "/ext/ntp/%s/%s/"
	UriExtErrorPrefix   = "/ext/error/%s/%s"
	UriOtaDeviceInform  = "/ota/device/inform/%s/%s"

	/* Model Raw From Cloud To Local Request And Response*/
	UriThingModelDownRaw      = "thing/model/downraw"
	UriThingModelDownRawReply = "thing/model/downrawreply"
	UriThingModelUpRaw        = "thing/model/upraw"
	UriThingModelUpRawReply   = "thing/model/uprawreply"

	UriRRPCRequestWildcard = "rrpc/request/+"

	/* service From Cloud To Local Request And Response*/
	UriThingServicePropertySet      = "thing/service/property/set"
	UriThingServicePropertySetReply = "thing/service/property/setreply"
	UriThingServicePropertyGet      = "thing/service/property/get"
	UriThingServicePropertyGetReply = "thing/service/property/getreply"
	UriThingServiceRequestWildcard  = "thing/service/+"
	UriThingServiceRequestWildcard2 = "thing/service/#"
	UriThingServiceRequest          = "thing/service/%s"
	UriThingServiceResponse         = "thing/service/%.*sreply"

	/* event From Local To Cloud Request And Response*/
	UriThingEventPropertyPost      = "thing/event/property/post"
	UriThingEventPropertyPostReply = "thing/event/property/postreply"
	UriThingEventPost              = "thing/event/%.*s/post"
	UriThingEventPostReply         = "thing/event/%s/postreply"
	UriThingEventPostReplyWildcard = "thing/event/+/postreply"

	/* device info From Local To Cloud Request And Response */
	UriThingDeviceInfoUpdate      = "thing/deviceinfo/update"
	UriThingDeviceInfoUpdateReply = "thing/deviceinfo/updatereply"
	UriThingDeviceInfoDelete      = "thing/deviceinfo/delete"
	UriThingDeviceInfoDeleteReply = "thing/deviceinfo/deletereply"

	/* dsl template From Local To Cloud Request And Response */
	UriThingDslTemplateGet      = "thing/dsltemplate/get"
	UriThingDslTemplateGetReply = "thing/dsltemplate/getreply"

	/* dynamic tsl From Local To Cloud Request And Response */
	UriThingDynamicTslGet      = "thing/dynamicTsl/get"
	UriThingDynamicTslGetReply = "thing/dynamicTsl/getreply"

	/* ntp From Local To Cloud Request And Response */
	UriNtpRequest  = "request"
	UriNtpResponse = "response"

	//! config
	UriThingConfigGet       = "thing/config/get"
	UriThingConfigGetReply  = "thing/config/getreply"
	UriThingConfigPush      = "thing/config/push"
	UriThingConfigPushReply = "thing/config/pushreply"
)
