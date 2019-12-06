package internal_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/fnv"
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"

	. "github.com/kamilsk/locker/internal"
)

func TestShardNumberCalculation(t *testing.T) {
	hashes := map[string]hash.Hash{
		"md5":     md5.New(),
		"sha1":    sha1.New(),
		"sha256":  sha256.New(),
		"sha512":  sha512.New(),
		"sum32":   fnv.New32(),
		"sum32a":  fnv.New32a(),
		"sum64":   fnv.New64(),
		"sum64a":  fnv.New64a(),
		"sum128":  fnv.New128(),
		"sum128a": fnv.New128a(),
	}
	keys := [...][]byte{[]byte(runtime.GOOS), []byte(runtime.GOARCH)}
	sizes := []uint64{
		math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64,
		math.MaxUint8 * (1 + math.MaxUint8), rand.New(rand.NewSource(time.Now().UnixNano())).Uint64(),
	}
	for name, algorithm := range hashes {
		t.Run(name, func(t *testing.T) {
			for _, key := range keys {
				algorithm.Write(key)
				checksum := algorithm.Sum(nil)
				algorithm.Reset()

				for _, size := range sizes {
					x, y := ShardNumberNaive(checksum, size), ShardNumberSimple(checksum, size)
					if x != y {
						t.Errorf("%d != %d", x, y)
						t.FailNow()
					}

					z := ShardNumberFast(checksum, size)
					if x != z {
						t.Errorf("%d != %d", x, z)
						t.FailNow()
					}
				}
			}
		})
	}
}

// BenchmarkShardNumberCalculation/naive:md5-4         	 1000000	      1011 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive:sha1-4        	 1000000	      1166 ns/op	     264 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive:sha256-4      	  583350	      1770 ns/op	     296 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive:sha512-4      	  382222	      2687 ns/op	     552 B/op	       9 allocs/op
// BenchmarkShardNumberCalculation/naive:sum32-4       	 2314318	       463 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive:sum32a-4      	 2239576	       466 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive:sum64-4       	 1904371	       616 ns/op	     136 B/op	       7 allocs/op
// BenchmarkShardNumberCalculation/naive:sum64a-4      	 1798017	       662 ns/op	     136 B/op	       7 allocs/op
// BenchmarkShardNumberCalculation/naive:sum128-4      	 1000000	      1076 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive:sum128a-4     	 1219659	       911 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/simple:md5-4        	 5383893	       223 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple:sha1-4       	 4214607	       283 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple:sha256-4     	 4009988	       287 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple:sha512-4     	 2707340	       417 ns/op	     208 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple:sum32-4      	 8844078	       129 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple:sum32a-4     	 8814894	       132 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple:sum64-4      	 9100696	       128 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple:sum64a-4     	 8227677	       130 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple:sum128-4     	 4863288	       227 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple:sum128a-4    	 5061507	       225 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/fast:md5-4          	 3878890	       304 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sha1-4         	 3090046	       381 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sha256-4       	 1827475	       651 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sha512-4       	  839475	      1368 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sum32-4        	18319796	        63.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sum32a-4       	19422997	        60.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sum64-4        	 9456823	       125 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sum64a-4       	 9441171	       128 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sum128-4       	 3846472	       310 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast:sum128a-4      	 3861493	       308 ns/op	       0 B/op	       0 allocs/op
func BenchmarkShardNumberCalculation(b *testing.B) {
	benchmarks := []struct {
		name           string
		implementation func([]byte, uint64) uint64
	}{
		{"naive", ShardNumberNaive},
		{"simple", ShardNumberSimple},
		{"fast", ShardNumberFast},
	}
	hashes := []struct {
		name string
		hash hash.Hash
	}{
		{"md5", md5.New()},
		{"sha1", sha1.New()},
		{"sha256", sha256.New()},
		{"sha512", sha512.New()},
		{"sum32", fnv.New32()},
		{"sum32a", fnv.New32a()},
		{"sum64", fnv.New64()},
		{"sum64a", fnv.New64a()},
		{"sum128", fnv.New128()},
		{"sum128a", fnv.New128a()},
	}
	key, size := []byte(runtime.Version()), rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	for _, benchmark := range benchmarks {
		for _, algorithm := range hashes {
			b.Run(benchmark.name+":"+algorithm.name, func(b *testing.B) {
				algorithm.hash.Write(key)
				checksum := algorithm.hash.Sum(nil)
				algorithm.hash.Reset()

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = benchmark.implementation(checksum, size)
				}
			})
		}
	}
}
