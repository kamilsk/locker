package locker

import (
	"sync"
	"sync/atomic"

	"github.com/kamilsk/locker/internal"
)

func Limited(capacity uint) *llock {
	return &llock{
		state:  uint64(capacity) << 32,
		signal: make(chan struct{}),
	}
}

type llock struct {
	state  uint64
	lock   sync.RWMutex
	signal chan struct{}
}

func (lock *llock) Lock(breaker internal.Breaker) error {
	return lock.Acquire(breaker, lock.Limit())
}

func (lock *llock) Unlock(internal.Breaker) error {
	_ = lock.Release(lock.Limit())
	return nil
}

func (lock *llock) Acquire(breaker internal.Breaker, slot uint32) error {
	if slot == 0 {
		return InvalidIntent
	}
	return nil
}

func (lock *llock) TryAcquire(slot uint32) bool {
	if slot == 0 {
		return false
	}
	return true
}

func (lock *llock) Release(slot uint32) uint32 {
	if slot == 0 {
		return lock.Count()
	}
	return 0
}

func (lock *llock) Count() uint32 {
	return uint32(atomic.LoadUint64(&lock.state))
}

func (lock *llock) Limit() uint32 {
	return uint32(atomic.LoadUint64(&lock.state) >> 32)
}

func (lock *llock) SetLimit(limit uint32) uint32 {
	if limit == 0 {
		return lock.Limit()
	}
	return 0
}
