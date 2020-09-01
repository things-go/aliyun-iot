package dataflow

import (
	"errors"
	"strings"
)

// 错误定义
var (
	ErrTopicInvalid = errors.New("topic invalid")
)

// 主题分隔符
const (
	SEP = "/"
)

// TopicType 主题类型
// see https://help.aliyun.com/document_detail/73736.html?spm=a2c4g.11186623.6.630.19485a105hJAhr
type TopicType int

// topic defined
const (
	// topic: /as/mqtt/status/{productKey}/{deviceName}
	TopicStatusWildcard = "/as/mqtt/status/+/+"
	// topic: /{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post
	// topic: /{productKey}/{deviceName}/thing/event/property/post
	TopicEventWildcard = "/+/+/thing/event/+/post"
	// topic: /{productKey}/{deviceName}/thing/lifecycle
	TopicLifecycleWildcard = "/+/+/thing/lifecycle"
	// topic: /{productKey}/{deviceName}/thing/topo/lifecycle
	TopicTopoLifecycleWildcard = "/+/+/thing/topo/lifecycle"
	// topic: /{productKey}/{deviceName}/thing/list/found
	TopicSubDeviceFoundWildcard = "/+/+/thing/list/found"
	// topic: /{productKey}/{deviceName}/thing/downlink/reply/message
	TopicDownLinkReplyWildcard = "/+/+/thing/downlink/reply/message"
	// topic: /sys/{productKey}/{deviceName}/thing/event/property/history/post
	// topic: /sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/history/post
	TopicHistoryEventWildcard = "/sys/+/+/thing/event/+/history/post"
	// topic: /sys/${productKey}/${deviceName}/ota/upgrade
	TopicOtaUpgrade = "/sys/+/+/ota/upgrade"
)

// TopicInfo 主题信息,解析topic后获得的信息
type TopicInfo struct {
	ProductKey string
	DeviceName string
	EventID    string // 仅TopicTypeEvent有效
}

// ParseTopicStatus parse status topic
// topic: /as/mqtt/status/{productKey}/{deviceName}
func ParseTopicStatus(topic string) (ti TopicInfo, err error) {
	if topic == "" {
		err = ErrTopicInvalid
		return
	}
	s := strings.Split(strings.TrimLeft(topic, SEP), SEP)
	if len(s) != 5 {
		err = ErrTopicInvalid
		return
	}
	return TopicInfo{ProductKey: s[3], DeviceName: s[4]}, nil
}

// ParseTopicEvent parse event topic
// topic: /{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/post
// topic: /{productKey}/{deviceName}/thing/event/property/post
func ParseTopicEvent(topic string) (ti TopicInfo, err error) {
	return parseTopic(topic, false, true, 6)
}

// ParseTopicLifecycle parse lifecycle topic
// topic: /{productKey}/{deviceName}/thing/lifecycle
func ParseTopicLifecycle(topic string) (ti TopicInfo, err error) {
	return parseTopic(topic, false, false, 4)
}

// ParseTopicTopoLifecycle parse topo lifecycle topic
// topic: /{productKey}/{deviceName}/thing/topo/lifecycle
func ParseTopicTopoLifecycle(topic string) (ti TopicInfo, err error) {
	return parseTopic(topic, false, false, 5)
}

// ParseTopicListFound parse listfound topic
// topic: /{productKey}/{deviceName}/thing/list/found
func ParseTopicListFound(topic string) (ti TopicInfo, err error) {
	return parseTopic(topic, false, false, 5)
}

// ParseTopicDownLinkReply parse downlink reply topic
// topic: /{productKey}/{deviceName}/thing/downlink/reply/message
func ParseTopicDownLinkReply(topic string) (ti TopicInfo, err error) {
	return parseTopic(topic, false, false, 6)
}

// ParseTopicHistoryEvent parse history event topic
// topic: /sys/{productKey}/{deviceName}/thing/event/property/history/post
// topic: /sys/{productKey}/{deviceName}/thing/event/{tsl.event.identifier}/history/post
func ParseTopicHistoryEvent(topic string) (ti TopicInfo, err error) {
	return parseTopic(topic, true, true, 8)
}

// ParseTopicOtaUpgrade parse ota upgrade topic
// topic: /sys/${productKey}/${deviceName}/ota/upgrade
func ParseTopicOtaUpgrade(topic string) (ti TopicInfo, err error) {
	return parseTopic(topic, true, false, 5)
}

func parseTopic(topic string, isSys, isEvent bool, match int) (ti TopicInfo, err error) {
	if topic == "" {
		err = ErrTopicInvalid
		return
	}
	s := strings.Split(strings.TrimLeft(topic, SEP), SEP)
	if len(s) != match {
		err = ErrTopicInvalid
		return
	}
	offset := 0
	if isSys {
		offset++
	}
	evt := ""
	if isEvent {
		evt = s[4+offset]
	}

	return TopicInfo{
		s[offset],
		s[1+offset],
		evt,
	}, nil
}
