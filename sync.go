package aiot

import (
	"strconv"
	"time"
)

// Message 回复的消息
type Message struct {
	ID   uint
	Data interface{}

	err error
}

// Token defines the interface for the tokens used to indicate when actions have completed.
type Token struct {
	message chan Message
}

// closedchan is a reusable closed channel.
var closedchan = make(chan Message)

func init() {
	close(closedchan)
}

// Wait the entry response,return ID,Data and error
func (sf *Token) Wait(timeout time.Duration) (m Message, err error) {
	tm := time.NewTimer(timeout)
	defer tm.Stop()
	select {
	case m, ok := <-sf.message:
		if ok {
			return m, m.err
		}
		return m, ErrEntryClosed
	case <-tm.C:
	}
	return m, ErrWaitTimeout
}

// putPending 缓存插入指定ID3
func (sf *Client) putPending(id uint) *Token {
	if sf.workOnWho != WorkOnMQTT {
		return &Token{closedchan}
	}
	entry := &Token{make(chan Message, 1)}
	sf.msgCache.SetDefault(strconv.FormatUint(uint64(id), 10), entry)
	return entry
}

// signalPending 指定缓存id收到回复,并发出同步通知
func (sf *Client) signalPending(msg Message) {
	key := strconv.FormatUint(uint64(msg.ID), 10)
	if v, ok := sf.msgCache.Get(key); ok {
		sf.msgCache.Delete(key)
		select {
		case v.(*Token).message <- msg:
		default:
		}
	}
}
