// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

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
// root: 网关设备
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
	defer sf.rw.RUnlock()
	return len(sf.nodes) + 1
}

// Add 增加一个子设备,子设备处于 DevStatusUnauthorized
func (sf *DevMgr) Add(meta infra.MetaTriad) error {
	if meta.ProductKey == "" || meta.DeviceName == "" {
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
func (sf *DevMgr) Delete(pk, dn string) {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	delete(sf.nodes, FormatKey(pk, dn))
}

func (sf *DevMgr) searchLocked(pk, dn string) (*DevNode, error) {
	if sf.root.productKey == pk && sf.root.deviceName == dn {
		return &sf.root, nil
	}
	node, ok := sf.nodes[FormatKey(pk, dn)]
	if !ok {
		return nil, ErrNotFound
	}
	return node, nil
}

// SearchAvail 使用productKey deviceName查找一个设备节点信息且avail = true
// 如果设备avail=false返回ErrNotAvail
func (sf *DevMgr) SearchAvail(pk, dn string) (*DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, err := sf.searchLocked(pk, dn)
	if err != nil {
		return nil, err
	}
	if node.avail {
		return node, nil
	}
	return nil, ErrNotAvail
}

// IsActive productKey deviceName的设备已使能且是否在线
func (sf *DevMgr) IsActive(pk, dn string) bool {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	node, err := sf.searchLocked(pk, dn)
	if err != nil {
		return false
	}
	return node.avail && node.status == DevStatusOnline
}

// Search 使用productKey deviceName查找一个设备节点信息
func (sf *DevMgr) Search(pk, dn string) (*DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	return sf.searchLocked(pk, dn)
}

// SetDeviceSecret 设置设备的密钥
func (sf *DevMgr) SetDeviceSecret(pk, dn, ds string) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	node, err := sf.searchLocked(pk, dn)
	if err != nil {
		return err
	}
	node.deviceSecret = ds
	return nil
}

// DeviceSecret 设备DeviceSecret
func (sf *DevMgr) DeviceSecret(pk, dn string) (string, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, err := sf.searchLocked(pk, dn)
	if err != nil {
		return "", err
	}
	return node.deviceSecret, nil
}

// SetDeviceAvail 设置avail
func (sf *DevMgr) SetDeviceAvail(pk, dn string, enable bool) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchLocked(pk, dn)
	if err != nil {
		return err
	}
	node.avail = enable
	return nil
}

// DeviceAvail 获取avail
func (sf *DevMgr) DeviceAvail(pk, dn string) (bool, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, err := sf.searchLocked(pk, dn)
	if err != nil {
		return false, err
	}
	return node.avail, nil
}

// SetDeviceStatus 设置设备的状态
func (sf *DevMgr) SetDeviceStatus(pk, dn string, status DevStatus) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchLocked(pk, dn)
	if err != nil {
		return err
	}
	node.status = status
	return nil
}

// DeviceStatus 获取设备的状态
func (sf *DevMgr) DeviceStatus(pk, dn string, status DevStatus) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchLocked(pk, dn)
	if err != nil {
		return err
	}
	node.status = status
	return nil
}

// FormatKey format pk dn --> {pk}.{dn}
func FormatKey(pk, dn string) string {
	return pk + "." + dn
}
