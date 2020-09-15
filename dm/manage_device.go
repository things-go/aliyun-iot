package dm

import (
	"sync"

	"github.com/thinkgos/aliyun-iot/infra"
)

// DevStatus 设备状态
type DevStatus byte

// 设备状态
const (
	DevStatusUnauthorized DevStatus = iota // Subdev Created
	DevStatusAuthorized                    // Receive Topo Add Notify
	DevStatusRegistered                    // Receive Subdev Registered
	DevStatusAttached                      // Receive Subdev Topo Add Reply
	DevStatusLogined                       // Receive Subdev Login Reply
	DevStatusOnline                        // After All Topic Subscribed
)

// DevMgr 设备管理
type DevMgr struct {
	root  DevNode // 网关设备节点或独立设备节点信息
	rw    sync.RWMutex
	nodes map[string]*DevNode
}

// DevNode 设备节点
type DevNode struct {
	productKey   string
	deviceName   string
	deviceSecret string
	avail        bool
	status       DevStatus
	ext          interface{}
}

// ProductKey 获得productKey
func (sf *DevNode) ProductKey() string { return sf.productKey }

// DeviceName 获得DeviceName
func (sf *DevNode) DeviceName() string { return sf.deviceName }

// DeviceSecret 获得DeviceSecret
func (sf *DevNode) DeviceSecret() string { return sf.deviceSecret }

// Avail 返回设备avail
func (sf *DevNode) Avail() bool { return sf.avail }

// Status 返回设备状态
func (sf *DevNode) Status() DevStatus { return sf.status }

// Extend 获得扩展参数值
func (sf *DevNode) Extend() interface{} { return sf.ext }

// NewDevMgr 设备管理是一个线程安全
func NewDevMgr(root infra.MetaTriad) *DevMgr {
	return &DevMgr{
		root: DevNode{
			root.ProductKey,
			root.DeviceName,
			root.DeviceSecret,
			true,
			DevStatusOnline,
			nil,
		},
		nodes: make(map[string]*DevNode),
	}
}

// Len 设备个数,含root设备
func (sf *DevMgr) Len() int {
	sf.rw.RLock()
	defer sf.rw.Unlock()
	return len(sf.nodes) + 1
}

// Create 创建一个子设备
func (sf *DevMgr) Create(meta infra.MetaTetrad) error {
	if meta.ProductKey == "" || meta.DeviceName == "" || meta.DeviceSecret == "" {
		return ErrInvalidParameter
	}

	sf.rw.Lock()
	defer sf.rw.Unlock()

	if meta.ProductKey == sf.root.productKey && meta.DeviceName == sf.root.deviceName {
		return ErrNotPermit
	}
	_, ok := sf.nodes[FormatKey(meta.ProductKey, meta.DeviceName)]
	if ok {
		return ErrDeviceHasExist
	}
	sf.nodes[FormatKey(meta.ProductKey, meta.DeviceName)] = &DevNode{
		meta.ProductKey,
		meta.DeviceName,
		meta.DeviceSecret,
		true,
		DevStatusUnauthorized,
		nil,
	}
	return nil
}

// Delete 删除一个子设备
func (sf *DevMgr) Delete(productKey, deviceName string) {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	delete(sf.nodes, FormatKey(productKey, deviceName))
}

func (sf *DevMgr) searchDevNodeLocked(productKey, deviceName string) (*DevNode, error) {
	if sf.root.productKey == productKey && sf.root.deviceName == deviceName {
		return &sf.root, nil
	}
	node, ok := sf.nodes[FormatKey(productKey, deviceName)]
	if !ok {
		return nil, ErrNotFound
	}
	return node, nil
}

// SearchDevNode 使用productKey deviceName查找一个设备节点信息
func (sf *DevMgr) SearchDevNode(productKey, deviceName string) (*DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	return sf.searchDevNodeLocked(productKey, deviceName)
}

// SetDeviceSecret 设置设备的密钥
func (sf *DevMgr) SetDeviceSecret(productKey, deviceName, deviceSecret string) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	node, err := sf.searchDevNodeLocked(productKey, deviceName)
	if err != nil {
		return err
	}
	node.deviceSecret = deviceSecret
	return nil
}

func (sf *DevMgr) DeviceSecret(productKey, deviceName string) (string, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, err := sf.searchDevNodeLocked(productKey, deviceName)
	if err != nil {
		return "", err
	}
	return node.deviceSecret, nil
}

// SetDevAvail 设置avail
func (sf *DevMgr) SetDevAvail(productKey, deviceName string, enable bool) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchDevNodeLocked(productKey, deviceName)
	if err != nil {
		return err
	}
	node.avail = enable
	return nil
}

// DevAvail 获取avail
func (sf *DevMgr) DevAvail(productKey, deviceName string) (bool, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, err := sf.searchDevNodeLocked(productKey, deviceName)
	if err != nil {
		return false, err
	}
	return node.avail, nil
}

// SetDevStatus 设置设备的状态
func (sf *DevMgr) SetDevStatus(productKey, deviceName string, status DevStatus) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchDevNodeLocked(productKey, deviceName)
	if err != nil {
		return err
	}
	node.status = status
	return nil
}

// SetDevStatus 获取设备的状态
func (sf *DevMgr) DevStatus(productKey, deviceName string, status DevStatus) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchDevNodeLocked(productKey, deviceName)
	if err != nil {
		return err
	}
	node.status = status
	return nil
}

func FormatKey(productKey, deviceName string) string {
	return productKey + "." + deviceName
}
