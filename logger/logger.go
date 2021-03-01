package logger

// Logger log interface
type Logger interface {
	Debugf(format string, arg ...interface{})
	Infof(format string, arg ...interface{})
	Warnf(format string, arg ...interface{})
	Errorf(format string, arg ...interface{})
	DPanicf(format string, arg ...interface{})
	Fatalf(format string, arg ...interface{})
}
