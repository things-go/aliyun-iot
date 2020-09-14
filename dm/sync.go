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

// Token defines the interface for the tokens used to indicate when actions have completed.
type Token struct {
	message chan message
}

// closedchan is a reusable closed channel.
var closedchan = make(chan message)

func init() {
	close(closedchan)
}

// WaitErr 等待同步回复,成功回复,只返回error
func (sf *Token) WaitErr(timeout time.Duration) error {
	_, _, err := sf.Wait(timeout)
	return err
}

// WaitData 等待同步回复,成功回复,返回data和error
func (sf *Token) WaitData(timeout time.Duration) (interface{}, error) {
	_, data, err := sf.Wait(timeout)
	return data, err
}

// WaitId 等待同步回复,成功回复,返回id和error
func (sf *Token) WaitID(timeout time.Duration) (uint, error) {
	id, _, err := sf.Wait(timeout)
	return id, err
}

// Wait the entry response,return id,data and error
func (sf *Token) Wait(timeout time.Duration) (uint, interface{}, error) {
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
	return 0, nil, ErrWaitTimeout
}

// Insert 缓存插入指定ID3
func (sf *Client) Insert(id uint) *Token {
	if sf.workOnWho != WorkOnMQTT {
		return &Token{closedchan}
	}
	entry := &Token{make(chan message, 1)}
	sf.msgCache.SetDefault(strconv.FormatUint(uint64(id), 10), entry)
	return entry
}

// signal 指定缓存id收到回复,并发出同步通知
func (sf *Client) signal(id uint, err error, data interface{}) {
	key := strconv.FormatUint(uint64(id), 10)
	if v, ok := sf.msgCache.Get(key); ok {
		sf.msgCache.Delete(key)
		select {
		case v.(*Token).message <- message{err, id, data}:
		default:
		}
	}
}
