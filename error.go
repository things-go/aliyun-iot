// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

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
	ErrNotActive         = errors.New("device not active")
	ErrNotAvail          = errors.New("device not avail")
)
