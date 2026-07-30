package main

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestTurnsRequired_ExactDivision(t *testing.T) {
	// 40:1 worm, 40 divisions: exactly one turn per division.
	wholeTurns, remNum, _ := turnsRequired(40, 40)
	if wholeTurns != 1 || remNum != 0 {
		t.Errorf("turnsRequired(40,40) = (%d, %d), want (1, 0)", wholeTurns, remNum)
	}
}

// TestTurnsRequired_ReconstructsRatio is self-verifying: whole turns
// plus the reduced remainder fraction must reconstruct the original
// ratio/divisions value.
func TestTurnsRequired_ReconstructsRatio(t *testing.T) {
	cases := []struct{ ratio, divisions int }{
		{40, 14}, {40, 7}, {90, 13}, {40, 24},
	}
	for _, c := range cases {
		wholeTurns, remNum, remDenom := turnsRequired(c.ratio, c.divisions)
		got := float64(wholeTurns) + float64(remNum)/float64(remDenom)
		want := float64(c.ratio) / float64(c.divisions)
		if diff := got - want; diff > 1e-9 || diff < -1e-9 {
			t.Errorf("turnsRequired(%d,%d) reconstructs to %v, want %v", c.ratio, c.divisions, got, want)
		}
	}
}

func TestRapidIndexSolution(t *testing.T) {
	if step, ok := rapidIndexSolution(24, 8); !ok || step != 3 {
		t.Errorf("rapidIndexSolution(24,8) = (%d,%v), want (3,true)", step, ok)
	}
	if _, ok := rapidIndexSolution(24, 7); ok {
		t.Error("rapidIndexSolution(24,7) ok = true, want false (7 does not divide 24)")
	}
	if _, ok := rapidIndexSolution(-1, 7); ok {
		t.Error("rapidIndexSolution(-1,7) ok = true, want false (no rapid index plate)")
	}
}

// TestPlateSolutions_KnownExample reproduces the real DIVHEAD.DAT
// worked example: 40:1 ratio, 14 divisions (no rapid index plate
// match, since 24 rapid-index holes is not a multiple of 14), against
// the plates DIVHEAD.DAT ships (15-20 on plate A, 21-33 on plate B,
// 37-49 on plate C). 40/14 = 2 whole turns plus a remainder of
// 12/14, which reduces to 6/7, so only a plate whose hole count is a
// multiple of 7 works.
func TestPlateSolutions_KnownExample(t *testing.T) {
	holeCircles := []int{15, 16, 17, 18, 19, 20, 21, 23, 27, 29, 31, 33, 37, 39, 41, 43, 47, 49}
	wholeTurns, remNum, remDenom := turnsRequired(40, 14)
	if wholeTurns != 2 {
		t.Fatalf("turnsRequired(40,14) whole turns = %d, want 2", wholeTurns)
	}
	if remNum != 6 || remDenom != 7 {
		t.Fatalf("turnsRequired(40,14) remainder = %d/%d, want 6/7", remNum, remDenom)
	}

	solutions := plateSolutions(holeCircles, remNum, remDenom)
	want := []plateSolution{{PlateHoles: 21, StepHoles: 18}, {PlateHoles: 49, StepHoles: 42}}
	if len(solutions) != len(want) {
		t.Fatalf("plateSolutions() = %v, want %v (21 and 49 are the multiples of 7 in this plate set)", solutions, want)
	}
	for i := range want {
		if solutions[i] != want[i] {
			t.Errorf("plateSolutions()[%d] = %+v, want %+v", i, solutions[i], want[i])
		}
	}
}

func TestPlateSolutions_NoMatch(t *testing.T) {
	solutions := plateSolutions([]int{15, 16, 17}, 1, 11)
	if len(solutions) != 0 {
		t.Errorf("plateSolutions() = %v, want none (no plate is a multiple of 11)", solutions)
	}
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	const schema = `
		CREATE TABLE divhead_settings (
			id                INTEGER PRIMARY KEY CHECK (id = 1),
			worm_gear_ratio   INTEGER NOT NULL,
			rapid_index_holes INTEGER NOT NULL
		);
		CREATE TABLE divhead_hole_circles (
			holes INTEGER NOT NULL
		);`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	return db
}

func TestLoadSettings_NotConfigured(t *testing.T) {
	db := newTestDB(t)
	_, ok, err := loadSettings(context.Background(), db)
	if err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if ok {
		t.Error("loadSettings() ok = true, want false for an empty table")
	}
}

func TestLoadSettingsAndHoleCircles_RoundTrip(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	if _, err := db.Exec(`INSERT INTO divhead_settings (id, worm_gear_ratio, rapid_index_holes) VALUES (1, 40, 24)`); err != nil {
		t.Fatalf("seed settings: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO divhead_hole_circles (holes) VALUES (15), (16), (21)`); err != nil {
		t.Fatalf("seed hole circles: %v", err)
	}

	s, ok, err := loadSettings(ctx, db)
	if err != nil || !ok {
		t.Fatalf("loadSettings() = (%+v, %v, %v)", s, ok, err)
	}
	if s.WormGearRatio != 40 || s.RapidIndexHoles != 24 {
		t.Errorf("loadSettings() = %+v, want {40 24}", s)
	}

	holes, err := loadHoleCircles(ctx, db)
	if err != nil {
		t.Fatalf("loadHoleCircles() error = %v", err)
	}
	want := []int{15, 16, 21}
	if len(holes) != len(want) {
		t.Fatalf("loadHoleCircles() = %v, want %v", holes, want)
	}
	for i := range want {
		if holes[i] != want[i] {
			t.Errorf("loadHoleCircles()[%d] = %d, want %d", i, holes[i], want[i])
		}
	}
}
