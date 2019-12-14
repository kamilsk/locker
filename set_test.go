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

func TestInterruptibleSet(t *testing.T) {
	t.Run("with custom hash option", func(t *testing.T) {
		container := InterruptibleSet(3, InterruptibleSetWithHash(sha1.New))
		if container.ByKey(runtime.GOOS) == container.ByKey(runtime.GOARCH) {
			t.Error("unexpected result")
			t.FailNow()
		}
	})
	t.Run("with custom mapping option", func(t *testing.T) {
		container := InterruptibleSet(3, InterruptibleSetWithMapping(func([]byte, uint64) uint64 { return 0 }))
		if container.ByKey(runtime.GOOS) != container.ByKey(runtime.GOARCH) {
			t.Error("unexpected result")
			t.FailNow()
		}
	})
}

func TestInterruptibleSet_ByFingerprint(t *testing.T) {
	fingerprints := [...][]byte{[]byte(runtime.GOOS), []byte(runtime.GOARCH)}
	set := InterruptibleSet(3)

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

func TestInterruptibleSet_ByKey(t *testing.T) {
	keys := [...]string{runtime.GOOS, runtime.GOARCH}
	set := InterruptibleSet(3)

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

func TestInterruptibleSet_ByVirtualShard(t *testing.T) {
	shards := [...]uint64{1, 5, 9}
	set := InterruptibleSet(3)

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

func TestInterruptibleSet_StressTest(t *testing.T) {
	if !*stress {
		t.SkipNow()
	}
	for range make([]struct{}, 1000) {
		TestInterruptibleSet_ByFingerprint(t)
		TestInterruptibleSet_ByKey(t)
		TestInterruptibleSet_ByVirtualShard(t)
	}
}

// BenchmarkInterruptibleSet_ByFingerprint/md5-4         	 2101126	       587 ns/op	     112 B/op	       2 allocs/op
// BenchmarkInterruptibleSet_ByFingerprint/sha1-4        	 1703542	       741 ns/op	     144 B/op	       2 allocs/op
func BenchmarkInterruptibleSet_ByFingerprint(b *testing.B) {
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
			set := InterruptibleSet(3, InterruptibleSetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByFingerprint(fingerprint)
			}
		})
	}
}

// BenchmarkInterruptibleSet_ByKey/md5-4         	 1851483	       678 ns/op	     120 B/op	       3 allocs/op
// BenchmarkInterruptibleSet_ByKey/sha1-4        	 1605897	       753 ns/op	     152 B/op	       3 allocs/op
func BenchmarkInterruptibleSet_ByKey(b *testing.B) {
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
			set := InterruptibleSet(3, InterruptibleSetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByKey(key)
			}
		})
	}
}

// BenchmarkInterruptibleSet_ByVirtualShard/md5-4         	123889586	         9.62 ns/op	       0 B/op	       0 allocs/op
// BenchmarkInterruptibleSet_ByVirtualShard/sha1-4        	123597364	         9.57 ns/op	       0 B/op	       0 allocs/op
func BenchmarkInterruptibleSet_ByVirtualShard(b *testing.B) {
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
			set := InterruptibleSet(3, InterruptibleSetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByVirtualShard(shard)
			}
		})
	}
}

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

func TestSet_ByFingerprint(t *testing.T) {
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

func TestSet_ByKey(t *testing.T) {
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

func TestSet_ByVirtualShard(t *testing.T) {
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

func TestSet_StressTest(t *testing.T) {
	if !*stress {
		t.SkipNow()
	}
	for range make([]struct{}, 1000) {
		TestSet_ByFingerprint(t)
		TestSet_ByKey(t)
		TestSet_ByVirtualShard(t)
	}
}

// BenchmarkSet_ByFingerprint/md5-4         	 2167080	       572 ns/op	     112 B/op	       2 allocs/op
// BenchmarkSet_ByFingerprint/sha1-4        	 2174104	       558 ns/op	     112 B/op	       2 allocs/op
func BenchmarkSet_ByFingerprint(b *testing.B) {
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

// BenchmarkSet_ByKey/md5-4         	 1968254	       596 ns/op	     120 B/op	       3 allocs/op
// BenchmarkSet_ByKey/sha1-4        	 2027947	       593 ns/op	     120 B/op	       3 allocs/op
func BenchmarkSet_ByKey(b *testing.B) {
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

// BenchmarkSet_ByVirtualShard/md5-4         	100000000	        10.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSet_ByVirtualShard/sha1-4        	100000000	        10.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSet_ByVirtualShard(b *testing.B) {
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

func TestRWSet_ByFingerprint(t *testing.T) {
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

func TestRWSet_ByKey(t *testing.T) {
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

func TestRWSet_ByVirtualShard(t *testing.T) {
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

func TestRWSet_StressTest(t *testing.T) {
	if !*stress {
		t.SkipNow()
	}
	for range make([]struct{}, 1000) {
		TestRWSet_ByFingerprint(t)
		TestRWSet_ByKey(t)
		TestRWSet_ByVirtualShard(t)
	}
}

// BenchmarkRWSet_ByFingerprint/md5-4         	 2084050	       587 ns/op	     112 B/op	       2 allocs/op
// BenchmarkRWSet_ByFingerprint/sha1-4        	 1712743	       706 ns/op	     144 B/op	       2 allocs/op
func BenchmarkRWSet_ByFingerprint(b *testing.B) {
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
			set := RWSet(3, RWSetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByFingerprint(fingerprint)
			}
		})
	}
}

// BenchmarkRWSet_ByKey/md5-4         	 1932391	       643 ns/op	     120 B/op	       3 allocs/op
// BenchmarkRWSet_ByKey/sha1-4        	 1603132	       792 ns/op	     152 B/op	       3 allocs/op
func BenchmarkRWSet_ByKey(b *testing.B) {
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
			set := RWSet(3, RWSetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByKey(key)
			}
		})
	}
}

// BenchmarkRWSet_ByVirtualShard/md5-4         	123782726	         9.68 ns/op	       0 B/op	       0 allocs/op
// BenchmarkRWSet_ByVirtualShard/sha1-4        	123174104	         9.71 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRWSet_ByVirtualShard(b *testing.B) {
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
			set := RWSet(3, RWSetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByVirtualShard(shard)
			}
		})
	}
}
