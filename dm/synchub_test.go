package dm

import (
	"errors"
	"testing"
	"time"
)

func TestSyncHub(t *testing.T) {
	hub := NewSyncHub()
	donefun := func(id int, tm time.Duration, err error) {
		time.Sleep(tm)
		hub.Done(id, err)
	}

	t.Run("done not block", func(t *testing.T) {
		donefun(0, time.Millisecond*10, nil)
		t.Log("done not block")
	})

	t.Run("timeout", func(t *testing.T) {
		go donefun(0, DefaultExpiration, nil)
		err := hub.Wait(0)
		if !errors.Is(err, ErrWaitMessageTimeout) {
			t.Errorf("hub wait should be timeout")
		}
	})

	t.Run("nil", func(t *testing.T) {
		go donefun(1, DefaultWaitTimeout/2, nil)
		err := hub.Wait(1)
		if !errors.Is(err, nil) {
			t.Errorf("hub wait should be nil")
		}
	})

	t.Run("by user send", func(t *testing.T) {
		go donefun(1, DefaultWaitTimeout/2, ErrNotFound)
		err := hub.Wait(1)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("hub wait should be %+v", ErrNotFound)
		}
	})
}
