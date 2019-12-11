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
	guard  sync.RWMutex
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
	for {
		select {
		case <-breaker.Done():
			return Interrupted
		default:
		}
		state, count, limit := lock.splitState()
		if newCount := count + slot; newCount <= limit {
			if atomic.CompareAndSwapUint64(&lock.state, state, uint64(limit<<32+newCount)) {
				return nil
			}
			continue
		}
		lock.guard.RLock()
		signal := lock.signal
		lock.guard.RUnlock()

		if atomic.LoadUint64(&lock.state) != state {
			continue
		}

		select {
		case <-breaker.Done():
			return Interrupted
		case <-signal:
			// potentially have a place
		}
	}
}

func (lock *llock) TryAcquire(slot uint32) bool {
	if slot == 0 {
		return false
	}
	for {
		state, count, limit := lock.splitState()
		if newCount := count + slot; newCount <= limit {
			if atomic.CompareAndSwapUint64(&lock.state, state, uint64(limit<<32+newCount)) {
				return true
			}
			continue
		}
		return false
	}
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

func (lock *llock) splitState() (uint64, uint32, uint32) {
	state := atomic.LoadUint64(&lock.state)
	return state, uint32(state), uint32(state >> 32)
}
