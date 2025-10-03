package rapidhash

import (
	"encoding/binary"
	"math/bits"
)

const SEED uint64 = 0xbdd89aa982704029

var secret = [3]uint64{
	0x2d358dccaa6c78a5,
	0x8bb84b93962eacc9,
	0x4b33a62ed433d4a3,
}

func Mul128(a, b uint64) (uint64, uint64) {
	hi, lo := bits.Mul64(a, b)
	return lo, hi
}

func Mix(a, b uint64) uint64 {
	lo, hi := Mul128(a, b)
	return lo ^ hi
}

// Hash implements the RapidHash algorithm
func Hash(data []byte, seed uint64) uint64 {
	length := uint64(len(data))
	seed ^= Mix(seed^secret[0], secret[1]) ^ length

	var a, b uint64

	if length <= 16 {
		if length >= 4 {
			a = (uint64(binary.LittleEndian.Uint32(data[0:4])) << 32) |
				uint64(binary.LittleEndian.Uint32(data[length-4:]))

			delta := ((length & 24) >> (length >> 3))
			b = (uint64(binary.LittleEndian.Uint32(data[delta:delta+4])) << 32) |
				uint64(binary.LittleEndian.Uint32(data[length-4-delta:length-delta]))
		} else if length > 0 {
			k := length
			a = (uint64(data[0]) << 56) | (uint64(data[k>>1]) << 32) | uint64(data[k-1])
			b = 0
		} else {
			a = 0
			b = 0
		}
	} else {
		i := length
		if i > 48 {
			see1 := seed
			see2 := seed

			for i >= 48 {
				p := data[length-i:]
				seed = Mix(binary.LittleEndian.Uint64(p[0:8])^secret[0],
					binary.LittleEndian.Uint64(p[8:16])^seed)
				see1 = Mix(binary.LittleEndian.Uint64(p[16:24])^secret[1],
					binary.LittleEndian.Uint64(p[24:32])^see1)
				see2 = Mix(binary.LittleEndian.Uint64(p[32:40])^secret[2],
					binary.LittleEndian.Uint64(p[40:48])^see2)
				i -= 48
			}
			seed ^= see1 ^ see2
		}

		if i > 16 {
			p := data[length-i:]
			seed = Mix(binary.LittleEndian.Uint64(p[0:8])^secret[2],
				binary.LittleEndian.Uint64(p[8:16])^seed^secret[1])
			if i > 32 {
				seed = Mix(binary.LittleEndian.Uint64(p[16:24])^secret[2],
					binary.LittleEndian.Uint64(p[24:32])^seed)
			}
		}

		a = binary.LittleEndian.Uint64(data[length-16 : length-8])
		b = binary.LittleEndian.Uint64(data[length-8:])
	}

	a ^= secret[1]
	b ^= seed
	lo, hi := Mul128(a, b)
	return Mix(lo^secret[0]^length, hi^secret[1])
}
