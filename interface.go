package locker

import "context"

// A Breaker carries a cancellation signal to break an action execution.
//
// It is a subset of `context.Context` and `github.com/kamilsk/breaker.Breaker`.
type Breaker interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
}

// A BreakCloser carries a cancellation signal to break an action execution
// and can release resources associated with it.
//
// It is a subset of `github.com/kamilsk/breaker.Breaker`.
type BreakCloser interface {
	Breaker
	// Close closes the Done channel and releases resources associated with it.
	Close()
}

// A DistributedLock carries of getting an exclusive lock to access
// a critical section in a distributed system.
type DistributedLock interface {
	// Lock locks a distributed mutex. If the lock is already in use,
	// the calling goroutine blocks until the mutex is available, timeout exited
	// or a network error occurred.
	Lock(context.Context) error
	// Unlock unlocks a distributed mutex. It returns an error if the mutex is not locked
	// on entry to Unlock, timeout exited or a network error occurred.
	Unlock(context.Context) error
}
