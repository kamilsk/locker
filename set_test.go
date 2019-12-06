package locker_test

import (
	"crypto/md5"
	"crypto/sha1"
	"hash"
	"hash/fnv"
	"math"
	"runtime"
	"testing"
	"time"

	. "github.com/kamilsk/locker"
)

func TestMutexSet_ByKey(t *testing.T) {
	keys := [...]string{runtime.GOOS, runtime.GOARCH}
	set := Set(3)
	for _, key := range keys {
		set.ByKey(key).Lock()
		go func(key string) {
			time.Sleep(time.Millisecond)
			set.ByKey(key).Unlock()
		}(key)
	}
	for _, key := range keys {
		set.ByKey(key).Lock()
	}
}

// BenchmarkMutexSet_ByFingerprint/md5-4         	 1893616	       612 ns/op	      16 B/op	       1 allocs/op
// BenchmarkMutexSet_ByFingerprint/sha1-4        	 2058212	       579 ns/op	      16 B/op	       1 allocs/op
// BenchmarkMutexSet_ByFingerprint/sum32-4       	 1978899	       614 ns/op	      16 B/op	       1 allocs/op
func BenchmarkMutexSet_ByFingerprint(b *testing.B) {
	benchmarks := []struct {
		name string
		hash hash.Hash
	}{
		{name: "md5", hash: md5.New()},
		{name: "sha1", hash: sha1.New()},
		{name: "sum32", hash: fnv.New32()},
	}
	fingerprint := []byte(runtime.GOOS)
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			set := Set(10, SetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByFingerprint(fingerprint)
			}
		})
	}
}

// BenchmarkMutexSet_ByKey/md5-4         	 1976034	       619 ns/op	      24 B/op	       2 allocs/op
// BenchmarkMutexSet_ByKey/sha1-4        	 1939647	       628 ns/op	      24 B/op	       2 allocs/op
// BenchmarkMutexSet_ByKey/sum32-4       	 1872367	       653 ns/op	      24 B/op	       2 allocs/op
func BenchmarkMutexSet_ByKey(b *testing.B) {
	benchmarks := []struct {
		name string
		hash hash.Hash
	}{
		{name: "md5", hash: md5.New()},
		{name: "sha1", hash: sha1.New()},
		{name: "sum32", hash: fnv.New32()},
	}
	key := runtime.GOOS
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			set := Set(10, SetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByKey(key)
			}
		})
	}
}

// BenchmarkMutexSet_ByVirtualShard/md5-4         	100000000	        10.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkMutexSet_ByVirtualShard/sha1-4        	100000000	        10.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkMutexSet_ByVirtualShard/sum32-4       	100000000	        10.6 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMutexSet_ByVirtualShard(b *testing.B) {
	benchmarks := []struct {
		name string
		hash hash.Hash
	}{
		{name: "md5", hash: md5.New()},
		{name: "sha1", hash: sha1.New()},
		{name: "sum32", hash: fnv.New32()},
	}
	shard := uint64(math.MaxUint64)
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			set := Set(10, SetWithHash(bm.hash))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.ByVirtualShard(shard)
			}
		})
	}
}
