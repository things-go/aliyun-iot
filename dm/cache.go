package dm

import (
	"strconv"

	"github.com/thinkgos/cache-go"
)

// MsgCacheEntry 消息缓存条目
type MsgCacheEntry struct {
	msgType MsgType // 消息类型
	id      int
	devID   int // 设备id
	data    string
}

// DevID 获得devID
func (sf *MsgCacheEntry) DevID() int {
	return sf.devID
}

// cacheInit 缓存初始化
func (sf *Client) cacheInit() {
	if !sf.cfg.hasCache {
		return
	}
	sf.pool = newPool()
	sf.msgCache = cache.New(sf.cfg.cacheExpiration, sf.cfg.cacheCleanupInterval)
	sf.msgCache.OnEvicted(func(s string, v interface{}) { // 超时处理
		sf.pool.Put(v)
		sf.debug("cache timeout - %s", s)
	})
	sf.msgCache.OnDeleted(func(s string, v interface{}) {
		sf.pool.Put(v)
	})
}

// CacheInsert 缓存插入指定ID
func (sf *Client) CacheInsert(id, devID int, msgType MsgType, data string) {
	if !sf.cfg.hasCache {
		return
	}
	entry := sf.pool.Get()
	entry.id = id
	entry.devID = devID
	entry.msgType = msgType
	entry.data = data
	sf.msgCache.SetDefault(strconv.Itoa(id), entry)
	sf.debug("cache insert - @%d", id)
}

// CacheGet 获取缓存消息
func (sf *Client) CacheGet(id int) (MsgCacheEntry, bool) {
	v, ok := sf.msgCache.Get(strconv.Itoa(id))
	if ok {
		return *(v.(*MsgCacheEntry)), true
	}
	return MsgCacheEntry{}, false
}

// CacheRemove 缓存删存指定ID
func (sf *Client) CacheRemove(id int) {
	if !sf.cfg.hasCache {
		return
	}
	sf.msgCache.Delete(strconv.Itoa(id))
	sf.debug("cache remove - @%d", id)
}
