package locker

import "github.com/kamilsk/locker/internal"

func Safe() *safe {
	lock := make(safe, 1)
	return &lock
}

type safe chan struct{}

func (lock safe) Lock(breaker internal.Breaker) error {
	select {
	case <-breaker.Done():
		return Interrupted
	case lock <- struct{}{}:
		return nil
	}
}

func (lock safe) Unlock(breaker internal.Breaker) error {
	select {
	case <-breaker.Done():
		return Interrupted
	case <-lock:
		return nil
	}
}
