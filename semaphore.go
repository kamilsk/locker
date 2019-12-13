package locker

import (
	"sync"
	"sync/atomic"

	"github.com/kamilsk/locker/internal"
)

// Limited returns a new instance of resizable semaphore
// and non-blocking mutex.
//
// Fully reworked of github.com/kamilsk/semaphore,
// inspired by github.com/marusama/semaphore.
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
	_, err := lock.Release(lock.Limit())
	return err
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

func (lock *llock) Release(slot uint32) (uint32, error) {
	if slot == 0 {
		return lock.Count(), nil
	}
	for {
		state, count, limit := lock.splitState()
		if count < slot {
			return count, InvalidIntent
		}

		if newCount := count - slot; atomic.CompareAndSwapUint64(&lock.state, state, uint64(limit<<32+newCount)) {
			signal := make(chan struct{})

			lock.guard.Lock()
			broadcast := lock.signal
			lock.signal = signal
			lock.guard.Unlock()

			close(broadcast)
			return count, nil
		}
	}
}

func (lock *llock) Count() uint32 {
	return uint32(atomic.LoadUint64(&lock.state))
}

func (lock *llock) Limit() uint32 {
	return uint32(atomic.LoadUint64(&lock.state) >> 32)
}

func (lock *llock) SetCapacity(capacity uint32) uint32 {
	if capacity == 0 {
		return lock.Limit()
	}

	for {
		state, count, limit := lock.splitState()
		if atomic.CompareAndSwapUint64(&lock.state, state, uint64(capacity<<32+count)) {
			signal := make(chan struct{})

			lock.guard.Lock()
			broadcast := lock.signal
			lock.signal = signal
			lock.guard.Unlock()

			close(broadcast)
			return limit
		}
	}
}

func (lock *llock) splitState() (state uint64, count uint32, limit uint32) {
	state = atomic.LoadUint64(&lock.state)
	return state, uint32(state), uint32(state >> 32)
}
