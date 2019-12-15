package locker_test

import (
	"context"
	"flag"
	"time"
)

var (
	stress  = flag.Bool("stress-test", false, "run stress tests")
	timeout = flag.Duration("timeout", time.Second, "use custom timeout, e.g. to debug")
)

// Wrap wraps the context and its cancel function into BreakCloser.
func Wrap(ctx context.Context, cancel context.CancelFunc) *wrapper {
	return &wrapper{ctx, cancel}
}

type wrapper struct {
	context.Context
	cancel context.CancelFunc
}

// Close closes the Done channel and releases resources associated with it.
func (breaker *wrapper) Close() {
	breaker.cancel()
}
