package dm

// DevInfoLabelCoordinateKey 地理位置标签
const DevInfoLabelCoordinateKey = "coordinate" //

// DeviceInfoLabel 更新设备标签的键值对
type DeviceInfoLabel struct {
	AttrKey   string `json:"attrKey"`
	AttrValue string `json:"attrValue"`
}

// DeviceLabelKey 删除设备标答的键
type DeviceLabelKey struct {
	AttrKey string `json:"attrKey"`
}
