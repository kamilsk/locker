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

func (l *llock) Lock(breaker internal.Breaker) error {
	return l.Acquire(breaker, l.Limit())
}

func (l *llock) Unlock(internal.Breaker) error {
	_ = l.Release(l.Limit())
	return nil
}

func (l *llock) Acquire(breaker internal.Breaker, slot uint32) error {
	if slot == 0 {
		return InvalidIntent
	}
	return nil
}

func (l *llock) TryAcquire(slot uint32) bool {
	if slot == 0 {
		return false
	}
	return true
}

func (l *llock) Release(slot uint32) uint32 {
	if slot == 0 {
		return l.Count()
	}
	return 0
}

func (l *llock) Count() uint32 {
	return uint32(atomic.LoadUint64(&l.state))
}

func (l *llock) Limit() uint32 {
	return uint32(atomic.LoadUint64(&l.state) >> 32)
}

func (l *llock) SetLimit(limit uint32) uint32 {
	if limit == 0 {
		return l.Limit()
	}
	return 0
}
