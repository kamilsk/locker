package locker

import (
	"time"

	"github.com/kamilsk/locker/internal"
)

func Distributed(ttl time.Duration) *dlock {
	return &dlock{}
}

type dlock struct{}

func (lock *dlock) Lock(internal.Breaker) error {
	return nil
}

func (lock *dlock) Unlock(internal.Breaker) error {
	return nil
}
