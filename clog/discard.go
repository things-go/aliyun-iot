package clog

import (
	"log"
)

// Discard is an Logger on which all Write calls succeed
// without doing anything.
type Discard struct{}

var _ LogProvider = (*Discard)(nil)

// NewDiscard a discard Logger
func NewDiscard() Discard { return Discard{} }

// Criticalf implement interface LogProvider
func (sf Discard) Criticalf(string, ...interface{}) {}

// Errorf implement interface LogProvider
func (sf Discard) Errorf(string, ...interface{}) {}

// Warnf implement interface LogProvider
func (sf Discard) Warnf(string, ...interface{}) {}

// Debugf implement interface LogProvider
func (sf Discard) Debugf(string, ...interface{}) {}

// Logger default log
type Logger struct {
	*log.Logger
}

var _ LogProvider = (*Logger)(nil)

// NewLogger new logger
func NewLogger(l *log.Logger) Logger {
	return Logger{l}
}

// Criticalf Log CRITICAL level message.
func (sf Logger) Criticalf(format string, v ...interface{}) {
	sf.Printf("[C]: "+format, v...)
}

// Errorf Log ERROR level message.
func (sf Logger) Errorf(format string, v ...interface{}) {
	sf.Printf("[E]: "+format, v...)
}

// Warnf Log WARN level message.
func (sf Logger) Warnf(format string, v ...interface{}) {
	sf.Printf("[W]: "+format, v...)
}

// Debugf Log DEBUG level message.
func (sf Logger) Debugf(format string, v ...interface{}) {
	sf.Printf("[D]: "+format, v...)
}
