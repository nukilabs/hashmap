package stringhasher

import "github.com/nukilabs/hashmap/internal/rapidhash"

const FlagCount = 8 // Save 8 bits to be used as flags

func ComputeHashAndMaskTop8Bits(data []byte, seed uint64) uint32 {
	return MaskTop8Bits(rapidhash.Hash(data, seed))
}

// MaskTop8Bits matches Chromium's StringHasher::MaskTop8Bits
// Keeps only the bottom 24 bits and avoids returning 0
func MaskTop8Bits(result uint64) uint32 {
	// Reserving space from the high bits for flags preserves most of the hash's
	// value, since hash lookup typically masks out the high bits anyway.
	result &= (1 << (32 - FlagCount)) - 1

	// This avoids ever returning a hash code of 0, since that is used to
	// signal "hash not computed yet". Setting the high bit maintains
	// reasonable fidelity to a hash code of 0 because it is likely to yield
	// exactly 0 when hash lookup masks out the high bits.
	if result == 0 {
		result = 0x80000000 >> FlagCount
	}

	return uint32(result)
}
