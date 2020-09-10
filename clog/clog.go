package clog

import (
	"sync/atomic"
)

// LogProvider RFC5424 log message levels only Debugf Warnf and Errorf
type LogProvider interface {
	Criticalf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

// Option logger option
type Option func(*Clog)

// WithLogger set logger provider
func WithLogger(p LogProvider) Option {
	return func(c *Clog) {
		if p != nil {
			c.provider = p
		}
	}
}

// WithEnableLogger enable
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

// Errorf Log ERROR level message.
func (sf Clog) Errorf(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Errorf(format, v...)
	}
}

// Warnf Log WARN level message.
func (sf Clog) Warnf(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Warnf(format, v...)
	}
}

// Debugf Log DEBUG level message.
func (sf Clog) Debugf(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Debugf(format, v...)
	}
}

// Criticalf Log CRITICAL level message.
func (sf Clog) Criticalf(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Criticalf(format, v...)
	}
}
