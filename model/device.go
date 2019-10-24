package model

import (
	"container/list"
	"errors"
	"sync"
)

// 设备本身, 对于网关,独立设备,就是指代本身
const DevItself = 0

// DevStatus 设备状态
type DevStatus byte

// 设备状态
const (
	DevStatusUnauthorized = iota // Subdev Created
	DevStatusAuthorized          // Receive Topo Add Notify
	DevStatusRegistered          // Receive Subdev Registered
	DevStatusAttached            // Receive Subdev Topo Add Reply
	DevStatusLogined             // Receive Subdev Login Reply
	DevStatusOnline              // After All Topic Subscribed
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
	listHead    *list.List
}

// 设备节点
type DevNode struct {
	id           int
	types        string
	ProductKey   string
	DeviceName   string
	DeviceSecret string
	avail        DevAvail
	status       DevStatus
}

func newDevMgr() *devMgr {
	return &devMgr{
		globalDevID: 0,
		listHead:    list.New(),
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
	return sf.listHead.Len()
}

// Create 创建设备,返回设备ID
func (sf *devMgr) Create(types, productKey, deviceName, deviceSecret string) (int, error) {
	if productKey == "" ||
		deviceName == "" ||
		deviceSecret == "" {
		return 0, errors.New("invalid parameter")
	}
	node, err := sf.searchByPkDn(productKey, deviceName)
	if err == nil {
		return node.id, nil
	}

	id := sf.nestDevID()
	sf.listHead.PushBack(&DevNode{
		id:           id,
		types:        types,
		ProductKey:   productKey,
		DeviceName:   deviceName,
		DeviceSecret: deviceSecret,
	})
	return id, nil
}

func (sf *devMgr) DeleteByID(devID int) error {
	if devID < 0 {
		return ErrInvalidParameter
	}

	sf.rw.Lock()
	defer sf.rw.Unlock()
	var next *list.Element
	for e := sf.listHead.Front(); e != nil; e = next {
		next = e.Next()
		node := e.Value.(*DevNode)
		if node.id == devID {
			sf.listHead.Remove(e)
			break
		}
	}
	return nil
}

func (sf *devMgr) DeleteByPkDn(productKey, deviceName string) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	var next *list.Element
	for e := sf.listHead.Front(); e != nil; e = next {
		next = e.Next()
		node := e.Value.(*DevNode)
		if node.ProductKey == productKey &&
			node.DeviceName == deviceName {
			sf.listHead.Remove(e)
		}
	}
	return nil
}

func (sf *devMgr) searchByID(devID int) (*DevNode, error) {
	for e := sf.listHead.Front(); e != nil; e = e.Next() {
		node := e.Value.(*DevNode)
		if node.id == devID {
			return node, nil
		}
	}
	return nil, ErrNotFound
}

func (sf *devMgr) SearchByID(devID int) (int, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, err := sf.searchByID(devID)
	if err != nil {
		return 0, err
	}
	return node.id, nil
}

func (sf *devMgr) SearchNodeByID(devID int) (DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, err := sf.searchByID(devID)
	if err != nil {
		return DevNode{}, err
	}
	return *node, err
}

func (sf *devMgr) searchByPkDn(productKey, deviceName string) (*DevNode, error) {
	for e := sf.listHead.Front(); e != nil; e = e.Next() {
		node := e.Value.(*DevNode)
		if node.ProductKey == productKey &&
			node.DeviceName == deviceName {
			return node, nil
		}
	}

	return nil, ErrNotFound
}

func (sf *devMgr) SearchByPkDn(productKey, deviceName string) (int, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, err := sf.searchByPkDn(productKey, deviceName)
	if err != nil {
		return 0, err
	}
	return node.id, nil
}

func (sf *devMgr) SearchNodeByPkDn(productKey, deviceName string) (DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, err := sf.searchByPkDn(productKey, deviceName)
	if err != nil {
		return DevNode{}, err
	}
	return *node, nil
}

func (sf *devMgr) SetDevAvail(devID int, enable bool) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchByID(devID)
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

func (sf *devMgr) DevAvail(productKey, deviceName string) (DevAvail, error) {
	sf.rw.RLock()
	defer sf.rw.Unlock()

	node, err := sf.searchByPkDn(productKey, deviceName)
	if err != nil {
		return DevAvailEnable, err
	}
	return node.avail, nil
}

func (sf *devMgr) SetDevStatus(devID int, status DevStatus) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchByID(devID)
	if err != nil {
		return ErrNotFound
	}
	node.status = status
	return nil
}

func (sf *devMgr) DevStatus(devID int) (DevStatus, error) {
	sf.rw.RLock()
	defer sf.rw.Unlock()
	node, err := sf.searchByID(devID)
	if err != nil {
		return DevStatusUnauthorized, err
	}
	return node.status, nil
}

func (sf *devMgr) SetDeviceSecret(devID int, deviceSecret string) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	node, err := sf.searchByID(devID)
	if err != nil {
		return err
	}
	node.DeviceSecret = deviceSecret
	return nil
}

func (sf *devMgr) DevTypes(devID int) (string, error) {
	if devID < 0 {
		return "", ErrInvalidParameter
	}

	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, err := sf.searchByID(devID)
	if err != nil {
		return "", err
	}
	return node.types, nil
}
