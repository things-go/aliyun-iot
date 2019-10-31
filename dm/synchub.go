package dm

import (
	"strconv"
	"time"

	"github.com/thinkgos/cache-go"
)

// 同步控制默认职
const (
	DefaultWaitTimeout     = 10 * time.Second
	DefaultExpiration      = DefaultWaitTimeout + 2*time.Second
	DefaultCleanUpInterval = 30 * time.Second
)

// SyncHub 同步控制
type SyncHub struct {
	c           *cache.Cache
	waitTimeout time.Duration
}

// NewSyncHub 新建同步控制
func NewSyncHub() *SyncHub {
	return &SyncHub{
		cache.New(DefaultExpiration, DefaultCleanUpInterval),
		DefaultWaitTimeout,
	}
}

// Done 发送同步通知
func (sf *SyncHub) Done(id int, err error) {
	v, ok := sf.c.Get(strconv.Itoa(id))
	if !ok {
		return
	}
	select {
	case v.(chan error) <- err:
	default:
	}
}

// Wait 等待同步
func (sf *SyncHub) Wait(id int, t ...time.Duration) error {
	tm := sf.waitTimeout
	if len(t) > 0 {
		tm = t[0]
	}
	node := make(chan error, 1)
	sf.c.SetDefault(strconv.Itoa(id), node)

	tk := time.NewTicker(tm)
	select {
	case v := <-node:
		return v
	case <-tk.C:
		tk.Stop()
	}
	return ErrWaitMessageTimeout
}
