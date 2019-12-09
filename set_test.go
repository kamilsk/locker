package locker_test

import (
	"crypto/md5"
	"crypto/sha1"
	"flag"
	"hash"
	"math"
	"runtime"
	"testing"

	. "github.com/kamilsk/locker"
)

var stress = flag.Bool("stress-test", false, "run stress tests")

func TestSet(t *testing.T) {
	t.Run("with custom hash option", func(t *testing.T) {
		container := Set(3, SetWithHash(sha1.New))
		if container.ByKey(runtime.GOOS) == container.ByKey(runtime.GOARCH) {
			t.Error("unexpected result")
			t.FailNow()
		}
	})
	t.Run("with custom mapping option", func(t *testing.T) {
		container := Set(3, SetWithMapping(func([]byte, uint64) uint64 { return 0 }))
		if container.ByKey(runtime.GOOS) != container.ByKey(runtime.GOARCH) {
			t.Error("unexpected result")
			t.FailNow()
		}
	})
}

func TestMutexSet_ByFingerprint(t *testing.T) {
	fingerprints := [...][]byte{[]byte(runtime.GOOS), []byte(runtime.GOARCH)}
	set := Set(3)

	for _, fingerprint := range fingerprints {
		origin := set.ByFingerprint(fingerprint)
		for range make([]struct{}, 1000) {
			current := set.ByFingerprint(fingerprint)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, fingerprint := range fingerprints {
		current := set.ByFingerprint(fingerprint)
		for _, fingerprint := range fingerprints[i+1:] {
			next := set.ByFingerprint(fingerprint)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestMutexSet_ByKey(t *testing.T) {
	keys := [...]string{runtime.GOOS, runtime.GOARCH}
	set := Set(3)

	for _, key := range keys {
		origin := set.ByKey(key)
		for range make([]struct{}, 1000) {
			current := set.ByKey(key)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, key := range keys {
		current := set.ByKey(key)
		for _, key := range keys[i+1:] {
			next := set.ByKey(key)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestMutexSet_ByVirtualShard(t *testing.T) {
	shards := [...]uint64{1, 5, 9}
	set := Set(3)

	for _, shard := range shards {
		origin := set.ByVirtualShard(shard)
		for range make([]struct{}, 1000) {
			current := set.ByVirtualShard(shard)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, shard := range shards {
		current := set.ByVirtualShard(shard)
		for _, shard := range shards[i+1:] {
			next := set.ByVirtualShard(shard)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestMutexSet_StressTest(t *testing.T) {
	if !*stress {
		t.SkipNow()
	}
	for range make([]struct{}, 1000) {
		TestMutexSet_ByFingerprint(t)
		TestMutexSet_ByKey(t)
		TestMutexSet_ByVirtualShard(t)
	}
}

// BenchmarkMutexSet_ByFingerprint/md5-4         	 2167080	       572 ns/op	     112 B/op	       2 allocs/op
// BenchmarkMutexSet_ByFingerprint/sha1-4        	 2174104	       558 ns/op	     112 B/op	       2 allocs/op
func BenchmarkMutexSet_ByFingerprint(b *testing.B) {
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
			set := Set(3, SetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByFingerprint(fingerprint)
			}
		})
	}
}

// BenchmarkMutexSet_ByKey/md5-4         	 1968254	       596 ns/op	     120 B/op	       3 allocs/op
// BenchmarkMutexSet_ByKey/sha1-4        	 2027947	       593 ns/op	     120 B/op	       3 allocs/op
func BenchmarkMutexSet_ByKey(b *testing.B) {
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
			set := Set(3, SetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByKey(key)
			}
		})
	}
}

// BenchmarkMutexSet_ByVirtualShard/md5-4         	100000000	        10.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkMutexSet_ByVirtualShard/sha1-4        	100000000	        10.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMutexSet_ByVirtualShard(b *testing.B) {
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
			set := Set(3, SetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByVirtualShard(shard)
			}
		})
	}
}

func TestRWSet(t *testing.T) {
	t.Run("with custom hash option", func(t *testing.T) {
		container := RWSet(3, RWSetWithHash(sha1.New))
		if container.ByKey(runtime.GOOS) == container.ByKey(runtime.GOARCH) {
			t.Error("unexpected result")
			t.FailNow()
		}
	})
	t.Run("with custom mapping option", func(t *testing.T) {
		container := RWSet(3, RWSetWithMapping(func([]byte, uint64) uint64 { return 0 }))
		if container.ByKey(runtime.GOOS) != container.ByKey(runtime.GOARCH) {
			t.Error("unexpected result")
			t.FailNow()
		}
	})
}

func TestRWMutexSet_ByFingerprint(t *testing.T) {
	fingerprints := [...][]byte{[]byte(runtime.GOOS), []byte(runtime.GOARCH)}
	set := RWSet(3)

	for _, fingerprint := range fingerprints {
		origin := set.ByFingerprint(fingerprint)
		for range make([]struct{}, 1000) {
			current := set.ByFingerprint(fingerprint)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, fingerprint := range fingerprints {
		current := set.ByFingerprint(fingerprint)
		for _, fingerprint := range fingerprints[i+1:] {
			next := set.ByFingerprint(fingerprint)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestRWMutexSet_ByKey(t *testing.T) {
	keys := [...]string{runtime.GOOS, runtime.GOARCH}
	set := RWSet(3)

	for _, key := range keys {
		origin := set.ByKey(key)
		for range make([]struct{}, 1000) {
			current := set.ByKey(key)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, key := range keys {
		current := set.ByKey(key)
		for _, key := range keys[i+1:] {
			next := set.ByKey(key)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestRWMutexSet_ByVirtualShard(t *testing.T) {
	shards := [...]uint64{1, 5, 9}
	set := RWSet(3)

	for _, shard := range shards {
		origin := set.ByVirtualShard(shard)
		for range make([]struct{}, 1000) {
			current := set.ByVirtualShard(shard)
			if origin != current {
				t.Error("non-deterministic result")
				t.FailNow()
			}
		}
	}

	for i, shard := range shards {
		current := set.ByVirtualShard(shard)
		for _, shard := range shards[i+1:] {
			next := set.ByVirtualShard(shard)
			if current == next {
				t.Error("has deadlock")
				t.FailNow()
			}
		}
	}
}

func TestRWMutexSet_StressTest(t *testing.T) {
	if !*stress {
		t.SkipNow()
	}
	for range make([]struct{}, 1000) {
		TestRWMutexSet_ByFingerprint(t)
		TestRWMutexSet_ByKey(t)
		TestRWMutexSet_ByVirtualShard(t)
	}
}
