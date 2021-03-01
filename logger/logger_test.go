package logger

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestStd(_ *testing.T) {
	l := New(log.New(ioutil.Discard, "testing", log.LstdFlags), WithEnable(true))
	l.Debugf("")
	l.Infof("")
	l.Errorf("")
	l.Warnf("")
	l.DPanicf("")
	l.Mode(false)
	l.Fatalf("")
}

func TestDiscard(_ *testing.T) {
	l := NewDiscard()
	l.Debugf("")
	l.Infof("")
	l.Errorf("")
	l.Warnf("")
	l.DPanicf("")
	l.Fatalf("")
}
