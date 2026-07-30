package main

import (
	"context"
	"math"
	"testing"
)

func realUnitsAndPrefixes(t *testing.T) ([]unitDef, []prefixDef) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	units, prefixes, err := loadUnitsAndPrefixes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return units, prefixes
}

func convertOnce(t *testing.T, units []unitDef, prefixes []prefixDef, x float64, from, to string) float64 {
	t.Helper()
	f, err := uanal(from, units, prefixes)
	if err != nil {
		t.Fatalf("uanal(%q) = %v", from, err)
	}
	tt, err := uanal(to, units, prefixes)
	if err != nil {
		t.Fatalf("uanal(%q) = %v", to, err)
	}
	if !compatible(f.PU, tt.PU) {
		t.Fatalf("%q and %q are not compatible: %v vs %v", from, to, f.PU, tt.PU)
	}
	y, _ := uconv(x, f, tt, units, prefixes)
	return y
}

func almostEqual(a, b, tol float64) bool { return math.Abs(a-b) < tol }

// The following tests reproduce UNIT.TXT's own worked tutorial
// examples exactly, against the real shipped reference.db.

func TestConvert_MilesPerHourToFeetPerSecond(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	got := convertOnce(t, units, prefixes, 60, "MILES/HOUR", "FEET/SEC")
	if !almostEqual(got, 88, 1e-6) {
		t.Errorf("60 MILES/HOUR = %v FEET/SEC, want 88", got)
	}
}

func TestConvert_MPHToFPSUsingAbbreviations(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	got := convertOnce(t, units, prefixes, 60, "MPH", "FPS")
	if !almostEqual(got, 88, 1e-6) {
		t.Errorf("60 MPH = %v FPS, want 88", got)
	}
}

func TestConvert_ImpliedUnitValueOfOne(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	got := convertOnce(t, units, prefixes, 1, "MILE/HOUR", "FEET/SEC")
	if !almostEqual(got, 1.46667, 1e-4) {
		t.Errorf("1 MILE/HOUR = %v FEET/SEC, want 1.46667", got)
	}
}

func TestConvert_DensityPoundsPerCubicFootToMilligramPerCubicCentimeter(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	got := convertOnce(t, units, prefixes, 10, "LB/FT^3", "MILLIGRAM/CENTIMETER^3")
	if !almostEqual(got, 160.185, 1e-2) {
		t.Errorf("10 LB/FT^3 = %v MILLIGRAM/CENTIMETER^3, want 160.185", got)
	}
}

func TestConvert_AreaSquareFeetToSquareInches(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	got := convertOnce(t, units, prefixes, 20, "FT^2", "IN^2")
	if !almostEqual(got, 2880, 1e-6) {
		t.Errorf("20 FT^2 = %v IN^2, want 2880", got)
	}
}

func TestConvert_PoundMassToKilogram(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	got := convertOnce(t, units, prefixes, 1, "POUND", "KG")
	if !almostEqual(got, 0.453592, 1e-5) {
		t.Errorf("1 POUND = %v KG, want 0.453592", got)
	}
}

func TestConvert_PoundforceToNewton(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	got := convertOnce(t, units, prefixes, 1, "POUNDFORCE", "NEWTON")
	if !almostEqual(got, 4.44822, 1e-4) {
		t.Errorf("1 POUNDFORCE = %v NEWTON, want 4.44822", got)
	}
}

func TestConvert_PoundMassToNewtonIsIncompatible(t *testing.T) {
	// Mass and force are different dimensions; UNIT.TXT documents this
	// exact "pound (mass) -> newton" example as INCOMPATIBLE UNITS.
	units, prefixes := realUnitsAndPrefixes(t)
	f, err := uanal("POUND", units, prefixes)
	if err != nil {
		t.Fatal(err)
	}
	tt, err := uanal("NEWTON", units, prefixes)
	if err != nil {
		t.Fatal(err)
	}
	if compatible(f.PU, tt.PU) {
		t.Error("POUND (mass) and NEWTON (force) should not be compatible")
	}
}

func TestUanal_HandlesPrefixedUnits(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	p, err := uanal("KILOMETER", units, prefixes)
	if err != nil {
		t.Fatal(err)
	}
	if p.NumPrefixIdx == -1 {
		t.Error("expected KILOMETER to resolve a KILO prefix")
	}
	if units[p.NumUnitIdx].Name != "METER" {
		t.Errorf("resolved unit = %q, want METER", units[p.NumUnitIdx].Name)
	}
}

func TestUanal_RejectsUnknownUnit(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	if _, err := uanal("NOTAREALUNIT", units, prefixes); err == nil {
		t.Error("expected an error for an unknown unit")
	}
}

