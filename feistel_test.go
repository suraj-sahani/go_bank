package main

import "testing"

func TestFeistelUniqueness(t *testing.T) {
	seen := make(map[int64]uint32)

	// Simulate 10,000 sequential values from the DB sequence
	for i := uint32(1); i <= 10_000; i++ {
		accNum := GenerateAccountNumber(i)

		if prev, exists := seen[accNum]; exists {
			t.Fatalf("collision: sequence values %d and %d both produced account number %d", prev, i, accNum)
		}
		seen[accNum] = i
	}
}

func TestFeistelDeterministic(t *testing.T) {
	// Same input must always produce the same output
	for i := uint32(1); i <= 100; i++ {
		a := GenerateAccountNumber(i)
		b := GenerateAccountNumber(i)
		if a != b {
			t.Fatalf("non-deterministic: sequence %d produced %d then %d", i, a, b)
		}
	}
}

func TestFeistelNonSequential(t *testing.T) {
	// Consecutive inputs should NOT produce consecutive outputs
	prev := GenerateAccountNumber(1)
	sequential := 0
	for i := uint32(2); i <= 100; i++ {
		curr := GenerateAccountNumber(i)
		if curr == prev+1 {
			sequential++
		}
		prev = curr
	}

	// Allow at most 2 accidental consecutive pairs out of 99
	if sequential > 2 {
		t.Fatalf("outputs look too sequential: %d consecutive pairs out of 99", sequential)
	}
}

