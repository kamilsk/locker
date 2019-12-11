package locker

import "github.com/kamilsk/locker/internal"

func Distributed() *dlock {
	return &dlock{}
}

type dlock struct{}

func (l *dlock) Lock(internal.Breaker) error {
	return nil
}

func (l *dlock) Unlock(internal.Breaker) error {
	return nil
}
