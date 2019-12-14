package locker_test

import (
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/kamilsk/locker"
	"github.com/kamilsk/locker/internal"
)

var timeout = flag.Duration("timeout", time.Second, "use custom timeout, e.g. to debug")

func ExampleInterruptible() {
	data, lock := []string{"1 ", "2 ", "3 ", "4"}, Interruptible()

	var handler http.HandlerFunc = func(rw http.ResponseWriter, req *http.Request) {
		if err := lock.Lock(req.Context()); err != nil {
			http.Error(rw, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
			return
		}
		_, _ = rw.Write([]byte(data[0]))
		data = data[1:]
		if err := lock.Unlock(req.Context()); err != nil {
			go func() { _ = lock.Unlock(context.Background()) }()
		}
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
	if err := lock.Unlock(ctx); err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}

	t.Run("try to unlock not-locked mutex", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()
		if err := lock.Unlock(ctx); err != CriticalIssue {
			t.Error("unexpected error value")
			t.FailNow()
		}
	})

	t.Run("try to use short-lived breaker", func(t *testing.T) {
		for range make([]struct{}, 10) {
			if err := lock.Lock(ctx); err != nil {
				t.Error("unexpected error")
				t.FailNow()
			}
			breaker := internal.Wrap(context.WithCancel(ctx))
			breaker.Close()
			if err := lock.Unlock(breaker); err != nil {
				if err != CriticalIssue {
					t.Error("unexpected error value")
					t.FailNow()
				}
				// unlock failed, success
				break
			}
			// else unlock won, repeat
		}
		if err := lock.Unlock(internal.Wrap(context.WithTimeout(ctx, time.Millisecond))); err != nil {
			t.Error("failed to implement test case")
			t.FailNow()
		}
	})

	t.Run("try to call lock multiple times", func(t *testing.T) {
		if err := lock.Lock(ctx); err != nil {
			t.Error("unexpected error")
			t.FailNow()
		}
		for range make([]struct{}, 10) {
			if err := lock.Lock(internal.Wrap(context.WithTimeout(ctx, time.Millisecond))); err != Interrupted {
				t.Error("unexpected error value")
				t.FailNow()
			}
		}
		if err := lock.Unlock(internal.Wrap(context.WithTimeout(ctx, time.Millisecond))); err != nil {
			t.Error("unexpected error")
			t.FailNow()
		}
	})

	t.Run("stress test", func(t *testing.T) {
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
	})
}

// BenchmarkInterruptible/interruptible_mutex-4         	 7655840	       164 ns/op	       0 B/op	       0 allocs/op
// BenchmarkInterruptible/built-in_mutex-4              	92805457	        12.6 ns/op	       0 B/op	       0 allocs/op
func BenchmarkInterruptible(b *testing.B) {
	ctx := context.Background()

	b.Run("interruptible mutex", func(b *testing.B) {
		lock := Interruptible()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = lock.Lock(ctx)
			_ = lock.Unlock(ctx)
		}
	})

	b.Run("built-in mutex", func(b *testing.B) {
		mx := sync.Mutex{}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mx.Lock()
			mx.Unlock()
		}
	})
}
