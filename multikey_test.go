package locker_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/fnv"
	"runtime"
	"testing"

	. "github.com/kamilsk/locker"
)

func TestMultiMutex(t *testing.T) {
	set := MultiMutex(3, md5.New())
	set.ByKey("test").Lock()
	set.ByKey("another").Lock()
	defer set.ByKey("test").Unlock()
	defer set.ByKey("another").Unlock()
	go func() { set.ByKey("test").Unlock() }()
	go func() { set.ByKey("another").Unlock() }()
	set.ByKey("test").Lock()
	set.ByKey("another").Lock()
}

// BenchmarkMultiMutex/md5-4         	  340592	      3822 ns/op	     712 B/op	      29 allocs/op
// BenchmarkMultiMutex/sha1-4        	  296442	      4391 ns/op	     912 B/op	      30 allocs/op
// BenchmarkMultiMutex/sha256-4      	  205178	      5874 ns/op	    1000 B/op	      29 allocs/op
// BenchmarkMultiMutex/sha512-4      	  131266	     10442 ns/op	    1864 B/op	      32 allocs/op
// BenchmarkMultiMutex/sum32-4       	  626733	      1696 ns/op	     288 B/op	      27 allocs/op
func BenchmarkMultiMutex(b *testing.B) {
	benchmarks := []struct {
		name string
		hash hash.Hash
	}{
		{name: "md5", hash: md5.New()},
		{name: "sha1", hash: sha1.New()},
		{name: "sha256", hash: sha256.New()},
		{name: "sha512", hash: sha512.New()},
		{name: "sum32", hash: fnv.New32()},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()

			keys := []string{runtime.GOOS, runtime.GOARCH, runtime.Version()}
			set := MultiMutex(10, bm.hash)
			for i := 0; i < b.N; i++ {
				for _, key := range keys {
					_ = set.ByKey(key)
				}
			}
		})
	}
}
