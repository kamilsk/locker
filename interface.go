package locker

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

// A SafeLock carries of getting an exclusive lock to access
// a critical section with a timeout specified by context.
type SafeLock interface {
	// Lock locks a mutex. If the lock is already in use,
	// the calling goroutine blocks until the mutex is available
	// or an error occurred.
	Lock(Breaker) error
	// Unlock unlocks a mutex. It returns an error if the mutex is not locked
	// on entry to Unlock or a timeout occurred.
	Unlock(Breaker) error
}
