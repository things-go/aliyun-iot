package dm

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/thinkgos/cache-go"
)

// MsgCacheEntry 消息缓存条目
type MsgCacheEntry struct {
	msgType MsgType    // 消息类型
	devID   int        // 设备id
	err     chan error // 用于wait通道
	done    uint32     // 表示消息接收到应答
}

// cacheInit 缓存初始化
func (sf *Client) cacheInit() {
	sf.msgCache = cache.New(sf.cfg.cacheExpiration, sf.cfg.cacheCleanupInterval)
	sf.pool = newPool()
	sf.msgCache.OnEvicted(func(id string, v interface{}) { // 超时处理
		entry := v.(*MsgCacheEntry)
		if atomic.LoadUint32(&entry.done) == 0 {
			if err := sf.ipcSendMessage(&ipcMessage{
				evt:     ipcEvtRequestTimeout,
				extend:  strconv.Itoa(entry.devID),
				payload: entry.msgType,
			}); err != nil {
				sf.warn("ipc send message cache timeout failed, %+v", err)
			}
		}
		sf.debug("cache evicted - @%s", id)
		sf.pool.Put(entry)
	})
	sf.msgCache.OnDeleted(func(s string, v interface{}) {
		sf.pool.Put(v.(*MsgCacheEntry))
	})
}

// CacheInsert 缓存插入指定ID
func (sf *Client) CacheInsert(id, devID int, msgType MsgType) {
	entry := sf.pool.Get()
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

// CacheWait 等待缓存ID的消息收到回复
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
	sf.debug("cache wait - @%d", id)
	select {
	case v := <-entry.err:
		return v
	case <-tk.C:
	}
	return ErrWaitMessageTimeout
}

// CacheDone 指定缓存id收到回复,并发出同步通知
func (sf *Client) CacheDone(id int, err error) {
	v, ok := sf.msgCache.Get(strconv.Itoa(id))
	if !ok {
		return
	}

	sf.debug("cache done - @%d", id)
	entry := v.(*MsgCacheEntry)
	atomic.StoreUint32(&entry.done, 1)
	select {
	case entry.err <- err:
	default:
	}
}
