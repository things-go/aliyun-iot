package dm

import (
	"sync"
)

type pool struct {
	pl sync.Pool
}

func newPool() *pool {
	return &pool{
		sync.Pool{
			New: func() interface{} { return new(MsgCacheEntry) },
		},
	}
}

func (sf *pool) Get() *MsgCacheEntry {
	return sf.pl.Get().(*MsgCacheEntry)
}

func (sf *pool) Put(entry interface{}) {
	sf.pl.Put(entry)
}
