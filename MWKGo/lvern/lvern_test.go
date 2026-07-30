package main

import (
	"math"
	"testing"
)

func TestVernierMainScale_FirstWorkedExample(t *testing.T) {
	// LVERN.TXT's first worked example: an inch main scale divided
	// into eighths, designed for a vernier resolving thirty-seconds.
	mainDivisionLength, vernierSubdivisions, exact := vernierMainScale(1.0, 8, 1.0/32.0)
	if !exact {
		t.Fatal("exact = false, want true")
	}
	if math.Abs(mainDivisionLength-0.125) > 1e-9 {
		t.Errorf("mainDivisionLength = %v, want 0.125", mainDivisionLength)
	}
	if vernierSubdivisions != 4 {
		t.Errorf("vernierSubdivisions = %d, want 4", vernierSubdivisions)
	}

	typicalSpan := 1.0 - mainDivisionLength
	if math.Abs(typicalSpan-0.875) > 1e-9 {
		t.Errorf("typicalSpan = %v, want 0.875", typicalSpan)
	}

	vernierDivision := vernierDivisionLength(typicalSpan, vernierSubdivisions)
	if math.Abs(vernierDivision-0.21875) > 1e-9 {
		t.Errorf("vernierDivision = %v, want 0.21875 (prints as 0.2188)", vernierDivision)
	}
}

func TestVernierMainScale_SecondWorkedExample(t *testing.T) {
	// LVERN.TXT's second worked example: a main scale divided into
	// tenths, designed for a vernier resolving to 0.01.
	mainDivisionLength, vernierSubdivisions, exact := vernierMainScale(1.0, 10, 0.01)
	if !exact {
		t.Fatal("exact = false, want true")
	}
	if math.Abs(mainDivisionLength-0.1) > 1e-9 {
		t.Errorf("mainDivisionLength = %v, want 0.1", mainDivisionLength)
	}
	if vernierSubdivisions != 10 {
		t.Errorf("vernierSubdivisions = %d, want 10", vernierSubdivisions)
	}

	typicalSpan := 1.0 - mainDivisionLength
	vernierDivision := vernierDivisionLength(typicalSpan, vernierSubdivisions)
	if math.Abs(vernierDivision-0.09) > 1e-9 {
		t.Errorf("vernierDivision = %v, want 0.09", vernierDivision)
	}
}

func TestNearestRationalFraction_WorkedExampleValues(t *testing.T) {
	cases := []struct {
		q                           float64
		wantWhole, wantNum, wantDen int64
	}{
		{1.0 / 32.0, 0, 1, 32},
		{0.125, 0, 1, 8},
		{0.01, 0, 1, 100},
		{0.1, 0, 1, 10},
		{0.09, 0, 9, 100},
	}
	for _, c := range cases {
		whole, num, den := nearestRationalFraction(c.q)
		if whole != c.wantWhole || num != c.wantNum || den != c.wantDen {
			t.Errorf("nearestRationalFraction(%v) = %d, %d, %d, want %d, %d, %d", c.q, whole, num, den, c.wantWhole, c.wantNum, c.wantDen)
		}
	}
}

func TestFormatNearestFraction(t *testing.T) {
	if got := formatNearestFraction(1.0 / 32.0); got != "1/32" {
		t.Errorf("formatNearestFraction(1/32) = %q, want %q", got, "1/32")
	}
	if got := formatNearestFraction(1.5); got != "1-1/2" {
		t.Errorf("formatNearestFraction(1.5) = %q, want %q", got, "1-1/2")
	}
}

func TestParseDecimalOrFraction(t *testing.T) {
	cases := []struct {
		s    string
		want float64
	}{
		{"1.5", 1.5},
		{"1/2", 0.5},
		{"1-1/2", 1.5},
		{".01", 0.01},
	}
	for _, c := range cases {
		got, err := parseDecimalOrFraction(c.s)
		if err != nil {
			t.Fatalf("parseDecimalOrFraction(%q) error = %v", c.s, err)
		}
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("parseDecimalOrFraction(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}
