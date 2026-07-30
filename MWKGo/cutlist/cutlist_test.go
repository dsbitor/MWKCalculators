package main

import (
	"os"
	"testing"
)

func loadExample(t *testing.T) (millis, []millis) {
	t.Helper()
	f, err := os.Open("testdata/example.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	stockLength, pieces, err := loadCutlist(f)
	if err != nil {
		t.Fatal(err)
	}
	return stockLength, pieces
}

func TestLoadCutlist_MatchesShippedExample(t *testing.T) {
	stockLength, pieces := loadExample(t)
	if stockLength != 6000 {
		t.Errorf("stockLength = %v, want 6000 (millis)", stockLength)
	}
	// 7+8+6+5+8+12+5 individual pieces, largest first.
	if len(pieces) != 51 {
		t.Fatalf("got %d pieces, want 51", len(pieces))
	}
	if pieces[0] != 2550 || pieces[len(pieces)-1] != 450 {
		t.Errorf("pieces[0]=%v, pieces[last]=%v, want 2550 and 450 (sorted descending)", pieces[0], pieces[len(pieces)-1])
	}
	for i := 1; i < len(pieces); i++ {
		if pieces[i] > pieces[i-1] {
			t.Fatalf("pieces not sorted descending at index %d: %v then %v", i, pieces[i-1], pieces[i])
		}
	}
}

// TestAssignCutlist_MatchesShippedExampleTotals confirms
// assignCutlist reaches the same theoretical-minimum totals CUTS.TXT
// reports for its own (different) algorithm on this dataset: waste
// 5.25 units across 12 standard lengths. Every piece must also be
// accounted for exactly once.
func TestAssignCutlist_MatchesShippedExampleTotals(t *testing.T) {
	stockLength, pieces := loadExample(t)
	stock := assignCutlist(stockLength, pieces)

	if len(stock) != 12 {
		t.Errorf("standard lengths used = %d, want 12", len(stock))
	}

	var waste millis
	var cutCount int
	for _, s := range stock {
		waste += s.Drop
		cutCount += len(s.Sizes)
	}
	if waste != 5250 {
		t.Errorf("waste = %v, want 5250 (millis)", waste)
	}
	if cutCount != len(pieces) {
		t.Errorf("pieces cut = %d, want %d (every piece assigned exactly once)", cutCount, len(pieces))
	}
}

// TestAssignCutlist_NeverExceedsStockLength confirms the basic
// correctness invariant of any cutting-stock assignment: no bar's
// cuts ever add up to more than the standard length it came from.
func TestAssignCutlist_NeverExceedsStockLength(t *testing.T) {
	stockLength, pieces := loadExample(t)
	stock := assignCutlist(stockLength, pieces)
	for i, s := range stock {
		var used millis
		for _, size := range s.Sizes {
			used += size
		}
		if used+s.Drop != stockLength {
			t.Errorf("bar %d: cuts (%v) + drop (%v) = %v, want %v", i, used, s.Drop, used+s.Drop, stockLength)
		}
		if s.Drop < 0 {
			t.Errorf("bar %d: drop = %v, want >= 0", i, s.Drop)
		}
	}
}

func TestToMillis_TruncatesRatherThanRounds(t *testing.T) {
	// CUTLIST.C's own f*1000 cast to an unsigned integer truncates
	// toward zero; a value with a 4th decimal digit should be cut
	// off, not rounded up.
	got := toMillis(1.2559)
	if got != 1255 {
		t.Errorf("toMillis(1.2559) = %v, want 1255 (truncated, not rounded to 1256)", got)
	}
}
