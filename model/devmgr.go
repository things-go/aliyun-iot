package model

import (
	"errors"
	"sync"
)

// 设备本身, 对于网关,独立设备,就是指代本身
const DevSelf = 0

type DevType byte

const (
	DevTypeSingle = 1 << iota
	DevTypeSubdev
	DevTypeGateway

	DevTypeMain = DevTypeSingle | DevTypeSubdev
	DevTypeALl
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

// 设备有效
type DevAvail byte

// 设备有效
const (
	DevAvailEnable DevAvail = iota
	DevAvailDisable
)

// devMgr 设备管理
type devMgr struct {
	globalDevID int
	rw          sync.RWMutex
	nodes       map[int]*DevNode
}

// 设备节点
type DevNode struct {
	id           int
	types        DevType
	ProductKey   string
	DeviceName   string
	DeviceSecret string
	avail        DevAvail
	status       DevStatus
}

func newDevMgr() *devMgr {
	return &devMgr{
		globalDevID: 1,
		nodes:       make(map[int]*DevNode),
	}
}

// 下一个设备id
func (sf *devMgr) nestDevID() int {
	id := sf.globalDevID
	sf.globalDevID++
	return id
}

// Len 设备个数
func (sf *devMgr) Len() int {
	sf.rw.RLock()
	defer sf.rw.Unlock()
	return len(sf.nodes)
}

// Create 创建一个设备,并返回设备ID
func (sf *devMgr) Create(types DevType, productKey, deviceName, deviceSecret string) (int, error) {
	if productKey == "" ||
		deviceName == "" ||
		deviceSecret == "" {
		return 0, errors.New("invalid parameter")
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchByPkDn(productKey, deviceName)
	if err == nil {
		return node.id, nil
	}

	id := sf.nestDevID()
	sf.nodes[id] = &DevNode{
		id:           id,
		types:        types,
		ProductKey:   productKey,
		DeviceName:   deviceName,
		DeviceSecret: deviceSecret,
	}
	return id, nil
}

func (sf *devMgr) insert(devID int, types DevType, productKey, deviceName, deviceSecret string) error {
	if productKey == "" ||
		deviceName == "" ||
		deviceSecret == "" {
		return errors.New("invalid parameter")
	}

	sf.rw.Lock()
	defer sf.rw.Unlock()

	if _, exist := sf.nodes[devID]; exist {
		return errors.New("device node has exist")
	}
	sf.nodes[devID] = &DevNode{
		id:           devID,
		types:        types,
		ProductKey:   productKey,
		DeviceName:   deviceName,
		DeviceSecret: deviceSecret,
	}
	return nil
}

// DeleteByID 删除一个设备
func (sf *devMgr) DeleteByID(devID int) {
	if devID < 0 {
		return
	}

	sf.rw.Lock()
	delete(sf.nodes, devID)
	sf.rw.Unlock()
}

// DeleteByPkDn 删除一个子设备
func (sf *devMgr) DeleteByPkDn(productKey, deviceName string) {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	for id, node := range sf.nodes {
		if node.ProductKey == productKey &&
			node.DeviceName == deviceName {
			delete(sf.nodes, id)
			return
		}
	}
}

// SearchByID 使用devID寻找一个节点
func (sf *devMgr) SearchByID(devID int) (int, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, exist := sf.nodes[devID]
	if !exist {
		return 0, ErrNotFound
	}
	return node.id, nil
}

// SearchNodeByID 使用devID查找一个设备节点信息
func (sf *devMgr) SearchNodeByID(devID int) (DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, exist := sf.nodes[devID]
	if !exist {
		return DevNode{}, ErrNotFound
	}
	return *node, nil
}

// searchByPkDn 使用productKey deviceName查找一个节点,需要带锁
func (sf *devMgr) searchByPkDn(productKey, deviceName string) (*DevNode, error) {
	for id, node := range sf.nodes {
		if node.ProductKey == productKey &&
			node.DeviceName == deviceName {
			delete(sf.nodes, id)
			return node, nil
		}
	}

	return nil, ErrNotFound
}

// SearchByPkDn 使用productKey deviceName查找一个设备,返回设备id
func (sf *devMgr) SearchByPkDn(productKey, deviceName string) (int, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, err := sf.searchByPkDn(productKey, deviceName)
	if err != nil {
		return 0, err
	}
	return node.id, nil
}

// SearchNodeByPkDn 使用productKey deviceName查找一个设备节点信息
func (sf *devMgr) SearchNodeByPkDn(productKey, deviceName string) (DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, err := sf.searchByPkDn(productKey, deviceName)
	if err != nil {
		return DevNode{}, err
	}
	return *node, nil
}

// SetDevAvail 设置avail
func (sf *devMgr) SetDevAvail(devID int, enable bool) error {
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

// DevAvail 获取avail
func (sf *devMgr) DevAvail(productKey, deviceName string) (DevAvail, error) {
	sf.rw.RLock()
	defer sf.rw.Unlock()

	node, err := sf.searchByPkDn(productKey, deviceName)
	if err != nil {
		return DevAvailEnable, err
	}
	return node.avail, nil
}

// SetDevStatus 设置设备的状态
func (sf *devMgr) SetDevStatus(devID int, status DevStatus) error {
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

// DevStatus 获取设备的状态
func (sf *devMgr) DevStatus(devID int) (DevStatus, error) {
	sf.rw.RLock()
	defer sf.rw.Unlock()

	node, exist := sf.nodes[devID]
	if !exist {
		return DevStatusUnauthorized, ErrNotFound
	}
	return node.status, nil
}

// SetDeviceSecret 设置设备密钥
func (sf *devMgr) SetDeviceSecret(devID int, deviceSecret string) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, exist := sf.nodes[devID]
	if !exist {
		return ErrNotFound
	}
	node.DeviceSecret = deviceSecret
	return nil
}

// 获得设备类型
func (sf *devMgr) DevTypes(devID int) (DevType, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, exist := sf.nodes[devID]
	if !exist {
		return 0, ErrNotFound
	}
	return node.types, nil
}

// ID 返回设备ID
func (sf *DevNode) ID() int {
	return sf.id
}

// ID 返回设备状态
func (sf *DevNode) Status() DevStatus {
	return sf.status
}

// ID 返回设备avail
func (sf *DevNode) Avail() DevAvail {
	return sf.avail
}

// ID 返回设备类型
func (sf *DevNode) Types() DevType {
	return sf.types
}
