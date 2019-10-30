package dm

import (
	"strconv"
	"time"

	"github.com/thinkgos/cache-go"
)

const (
	DefaultWaitTimeout     = 10 * time.Second
	DefaultExpiration      = DefaultWaitTimeout + 2*time.Second
	DefaultCleanUpInterval = 30 * time.Second
)

type SyncHub struct {
	c           *cache.Cache
	waitTimeout time.Duration
}

func NewSyncHub() *SyncHub {
	return &SyncHub{
		cache.New(DefaultExpiration, DefaultCleanUpInterval),
		DefaultWaitTimeout,
	}
}

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
