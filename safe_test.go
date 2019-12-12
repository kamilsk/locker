package locker_test

import (
	"context"
	"sync"
	"testing"
	"time"

	. "github.com/kamilsk/locker"
)

func TestSafe(t *testing.T) {
	lock := Safe()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := lock.Lock(ctx); err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}

	if err := lock.Unlock(ctx); err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}

	t.Run("unlock not-locked mutex", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()

		if err := lock.Unlock(ctx); err != Interrupted {
			t.Error("unexpected error value")
			t.FailNow()
		}
	})

	t.Run("double lock", func(t *testing.T) {
		if err := lock.Lock(ctx); err != nil {
			t.Error("unexpected error")
			t.FailNow()
		}

		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()

		if err := lock.Lock(ctx); err != Interrupted {
			t.Error("unexpected error value")
			t.FailNow()
		}
	})

	t.Run("stress test", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for range make([]struct{}, 1000) {
			wg.Add(2)
			go func() {
				if err := lock.Lock(ctx); err != nil {
					t.Error("unexpected error")
					t.FailNow()
				}
				wg.Done()
			}()
			go func() {
				if err := lock.Unlock(ctx); err != nil {
					t.Error("unexpected error")
					t.FailNow()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
}
