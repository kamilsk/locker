package internal

import (
	"encoding/hex"
	"math/big"
)

func ShardNumberNaive(in []byte, size uint64) uint64 {
	base, success := big.NewInt(0).SetString(hex.EncodeToString(in), 16)
	if !success {
		panic("invalid input")
	}
	return big.NewInt(0).Mod(base, big.NewInt(0).SetUint64(size)).Uint64()
}

func ShardNumberSimple(in []byte, size uint64) uint64 {
	return big.NewInt(0).Mod(big.NewInt(0).SetBytes(in), big.NewInt(0).SetUint64(size)).Uint64()
}

func ShardNumberFast(in []byte, size uint64) uint64 {
	shard := uint64(0)
	for f, i := uint64(1), len(in)-1; i >= 0; i-- {
		curByte := in[i]
		for j := 0; j < 8; j++ {
			// we iterate over bits to use sum instead of multiplication for the factor:
			// on each iteration: factor = factor * 2 % size <=> (factor + factor) % size
			bit := curByte % 2
			curByte /= 2
			if bit == 1 {
				shard = sumModWithoutOverflow(shard, f, size)
			}
			f = sumModWithoutOverflow(f, f, size)
		}
	}

	return shard
}

// sumModWithoutOverflow – overflow-safe (a + b) % d operation
func sumModWithoutOverflow(a, b, d uint64) uint64 {
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
