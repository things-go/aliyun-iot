package ahttp

import (
	"errors"
)

var (
	ErrUnknown             = errors.New("code<10000>: unknown error")
	ErrParamException      = errors.New("code<10001>: param exception")
	ErrAuthFailed          = errors.New("code<20000>: auth failed")
	ErrUpdateSessionFailed = errors.New("code<20000>: update session failed")
	ErrRequestTooMany      = errors.New("code<40000>: request too many")
)
