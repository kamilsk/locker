package locker

import (
	"encoding/hex"
	"hash"
	"math/big"
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
	base := big.NewInt(0)
	_, _ = base.SetString(hex.EncodeToString(mx.hash.Sum(nil)), 16)
	mx.hash.Reset()
	shard := big.NewInt(0).Mod(base, big.NewInt(int64(len(mx.set))))
	return mx.ByVirtualShard(shard.Uint64())
}

func (mx multimx) ByKey(key string) *sync.Mutex {
	return mx.ByFingerprint([]byte(key))
}

func (mx multimx) ByVirtualShard(shard uint64) *sync.Mutex {
	return &mx.set[shard%uint64(len(mx.set))]
}
