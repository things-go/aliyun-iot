package ahttp

import (
	"fmt"
)

// 错误码
const (
	CodeSuccess              = 0
	CodeUnknown              = 10000
	CodeParamException       = 10001
	CodeAuthFailed           = 20000
	CodeTokenExpired         = 20001 // 需重新调用auth进行鉴权，获取token
	CodeTokenIsNull          = 20002 // 需重新调用auth进行鉴权，获取token
	CodeTokenCheckFailed     = 20003 // 根据token获取identify信息失败。需重新调用auth进行鉴权，获取token
	CodePublishMessageFailed = 30001
	CodeUpdateSessionFailed  = 20004
	CodeRequestTooMany       = 40000
)

// 错误定义
var (
	ErrUnknown              = fmt.Errorf("code<%d>: unknown error", CodeUnknown)
	ErrParamException       = fmt.Errorf("code<%d>: param exception", CodeParamException)
	ErrAuthFailed           = fmt.Errorf("code<%d>: auth failed", CodeAuthFailed)
	ErrUpdateSessionFailed  = fmt.Errorf("code<%d>: update session failed", CodeUpdateSessionFailed)
	ErrRequestTooMany       = fmt.Errorf("code<%d>: request too many", CodeRequestTooMany)
	ErrTokenExpired         = fmt.Errorf("code<%d>: token is expired", CodeTokenExpired)
	ErrTokenIsNull          = fmt.Errorf("code<%d>: token is null", CodeTokenIsNull)
	ErrTokenCheckFailed     = fmt.Errorf("code<%d>: check token failed", CodeTokenCheckFailed)
	ErrPublishMessageFailed = fmt.Errorf("code<%d>: publish message error", CodePublishMessageFailed)
)
