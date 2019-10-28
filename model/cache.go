package model

import (
	"strconv"

	"github.com/thinkgos/cache-go"
)

type messageCacheEntry struct {
	msgType MsgType // 消息类型
	id      int
	devID   int // 设备id
	data    string
}

func (sf *Manager) CacheInit() {
	if !sf.cfg.enableCache {
		return
	}
	sf.pool = newPool()
	sf.msgCache = cache.New(sf.cfg.cacheExpiration, sf.cfg.cacheCleanupInterval)
	sf.msgCache.OnEvicted(func(s string, v interface{}) { // 超时处理
		sf.pool.Put(v)
		sf.debug("cache timeout - %s", s)
	})
}

func (sf *Manager) CacheInsert(id, devID int, msgType MsgType, data string) {
	if !sf.cfg.enableCache {
		return
	}
	entry := sf.pool.Get()
	entry.id = id
	entry.devID = devID
	entry.msgType = msgType
	entry.data = data
	sf.msgCache.SetDefault(strconv.Itoa(id), entry)
	sf.debug("cache insert - %d", id)
}

func (sf *Manager) CacheRemove(id int) {
	if !sf.cfg.enableCache {
		return
	}
	sf.msgCache.Delete(strconv.Itoa(id))
	sf.debug("cache remove - %d", id)
}
