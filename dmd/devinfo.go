package dmd

// DevInfoLabelCoordinateKey 地理位置标签
const DevInfoLabelCoordinateKey = "coordinate" //

// DevInfoLabelUpdate 更新设备标签的键值对
type DevInfoLabelUpdate struct {
	AttrKey   string `json:"attrKey"`
	AttrValue string `json:"attrValue"`
}

// DevInfoLabelDelete 删除设备标答的键
type DevInfoLabelDelete struct {
	AttrKey string `json:"attrKey"`
}
