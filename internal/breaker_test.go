package internal_test

import (
	"context"
	"testing"

	. "github.com/kamilsk/locker/internal"
)

func TestWrap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	breaker := Wrap(ctx, cancel)
	if ctx.Err() != nil {
		t.Error("unexpected error")
		t.FailNow()
	}
	breaker.Close()
	if ctx.Err() != context.Canceled {
		t.Error("canceled context is expected")
		t.FailNow()
	}
}
