package main

import (
	"context"
	"math"
	"strings"
	"testing"

	"mwkgo/internal/refdata"
)

const eps = 1e-9

func almostEqual(a, b float64) bool { return math.Abs(a-b) < eps }

// TestLoadMaterials_RealReferenceData checks loadMaterials against the
// actual shipped reference.db (not a fake), the same way expand's own
// loadMaterials test does: a well known material (Aluminum, density
// close to the widely cited 0.0924 lb/in^3) is present, and the list
// comes back sorted alphabetically case-insensitively.
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
			if math.Abs(m.DensityLbPerIn3-0.0924) > 1e-9 {
				t.Errorf("Aluminum density = %v, want 0.0924", m.DensityLbPerIn3)
			}
		}
	}
	if !found {
		t.Error(`materials does not contain "Aluminum"`)
	}

	for i := 1; i < len(materials); i++ {
		if strings.ToLower(materials[i-1].Name) > strings.ToLower(materials[i].Name) {
			t.Errorf("materials not sorted case-insensitively: %q before %q", materials[i-1].Name, materials[i].Name)
		}
	}
}

// TestWeight_MaterialTimesVolumeTimesUnitFactor exercises the actual
// weight computation main() performs (density * volume * unit
// conversion factor), tying a real material's density from
// reference.db to a known shape's volume - the end-to-end path the
// pure per-shape volume-formula tests elsewhere in this file don't
// cover on their own.
func TestWeight_MaterialTimesVolumeTimesUnitFactor(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	ctx := context.Background()
	db, err := refdata.Open(ctx)
	if err != nil {
		t.Fatalf("open real reference database: %v", err)
	}
	defer db.Close()

	materials, err := loadMaterials(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	var aluminum material
	for _, m := range materials {
		if m.Name == "Aluminum" {
			aluminum = m
		}
	}
	if aluminum.Name == "" {
		t.Fatal(`reference data does not contain "Aluminum"`)
	}

	// A 1 in diameter sphere of Aluminum, in inches: volume =
	// PI/6 = 0.523599 in^3, weight = density * volume.
	volume := sphereVolume(1)
	_, factIn := unitConversion("i")
	gotWeight := aluminum.DensityLbPerIn3 * volume * factIn
	wantWeight := 0.0924 * (math.Pi / 6)
	if !almostEqual(gotWeight, wantWeight) {
		t.Errorf("weight (in) = %v, want %v", gotWeight, wantWeight)
	}

	// The same sphere, but with its diameter given in millimeters
	// (25.4 mm = 1 in): the unit conversion factor must bring the
	// computed weight back to the same value as the inches case.
	volumeMM := sphereVolume(25.4)
	_, factMM := unitConversion("m")
	gotWeightMM := aluminum.DensityLbPerIn3 * volumeMM * factMM
	if !almostEqual(gotWeightMM, wantWeight) {
		t.Errorf("weight (25.4mm sphere, mm units) = %v, want %v (same physical sphere as the 1in case)", gotWeightMM, wantWeight)
	}
}

func TestRectangularVolume(t *testing.T) {
	if got := rectangularVolume(2, 3, 4); got != 24 {
		t.Errorf("rectangularVolume(2,3,4) = %v, want 24", got)
	}
}

func TestCylinderVolume_MatchesStandardFormula(t *testing.T) {
	d, h := 4.0, 10.0
	r := d / 2
	want := math.Pi * r * r * h
	if got := cylinderVolume(d, h); !almostEqual(got, want) {
		t.Errorf("cylinderVolume(%v,%v) = %v, want %v", d, h, got, want)
	}
}

func TestAnnulusVolume_MatchesPipeVolume(t *testing.T) {
	// The washer-shaped annulus and the cylindrical pipe are the same
	// solid described two ways; WEIGHT.C's own formulas for them are
	// algebraically identical.
	d1, d2, h := 1.5, 0.75, 0.25
	if got, want := annulusVolume(d1, d2, h), pipeVolume(d1, d2, h); !almostEqual(got, want) {
		t.Errorf("annulusVolume(%v,%v,%v) = %v, want %v (== pipeVolume)", d1, d2, h, got, want)
	}
}

func TestFrustumVolume_EqualDiametersMatchesCylinder(t *testing.T) {
	// A frustum whose two end diameters are equal is just a cylinder.
	d, h := 3.0, 8.0
	if got, want := frustumVolume(d, d, h), cylinderVolume(d, h); !almostEqual(got, want) {
		t.Errorf("frustumVolume(%v,%v,%v) = %v, want %v (== cylinderVolume)", d, d, h, got, want)
	}
}

func TestSphereVolume_MatchesStandardFormula(t *testing.T) {
	d := 6.0
	r := d / 2
	want := (4.0 / 3.0) * math.Pi * r * r * r
	if got := sphereVolume(d); !almostEqual(got, want) {
		t.Errorf("sphereVolume(%v) = %v, want %v", d, got, want)
	}
}

func TestSphericalCapVolume_HemisphereIsHalfSphere(t *testing.T) {
	// A cap of height equal to the radius is exactly a hemisphere.
	d := 5.0
	got := sphericalCapVolume(d, d/2)
	want := sphereVolume(d) / 2
	if !almostEqual(got, want) {
		t.Errorf("sphericalCapVolume(%v, %v) = %v, want %v (half of sphereVolume)", d, d/2, got, want)
	}
}

func TestHemisphericalShellVolume_ZeroThicknessIsZero(t *testing.T) {
	if got := hemisphericalShellVolume(2, 2); got != 0 {
		t.Errorf("hemisphericalShellVolume(2,2) = %v, want 0", got)
	}
}

func TestHemisphericalShellVolume_NoInsideIsHalfSphere(t *testing.T) {
	// A shell with no inside cavity is just a solid hemisphere.
	d := 4.0
	got := hemisphericalShellVolume(d, 0)
	want := sphereVolume(d) / 2
	if !almostEqual(got, want) {
		t.Errorf("hemisphericalShellVolume(%v,0) = %v, want %v", d, got, want)
	}
}

func TestTorusVolume_MatchesPappusTheorem(t *testing.T) {
	crossSection, outside := 1.5, 6.5
	centerlineRadius := (outside - crossSection) / 2
	crossSectionArea := math.Pi * (crossSection / 2) * (crossSection / 2)
	want := 2 * math.Pi * centerlineRadius * crossSectionArea
	if got := torusVolume(crossSection, outside); !almostEqual(got, want) {
		t.Errorf("torusVolume(%v,%v) = %v, want %v (Pappus's theorem)", crossSection, outside, got, want)
	}
}

func TestConicalWedgeVolume_EqualSagittaHalvesTheCone(t *testing.T) {
	h, r := 10.0, 3.0
	fullCone := math.Pi * r * r * h / 3
	got := conicalWedgeVolume(h, r, r)
	if !almostEqual(got, fullCone/2) {
		t.Errorf("conicalWedgeVolume(%v,%v,%v) = %v, want %v (half the full cone)", h, r, r, got, fullCone/2)
	}
}

func TestConicalWedgeVolume_ZeroSagittaIsZeroVolume(t *testing.T) {
	// A sagitta of 0 means the cutting plane is tangent to the base
	// circle, slicing off nothing.
	h, r := 10.0, 3.0
	if got := conicalWedgeVolume(h, 0, r); !almostEqual(got, 0) {
		t.Errorf("conicalWedgeVolume(%v,0,%v) = %v, want 0", h, r, got)
	}
}

func TestConicalWedgeVolume_SagittaOfDiameterIsWholeCone(t *testing.T) {
	// A sagitta equal to the base diameter (2r) means the cutting
	// plane just clears the far side of the base circle, leaving the
	// whole cone.
	h, r := 10.0, 3.0
	fullCone := math.Pi * r * r * h / 3
	got := conicalWedgeVolume(h, 2*r, r)
	if !almostEqual(got, fullCone) {
		t.Errorf("conicalWedgeVolume(%v,%v,%v) = %v, want %v (the whole cone)", h, 2*r, r, got, fullCone)
	}
}

func TestTangentOgiveVolume_ZeroLengthIsZero(t *testing.T) {
	if got := tangentOgiveVolume(2, 0); !almostEqual(got, 0) {
		t.Errorf("tangentOgiveVolume(2,0) = %v, want 0", got)
	}
}

func TestPolygonVolume_RejectsFewerThanThreeSides(t *testing.T) {
	if _, err := polygonVolume(2, 1, 1); err == nil {
		t.Error("expected an error for a 2-sided \"polygon\"")
	}
}

func TestPolygonVolume_ApproachesCylinderAsSidesIncrease(t *testing.T) {
	// A regular polygon with many sides, measured across flats,
	// approximates a circle of the same diameter.
	d, h := 4.0, 10.0
	got, err := polygonVolume(2000, d, h)
	if err != nil {
		t.Fatal(err)
	}
	want := cylinderVolume(d, h)
	if diff := math.Abs(got-want) / want; diff > 1e-4 {
		t.Errorf("polygonVolume(2000,%v,%v) = %v, want ~%v (cylinder limit), relative diff %v", d, h, got, want, diff)
	}
}

func TestPolygonVolume_HexagonMatchesKnownArea(t *testing.T) {
	// A regular hexagon with distance 1 across flats has side length
	// 1/sqrt(3) and area (3*sqrt(3)/2)*side^2 = sqrt(3)/2, a standard
	// closed-form result independent of WEIGHT.C's own trig formula.
	across, h := 1.0, 1.0
	side := across / math.Sqrt(3)
	wantArea := (3 * math.Sqrt(3) / 2) * side * side
	got, err := polygonVolume(6, across, h)
	if err != nil {
		t.Fatal(err)
	}
	if !almostEqual(got, wantArea*h) {
		t.Errorf("polygonVolume(6,%v,%v) = %v, want %v", across, h, got, wantArea*h)
	}
}

func TestUnitConversion_KnownFactors(t *testing.T) {
	cases := []struct {
		choice   string
		wantName string
		wantFact float64
	}{
		{"i", "in", 1},
		{"f", "ft", 1728},
		{"m", "mm", 1 / math.Pow(25.4, 3)},
		{"c", "cm", 1 / math.Pow(2.54, 3)},
		{"x", "in", 1}, // anything unrecognized defaults to inches, matching WEIGHT.C
	}
	for _, c := range cases {
		name, fact := unitConversion(c.choice)
		if name != c.wantName || !almostEqual(fact, c.wantFact) {
			t.Errorf("unitConversion(%q) = (%q, %v), want (%q, %v)", c.choice, name, fact, c.wantName, c.wantFact)
		}
	}
}
