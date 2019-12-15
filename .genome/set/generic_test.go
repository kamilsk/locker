// +build genome

package set_test

import (
	"crypto/md5"
	"crypto/sha1"
	"flag"
	"hash"
	"math"
	"runtime"
	"testing"

	. "github.com/kamilsk/locker/genome/set"
)

var stress = flag.Bool("stress-test", false, "run stress tests")

func TestNewContainer(t *testing.T) {
	t.Run("with custom hash option", func(t *testing.T) {
		container := NewContainer(3, ContainerWithHash(sha1.New))
		if container.ByKey(runtime.GOOS) == container.ByKey(runtime.GOARCH) {
			t.Error("unexpected result")
			t.FailNow()
		}
	})
	t.Run("with custom mapping option", func(t *testing.T) {
		container := NewContainer(3, ContainerWithMapping(func([]byte, uint64) uint64 { return 0 }))
		if container.ByKey(runtime.GOOS) != container.ByKey(runtime.GOARCH) {
			t.Error("unexpected result")
			t.FailNow()
		}
	})
}

func TestContainer_ByFingerprint(t *testing.T) {
	t.Parallel()

	fingerprints := [...][]byte{[]byte(runtime.GOOS), []byte(runtime.GOARCH)}
	container := NewContainer(3)

	for _, fingerprint := range fingerprints {
		origin := container.ByFingerprint(fingerprint)
		for range make([]struct{}, 1000) {
			current := container.ByFingerprint(fingerprint)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, fingerprint := range fingerprints {
		current := container.ByFingerprint(fingerprint)
		for _, fingerprint := range fingerprints[i+1:] {
			next := container.ByFingerprint(fingerprint)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestContainer_ByKey(t *testing.T) {
	t.Parallel()

	keys := [...]string{runtime.GOOS, runtime.GOARCH}
	container := NewContainer(3)

	for _, key := range keys {
		origin := container.ByKey(key)
		for range make([]struct{}, 1000) {
			current := container.ByKey(key)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, key := range keys {
		current := container.ByKey(key)
		for _, key := range keys[i+1:] {
			next := container.ByKey(key)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestContainer_ByVirtualShard(t *testing.T) {
	t.Parallel()

	shards := [...]uint64{1, 5, 9}
	container := NewContainer(3)

	for _, shard := range shards {
		origin := container.ByVirtualShard(shard)
		for range make([]struct{}, 1000) {
			current := container.ByVirtualShard(shard)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, shard := range shards {
		current := container.ByVirtualShard(shard)
		for _, shard := range shards[i+1:] {
			next := container.ByVirtualShard(shard)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestContainer_StressTest(t *testing.T) {
	if *stress {
		for range make([]struct{}, 1000) {
			t.Run("by fingerprint", TestContainer_ByFingerprint)
			t.Run("by key", TestContainer_ByKey)
			t.Run("by virtual shard", TestContainer_ByVirtualShard)
		}
	}
}

func BenchmarkContainer_ByFingerprint(b *testing.B) {
	benchmarks := []struct {
		name string
		hash func() hash.Hash
	}{
		{name: "md5", hash: md5.New},
		{name: "sha1", hash: sha1.New},
	}
	fingerprint := []byte(runtime.GOOS)
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			container := NewContainer(3, ContainerWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = container.ByFingerprint(fingerprint)
			}
		})
	}
}

func BenchmarkContainer_ByKey(b *testing.B) {
	benchmarks := []struct {
		name string
		hash func() hash.Hash
	}{
		{name: "md5", hash: md5.New},
		{name: "sha1", hash: sha1.New},
	}
	key := runtime.GOOS
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			container := NewContainer(3, ContainerWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = container.ByKey(key)
			}
		})
	}
}

func BenchmarkContainer_ByVirtualShard(b *testing.B) {
	benchmarks := []struct {
		name string
		hash func() hash.Hash
	}{
		{name: "md5", hash: md5.New},
		{name: "sha1", hash: sha1.New},
	}
	shard := uint64(math.MaxUint64)
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			container := NewContainer(3, ContainerWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = container.ByVirtualShard(shard)
			}
		})
	}
}
