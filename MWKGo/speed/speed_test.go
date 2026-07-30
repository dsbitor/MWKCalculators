package main

import (
	"context"
	"database/sql"
	"math"
	"testing"

	_ "modernc.org/sqlite"

	"mwkgo/internal/refdata"
)

// TestRPMRange_MatchesSurfaceSpeedFormula is self-verifying: it
// recomputes the standard cutting-speed relation independently
// (surface speed = pi * diameter * rpm / 12) and checks rpmRange's
// output satisfies it, rather than asserting a specific number.
func TestRPMRange_MatchesSurfaceSpeedFormula(t *testing.T) {
	m := material{Name: "TEST", LowSFPM: 200, HighSFPM: 300}
	diameter := 1.5

	low, high := rpmRange(m, diameter)

	gotLowSFPM := math.Pi * diameter * low / 12
	gotHighSFPM := math.Pi * diameter * high / 12
	if math.Abs(gotLowSFPM-float64(m.LowSFPM)) > 1e-9 {
		t.Errorf("low rpm %v implies %v sfpm, want %v", low, gotLowSFPM, m.LowSFPM)
	}
	if math.Abs(gotHighSFPM-float64(m.HighSFPM)) > 1e-9 {
		t.Errorf("high rpm %v implies %v sfpm, want %v", high, gotHighSFPM, m.HighSFPM)
	}
	if low >= high {
		t.Errorf("low rpm %v is not less than high rpm %v", low, high)
	}
}

func TestRPMRange_HalvingDiameterDoublesRPM(t *testing.T) {
	// Surface speed is proportional to diameter * rpm, so at constant
	// surface speed, halving the diameter must double the rpm.
	m := material{Name: "TEST", LowSFPM: 100, HighSFPM: 100}
	low1, _ := rpmRange(m, 2.0)
	low2, _ := rpmRange(m, 1.0)
	if math.Abs(low2-2*low1) > 1e-9 {
		t.Errorf("rpm at diameter 1.0 = %v, want double rpm at diameter 2.0 (%v)", low2, low1)
	}
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	const createTable = `CREATE TABLE machining_speeds (
		list_position INTEGER PRIMARY KEY,
		material      TEXT    NOT NULL,
		low_sfpm      INTEGER NOT NULL,
		high_sfpm     INTEGER NOT NULL
	)`
	if _, err := db.Exec(createTable); err != nil {
		t.Fatalf("create machining_speeds table: %v", err)
	}
	return db
}

func TestLoadMaterials_PreservesFileOrder(t *testing.T) {
	db := newTestDB(t)
	// Deliberately out of alphabetical order, matching SPEED.DAT (and
	// SPEED.C, which never sorts): loadMaterials must return this
	// insertion order, not resorted.
	if _, err := db.Exec(`INSERT INTO machining_speeds (material, low_sfpm, high_sfpm) VALUES
		('ALUMINUM AND ALLOYS', 200, 300), ('BRASS AND SOFT BRONZE', 100, 300), ('LOW CARBON STEEL', 80, 150)`); err != nil {
		t.Fatalf("seed machining_speeds: %v", err)
	}

	got, err := loadMaterials(context.Background(), db)
	if err != nil {
		t.Fatalf("loadMaterials() error = %v", err)
	}

	want := []string{"ALUMINUM AND ALLOYS", "BRASS AND SOFT BRONZE", "LOW CARBON STEEL"}
	if len(got) != len(want) {
		t.Fatalf("loadMaterials() returned %d materials, want %d", len(got), len(want))
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Errorf("loadMaterials()[%d].Name = %q, want %q", i, got[i].Name, name)
		}
	}
}

// TestLoadMaterials_RealReferenceData is an integration check against
// the actual embedded reference.db: it confirms the fact this
// program's default answer depends on, that "ALUMINUM AND ALLOYS"
// really is material number 1 in the shipped table, matching
// SPEED.DAT's own file order.
func TestLoadMaterials_RealReferenceData(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	ctx := context.Background()
	db, err := refdata.Open(ctx)
	if err != nil {
		t.Fatalf("open real reference database: %v", err)
	}
	defer db.Close()

	materials, err := loadMaterials(ctx, db)
	if err != nil {
		t.Fatalf("loadMaterials() error = %v", err)
	}
	if len(materials) == 0 {
		t.Fatal("loadMaterials() returned no materials")
	}
	if materials[0].Name != "ALUMINUM AND ALLOYS" {
		t.Errorf("materials[0].Name = %q, want %q (material number 1, 1-based)", materials[0].Name, "ALUMINUM AND ALLOYS")
	}
}
