package locker

import (
	"encoding/hex"
	"math/big"
)

func OldBytesToUint64Mod(in []byte, divisor uint64) uint64 {
	base := big.NewInt(0)
	_, _ = base.SetString(hex.EncodeToString(in), 16)
	shard := big.NewInt(0).Mod(base, big.NewInt(int64(divisor)))
	return shard.Uint64()

}

func BytesToUint64Mod(in []byte, divisor uint64) uint64 {
	res := uint64(0)
	multiplier := uint64(1)
	for i := len(in) - 1; i >= 0; i-- {
		res += uint64(in[i]) * multiplier
		res %= divisor
		multiplier *= 256
		multiplier %= divisor
	}

	return res
}
