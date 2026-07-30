// combi enumerates every combination of N things taken M at a time,
// printing each one as a string of letters rather than just counting
// them.
//
// Converted from COMBI.C (M. W. Klotz), Math/combi. The Tier 2
// preparatory scan found this program misclassified: it only ever
// opens an *output* report file, never reads a data file, so it
// belongs with the Tier 1 no-data-file calculators, not Tier 2; see
// ai/plans/c-to-go-conversion-plan.md's clean-up phase note. Like
// several other Tier 1 programs (ballcut, boltcirc, flywheel, ...),
// this drops the original's file-save feature and prints straight to
// stdout instead.
package main

import (
	"fmt"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

// maxN matches COMBI.C's own nmax/mmax: the number of letters
// available (26 upper-case plus 26 lower-case) to label each element
// of the enumeration.
const maxN = 52

// alphabet mirrors COMBI.C's own alpha[] array exactly, including its
// unused leading '.' placeholder at index 0 (element indices are
// 1-based throughout, matching the original).
const alphabet = ".ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// combinationState is the revolving-door combination generator's own
// mutable state: c[i] (1-indexed, index 0 unused) is 1 if element i
// is in the current combination.
type combinationState struct {
	n int
	c []int
}

// newCombinationState starts at the lexicographically-first
// combination (elements 1..m), matching COMBI.C's own initialization.
func newCombinationState(n, m int) *combinationState {
	c := make([]int, n+1)
	for i := 1; i <= m; i++ {
		c[i] = 1
	}
	return &combinationState{n: n, c: c}
}

// next advances to the next combination in place, matching COMBI.C's
// own combo() function exactly: the classic Nijenhuis & Wilf
// "revolving door" algorithm, which steps through every combination
// of an n-element set by complementing exactly two bits per step.
// Its branching depends on whether m is odd or even and a scan for
// the first position (from the right) where the current combination
// differs from a fixed reference pattern; two of its cases are
// reachable from either parity, so they are written once here and
// selected by the useScan/useC120 flags rather than duplicated.
func (s *combinationState) next(m int) {
	n := s.n
	c := s.c

	j1 := 0
	for j := 1; j <= n-1; j++ {
		if c[n] != c[n-j] {
			j1 = j
			break
		}
	}

	var k1, k2 int
	useScan, useC120 := false, false

	if m%2 == 0 {
		switch {
		case c[n] != 1:
			k1, k2 = n-j1, n-j1+1
		case j1%2 == 0:
			useC120 = true
		default:
			useScan = true
		}
	} else {
		if c[n] == 1 {
			if j1%2 == 0 {
				useScan = true
			} else {
				useC120 = true
			}
		} else {
			k2 = n - j1 - 1
			switch {
			case k2 == 0:
				k1, k2 = 1, n+1-m
			case c[k2+1] == 1 && c[k2] == 1:
				k1 = n
			default:
				k1 = k2 + 1
			}
		}
	}

	if useC120 {
		k1 = n - j1
		k2 = min(k1+2, n)
	}

	if useScan {
		jp := n - j1 - 1
		found := false
		for i := 1; i <= jp; i++ {
			i1 := jp + 2 - i
			if c[i1] == 0 {
				continue
			}
			if c[i1-1] == 1 {
				k1, k2 = i1-1, n-j1
			} else {
				k1, k2 = i1-1, n+1-j1
			}
			found = true
			break
		}
		if !found {
			k1, k2 = 1, n+1-m
		}
	}

	c[k1] = 1 - c[k1]
	c[k2] = 1 - c[k2]
}

// label renders the current combination as a string of letters,
// matching COMBI.C's own output loop.
func (s *combinationState) label() string {
	var b strings.Builder
	for j := 1; j <= s.n; j++ {
		if s.c[j] != 0 {
			b.WriteByte(alphabet[j])
		}
	}
	return b.String()
}

// binomial computes C(n,m), matching COMBI.C's own computation
// exactly: multiply the full descending range (m+1..n), then divide
// by (n-m)! one factor at a time (rather than the more commonly seen
// per-step-interleaved multiplicative formula, which stays exact at
// every step) — reproduced as-is since it is what the original
// computes. This intermediate product can be astronomically larger
// than the final result (n=52, m=5 multiplies out to roughly 8e64
// before ever dividing), so for n close to the documented maximum of
// 52 combined with a small m, this overflows even int64 — the same
// failure COMBI.C's own 32-bit long already had for such inputs, not
// a regression introduced by this port.
func binomial(n, m int) int64 {
	ncom := int64(1)
	for i := m + 1; i <= n; i++ {
		ncom *= int64(i)
	}
	for i := 1; i <= n-m; i++ {
		ncom /= int64(i)
	}
	return ncom
}

// enumerateCombinations returns every combination of n things taken m
// at a time, in the revolving-door algorithm's own generation order.
// Matching COMBI.C's own main() loop, next() is called once before
// recording each combination — including the first — so the initial
// state (elements 1..m) is never itself reported; the sequence
// instead cycles all the way back around to it as the *last* entry.
func enumerateCombinations(n, m int) []string {
	state := newCombinationState(n, m)
	ncom := binomial(n, m)
	results := make([]string, ncom)
	for i := range results {
		state.next(m)
		results[i] = state.label()
	}
	return results
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "combi:", err)
		os.Exit(1)
	}

	fmt.Printf("Combinations of N(<=%d) things taken M(<=%d) at a time\n\n", maxN, maxN)

	var n, m int
	for {
		n = prompter.Int("Number of things", 6)
		if n <= maxN {
			break
		}
		fmt.Printf("N <= %d\n", maxN)
	}
	for {
		m = prompter.Int("Taken m at a time", 4)
		if m > maxN {
			fmt.Printf("M <= %d\n", maxN)
			continue
		}
		if m >= n {
			fmt.Printf("M < %d\n", n)
			continue
		}
		break
	}

	fmt.Printf("%d things taken %d at a time\n", n, m)
	combos := enumerateCombinations(n, m)
	fmt.Printf("Number of combinations = %d\n", len(combos))
	for i, combo := range combos {
		fmt.Printf("%d: %s\n", i+1, combo)
	}
}
