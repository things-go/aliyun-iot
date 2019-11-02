package dm

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/thinkgos/cache-go"
)

// MsgCacheEntry 消息缓存条目
type MsgCacheEntry struct {
	msgType MsgType // 消息类型
	id      int
	devID   int // 设备id
	err     chan error
	done    uint32
}

// cacheInit 缓存初始化
func (sf *Client) cacheInit() {
	sf.msgCache = cache.New(sf.cfg.cacheExpiration, sf.cfg.cacheCleanupInterval)
	sf.pool = newPool()
	sf.msgCache.OnEvicted(func(s string, v interface{}) { // 超时处理
		entry := v.(*MsgCacheEntry)
		if atomic.LoadUint32(&entry.done) == 0 {
			if err := sf.ipcSendMessage(&ipcMessage{
				evt:     ipcEvtRRPCRequest,
				extend:  strconv.Itoa(entry.devID),
				payload: entry.msgType,
			}); err != nil {
				sf.warn("ipc send message cache timeout failed, %+v", err)
			}
		}

		sf.pool.Put(entry)
	})
	sf.msgCache.OnDeleted(func(s string, v interface{}) {
		sf.pool.Put(v.(*MsgCacheEntry))
	})
}

// CacheInsert 缓存插入指定ID
func (sf *Client) CacheInsert(id, devID int, msgType MsgType) {
	entry := sf.pool.Get()
	entry.id = id
	entry.devID = devID
	entry.msgType = msgType
	entry.done = 0
	sf.msgCache.SetDefault(strconv.Itoa(id), entry)
	sf.debug("cache insert - @%d", id)
}

// CacheGet 获取缓存消息设备ID
func (sf *Client) CacheGet(id int) (int, bool) {
	v, ok := sf.msgCache.Get(strconv.Itoa(id))
	if ok {
		return v.(*MsgCacheEntry).devID, true
	}
	return 0, false
}

func (sf *Client) CacheWait(id int, t ...time.Duration) error {
	v, ok := sf.msgCache.Get(strconv.Itoa(id))
	if !ok {
		return ErrNotFound
	}

	entry := v.(*MsgCacheEntry)

	tm := 10 * time.Second
	if len(t) > 0 {
		tm = t[0]
	}

	tk := time.NewTicker(tm)
	defer tk.Stop()
	select {
	case v := <-entry.err:
		return v
	case <-tk.C:
	}
	return ErrWaitMessageTimeout
}

// CacheDone 发送同步通知
func (sf *Client) CacheDone(id int, err error) {
	v, ok := sf.msgCache.Get(strconv.Itoa(id))
	if !ok {
		return
	}
	entry := v.(*MsgCacheEntry)
	select {
	case entry.err <- err:
		atomic.StoreUint32(&entry.done, 1)
	default:
	}
}
