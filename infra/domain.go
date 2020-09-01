package infra

// HTTPCloudDomain http 域名
var HTTPCloudDomain = []string{
	"iot-auth.cn-shanghai.aliyuncs.com",    /* Shanghai */
	"iot-auth.ap-southeast-1.aliyuncs.com", /* Singapore */
	"iot-auth.ap-northeast-1.aliyuncs.com", /* Japan */
	"iot-auth.us-west-1.aliyuncs.com",      /* America */
	"iot-auth.eu-central-1.aliyuncs.com",   /* Germany */
}

// MQTTCloudDomain mqtt 域名
var MQTTCloudDomain = []string{
	"iot-as-mqtt.cn-shanghai.aliyuncs.com",    /* Shanghai */
	"iot-as-mqtt.ap-southeast-1.aliyuncs.com", /* Singapore */
	"iot-as-mqtt.ap-northeast-1.aliyuncs.com", /* Japan */
	"iot-as-mqtt.us-west-1.aliyuncs.com",      /* America */
	"iot-as-mqtt.eu-central-1.aliyuncs.com",   /* Germany */
}

// CloudRegion MQTT,HTPP云端地域
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

// CloudRegionDomain 云端域信息
type CloudRegionDomain struct {
	Region       CloudRegion
	CustomDomain string // address:port,采用CloudRegionCustom需要定义此字段
}

// MetaInfo 产品与设备三元组
type MetaInfo struct {
	ProductKey    string
	ProductSecret string
	DeviceName    string
	DeviceSecret  string
}
