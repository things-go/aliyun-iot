package aHttp

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

// CodeError code 错误
type CodeError struct {
	code    int
	message string
}

// NewCodeError 生成code错误
func NewCodeError(code int, message string) error {
	return &CodeError{code, message}
}

// Error 实现error接口
func (sf *CodeError) Error() string {
	return fmt.Sprintf("code<%d>: %s", sf.code, sf.message)
}

// Code unwrap then got code
func (sf *CodeError) Code() int {
	return sf.code
}
