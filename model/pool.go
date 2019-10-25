package model

import (
	"sync"
)

type pool struct {
	pl sync.Pool
}

func newPool() *pool {
	return &pool{
		sync.Pool{
			New: func() interface{} { return new(messageCacheEntry) },
		},
	}
}

func (sf *pool) Get() *messageCacheEntry {
	return sf.pl.Get().(*messageCacheEntry)
}

func (sf *pool) Put(entry interface{}) {
	sf.pl.Put(entry)
}
