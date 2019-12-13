package internal

import "context"

// Wrap wraps the context and its cancel function into BreakCloser.
func Wrap(ctx context.Context, cancel context.CancelFunc) BreakCloser {
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
