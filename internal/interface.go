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

// A Locker carries of getting an exclusive lock to access
// a critical section with the ability to interrupt the action.
type Locker interface {
	// Lock takes an exclusive lock. If the lock is already in use,
	// the calling goroutine blocks until the mutex is available or
	// an error occurred, e.g. if the Breaker is done.
	Lock(Breaker) error
	// Unlock releases an exclusive lock. It could return an error
	// if the mutex is not locked on entry to Unlock or
	// the Breaker is done.
	Unlock(Breaker) error
}

// A FastLocker is a Locker with a possibility to take an exclusive lock
// fast or failure if it not possible at that moment.
type FastLocker interface {
	Locker
	// TryLock is a fail-fast version of the Lock method.
	// It returns true if the mutex is locked by the calling goroutine
	// or false otherwise.
	TryLock() bool
}

// A SafeLocker is a Locker with a guarantee that the unlock operation
// is achievable in a finite time.
type SafeLocker interface {
	Locker
	// MustUnlock is a fail-fast version of the Unlock method.
	// It is a runtime error if the mutex is not locked on entry to Unlock.
	MustUnlock()
}

type Semaphore interface {
	Acquire(Breaker, uint32) error
	Release(uint32) (uint32, error)
}

type FastSemaphore interface {
	Semaphore
	TryAcquire(uint32) bool
}

type Observable interface {
	Count() uint32
	Limit() uint32
}

type Resizable interface {
	SetCapacity(uint32) uint32
}
