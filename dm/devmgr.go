package dm

import (
	"errors"
	"sync"
)

// DevNodeLocal 设备本身, 对于网关,独立设备,就是指代本身
const DevNodeLocal = 0

// DevType 设备类型
type DevType byte

// 设备类型定义
const (
	DevTypeSingle = 1 << iota
	DevTypeSubDev
	DevTypeGateway

	// DevTypeMain = DevTypeSingle | DevTypeSubDev
	// DevTypeALl
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

// DevAvail 设备有效
type DevAvail byte

// 设备有效
const (
	DevAvailEnable DevAvail = iota
	DevAvailDisable
)

// DevMgr 设备管理
type DevMgr struct {
	globalDevID int
	rw          sync.RWMutex
	nodes       map[int]*DevNode
}

// DevNode 设备节点
type DevNode struct {
	id           int
	types        DevType
	productKey   string
	deviceName   string
	deviceSecret string
	avail        DevAvail
	status       DevStatus
}

// NewDevMgr 设备管理是一个线程安全
func NewDevMgr() *DevMgr {
	return &DevMgr{
		globalDevID: 1,
		nodes:       make(map[int]*DevNode),
	}
}

// 下一个设备id
func (sf *DevMgr) nextDevID() int {
	id := sf.globalDevID
	sf.globalDevID++
	return id
}

// Len 设备个数
func (sf *DevMgr) Len() int {
	sf.rw.RLock()
	defer sf.rw.Unlock()
	return len(sf.nodes)
}

// Create 创建一个设备,并返回设备ID
func (sf *DevMgr) Create(types DevType, productKey, deviceName, deviceSecret string) (int, error) {
	if productKey == "" ||
		deviceName == "" ||
		deviceSecret == "" {
		return 0, ErrInvalidParameter
	}

	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchNodeByPkDn(productKey, deviceName)
	if err != nil {
		id := sf.nextDevID()
		sf.nodes[id] = &DevNode{
			id:           id,
			types:        types,
			productKey:   productKey,
			deviceName:   deviceName,
			deviceSecret: deviceSecret,
		}
		return id, nil
	}
	return node.id, nil
}

// insert
func (sf *DevMgr) insert(devID int, types DevType, productKey, deviceName, deviceSecret string) error {
	if productKey == "" ||
		deviceName == "" ||
		deviceSecret == "" ||
		devID < 0 {
		return ErrInvalidParameter
	}

	sf.rw.Lock()
	defer sf.rw.Unlock()

	if _, exist := sf.nodes[devID]; exist {
		return errors.New("device node has exist")
	}
	sf.nodes[devID] = &DevNode{
		id:           devID,
		types:        types,
		productKey:   productKey,
		deviceName:   deviceName,
		deviceSecret: deviceSecret,
	}
	return nil
}

// DeleteByID 删除一个设备, DevSelf不可删除
func (sf *DevMgr) DeleteByID(devID int) {
	if devID < 0 || devID == DevNodeLocal {
		return
	}
	sf.rw.Lock()
	delete(sf.nodes, devID)
	sf.rw.Unlock()
}

// DeleteByPkDn 删除一个子设备, DevNodeLocal不可删除
func (sf *DevMgr) DeleteByPkDn(productKey, deviceName string) {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	for id, node := range sf.nodes {
		if node.productKey == productKey &&
			node.deviceName == deviceName {
			if id != DevNodeLocal {
				delete(sf.nodes, id)
			}
			return
		}
	}
}

// searchNodeByPkDn 使用productKey deviceName查找一个节点
// 需要带锁操作
func (sf *DevMgr) searchNodeByPkDn(productKey, deviceName string) (*DevNode, error) {
	for id, node := range sf.nodes {
		if node.productKey == productKey &&
			node.deviceName == deviceName {
			delete(sf.nodes, id)
			return node, nil
		}
	}
	return nil, ErrNotFound
}

// SearchNodeByID 使用devID查找一个设备节点信息
func (sf *DevMgr) SearchNodeByID(devID int) (DevNode, error) {
	if devID < 0 {
		return DevNode{}, ErrInvalidParameter
	}
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, exist := sf.nodes[devID]
	if !exist {
		return DevNode{}, ErrNotFound
	}
	return *node, nil
}

// SearchNodeByPkDn 使用productKey deviceName查找一个设备节点信息
func (sf *DevMgr) SearchNodeByPkDn(productKey, deviceName string) (DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, err := sf.searchNodeByPkDn(productKey, deviceName)
	if err != nil {
		return DevNode{}, err
	}
	return *node, nil
}

// SetDevAvailByID 设置avail
func (sf *DevMgr) SetDevAvailByID(devID int, enable bool) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, exist := sf.nodes[devID]
	if !exist {
		return ErrNotFound
	}
	if enable {
		node.avail = DevAvailEnable
	} else {
		node.avail = DevAvailDisable
	}
	return nil
}

