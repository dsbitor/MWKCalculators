package main

import (
	"context"
	"math"
	"testing"

	"mwkgo/internal/refdata"
)

var testWireGages = []gageEntry{
	{Gage: "10", SizeIn: 0.1019},
	{Gage: "12", SizeIn: 0.0808},
	{Gage: "14", SizeIn: 0.0641},
}

func TestFindByDesignation_ExactMatch(t *testing.T) {
	entry, ok := findByDesignation(testWireGages, "12")
	if !ok {
		t.Fatal("findByDesignation(12) ok = false, want true")
	}
	if entry.SizeIn != 0.0808 {
		t.Errorf("findByDesignation(12).SizeIn = %v, want 0.0808", entry.SizeIn)
	}
}

func TestFindByDesignation_NotFound(t *testing.T) {
	if _, ok := findByDesignation(testWireGages, "999"); ok {
		t.Error("findByDesignation(999) ok = true, want false")
	}
}

func TestFindClosest_ExactAndNearMatches(t *testing.T) {
	cases := []struct {
		size float64
		want string
	}{
		{0.0808, "12"}, // exact match
		{0.081, "12"},  // closer to 12 than to 10 or 14
		{0.09, "12"},   // between 12 (0.0808) and 10 (0.1019), closer to 12
		{1.0, "10"},    // far outside the range: closest is the nearest endpoint
	}
	for _, c := range cases {
		got := findClosest(testWireGages, c.size)
		if got.Gage != c.want {
			t.Errorf("findClosest(%v) = %q, want %q", c.size, got.Gage, c.want)
		}
	}
}

// TestFindClosest_IsAlwaysAtLeastAsCloseAsAnyOtherEntry is
// self-verifying: for a range of query sizes, the chosen entry's
// distance must be less than or equal to every other entry's
// distance, which is the actual defining property of "closest"
// rather than a re-statement of the algorithm.
func TestFindClosest_IsAlwaysAtLeastAsCloseAsAnyOtherEntry(t *testing.T) {
	for _, size := range []float64{0.05, 0.0808, 0.095, 0.2, 0.0} {
		got := findClosest(testWireGages, size)
		gotDiff := math.Abs(size - got.SizeIn)
		for _, e := range testWireGages {
			if diff := math.Abs(size - e.SizeIn); diff < gotDiff {
				t.Errorf("findClosest(%v) = %q (diff %v), but %q is closer (diff %v)", size, got.Gage, gotDiff, e.Gage, diff)
			}
		}
	}
}

// TestLoadGages_RealReferenceData is an integration check against the
// actual embedded reference.db, confirming known standard wire and
// sheet gage sizes are present (AWG 0000, 0.4600 in diameter, is a
// widely documented reference value).
func TestLoadGages_RealReferenceData(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	ctx := context.Background()
	db, err := refdata.Open(ctx)
	if err != nil {
		t.Fatalf("open real reference database: %v", err)
	}
	defer db.Close()

	wire, err := loadGages(ctx, db, "wire_gages")
	if err != nil {
		t.Fatalf("loadGages(wire_gages) error = %v", err)
	}
	entry, ok := findByDesignation(wire, "0000")
	if !ok {
		t.Fatal(`findByDesignation(wire, "0000") ok = false, want true`)
	}
	if math.Abs(entry.SizeIn-0.4600) > 1e-9 {
		t.Errorf(`wire gage "0000" size = %v, want 0.4600`, entry.SizeIn)
	}

	sheet, err := loadGages(ctx, db, "sheet_gages")
	if err != nil {
		t.Fatalf("loadGages(sheet_gages) error = %v", err)
	}
	if len(sheet) == 0 {
		t.Error("sheet_gages is empty")
	}
}
