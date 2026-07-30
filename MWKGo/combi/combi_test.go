package main

import "testing"

// TestEnumerateCombinations_MatchesCombiTXTWorkedExample reproduces
// COMBI.TXT's own worked example (6 things taken 4 at a time) exactly,
// including the generation order: 15 combinations, ending on ABCD (the
// starting pattern), which the revolving-door algorithm only reaches
// again after cycling through every other combination.
func TestEnumerateCombinations_MatchesCombiTXTWorkedExample(t *testing.T) {
	want := []string{
		"ABCE", "ABCF", "ACEF", "ACDE", "ACDF", "ADEF", "CDEF", "BCEF",
		"BCDE", "BCDF", "BDEF", "ABEF", "ABDE", "ABDF", "ABCD",
	}
	got := enumerateCombinations(6, 4)
	if len(got) != len(want) {
		t.Fatalf("got %d combinations, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("combination %d = %q, want %q", i+1, got[i], want[i])
		}
	}
}

func TestBinomial_MatchesKnownValues(t *testing.T) {
	cases := []struct {
		n, m int
		want int64
	}{
		{6, 4, 15},
		{10, 3, 120},
		{5, 0, 1},
	}
	for _, c := range cases {
		if got := binomial(c.n, c.m); got != c.want {
			t.Errorf("binomial(%d,%d) = %d, want %d", c.n, c.m, got, c.want)
		}
	}
}

// TestEnumerateCombinations_EveryEntryIsDistinctAndCorrectLength
// confirms the enumeration is a genuine listing of unique
// combinations (no repeats, no combination of the wrong size) for a
// case with no independently published worked example.
func TestEnumerateCombinations_EveryEntryIsDistinctAndCorrectLength(t *testing.T) {
	n, m := 8, 3
	combos := enumerateCombinations(n, m)
	want := binomial(n, m)
	if int64(len(combos)) != want {
		t.Fatalf("got %d combinations, want %d", len(combos), want)
	}

	seen := make(map[string]bool, len(combos))
	for _, c := range combos {
		if len(c) != m {
			t.Errorf("combination %q has length %d, want %d", c, len(c), m)
		}
		if seen[c] {
			t.Errorf("combination %q appeared more than once", c)
		}
		seen[c] = true
	}
}

func TestEnumerateCombinations_SmallestCase(t *testing.T) {
	// 2 things taken 1 at a time has exactly 2 combinations: "A" and "B".
	got := enumerateCombinations(2, 1)
	want := []string{"B", "A"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("enumerateCombinations(2,1) = %v, want %v", got, want)
	}
}
