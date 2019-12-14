package locker

import "github.com/kamilsk/locker/internal"

// Interruptible returns a new instance of safe interruptible mutex.
//
//  lock := locker.Interruptible()
//
//  var handler http.HandlerFunc = func(rw http.ResponseWriter, req *http.Request) {
//  	if err := lock.Lock(req.Context()); err != nil {
//  		http.Error(rw, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
//  		return
//  	}
//  	// critical section with lock protection
//  	if err := lock.Unlock(req.Context()); err != nil {
//  		// timeout occurred or connection is broken, but we must unlock it anyway
//  		go func() { _ = lock.Unlock(context.Background()) }()
//  	}
//  }
//
func Interruptible() *ilock {
	lock := make(ilock, 1)
	return &lock
}

type ilock chan struct{}

// Lock takes an exclusive lock. If the lock is already in use,
// the calling goroutine blocks until the mutex is available or
// an error occurred, e.g. if the Breaker is done.
func (lock *ilock) Lock(breaker internal.Breaker) error {
	select {
	case <-breaker.Done():
		return Interrupted
	case *lock <- struct{}{}:
		return nil
	}
}

// Unlock releases an exclusive lock. It could return an error
// if the mutex is not locked on entry to Unlock or
// the Breaker is done.
//
// The CriticalIssue error is an important part of the mutex API.
func (lock *ilock) Unlock(breaker internal.Breaker) error {
	select {
	case <-breaker.Done():
		return CriticalIssue
	case <-*lock:
		return nil
	}
}
