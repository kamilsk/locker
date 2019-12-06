package internal_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/fnv"
	"math"
	"runtime"
	"testing"

	. "github.com/kamilsk/locker/internal"
)

var keys = [...]string{runtime.GOOS, runtime.GOARCH}

func TestShardNumberCalculation(t *testing.T) {
	sources := map[string]hash.Hash{
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
	size := uint64(math.MaxUint64)
	for name, checksum := range sources {
		t.Run(name, func(t *testing.T) {
			for _, key := range keys {
				checksum.Write([]byte(key))
				in := checksum.Sum(nil)
				checksum.Reset()

				x, y := ShardNumberNaive(in, size), ShardNumberSimple(in, size)
				if x != y {
					t.Errorf("%d != %d", x, y)
					t.FailNow()
				}

				z := ShardNumberFast(in, size)
				if x != z {
					t.Errorf("%d != %d", x, z)
					t.FailNow()
				}
			}
		})
	}
}

// BenchmarkShardNumberCalculation/naive,md5:darwin-12         	 3000000	       575 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,md5:amd64-12          	 3000000	       560 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha1:darwin-12        	 2000000	       644 ns/op	     264 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha1:amd64-12         	 2000000	       647 ns/op	     264 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha256:darwin-12      	 2000000	       868 ns/op	     296 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha256:amd64-12       	 2000000	       849 ns/op	     296 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha512:darwin-12      	 1000000	      1512 ns/op	     552 B/op	       9 allocs/op
// BenchmarkShardNumberCalculation/naive,sha512:amd64-12       	 1000000	      1473 ns/op	     552 B/op	       9 allocs/op
// BenchmarkShardNumberCalculation/naive,sum32:darwin-12       	 5000000	       264 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive,sum32:amd64-12        	 5000000	       265 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive,sum32a:darwin-12      	 5000000	       263 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive,sum32a:amd64-12       	 5000000	       269 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive,sum64:darwin-12       	 5000000	       369 ns/op	     136 B/op	       7 allocs/op
// BenchmarkShardNumberCalculation/naive,sum64:amd64-12        	 5000000	       370 ns/op	     136 B/op	       7 allocs/op
// BenchmarkShardNumberCalculation/naive,sum64a:darwin-12      	 3000000	       368 ns/op	     136 B/op	       7 allocs/op
// BenchmarkShardNumberCalculation/naive,sum64a:amd64-12       	 5000000	       373 ns/op	     136 B/op	       7 allocs/op
// BenchmarkShardNumberCalculation/naive,sum128:darwin-12      	 3000000	       560 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum128:amd64-12       	 3000000	       564 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum128a:darwin-12     	 3000000	       557 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum128a:amd64-12      	 3000000	       586 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/simple,md5:darwin-12        	10000000	       174 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,md5:amd64-12         	10000000	       166 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha1:darwin-12       	10000000	       196 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha1:amd64-12        	10000000	       203 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha256:darwin-12     	10000000	       222 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha256:amd64-12      	10000000	       233 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha512:darwin-12     	 5000000	       346 ns/op	     208 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha512:amd64-12      	 5000000	       354 ns/op	     208 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum32:darwin-12      	20000000	        97.8 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum32:amd64-12       	20000000	        95.9 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum32a:darwin-12     	20000000	        93.0 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum32a:amd64-12      	20000000	        95.1 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum64:darwin-12      	20000000	        87.2 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum64:amd64-12       	20000000	        86.5 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum64a:darwin-12     	20000000	        86.3 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum64a:amd64-12      	20000000	        86.2 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum128:darwin-12     	10000000	       157 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum128:amd64-12      	10000000	       158 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum128a:darwin-12    	10000000	       157 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum128a:amd64-12     	10000000	       157 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/fast,md5:darwin-12          	10000000	       177 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,md5:amd64-12           	10000000	       175 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha1:darwin-12         	10000000	       217 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha1:amd64-12          	10000000	       221 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha256:darwin-12       	 5000000	       348 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha256:amd64-12        	 5000000	       346 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha512:darwin-12       	 2000000	       701 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha512:amd64-12        	 2000000	       701 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum32:darwin-12        	30000000	        45.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum32:amd64-12         	30000000	        45.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum32a:darwin-12       	30000000	        44.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum32a:amd64-12        	30000000	        45.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum64:darwin-12        	20000000	        89.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum64:amd64-12         	20000000	        91.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum64a:darwin-12       	20000000	        88.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum64a:amd64-12        	20000000	        88.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum128:darwin-12       	10000000	       174 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum128:amd64-12        	10000000	       174 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum128a:darwin-12      	10000000	       175 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum128a:amd64-12       	10000000	       174 ns/op	       0 B/op	       0 allocs/op
func BenchmarkShardNumberCalculation(b *testing.B) {
	hashes := []struct {
		name     string
		checksum hash.Hash
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
	benchmarks := []struct {
		name           string
		implementation func([]byte, uint64) uint64
	}{
		{"naive", ShardNumberNaive},
		{"simple", ShardNumberSimple},
		{"fast", ShardNumberFast},
	}
	size := uint64(math.MaxUint64)
	for _, benchmark := range benchmarks {
		for _, algorithm := range hashes {
			for _, key := range keys {
				b.Run(benchmark.name+","+algorithm.name+":"+key, func(b *testing.B) {
					algorithm.checksum.Write([]byte(key))
					in := algorithm.checksum.Sum(nil)
					algorithm.checksum.Reset()

					b.ResetTimer()
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						_ = benchmark.implementation(in, size)
					}
				})
			}
		}
	}
}
