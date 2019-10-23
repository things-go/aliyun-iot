package infra

// http 域名
var HTTPCloudDomain = []string{
	"iot-auth.cn-shanghai.aliyuncs.com",    /* Shanghai */
	"iot-auth.ap-southeast-1.aliyuncs.com", /* Singapore */
	"iot-auth.ap-northeast-1.aliyuncs.com", /* Japan */
	"iot-auth.us-west-1.aliyuncs.com",      /* America */
	"iot-auth.eu-central-1.aliyuncs.com",   /* Germany */
}

// mqtt 域名
var MQTTCloudDomain = []string{
	"iot-as-mqtt.cn-shanghai.aliyuncs.com",    /* Shanghai */
	"iot-as-mqtt.ap-southeast-1.aliyuncs.com", /* Singapore */
	"iot-as-mqtt.ap-northeast-1.aliyuncs.com", /* Japan */
	"iot-as-mqtt.us-west-1.aliyuncs.com",      /* America */
	"iot-as-mqtt.eu-central-1.aliyuncs.com",   /* Germany */
}

// CloudRegionRegion MQTT,HTPP云端地域
type CloudRegion byte

// 云平台地域定义CloudRegionRegion
const (
	CloudRegionShangHai CloudRegion = iota
	CloudRegionSingapore
	CloudRegionJapan
	CloudRegionAmerica
	CloudRegionGermany
	CloudRegionCustom
)
