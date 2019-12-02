package locker

import (
	"crypto/md5"
	"hash"
	"sync"
)

func MultiMutex(capacity int) multimx {
	return multimx{hash: md5.New(), set: make([]sync.Mutex, capacity)}
}

type multimx struct {
	hash hash.Hash
	set  []sync.Mutex
}

func (mx multimx) Mutex(key string) *sync.Mutex {
	var spot int
	for _, point := range mx.hash.Sum([]byte(key)) {
		spot += int(point)
	}
	spot %= len(mx.set)
	return &mx.set[spot]
}
