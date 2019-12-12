package internal

// A Breaker carries a cancellation signal to break an action execution.
type Breaker interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
}

// A BreakCloser carries a cancellation signal to break an action execution
// and can release resources associated with it.
type BreakCloser interface {
	Breaker
	// Close closes the Done channel and releases resources associated with it.
	Close()
}

// A SafeLock carries of getting an exclusive lock to access
// a critical section with the ability to interrupt the action.
type SafeLock interface {
	// Lock locks a mutex. If the lock is already in use,
	// the calling goroutine blocks until the mutex is available
	// or an error occurred.
	Lock(Breaker) error
	// Unlock unlocks a mutex. It could return an error if the mutex
	// is not locked on entry to Unlock.
	Unlock(Breaker) error
}

type Semaphore interface {
	Acquire(Breaker, uint32) error
	TryAcquire(uint32) bool
	Release(uint32) (uint32, error)
}

type Observable interface {
	Count() uint32
	Limit() uint32
}

type Resizable interface {
	SetCapacity(uint32) uint32
}
