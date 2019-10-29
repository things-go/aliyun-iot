package dm

import (
	"errors"
)

var ErrInvalidURI = errors.New("invalid URI")
var ErrNotFound = errors.New("not found")
var ErrInvalidParameter = errors.New("invalid parameter")
var ErrNotSupportMsgType = errors.New("not support message type")
var ErrNotSupportFeature = errors.New(" not support feature")