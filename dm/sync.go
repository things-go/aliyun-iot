package dm

import (
	"strconv"
	"time"
)

// message message
type message struct {
	err  error
	id   uint
	data interface{}
}

// Entry entry
type Entry struct {
	message chan message
}

// wait the entry response,return id,data and error
func (sf *Entry) Wait(t time.Duration) (uint, interface{}, error) {
	tm := time.NewTimer(t)
	defer tm.Stop()
	select {
	case v := <-sf.message:
		return v.id, v.data, v.err
	case <-tm.C:
	}
	return 0, nil, ErrWaitMessageTimeout
}

// Insert 缓存插入指定ID
func (sf *Client) Insert(id uint) *Entry {
	entry := &Entry{make(chan message, 1)}
	sf.msgCache.SetDefault(strconv.FormatUint(uint64(id), 10), entry)
	sf.debugf("cache Insert - @%d", id)
	return entry
}

// done 指定缓存id收到回复,并发出同步通知
func (sf *Client) done(id uint, err error, data interface{}) {
	key := strconv.FormatUint(uint64(id), 10)
	if v, ok := sf.msgCache.Get(key); ok {
		sf.msgCache.Delete(key)
		sf.debugf("cache done - @%d", id)
		select {
		case v.(*Entry).message <- message{err, id, data}:
		default:
		}
	}
}
