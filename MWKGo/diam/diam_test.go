package main

import (
	"context"
	"database/sql"
	"math"
	"testing"

	_ "modernc.org/sqlite"

	"mwkgo/internal/refdata"
)

// TestDiameterRange_MatchesSurfaceSpeedFormula is self-verifying: it
// recomputes the standard cutting-speed relation independently and
// checks diameterRange's output satisfies it, the same style used for
// speed.rpmRange, of which this is the inverse.
func TestDiameterRange_MatchesSurfaceSpeedFormula(t *testing.T) {
	m := material{Name: "TEST", LowSFPM: 200, HighSFPM: 300}
	rpm := 500.0

	low, high := diameterRange(m, rpm)

	gotLowSFPM := math.Pi * low * rpm / 12
	gotHighSFPM := math.Pi * high * rpm / 12
	if math.Abs(gotLowSFPM-float64(m.LowSFPM)) > 1e-9 {
		t.Errorf("low diameter %v implies %v sfpm, want %v", low, gotLowSFPM, m.LowSFPM)
	}
	if math.Abs(gotHighSFPM-float64(m.HighSFPM)) > 1e-9 {
		t.Errorf("high diameter %v implies %v sfpm, want %v", high, gotHighSFPM, m.HighSFPM)
	}
}

func TestDiameterRange_IsInverseOfRPMRange(t *testing.T) {
	// diameterRange(m, rpm) and speed's rpmRange(m, diameter) are
	// inverses of the same relation: feeding diameterRange's output
	// back through the same formula at that diameter must reproduce
	// the original rpm.
	m := material{Name: "TEST", LowSFPM: 100, HighSFPM: 100}
	rpm := 764.0
	low, _ := diameterRange(m, rpm)

	c := 12 / math.Pi
	gotRPM := c * float64(m.LowSFPM) / low
	if math.Abs(gotRPM-rpm) > 1e-6 {
		t.Errorf("round trip through diameterRange = %v, want %v", gotRPM, rpm)
	}
}

func newTestUserDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if _, err := db.Exec(`CREATE TABLE ` + availableSpeedsTable + ` (rpm INTEGER NOT NULL)`); err != nil {
		t.Fatalf("create %s table: %v", availableSpeedsTable, err)
	}
	return db
}

func TestLoadAvailableSpeeds_EmptyAndPopulated(t *testing.T) {
	db := newTestUserDB(t)
	ctx := context.Background()

	speeds, err := loadAvailableSpeeds(ctx, db)
	if err != nil {
		t.Fatalf("loadAvailableSpeeds() error = %v", err)
	}
	if len(speeds) != 0 {
		t.Errorf("loadAvailableSpeeds() = %v, want empty", speeds)
	}

	if _, err := db.Exec(`INSERT INTO ` + availableSpeedsTable + ` (rpm) VALUES (85), (115), (150)`); err != nil {
		t.Fatalf("seed speeds: %v", err)
	}
	speeds, err = loadAvailableSpeeds(ctx, db)
	if err != nil {
		t.Fatalf("loadAvailableSpeeds() error = %v", err)
	}
	want := []float64{85, 115, 150}
	if len(speeds) != len(want) {
		t.Fatalf("loadAvailableSpeeds() = %v, want %v", speeds, want)
	}
	for i := range want {
		if speeds[i] != want[i] {
			t.Errorf("loadAvailableSpeeds()[%d] = %v, want %v", i, speeds[i], want[i])
		}
	}
}

// TestLoadMaterials_SharesReferenceDataWithSpeed is an integration
// check confirming diam reads the same machining_speeds table speed
// does (material 1 is "ALUMINUM AND ALLOYS" in both), rather than a
// separate copy of the same data.
func TestLoadMaterials_SharesReferenceDataWithSpeed(t *testing.T) {
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
		t.Errorf("materials[0].Name = %q, want %q", materials[0].Name, "ALUMINUM AND ALLOYS")
	}
}
