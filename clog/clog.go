package clog

import (
	"sync/atomic"
)

// LogProvider RFC5424 log message levels only Debug Warn and Error
type LogProvider interface {
	Critical(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Option func(*Clog)

func WithLogger(p LogProvider) Option {
	return func(c *Clog) {
		if p != nil {
			c.provider = p
		}
	}
}

// SetLogProvider set provider provider
func WithEnableLogger() Option {
	return func(c *Clog) {
		c.LogMode(true)
	}
}

// Clog Logger internal implement
type Clog struct {
	provider LogProvider
	// is log output enabled,1: enable, 0: disable
	has uint32
}

// New new a Logger, default use Discard
func New(opts ...Option) *Clog {
	c := &Clog{NewDiscard(), 0}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// LogMode set enable or disable log output when you has set provider
func (sf *Clog) LogMode(enable bool) {
	if enable {
		atomic.StoreUint32(&sf.has, 1)
	} else {
		atomic.StoreUint32(&sf.has, 0)
	}
}

// Critical Log CRITICAL level message.
func (sf Clog) Critical(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Critical(format, v...)
	}
}

// Error Log ERROR level message.
func (sf Clog) Error(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Error(format, v...)
	}
}

// Warn Log WARN level message.
func (sf Clog) Warn(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Warn(format, v...)
	}
}

// Debug Log DEBUG level message.
func (sf Clog) Debug(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Debug(format, v...)
	}
}
