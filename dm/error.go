package dm

import (
	"errors"
)

// 错误相关定义
var (
	ErrInvalidURI        = errors.New("invalid URI")
	ErrNotFound          = errors.New("not found")
	ErrInvalidParameter  = errors.New("invalid parameter")
	ErrNotSupportFeature = errors.New("not support feature")
	ErrWaitTimeout       = errors.New("wait timeout")
	ErrEntryClosed       = errors.New("entry has closed")
	ErrDeviceHasExist    = errors.New("device has exist")
	ErrNotPermit         = errors.New("not permit")
)
