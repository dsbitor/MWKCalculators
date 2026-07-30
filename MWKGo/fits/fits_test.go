package main

import (
	"context"
	"database/sql"
	"math"
	"testing"

	_ "modernc.org/sqlite"

	"mwkgo/internal/refdata"
)

// TestShaftDiameter_PushFitWorkedExample reproduces FITS.DAT's own
// documented worked example verbatim: "For a push fit on a nominal
// 1" shaft, machine the hole to exactly 1.0000", and machine the
// shaft to -0.35*(1.0)-0.15 = -0.5 thou less than the nominal size
// (0.9995")."
func TestShaftDiameter_PushFitWorkedExample(t *testing.T) {
	push := fit{Name: "Push", ConstantThou: -0.15, AllowanceThouPerIn: -0.35}
	got := shaftDiameter(push, 1.0)
	want := 0.9995
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("shaftDiameter(push, 1.0) = %v, want %v", got, want)
	}
}

func TestShaftDiameter_ZeroFitLeavesDiameterUnchanged(t *testing.T) {
	// A fit with zero constant and zero allowance is a self-verifying
	// identity: the shaft diameter must equal the nominal diameter.
	zero := fit{Name: "Zero", ConstantThou: 0, AllowanceThouPerIn: 0}
	for _, d := range []float64{0.5, 1.0, 3.25} {
		if got := shaftDiameter(zero, d); got != d {
			t.Errorf("shaftDiameter(zero, %v) = %v, want %v", d, got, d)
		}
	}
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	const createTable = `CREATE TABLE fits (
		list_position         INTEGER PRIMARY KEY,
		name                  TEXT NOT NULL,
		constant_thou         REAL NOT NULL,
		allowance_thou_per_in REAL NOT NULL
	)`
	if _, err := db.Exec(createTable); err != nil {
		t.Fatalf("create fits table: %v", err)
	}
	return db
}

func TestLoadFits_PreservesFileOrder(t *testing.T) {
	db := newTestDB(t)
	// Insert deliberately out of alphabetical order, matching how the
	// original FITS.DAT (and FITS.C, which never sorts) lists fits:
	// the point under test is that loadFits returns them in this
	// insertion order, not resorted.
	if _, err := db.Exec(`INSERT INTO fits (name, constant_thou, allowance_thou_per_in) VALUES
		('Shrink', 0.5, 1.5), ('Force', 0.5, 0.75), ('Push', -0.15, -0.35)`); err != nil {
		t.Fatalf("seed fits: %v", err)
	}

	got, err := loadFits(context.Background(), db)
	if err != nil {
		t.Fatalf("loadFits() error = %v", err)
	}

	want := []string{"Shrink", "Force", "Push"}
	if len(got) != len(want) {
		t.Fatalf("loadFits() returned %d fits, want %d", len(got), len(want))
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Errorf("loadFits()[%d].Name = %q, want %q", i, got[i].Name, name)
		}
	}
}

// TestLoadFits_RealReferenceData is an integration check against the
// actual embedded reference.db: it confirms the fact this program's
// default answer depends on, that "Push" really is fit number 4 in
// the shipped table, matching FITS.DAT's own file order.
func TestLoadFits_RealReferenceData(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	ctx := context.Background()
	db, err := refdata.Open(ctx)
	if err != nil {
		t.Fatalf("open real reference database: %v", err)
	}
	defer db.Close()

	fits, err := loadFits(ctx, db)
	if err != nil {
		t.Fatalf("loadFits() error = %v", err)
	}
	if len(fits) < 4 {
		t.Fatalf("loadFits() returned %d fits, want at least 4", len(fits))
	}
	if fits[3].Name != "Push" {
		t.Errorf("fits[3].Name = %q, want %q (fit number 4, 1-based)", fits[3].Name, "Push")
	}
}
