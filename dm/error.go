package dm

import (
	"errors"
)

// 错误相关定义
var (
	ErrInvalidURI         = errors.New("invalid URI")
	ErrNotFound           = errors.New("not found")
	ErrInvalidParameter   = errors.New("invalid parameter")
	ErrNotSupportMsgType  = errors.New("not support message type")
	ErrNotSupportFeature  = errors.New("not support feature")
	ErrWaitMessageTimeout = errors.New("wait message timeout")
	ErrIPCMessageBuffFull = errors.New("ipc message buff is full")
)
