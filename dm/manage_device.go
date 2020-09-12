package dm

import (
	"errors"
	"sync"

	"github.com/thinkgos/aliyun-iot/infra"
)

// DevNodeLocal 设备本身, 对于网关,独立设备,就是指代本身
const DevNodeLocal = 0

// DevType 设备类型
type DevType byte

// 设备类型定义
const (
	DevTypeSingle  = 1 << iota // 独立设备
	DevTypeSubDev              // 子设备
	DevTypeGateway             // 网关设备

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

// DevMgr 设备管理
type DevMgr struct {
	nextID int
	rw     sync.RWMutex
	nodes  map[int]*DevNode
}

// DevNode 设备节点
type DevNode struct {
	id           int
	types        DevType
	productKey   string
	deviceName   string
	deviceSecret string
	avail        bool
	status       DevStatus
}

// NewDevMgr 设备管理是一个线程安全
func NewDevMgr() *DevMgr {
	return &DevMgr{
		nextID: 1,
		nodes:  make(map[int]*DevNode),
	}
}

// 下一个设备id
func (sf *DevMgr) nextDevID() int {
	for {
		if sf.nextID <= 0 {
			sf.nextID = 1
		}
		id := sf.nextID
		sf.nextID++
		if _, exist := sf.nodes[id]; !exist {
			return id
		}
	}
}

// Len 设备个数
func (sf *DevMgr) Len() int {
	sf.rw.RLock()
	defer sf.rw.Unlock()
	return len(sf.nodes)
}

// Create 创建一个设备,并返回设备ID
func (sf *DevMgr) Create(types DevType, meta infra.MetaInfo) (int, error) {
	if meta.ProductKey == "" || meta.DeviceName == "" || meta.DeviceSecret == "" {
		return 0, ErrInvalidParameter
	}

	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchNodeByPkDnLocked(meta.ProductKey, meta.DeviceName)
	if err != nil {
		id := sf.nextDevID()
		sf.nodes[id] = &DevNode{
			id:           id,
			types:        types,
			productKey:   meta.ProductKey,
			deviceName:   meta.DeviceName,
			deviceSecret: meta.DeviceSecret,
			avail:        true,
		}
		return id, nil
	}
	return node.id, nil
}

// Insert
func (sf *DevMgr) insert(devID int, types DevType, meta infra.MetaInfo) error {
	if meta.ProductKey == "" || meta.DeviceName == "" ||
		meta.DeviceSecret == "" || devID < 0 {
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
		productKey:   meta.ProductKey,
		deviceName:   meta.DeviceName,
		deviceSecret: meta.DeviceSecret,
		avail:        true,
	}
	return nil
}

// Delete 删除一个设备, DevNodeLocal 不可删除
func (sf *DevMgr) Delete(devID int) {
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

	node, err := sf.searchNodeByPkDnLocked(productKey, deviceName)
	if err == nil && node.id != DevNodeLocal {
		delete(sf.nodes, node.id)
	}
}

// searchNodeByPkDnLocked 使用productKey deviceName查找一个节点
// 需要带锁操作
func (sf *DevMgr) searchNodeByPkDnLocked(productKey, deviceName string) (*DevNode, error) {
	for _, node := range sf.nodes {
		if node.productKey == productKey && node.deviceName == deviceName {
			return node, nil
		}
	}
	return nil, ErrNotFound
}

// SearchNode 使用devID查找一个设备节点信息
func (sf *DevMgr) SearchNode(devID int) (*DevNode, error) {
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	if node, exist := sf.nodes[devID]; exist {
		return node, nil
	}
	return nil, ErrNotFound
}

// SearchNodeByPkDn 使用productKey deviceName查找一个设备节点信息
func (sf *DevMgr) SearchNodeByPkDn(productKey, deviceName string) (*DevNode, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	return sf.searchNodeByPkDnLocked(productKey, deviceName)
}

// SetDevAvail 设置avail
func (sf *DevMgr) SetDevAvail(devID int, enable bool) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	if node, exist := sf.nodes[devID]; exist {
		node.avail = enable
		return nil
	}
	return ErrNotFound
}

// SetDevAvailByPkDN 设置avail
func (sf *DevMgr) SetDevAvailByPkDN(productKey, deviceName string, enable bool) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchNodeByPkDnLocked(productKey, deviceName)
	if err != nil {
		return err
	}
	node.avail = enable
	return nil
}

// DevAvail 获取devAvail
func (sf *DevMgr) DevAvail(devID int) (bool, error) {
	if devID < 0 {
		return false, ErrInvalidParameter
	}
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	if node, exist := sf.nodes[devID]; exist {
		return node.avail, nil
	}
	return false, ErrNotFound
}

// DevAvailByPkDn 获取avail
func (sf *DevMgr) DevAvailByPkDn(productKey, deviceName string) (bool, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()

	node, err := sf.searchNodeByPkDnLocked(productKey, deviceName)
	if err != nil {
		return false, err
	}
	return node.avail, nil
}

// SetDevStatus 设置设备的status
func (sf *DevMgr) SetDevStatus(devID int, status DevStatus) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()

	if node, exist := sf.nodes[devID]; exist {
		node.status = status
		return nil
	}
	return ErrNotFound
}

// SetDevStatusByPkDn 设置设备的状态
func (sf *DevMgr) SetDevStatusByPkDn(productKey, deviceName string, status DevStatus) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchNodeByPkDnLocked(productKey, deviceName)
	if err != nil {
		return err
	}
	node.status = status
	return nil
}

// SetDeviceSecret 设置设备密钥
func (sf *DevMgr) SetDeviceSecret(devID int, deviceSecret string) error {
	if devID < 0 {
		return ErrInvalidParameter
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()

	if node, exist := sf.nodes[devID]; exist {
		node.deviceSecret = deviceSecret
		return nil
	}
	return ErrNotFound
}

// SetDeviceSecretByPkDn 设置设备的密钥
func (sf *DevMgr) SetDeviceSecretByPkDn(productKey, deviceName, deviceSecret string) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()

	node, err := sf.searchNodeByPkDnLocked(productKey, deviceName)
	if err != nil {
		return err
	}
	node.deviceSecret = deviceSecret
	return nil
}

// ID 返回设备ID
func (sf *DevNode) ID() int { return sf.id }

// Types 返回设备类型
func (sf *DevNode) Types() DevType { return sf.types }

// Status 返回设备状态
func (sf *DevNode) Status() DevStatus { return sf.status }

// Avail 返回设备avail
func (sf *DevNode) Avail() bool { return sf.avail }

// ProductKey 获得productKey
func (sf *DevNode) ProductKey() string { return sf.productKey }

// DeviceName 获得DeviceName
func (sf *DevNode) DeviceName() string { return sf.deviceName }

// DeviceSecret 获得DeviceSecret
func (sf *DevNode) DeviceSecret() string { return sf.deviceSecret }
