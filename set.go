package locker

import (
	"crypto/md5"
	"hash"
	"sync"

	"github.com/kamilsk/locker/internal"
)

func Set(capacity uint, options ...SetOption) *MutexSet {
	set := &MutexSet{set: make([]sync.Mutex, capacity), size: uint64(capacity)}
	for _, option := range options {
		option(set)
	}
	if set.hash == nil {
		set.hash = md5.New
	}
	if set.idx == nil {
		set.idx = internal.ShardNumberFast
	}
	return set
}

type SetOption func(*MutexSet)

func SetWithHash(builder func() hash.Hash) SetOption {
	return func(c *MutexSet) { c.hash = builder }
}

func SetWithMapping(index func([]byte, uint64) uint64) SetOption {
	return func(c *MutexSet) { c.idx = index }
}

type MutexSet struct {
	hash func() hash.Hash
	idx  func([]byte, uint64) uint64
	set  []sync.Mutex
	size uint64
}

func (c MutexSet) ByFingerprint(fingerprint []byte) *sync.Mutex {
	h := c.hash()
	_, _ = h.Write(fingerprint)
	shard := c.idx(h.Sum(nil), c.size)
	h.Reset()
	return &c.set[shard]
}

func (c MutexSet) ByKey(key string) *sync.Mutex {
	return c.ByFingerprint([]byte(key))
}

func (c MutexSet) ByVirtualShard(shard uint64) *sync.Mutex {
	return &c.set[shard%c.size]
}

func RWSet(capacity uint, options ...RWSetOption) *RWMutexSet {
	set := &RWMutexSet{set: make([]sync.RWMutex, capacity), size: uint64(capacity)}
	for _, option := range options {
		option(set)
	}
	if set.hash == nil {
		set.hash = md5.New
	}
	if set.idx == nil {
		set.idx = internal.ShardNumberFast
	}
	return set
}

type RWSetOption func(*RWMutexSet)

func RWSetWithHash(builder func() hash.Hash) RWSetOption {
	return func(c *RWMutexSet) { c.hash = builder }
}

func RWSetWithMapping(index func([]byte, uint64) uint64) RWSetOption {
	return func(c *RWMutexSet) { c.idx = index }
}

type RWMutexSet struct {
	hash func() hash.Hash
	idx  func([]byte, uint64) uint64
	set  []sync.RWMutex
	size uint64
}

func (c RWMutexSet) ByFingerprint(fingerprint []byte) *sync.RWMutex {
	h := c.hash()
	_, _ = h.Write(fingerprint)
	shard := c.idx(h.Sum(nil), c.size)
	h.Reset()
	return &c.set[shard]
}

func (c RWMutexSet) ByKey(key string) *sync.RWMutex {
	return c.ByFingerprint([]byte(key))
}

func (c RWMutexSet) ByVirtualShard(shard uint64) *sync.RWMutex {
	return &c.set[shard%c.size]
}
