package locker_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/kamilsk/locker"
	"github.com/kamilsk/locker/internal"
)

func ExampleInterruptible() {
	data, lock := []string{"1 ", "2 ", "3 ", "4"}, Interruptible()

	var handler http.HandlerFunc = func(rw http.ResponseWriter, req *http.Request) {
		if err := lock.Lock(req.Context()); err != nil {
			http.Error(rw, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
			return
		}
		defer lock.MustUnlock()
		// critical section with lock protection
		// only one goroutine can be here one moment in time
		_, _ = rw.Write([]byte(data[0]))
		data = data[1:]
	}

	rec, wg := httptest.NewRecorder(), &sync.WaitGroup{}
	wg.Add(len(data))
	for range data {
		go func() {
			handler.ServeHTTP(rec, &http.Request{})
			wg.Done()
		}()
	}
	wg.Wait()

	_, _ = rec.Body.WriteTo(os.Stdout)
	// output: 1 2 3 4
}

func TestInterruptible(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	lock := Interruptible()
	if err := lock.Lock(ctx); err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}
	if lock.TryLock() {
		t.Error("unexpected double lock")
		t.FailNow()
	}
	if err := lock.Unlock(ctx); err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}

	t.Run("try to unlock not-locked mutex", func(t *testing.T) {
		defer func() {
			if r := recover(); r != CriticalIssue {
				t.Error("panic with CriticalIssue is expected")
			}
		}()

		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()
		if err := lock.Unlock(ctx); err != InvalidIntent {
			t.Error("unexpected error value")
			t.FailNow()
		}
		lock.MustUnlock()
	})

	t.Run("try to use short-lived breaker", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("failed to implement test case: %+v", r)
				t.FailNow()
			}
		}()
		for range make([]struct{}, 10) {
			if !lock.TryLock() {
				t.Error("lock is expected")
				t.FailNow()
			}
			breaker := Wrap(context.WithCancel(ctx))
			breaker.Close()
			if err := lock.Unlock(breaker); err != nil {
				if err != InvalidIntent {
					lock.MustUnlock()
					t.Error("unexpected error value")
					t.FailNow()
				}
				// unlock failed, success
				break
			}
			// else unlock won, repeat
		}
		lock.MustUnlock()
	})

	t.Run("try to call lock multiple times", func(t *testing.T) {
		if err := lock.Lock(ctx); err != nil {
			t.Error("unexpected error")
			t.FailNow()
		}
		for range make([]struct{}, 10) {
			if err := lock.Lock(Wrap(context.WithTimeout(ctx, time.Millisecond))); err != Interrupted {
				t.Error("unexpected error value")
				t.FailNow()
			}
		}
		if err := lock.Unlock(Wrap(context.WithTimeout(ctx, time.Millisecond))); err != nil {
			t.Error("unexpected error")
			t.FailNow()
		}
	})
}

func TestInterruptible_StressTest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	lock := Interruptible()
	if *stress {
		wg := sync.WaitGroup{}
		for range make([]struct{}, 1000) {
			wg.Add(1)
			go func() {
				if err := lock.Lock(ctx); err != nil {
					t.Error("unexpected error")
					t.FailNow()
				}
				if err := lock.Unlock(ctx); err != nil {
					t.Error("unexpected error")
					t.FailNow()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

// BenchmarkInterruptible/interruptible_locker-4         	 7418108	       156 ns/op	       0 B/op	       0 allocs/op
// BenchmarkInterruptible/built-in_locker-4              	68697568	        17.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkInterruptible(b *testing.B) {
	ctx := context.Background()

	b.Run("interruptible locker", func(b *testing.B) {
		var lock internal.Locker = Interruptible()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = lock.Lock(ctx)
			_ = lock.Unlock(ctx)
		}
	})

	b.Run("built-in locker", func(b *testing.B) {
		var lock sync.Locker = &sync.Mutex{}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			lock.Lock()
			lock.Unlock()
		}
	})
}
