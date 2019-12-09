package set

import (
	"crypto/md5"
	"hash"
)

func NewContainer(capacity uint, options ...ContainerOption) *Container {
	set := &Container{set: make([]T, capacity), size: uint64(capacity)}
	for _, option := range options {
		option(set)
	}
	if set.hash == nil {
		set.hash = md5.New
	}
	if set.idx == nil {
		// TODO:debt how to use internal?
		set.idx = func(checksum []byte, size uint64) uint64 {
			sumModWithoutOverflow := func(a, b, d uint64) uint64 {
				if a < d-b {
					// 1. a + b < d ?
					//    a < d - b
					return a + b
				} else if a > d-b {
					// 2. a + b > d ?
					//      a > d -b
					return d - ((d - a) + (d - b))
				} else {
					// 3. a + b == d
					return 0
				}
			}

			var shard, factor uint64 = 0, 1
			for i := len(checksum) - 1; i >= 0; i-- {
				curByte := checksum[i]
				for j := 0; j < 8; j++ {
					// we iterate over bits to use sum instead of multiplication for the factor:
					// on each iteration: factor = factor * 2 % size <=> (factor + factor) % size
					bit := curByte % 2
					curByte /= 2
					if bit == 1 {
						shard = sumModWithoutOverflow(shard, factor, size)
					}
					factor = sumModWithoutOverflow(factor, factor, size)
				}
			}
			return shard
		}
	}
	return set
}

type T interface{}

type Container struct {
	hash func() hash.Hash
	idx  func([]byte, uint64) uint64
	set  []T
	size uint64
}

type ContainerOption func(*Container)

func ContainerWithHash(builder func() hash.Hash) ContainerOption {
	return func(c *Container) { c.hash = builder }
}

func ContainerWithMapping(index func([]byte, uint64) uint64) ContainerOption {
	return func(c *Container) { c.idx = index }
}

func (c Container) ByFingerprint(fingerprint []byte) *T {
	h := c.hash()
	_, _ = h.Write(fingerprint)
	shard := c.idx(h.Sum(nil), c.size)
	h.Reset()
	return &c.set[shard]
}

func (c Container) ByKey(key string) *T {
	return c.ByFingerprint([]byte(key))
}

func (c Container) ByVirtualShard(shard uint64) *T {
	return &c.set[shard%c.size]
}
