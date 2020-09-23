// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package infra

// HTTPCloudDomain http 域名
var HTTPCloudDomain = []string{
	"iot-auth.cn-shanghai.aliyuncs.com",    // Shanghai
	"iot-auth.ap-southeast-1.aliyuncs.com", // Singapore
	"iot-auth.ap-northeast-1.aliyuncs.com", // Japan
	"iot-auth.us-west-1.aliyuncs.com",      // America
	"iot-auth.eu-central-1.aliyuncs.com",   // Germany
}

// MQTTCloudDomain mqtt 域名
var MQTTCloudDomain = []string{
	"iot-as-mqtt.cn-shanghai.aliyuncs.com",    // Shanghai
	"iot-as-mqtt.ap-southeast-1.aliyuncs.com", // Singapore
	"iot-as-mqtt.ap-northeast-1.aliyuncs.com", // Japan
	"iot-as-mqtt.us-west-1.aliyuncs.com",      // America
	"iot-as-mqtt.eu-central-1.aliyuncs.com",   // Germany
}

// CloudRegion MQTT, HTTP云端地域定义
type CloudRegion byte

// CloudRegionRegion 云平台地域定义
const (
	CloudRegionShangHai CloudRegion = iota
	CloudRegionSingapore
	CloudRegionJapan
	CloudRegionAmerica
	CloudRegionGermany
	CloudRegionCustom
)

// CloudRegionDomain 云端域信息
type CloudRegionDomain struct {
	Region       CloudRegion
	CustomDomain string // address:port,当Region为CloudRegionCustom需要定义此字段,其它无效
}
