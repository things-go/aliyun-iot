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

// closedchan is a reusable closed channel.
var closedchan = make(chan message)

func init() {
	close(closedchan)
}

// Wait the entry response,return id,data and error
func (sf *Entry) Wait(timeout time.Duration) (uint, interface{}, error) {
	tm := time.NewTimer(timeout)
	defer tm.Stop()
	select {
	case v, ok := <-sf.message:
		if ok {
			return v.id, v.data, v.err
		}
		return 0, nil, ErrEntryClosed
	case <-tm.C:
	}
	return 0, nil, ErrWaitMessageTimeout
}

// Insert 缓存插入指定ID3
func (sf *Client) Insert(id uint) *Entry {
	if sf.workOnWho != WorkOnMQTT {
		return &Entry{closedchan}
	}
	entry := &Entry{make(chan message, 1)}
	sf.msgCache.SetDefault(strconv.FormatUint(uint64(id), 10), entry)
	sf.log.Debugf("cache Insert - @%d", id)
	return entry
}

// signal 指定缓存id收到回复,并发出同步通知
func (sf *Client) signal(id uint, err error, data interface{}) {
	key := strconv.FormatUint(uint64(id), 10)
	if v, ok := sf.msgCache.Get(key); ok {
		sf.msgCache.Delete(key)
		sf.log.Debugf("cache signal - @%d", id)
		select {
		case v.(*Entry).message <- message{err, id, data}:
		default:
		}
	}
}
