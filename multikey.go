package locker

import (
	"hash"
	"sync"
)

func MultiMutex(capacity int, hash hash.Hash) multimx {
	return multimx{hash: hash, set: make([]sync.Mutex, capacity)}
}

type multimx struct {
	hash hash.Hash
	set  []sync.Mutex
}

func (mx multimx) ByFingerprint(fingerprint []byte) *sync.Mutex {
	_, _ = mx.hash.Write(fingerprint)
	shard := BytesToUint64Mod(mx.hash.Sum(nil), uint64(len(mx.set)))
	mx.hash.Reset()
	return mx.ByVirtualShard(shard)
}

func (mx multimx) ByKey(key string) *sync.Mutex {
	return mx.ByFingerprint([]byte(key))
}

func (mx multimx) ByVirtualShard(shard uint64) *sync.Mutex {
	return &mx.set[shard%uint64(len(mx.set))]
}
