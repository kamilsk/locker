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
//  	defer lock.MustUnlock()
//  	// critical section with lock protection
//  	// only one goroutine can be here one moment in time
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

// TryLock is a fail-fast version of the Lock method.
// It returns true if the mutex is locked by the calling goroutine
// or false otherwise.
func (lock *ilock) TryLock() bool {
	select {
	case *lock <- struct{}{}:
		return true
	default:
		return false
	}
}

// Unlock releases an exclusive lock. It could return an error
// if the mutex is not locked on entry to Unlock or
// the Breaker is done. In this case the calling goroutine
// needs to release the mutex in background.
//
//  var handler http.HandlerFunc = func(rw http.ResponseWriter, req *http.Request) {
//  	if err := lock.Lock(req.Context()); err != nil {
//  		http.Error(rw, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
//  		return
//  	}
//  	// critical section with lock protection
//  	// only one goroutine can be here one moment in time
//  	if err := lock.Unlock(req.Context()); err != nil {
//  		// timeout occurred or connection is broken,
//  		// but we need to release the mutex anyway
//  		go lock.Unlock(context.Background())
//  		// or go lock.MustUnlock()
//  	}
//  }
//
func (lock *ilock) Unlock(breaker internal.Breaker) error {
	select {
	case <-breaker.Done():
		return InvalidIntent
	case <-*lock:
		return nil
	}
}

// MustUnlock is a fail-fast version of the Unlock method.
// It is a runtime error if the mutex is not locked on entry to Unlock.
func (lock *ilock) MustUnlock() {
	select {
	case <-*lock:
		return
	default:
		panic(CriticalIssue)
	}
}
