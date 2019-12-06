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
	set := MutexSet{set: make([]sync.Mutex, capacity), size: uint64(capacity)}
	for _, option := range options {
		option(set)
	}
	if set.hash == nil {
		set.hash = md5.New()
	}
	if set.idx == nil {
		set.idx = internal.ShardNumberFast
	}
	return &set
}

type SetOption func(MutexSet)

func SetWithHash(hash hash.Hash) SetOption {
	return func(set MutexSet) { set.hash = hash }
}

func SetWithMapping(index func([]byte, uint64) uint64) SetOption {
	return func(set MutexSet) { set.idx = index }
}

type MutexSet struct {
	hash hash.Hash
	idx  func([]byte, uint64) uint64
	set  []sync.Mutex
	size uint64
}

func (mx MutexSet) ByFingerprint(fingerprint []byte) *sync.Mutex {
	_, _ = mx.hash.Write(fingerprint)
	shard := mx.idx(mx.hash.Sum(nil), mx.size)
	mx.hash.Reset()
	return &mx.set[shard]
}

func (mx MutexSet) ByKey(key string) *sync.Mutex {
	return mx.ByFingerprint([]byte(key))
}

func (mx MutexSet) ByVirtualShard(shard uint64) *sync.Mutex {
	return &mx.set[shard%mx.size]
}
