package locker

import (
	"crypto/md5"
	"hash"
	"sync"

	"github.com/kamilsk/locker/internal"
)

func Set(capacity uint, options ...SetOption) *mset {
	set := &mset{set: make([]sync.Mutex, capacity), size: uint64(capacity)}
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

type SetOption func(*mset)

func SetWithHash(builder func() hash.Hash) SetOption {
	return func(c *mset) { c.hash = builder }
}

func SetWithMapping(index func([]byte, uint64) uint64) SetOption {
	return func(c *mset) { c.idx = index }
}

type mset struct {
	hash func() hash.Hash
	idx  func([]byte, uint64) uint64
	set  []sync.Mutex
	size uint64
}

func (c *mset) ByFingerprint(fingerprint []byte) *sync.Mutex {
	h := c.hash()
	_, _ = h.Write(fingerprint)
	shard := c.idx(h.Sum(nil), c.size)
	h.Reset()
	return &c.set[shard]
}

func (c *mset) ByKey(key string) *sync.Mutex {
	return c.ByFingerprint([]byte(key))
}

func (c *mset) ByVirtualShard(shard uint64) *sync.Mutex {
	return &c.set[shard%c.size]
}

func RWSet(capacity uint, options ...RWSetOption) *rwset {
	set := &rwset{set: make([]sync.RWMutex, capacity), size: uint64(capacity)}
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

type RWSetOption func(*rwset)

func RWSetWithHash(builder func() hash.Hash) RWSetOption {
	return func(c *rwset) { c.hash = builder }
}

func RWSetWithMapping(index func([]byte, uint64) uint64) RWSetOption {
	return func(c *rwset) { c.idx = index }
}

type rwset struct {
	hash func() hash.Hash
	idx  func([]byte, uint64) uint64
	set  []sync.RWMutex
	size uint64
}

func (c *rwset) ByFingerprint(fingerprint []byte) *sync.RWMutex {
	h := c.hash()
	_, _ = h.Write(fingerprint)
	shard := c.idx(h.Sum(nil), c.size)
	h.Reset()
	return &c.set[shard]
}

func (c *rwset) ByKey(key string) *sync.RWMutex {
	return c.ByFingerprint([]byte(key))
}

func (c *rwset) ByVirtualShard(shard uint64) *sync.RWMutex {
	return &c.set[shard%c.size]
}
