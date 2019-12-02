package locker_test

import (
	"github.com/kamilsk/locker"
	"testing"
)

func TestBytesToUint64Mod(t *testing.T) {
	tests := []string{
		"",
		"a",
		"testtesttesttesttesttesttesttesttesttesttesttesttesttesttest",
	}

	for _, test := range tests {
		t.Run("test", func(t *testing.T) {
			a := locker.OldBytesToUint64Mod([]byte(test), 100500)
			b := locker.BytesToUint64Mod([]byte(test), 100500)
			if a != b {
				t.Error(test, a, b)
			}
		})
	}
}

// BenchmarkBytesToUint64Mod/using_bigint-12   	 3000000	       585 ns/op	     216 B/op	       8 allocs/op
// BenchmarkBytesToUint64Mod/without_bigint-12 	10000000	       215 ns/op	       0 B/op	       0 allocs/op
func BenchmarkBytesToUint64Mod(b *testing.B) {
	benchmarks := []struct {
		name      string
		converter func(in []byte, divisor uint64) uint64
	}{
		{
			name:      "using bigint",
			converter: locker.OldBytesToUint64Mod,
		},
		{
			name:      "without bigint",
			converter: locker.BytesToUint64Mod,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			str := []byte("testtesttesttest")

			for i := 0; i < b.N; i++ {
				_ = bm.converter(str, 100500)
			}
		})
	}
}
