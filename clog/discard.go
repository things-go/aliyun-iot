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

func (sf Discard) Critical(string, ...interface{}) {}
func (sf Discard) Error(string, ...interface{})    {}
func (sf Discard) Warn(string, ...interface{})     {}
func (sf Discard) Debug(string, ...interface{})    {}

// default log
type Logger struct {
	*log.Logger
}

var _ LogProvider = (*Logger)(nil)

func NewLogger(l *log.Logger) Logger {
	return Logger{l}
}

// Critical Log CRITICAL level message.
func (sf Logger) Critical(format string, v ...interface{}) {
	sf.Printf("[C]: "+format, v...)
}

// Error Log ERROR level message.
func (sf Logger) Error(format string, v ...interface{}) {
	sf.Printf("[E]: "+format, v...)
}

// Warn Log WARN level message.
func (sf Logger) Warn(format string, v ...interface{}) {
	sf.Printf("[W]: "+format, v...)
}

// Debug Log DEBUG level message.
func (sf Logger) Debug(format string, v ...interface{}) {
	sf.Printf("[D]: "+format, v...)
}
