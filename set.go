package locker

import (
	"crypto/md5"
	"hash"
	"sync"

	"github.com/kamilsk/locker/internal"
)

func Set(capacity int, options ...SetOption) *MutexSet {
	if capacity < 1 {
		panic("capacity must be greater than zero")
	}
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
	return func(set *MutexSet) { set.hash = builder }
}

func SetWithMapping(index func([]byte, uint64) uint64) SetOption {
	return func(set *MutexSet) { set.idx = index }
}

type MutexSet struct {
	hash func() hash.Hash
	idx  func([]byte, uint64) uint64
	set  []sync.Mutex
	size uint64
}

func (mx MutexSet) ByFingerprint(fingerprint []byte) *sync.Mutex {
	h := mx.hash()
	_, _ = h.Write(fingerprint)
	shard := mx.idx(h.Sum(nil), mx.size)
	h.Reset()
	return &mx.set[shard]
}

func (mx MutexSet) ByKey(key string) *sync.Mutex {
	return mx.ByFingerprint([]byte(key))
}

func (mx MutexSet) ByVirtualShard(shard uint64) *sync.Mutex {
	return &mx.set[shard%mx.size]
}

func RWSet(capacity int, options ...RWSetOption) *RWMutexSet {
	if capacity < 1 {
		panic("capacity must be greater than zero")
	}
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
	return func(set *RWMutexSet) { set.hash = builder }
}

func RWSetWithMapping(index func([]byte, uint64) uint64) RWSetOption {
	return func(set *RWMutexSet) { set.idx = index }
}

type RWMutexSet struct {
	hash func() hash.Hash
	idx  func([]byte, uint64) uint64
	set  []sync.RWMutex
	size uint64
}

func (mx RWMutexSet) ByFingerprint(fingerprint []byte) *sync.RWMutex {
	h := mx.hash()
	_, _ = h.Write(fingerprint)
	shard := mx.idx(h.Sum(nil), mx.size)
	h.Reset()
	return &mx.set[shard]
}

func (mx RWMutexSet) ByKey(key string) *sync.RWMutex {
	return mx.ByFingerprint([]byte(key))
}

func (mx RWMutexSet) ByVirtualShard(shard uint64) *sync.RWMutex {
	return &mx.set[shard%mx.size]
}
