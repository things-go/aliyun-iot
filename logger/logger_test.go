package logger

import (
	"testing"
)

func TestDiscard(_ *testing.T) {
	l := NewDiscard()
	l.Debugf("")
	l.Infof("")
	l.Errorf("")
	l.Warnf("")
	l.DPanicf("")
	l.Fatalf("")
}
