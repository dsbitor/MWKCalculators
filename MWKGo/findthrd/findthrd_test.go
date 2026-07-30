package main

import (
	"context"
	"testing"

	"mwkgo/internal/refdata"
)

var testThreads = []thread{
	{Name: "1/4 UNEF", DiamIn: 0.25, DiamMM: 6.35, PitchTPI: 32, PitchMM: 0.794},
	{Name: "5/16 UNEF", DiamIn: 0.3125, DiamMM: 7.938, PitchTPI: 32, PitchMM: 0.794},
	{Name: "10 W.INS", DiamIn: 0.01, DiamMM: 0.254, PitchTPI: 400, PitchMM: 0.064},
}

func TestFindMatches_ExactMatch(t *testing.T) {
	matches := findMatches(testThreads, searchDiamIn, 0.25, 1.0)
	if len(matches) != 1 || matches[0].Name != "1/4 UNEF" {
		t.Errorf("findMatches(diam=0.25) = %v, want [1/4 UNEF]", matches)
	}
}

func TestFindMatches_WithinTolerance(t *testing.T) {
	// 0.252 is within 1% of 0.25 (0.25 * 1.01 = 0.2525).
	matches := findMatches(testThreads, searchDiamIn, 0.252, 1.0)
	if len(matches) != 1 || matches[0].Name != "1/4 UNEF" {
		t.Errorf("findMatches(diam=0.252, tol=1%%) = %v, want [1/4 UNEF]", matches)
	}
}

func TestFindMatches_OutsideTolerance(t *testing.T) {
	matches := findMatches(testThreads, searchDiamIn, 0.3, 1.0)
	if len(matches) != 0 {
		t.Errorf("findMatches(diam=0.3, tol=1%%) = %v, want none", matches)
	}
}

func TestFindMatches_ToleranceScalesFromTableValueNotTarget(t *testing.T) {
	// FINDTHRD.C brackets the search value between
	// tableValue*(1-tol) and tableValue*(1+tol) -- the tolerance
	// window's width depends on the table entry, not the search
	// value. For an asymmetric case this makes a real difference at
	// the edge: a target 10% below a table value of 400 sits right at
	// the edge of a 10% tolerance window (400*0.9 = 360), which would
	// NOT be within 10% of a target-scaled window (360*(1-0.1)=324 to
	// 360*(1.1)=396, excluding 360's own source value 400... this
	// construction instead directly checks the documented formula).
	th := []thread{{Name: "T", PitchTPI: 400}}
	target := 360.0 // exactly tableValue * (1 - 0.10)
	if len(findMatches(th, searchPitchTPI, target, 10.0)) != 1 {
		t.Error("findMatches() did not match a target exactly at the table-scaled tolerance boundary")
	}
	if len(findMatches(th, searchPitchTPI, target-0.01, 10.0)) != 0 {
		t.Error("findMatches() matched a target just outside the table-scaled tolerance boundary")
	}
}

func TestFindMatches_SearchesEachFieldIndependently(t *testing.T) {
	if len(findMatches(testThreads, searchPitchTPI, 400, 1.0)) != 1 {
		t.Error("findMatches(pitchTPI=400) did not find the expected match")
	}
	if len(findMatches(testThreads, searchDiamMM, 6.35, 1.0)) != 1 {
		t.Error("findMatches(diamMM=6.35) did not find the expected match")
	}
	if len(findMatches(testThreads, searchPitchMM, 0.794, 1.0)) != 2 {
		t.Error("findMatches(pitchMM=0.794) did not find both expected matches")
	}
}

// TestLoadThreads_RealReferenceData is an integration check against
// the actual embedded reference.db, confirming a well known standard
// thread (1/4 UNEF: 0.25 in diameter, 32 tpi) is present.
func TestLoadThreads_RealReferenceData(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	ctx := context.Background()
	db, err := refdata.Open(ctx)
	if err != nil {
		t.Fatalf("open real reference database: %v", err)
	}
	defer db.Close()

	threads, err := loadThreads(ctx, db)
	if err != nil {
		t.Fatalf("loadThreads() error = %v", err)
	}
	if len(threads) < 400 {
		t.Errorf("loadThreads() returned %d threads, want at least 400", len(threads))
	}

	matches := findMatches(threads, searchDiamIn, 0.25, 0.5)
	found := false
	for _, m := range matches {
		if m.Name == "1/4 UNEF" {
			found = true
		}
	}
	if !found {
		t.Error(`findMatches(diamIn=0.25) does not include "1/4 UNEF"`)
	}
}
