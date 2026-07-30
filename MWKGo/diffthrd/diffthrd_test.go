package main

import (
	"context"
	"database/sql"
	"math"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"mwkgo/internal/csvtable"
)

func TestExactMatch(t *testing.T) {
	pitches := []float64{4, 8, 20, 40}
	if !exactMatch(pitches, 20) {
		t.Error("exactMatch(20) = false, want true")
	}
	if exactMatch(pitches, 21) {
		t.Error("exactMatch(21) = true, want false")
	}
}

// TestBestPair_EffectivePitchFormula is self-verifying: for the pair
// bestPair returns, the differential effective-pitch formula
// (1/coarse - 1/fine, inverted) must reproduce the reported effective
// pitch.
func TestBestPair_EffectivePitchFormula(t *testing.T) {
	pitches := []float64{8, 10, 12, 16, 20, 40}
	coarse, fine, effective, found := bestPair(pitches, 100)
	if !found {
		t.Fatal("bestPair() found = false, want true")
	}
	if coarse > fine {
		t.Errorf("coarse (%v) > fine (%v), want coarse <= fine", coarse, fine)
	}
	want := 1 / (1/coarse - 1/fine)
	if math.Abs(effective-want) > 1e-9 {
		t.Errorf("effective = %v, want %v (recomputed from coarse=%v, fine=%v)", effective, want, coarse, fine)
	}
}

// TestBestPair_IsActuallyClosestAmongAllPairs is self-verifying: it
// independently recomputes every pair's effective pitch and checks
// none is closer to the desired pitch than the one bestPair chose.
func TestBestPair_IsActuallyClosestAmongAllPairs(t *testing.T) {
	pitches := []float64{4, 4.5, 8, 10, 20, 40, 80}
	desired := 250.0
	_, _, effective, found := bestPair(pitches, desired)
	if !found {
		t.Fatal("bestPair() found = false, want true")
	}
	gotError := math.Abs(effective - desired)

	for i, pi := range pitches {
		for j, pj := range pitches {
			if i == j {
				continue
			}
			pc, pf := pi, pj
			if pc > pf {
				pc, pf = pf, pc
			}
			d := 1/pc - 1/pf
			if d == 0 {
				continue
			}
			eff := 1 / d
			if err := math.Abs(eff - desired); err < gotError-1e-9 {
				t.Errorf("pair (%v,%v) has effective pitch %v (error %v), closer than bestPair's result (error %v)", pc, pf, eff, err, gotError)
			}
		}
	}
}

func TestBestPair_NoUsablePairs(t *testing.T) {
	// A single pitch can't form a pair at all.
	if _, _, _, found := bestPair([]float64{10}, 100); found {
		t.Error("bestPair() with one pitch found = true, want false")
	}
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if _, err := db.Exec(`CREATE TABLE ` + pitchesTable + ` (pitch_tpi REAL NOT NULL)`); err != nil {
		t.Fatalf("create %s table: %v", pitchesTable, err)
	}
	return db
}

func TestLoadPitches_EmptyTable(t *testing.T) {
	db := newTestDB(t)
	pitches, err := loadPitches(context.Background(), db)
	if err != nil {
		t.Fatalf("loadPitches() error = %v", err)
	}
	if len(pitches) != 0 {
		t.Errorf("loadPitches() = %v, want empty", pitches)
	}
}

func TestImportPitches_RoundTripsThroughCSV(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	csvData := "pitch_tpi\n4\n8\n20\n"
	if err := csvtable.Import(ctx, db, pitchesTable, strings.NewReader(csvData)); err != nil {
		t.Fatalf("csvtable.Import() error = %v", err)
	}

	pitches, err := loadPitches(ctx, db)
	if err != nil {
		t.Fatalf("loadPitches() error = %v", err)
	}
	want := []float64{4, 8, 20}
	if len(pitches) != len(want) {
		t.Fatalf("loadPitches() = %v, want %v", pitches, want)
	}
	for i := range want {
		if pitches[i] != want[i] {
			t.Errorf("loadPitches()[%d] = %v, want %v", i, pitches[i], want[i])
		}
	}
}
