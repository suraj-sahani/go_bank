package main

// feistelKey is the fixed key used to scramble sequential account numbers
// into unique, non-sequential values via a Feistel cipher.
// This is for obfuscation (hiding sequential patterns), not cryptographic security.
const feistelKey uint32 = 0xDEADBEEF

// GenerateAccountNumber takes a sequential value from a database sequence
// and returns a unique, scrambled account number using a Feistel cipher.
//
// The Feistel cipher is a bijection — every unique input maps to a unique
// output — so uniqueness is guaranteed by the underlying database sequence.
func GenerateAccountNumber(seqVal uint32) int64 {
	return int64(feistelEncrypt(seqVal, feistelKey))
}

// feistelEncrypt applies a 4-round Feistel network to a 32-bit value.
// It splits the input into two 16-bit halves and alternately applies
// a keyed round function, producing a permutation of the full uint32 space.
func feistelEncrypt(value, key uint32) uint32 {
	left := (value >> 16) & 0xFFFF
	right := value & 0xFFFF

	for i := uint32(0); i < 4; i++ {
		// Derive a unique sub-key per round using the golden ratio constant
		roundKey := key ^ (i * 0x9E3779B9)
		f := feistelRound(right, roundKey)
		left, right = right, left^f
	}

	return (left << 16) | right
}

// feistelRound is the mixing function for one Feistel round.
// It takes a 16-bit half-block and a round key, and produces
// a 16-bit pseudo-random output using MurmurHash3-style mixing.
func feistelRound(value, key uint32) uint32 {
	v := value * 0xCC9E2D51
	v = (v << 15) | (v >> 17) // rotate left 15
	v *= 0x1B873593
	v ^= key
	v = (v << 13) | (v >> 19) // rotate left 13
	v = v*5 + 0xE6546B64
	return v & 0xFFFF // mask to 16-bit half-block width
}

