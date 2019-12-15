package internal

import (
	"encoding/hex"
	"math/big"
)

// ShardNumberNaive calculates a shard number by a hash/fingerprint of a custom payload
// and a size of an available mapping set.
//
// This is a naive implementation to control the correctness of the simple implementation.
func ShardNumberNaive(checksum []byte, size uint64) uint64 {
	base, _ := big.NewInt(0).SetString(hex.EncodeToString(checksum), 16)
	return big.NewInt(0).Mod(base, big.NewInt(0).SetUint64(size)).Uint64()
}

// ShardNumberNaive calculates a shard number by a hash/fingerprint of a custom payload
// and a size of an available mapping set.
//
// This is a simple implementation to control the correctness of the fast implementation.
func ShardNumberSimple(checksum []byte, size uint64) uint64 {
	return big.NewInt(0).Mod(big.NewInt(0).SetBytes(checksum), big.NewInt(0).SetUint64(size)).Uint64()
}

// ShardNumberNaive calculates a shard number by a hash/fingerprint of a custom payload
// and a size of an available mapping set.
//
// This is a fast implementation without disadvantages of the simple implementation.
func ShardNumberFast(checksum []byte, size uint64) uint64 {
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
	}
	// 3. a + b == d
	return 0
}
