package main

import (
	"context"
	"math"
	"testing"

	"mwkgo/internal/refdata"
)

// TestLengthChangeAndTemperatureChange_AreInverses is self-verifying:
// computing a length change and then recovering the temperature
// change from it must reproduce the original temperature change, for
// any material and length.
func TestLengthChangeAndTemperatureChange_AreInverses(t *testing.T) {
	cases := []struct {
		cte, length, deltaT float64
	}{
		{12.44e-6, 10.0, 100.0}, // aluminum, 10 in, 100 degF
		{6.3e-6, 1.0, -50.0},    // carbon steel, 1 in, cooling
		{17e-6, 24.0, 20.0},     // zinc, 24 in, small change
	}
	for _, c := range cases {
		dl := lengthChange(c.cte, c.length, c.deltaT)
		gotDeltaT := temperatureChange(c.cte, c.length, dl)
		if math.Abs(gotDeltaT-c.deltaT) > 1e-9 {
			t.Errorf("temperatureChange(lengthChange(%v,%v,%v)) = %v, want %v", c.cte, c.length, c.deltaT, gotDeltaT, c.deltaT)
		}
	}
}

func TestLengthChange_ZeroTemperatureChangeIsZeroLengthChange(t *testing.T) {
	if got := lengthChange(12.44e-6, 10.0, 0); got != 0 {
		t.Errorf("lengthChange with zero deltaT = %v, want 0", got)
	}
}

// TestLoadMaterials_RealReferenceData is an integration check against
// the actual embedded reference.db, confirming a well known material
// (aluminum, CTE close to 12.44 ppm/degF, a widely cited reference
// value) is present, and that materials come back sorted
// alphabetically case-insensitively.
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

	var found bool
	for _, m := range materials {
		if m.Name == "Aluminum" {
			found = true
			if math.Abs(m.CTEPPMPerDegF-12.44) > 1e-9 {
				t.Errorf("Aluminum CTE = %v, want 12.44", m.CTEPPMPerDegF)
			}
		}
	}
	if !found {
		t.Error(`materials does not contain "Aluminum"`)
	}

	for i := 1; i < len(materials); i++ {
		prev, cur := materials[i-1].Name, materials[i].Name
		if compareFold(prev, cur) > 0 {
			t.Errorf("materials not sorted case-insensitively: %q before %q", prev, cur)
		}
	}
}

// compareFold does a simple case-insensitive ASCII comparison,
// avoiding a dependency on any particular collation library for this
// test's own ordering check.
func compareFold(a, b string) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		ca, cb := toLowerByte(a[i]), toLowerByte(b[i])
		if ca != cb {
			return int(ca) - int(cb)
		}
	}
	return len(a) - len(b)
}

func toLowerByte(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
