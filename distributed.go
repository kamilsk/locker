package locker

import "github.com/kamilsk/locker/internal"

func Distributed() *dlock {
	return &dlock{}
}

type dlock struct{}

func (lock *dlock) Lock(internal.Breaker) error {
	return nil
}

func (lock *dlock) Unlock(internal.Breaker) error {
	return nil
}
