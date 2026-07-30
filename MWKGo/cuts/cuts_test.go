package main

import (
	"math"
	"os"
	"testing"
)

func loadExample(t *testing.T, path string) (float64, []pieceReq) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	bar, pieces, err := loadCuts(f)
	if err != nil {
		t.Fatal(err)
	}
	return bar, pieces
}

func TestLoadCuts_MatchesShippedExample(t *testing.T) {
	bar, pieces := loadExample(t, "testdata/example.dat")
	if bar != 6.0 {
		t.Errorf("bar = %v, want 6.0", bar)
	}
	// Sorted descending by size, matching CUTS.C's own bubble sort
	// requirement, regardless of the data file's own order.
	want := []pieceReq{
		{Num: 7, Size: 2.55},
		{Num: 8, Size: 1.95},
		{Num: 6, Size: 1.65},
		{Num: 5, Size: 1.35},
		{Num: 8, Size: 0.90},
		{Num: 12, Size: 0.60},
		{Num: 5, Size: 0.45},
	}
	if len(pieces) != len(want) {
		t.Fatalf("got %d pieces, want %d", len(pieces), len(want))
	}
	for i, w := range want {
		if pieces[i] != w {
			t.Errorf("piece %d = %+v, want %+v", i, pieces[i], w)
		}
	}
}

// TestBestCombination_MatchesShippedExampleFirstBar confirms the
// first bar of CUTS.TXT's own worked example: two pieces of length
// 2.55 and one of 0.90 from a 6-unit bar, with zero waste.
func TestBestCombination_MatchesShippedExampleFirstBar(t *testing.T) {
	bar, pieces := loadExample(t, "testdata/example.dat")
	waste, best := bestCombination(bar, pieces, false)
	if waste != 0 {
		t.Errorf("waste = %v, want 0", waste)
	}
	want := []int{2, 0, 0, 0, 1, 0, 0}
	if !intSlicesEqual(best, want) {
		t.Errorf("best = %v, want %v", best, want)
	}
}

// TestCutsMain_MatchesShippedExampleTotals is a coarser end-to-end
// check on the same dataset: after every bar is cut, total waste and
// bar count must match CUTS.TXT's own reported totals exactly (both
// are the theoretical minimum, so this dataset is fully solvable).
func TestCutsMain_MatchesShippedExampleTotals(t *testing.T) {
	bar, pieces := loadExample(t, "testdata/example.dat")
	var twaste float64
	nbar := 0
	for {
		waste, best := bestCombination(bar, pieces, false)
		nbar++
		twaste += waste
		more := false
		for i := range pieces {
			pieces[i].Num -= best[i]
			if pieces[i].Num > 0 {
				more = true
			}
		}
		if waste == bar {
			more = false
		}
		if !more {
			break
		}
	}
	if math.Abs(twaste-5.25) > 1e-6 {
		t.Errorf("total waste = %v, want 5.25", twaste)
	}
	if nbar != 12 {
		t.Errorf("standard lengths used = %d, want 12", nbar)
	}
}

// TestBestCombination_ZeroWasteFlagMatchesUpdateNote confirms
// CUTS.TXT's own "Update 2/02" case: cutting 1,10 / 2,7 / 1,6 / 2,4
// from 20-unit stock. Without restricting to zero waste, the greedy
// search's first bar wastes 3 (a 1+1+... combination); restricting to
// zero waste first finds the single bar that fits a 10, a 6, and a 4
// with none left over.
func TestBestCombination_ZeroWasteFlagMatchesUpdateNote(t *testing.T) {
	bar, pieces := loadExample(t, "testdata/zerowaste.dat")

	_, greedyBest := bestCombination(bar, pieces, false)
	if intSlicesEqual(greedyBest, []int{1, 0, 1, 1}) {
		t.Error("unrestricted search unexpectedly found the zero-waste combination first")
	}

	zeroWaste, zeroBest := bestCombination(bar, pieces, true)
	if zeroWaste != 0 {
		t.Errorf("zero-waste search waste = %v, want 0", zeroWaste)
	}
	want := []int{1, 0, 1, 1} // 10 + 6 + 4 = 20 exactly
	if !intSlicesEqual(zeroBest, want) {
		t.Errorf("zero-waste search best = %v, want %v", zeroBest, want)
	}
}

func TestDiff_WithinPrecisionIsZero(t *testing.T) {
	if diff(1.0001, 1.0) != 0 {
		t.Error("diff(1.0001, 1.0) should be 0 (within precision)")
	}
	if diff(1.1, 1.0) != 1 {
		t.Error("diff(1.1, 1.0) should be 1")
	}
	if diff(0.9, 1.0) != -1 {
		t.Error("diff(0.9, 1.0) should be -1")
	}
}
