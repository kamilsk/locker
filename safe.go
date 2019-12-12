package locker

import "github.com/kamilsk/locker/internal"

// Safe returns a new instance of safe lock.
func Safe() *safe {
	lock := make(safe, 1)
	return &lock
}

type safe chan struct{}

// Lock locks a mutex. If the lock is already in use,
// the calling goroutine blocks until the mutex is available
// or an error occurred.
func (lock safe) Lock(breaker internal.Breaker) error {
	select {
	case <-breaker.Done():
		return Interrupted
	case lock <- struct{}{}:
		return nil
	}
}

// Unlock unlocks a mutex. It could return an error if the mutex
// is not locked on entry to Unlock.
func (lock safe) Unlock(breaker internal.Breaker) error {
	select {
	case <-breaker.Done():
		return Interrupted
	case <-lock:
		return nil
	}
}