// SetDevAvailByPkDN 设置avail
func (sf *DevMgr) SetDevAvailByPkDN(productKey, deviceName string, enable bool) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchNodeByPkDn(productKey, deviceName)
	if err != nil {
		return err
	}
	if enable {
		node.avail = DevAvailEnable
	} else {
		node.avail = DevAvailDisable
	}
	return nil
}

// DevAvailByID 获取devAvail
func (sf *DevMgr) DevAvailByID(devID int) (DevAvail, error) {
	if devID < 0 {
		return 0, ErrInvalidParameter
	}
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	if node, exist := sf.nodes[devID]; exist {
		return node.avail, nil
	}
	return 0, ErrNotFound
}

// DevAvailByPkDn 获取avail
func (sf *DevMgr) DevAvailByPkDn(productKey, deviceName string) (DevAvail, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, err := sf.searchNodeByPkDn(productKey, deviceName)
	if err != nil {
		return DevAvailEnable, err
	}
	return node.avail, nil
}

// SetDevStatusByID 设置设备的status
func (sf *DevMgr) SetDevStatusByID(devID int, status DevStatus) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, exist := sf.nodes[devID]
	if !exist {
		return ErrNotFound
	}
	node.status = status
	return nil
}

// SetDevStatusByPkDn 设置设备的状态
func (sf *DevMgr) SetDevStatusByPkDn(productKey, deviceName string, status DevStatus) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchNodeByPkDn(productKey, deviceName)
	if err != nil {
		return err
	}
	node.status = status
	return nil
}

// SetDeviceSecretByID 设置设备密钥
func (sf *DevMgr) SetDeviceSecretByID(devID int, deviceSecret string) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, exist := sf.nodes[devID]
	if !exist {
		return ErrNotFound
	}
	node.deviceSecret = deviceSecret
	return nil
}

// SetDeviceSecretByPkDn 设置设备的密钥
func (sf *DevMgr) SetDeviceSecretByPkDn(productKey, deviceName, deviceSecret string) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchNodeByPkDn(productKey, deviceName)
	if err != nil {
		return err
	}
	node.deviceSecret = deviceSecret
	return nil
}

// ID 返回设备ID
func (sf *DevNode) ID() int {
	return sf.id
}

// Types 返回设备类型
func (sf *DevNode) Types() DevType {
	return sf.types
}

// Status 返回设备状态
func (sf *DevNode) Status() DevStatus {
	return sf.status
}

// Avail 返回设备avail
func (sf *DevNode) Avail() DevAvail {
	return sf.avail
}

// ProductKey 获得productKey
func (sf *DevNode) ProductKey() string {
	return sf.productKey
}

// DeviceName 获得DeviceName
func (sf *DevNode) DeviceName() string {
	return sf.deviceName
}

// DeviceSecret 获得DeviceSecret
func (sf *DevNode) DeviceSecret() string {
	return sf.deviceSecret
}