func TestUanal_RejectsEmptyExpression(t *testing.T) {
	units, prefixes := realUnitsAndPrefixes(t)
	if _, err := uanal("", units, prefixes); err != errNoUnitSpecified {
		t.Errorf("uanal(\"\") error = %v, want errNoUnitSpecified", err)
	}
}

// TestDeplural_PluralUnitNamesResolve confirms deplural strips a
// trailing "S" or "ES" from a made-up plural form, but leaves an
// exact, already-known unit name (like MILES, itself defined in
// UNIT.DAT) completely alone.
func TestDeplural_PluralUnitNamesResolve(t *testing.T) {
	units, _ := realUnitsAndPrefixes(t)
	if got := deplural("MILES", units); got != "MILES" {
		t.Errorf("deplural(MILES) = %q, want MILES unchanged (it is itself a defined unit)", got)
	}
	if got := deplural("INCHES", units); got != "INCH" {
		t.Errorf("deplural(INCHES) = %q, want INCH", got)
	}
}

func TestPprim_FormatsPositiveAndNegativeDimensions(t *testing.T) {
	// Force: LENGTH^1 * MASS^1 * TIME^-2
	got := pprim(primaryNames, [nPrim]int{1, 1, -2, 0, 0, 0, 0})
	want := " (LENGTH * MASS) / (TIME^2)"
	if got != want {
		t.Errorf("pprim(force dims) = %q, want %q", got, want)
	}
}

func TestPprim_NoNegativeDimensionsOmitsSecondGroup(t *testing.T) {
	got := pprim(primaryNames, [nPrim]int{2, 0, 0, 0, 0, 0, 0})
	want := " (LENGTH^2)"
	if got != want {
		t.Errorf("pprim(area dims) = %q, want %q", got, want)
	}
}

func TestEngprnt_ForcesExponentToMultipleOfThree(t *testing.T) {
	// The exponent field is always 3 digits wide (matching UNIT.TXT's
	// own worked example, "222.886230E-012") and always a multiple of
	// three: 12345678 = 12.345678E+006, not the "natural" +007.
	got := engprnt(12345678.0)
	want := " 12.345678E+006"
	if got != want {
		t.Errorf("engprnt(12345678) = %q, want %q", got, want)
	}
}

func TestEngprnt_MatchesUnitTXTWorkedExample(t *testing.T) {
	// UNIT.TXT's own worked example: 10 XX = 0.222886 YY = 222.886230E-012.
	got := engprnt(222.88623e-12)
	want := "222.886230E-012"
	if got != want {
		t.Errorf("engprnt(222.88623e-12) = %q, want %q", got, want)
	}
}

func TestEngprnt_PlainDecimalUnchanged(t *testing.T) {
	// A value %.6G renders without an exponent should pass through as-is.
	got := engprnt(88.0)
	if got != "88" {
		t.Errorf("engprnt(88) = %q, want 88", got)
	}
}

func TestApplyPrimaryUnit_BuildsCompoundUnit(t *testing.T) {
	// UNIT.TXT's own worked "xx = ft^2/mile^3" example: exponent 2 on
	// FT, then exponent -3 on MILE.
	units, _ := realUnitsAndPrefixes(t)
	ft := units[ufind("FT", units)]
	mile := units[ufind("MILE", units)]

	fact, pfact := 1.0, [nPrim]int{}
	fact, pfact = applyPrimaryUnit(fact, pfact, ft.Fact, ft.PFact, 2)
	fact, pfact = applyPrimaryUnit(fact, pfact, mile.Fact, mile.PFact, -3)

	wantFact := 4.48659e10
	if diff := math.Abs(fact-wantFact) / wantFact; diff > 1e-3 {
		t.Errorf("defined unit factor = %v, want ~%v (UNIT.TXT's own XX example)", fact, wantFact)
	}
	wantPFact := [nPrim]int{-1, 0, 0, 0, 0, 0, 0}
	if pfact != wantPFact {
		t.Errorf("defined unit dims = %v, want %v", pfact, wantPFact)
	}
}

func TestAlphabeticalOrder_IsCaseInsensitiveAndStable(t *testing.T) {
	units, _ := realUnitsAndPrefixes(t)
	order := alphabeticalOrder(units)
	for i := 1; i < len(order); i++ {
		a, b := units[order[i-1]].Name, units[order[i]].Name
		if compareFold(a, b) > 0 {
			t.Fatalf("units not sorted case-insensitively: %q before %q", a, b)
		}
	}
}

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
