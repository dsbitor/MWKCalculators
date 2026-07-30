package main

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestParseValue(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want float64
	}{
		{"plain decimal", "7", 7},
		{"leading dot", ".2", 0.2},
		{"simple fraction", "3/4", 0.75},
		{"mixed fraction", "1&1/2", 1.5},
		{"mixed fraction with larger numbers", "4&3/8", 4.375},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseValue(c.in)
			if err != nil {
				t.Fatalf("parseValue(%q) error = %v", c.in, err)
			}
			if math.Abs(got-c.want) > 1e-12 {
				t.Errorf("parseValue(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestParseValue_Errors(t *testing.T) {
	for _, in := range []string{"", "abc", "3/0"} {
		if _, err := parseValue(in); err == nil {
			t.Errorf("parseValue(%q) error = nil, want an error", in)
		}
	}
}

func TestGCDFloat(t *testing.T) {
	cases := []struct {
		x, y, want float64
	}{
		{12, 18, 6},
		{7, 13, 1},
		{64, 64, 64},
		{1, 64, 1},
	}
	for _, c := range cases {
		if got := gcdFloat(c.x, c.y); got != c.want {
			t.Errorf("gcdFloat(%v, %v) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
}

// TestParseDataEntry_WorkedExample reproduces the worked example from
// MIX.TXT: entering "+3.2in" then "+1.5cm" leaves the inch
// accumulator holding 3.2 + 1.5/2.54 inches (the independently
// computable sum of both entries converted to a common unit), and
// every other accumulator holding that same physical length in its
// own unit.
func TestApplyEntry_WorkedExample(t *testing.T) {
	state := newCalcState()

	op, value, unitIndex, err := parseDataEntry("+3.2in", state.defaultUnit)
	if err != nil {
		t.Fatalf("parseDataEntry(+3.2in) error = %v", err)
	}
	if op != opAdd || unitIndex != unitIN {
		t.Fatalf("parseDataEntry(+3.2in) = op %d, unit %d, want opAdd, unitIN", op, unitIndex)
	}
	if err := state.applyEntry(op, unitIndex, value); err != nil {
		t.Fatalf("applyEntry(+3.2in) error = %v", err)
	}

	op, value, unitIndex, err = parseDataEntry("+1.5cm", state.defaultUnit)
	if err != nil {
		t.Fatalf("parseDataEntry(+1.5cm) error = %v", err)
	}
	if err := state.applyEntry(op, unitIndex, value); err != nil {
		t.Fatalf("applyEntry(+1.5cm) error = %v", err)
	}

	wantInches := 3.2 + 1.5/2.54
	gotInches := state.acc[unitIN-1]
	if math.Abs(gotInches-wantInches) > 1e-9 {
		t.Errorf("inch accumulator = %v, want %v", gotInches, wantInches)
	}

	wantMeters := wantInches * unitPerInch[unitM]
	gotMeters := state.acc[unitM-1]
	if math.Abs(gotMeters-wantMeters) > 1e-12 {
		t.Errorf("meter accumulator = %v, want %v", gotMeters, wantMeters)
	}

	wantFeet := wantInches * unitPerInch[unitFT]
	gotFeet := state.acc[unitFT-1]
	if math.Abs(gotFeet-wantFeet) > 1e-12 {
		t.Errorf("foot accumulator = %v, want %v", gotFeet, wantFeet)
	}

	if state.exponent != 1 {
		t.Errorf("exponent = %d, want 1", state.exponent)
	}
}

// TestApplyEntry_AreaExponent reproduces the second half of the
// MIX.TXT worked example: squaring a length by multiplying it by
// itself should double the dimension exponent to 2 (area), and the
// resulting accumulator values must equal the length squared,
// converted to each unit's own squared-unit scale factor.
func TestApplyEntry_AreaExponent(t *testing.T) {
	state := newCalcState()
	state.defaultUnit = unitIN

	if _, err := enter(state, "5.25"); err != nil {
		t.Fatalf("enter(5.25) error = %v", err)
	}
	if _, err := enter(state, "*5&1/4"); err != nil {
		t.Fatalf("enter(*5&1/4) error = %v", err)
	}

	if state.exponent != 2 {
		t.Fatalf("exponent = %d, want 2 (area)", state.exponent)
	}

	wantSquareInches := 5.25 * 5.25
	gotSquareInches := state.acc[unitIN-1]
	if math.Abs(gotSquareInches-wantSquareInches) > 1e-9 {
		t.Errorf("square-inch accumulator = %v, want %v", gotSquareInches, wantSquareInches)
	}

	wantSquareMeters := wantSquareInches * unitPerInch[unitM] * unitPerInch[unitM]
	gotSquareMeters := state.acc[unitM-1]
	if math.Abs(gotSquareMeters-wantSquareMeters) > 1e-12 {
		t.Errorf("square-meter accumulator = %v, want %v", gotSquareMeters, wantSquareMeters)
	}
}

func enter(state *calcState, input string) (bool, error) {
	op, value, unitIndex, err := parseDataEntry(input, state.defaultUnit)
	if err != nil {
		return false, err
	}
	return true, state.applyEntry(op, unitIndex, value)
}

func TestApplyEntry_DivideByZero(t *testing.T) {
	state := newCalcState()
	if err := processEntry(state, "/0in"); err != errDivideByZero {
		t.Errorf("processEntry(/0in) error = %v, want errDivideByZero", err)
	}
}

func TestApplyEntry_UnknownUnit(t *testing.T) {
	state := newCalcState()
	if err := processEntry(state, "3xyz"); err != errUnknownUnit {
		t.Errorf("processEntry(3xyz) error = %v, want errUnknownUnit", err)
	}
}

func TestApplyEntry_RejectsNonDimensionalAddToLength(t *testing.T) {
	state := newCalcState()
	if err := processEntry(state, "3in"); err != nil {
		t.Fatalf("processEntry(3in) error = %v", err)
	}
	if err := processEntry(state, "+2nd"); err != errMixedNonDim {
		t.Errorf("processEntry(+2nd) error = %v, want errMixedNonDim", err)
	}
}

func TestApplyEntry_NonDimensionalMultiplyIsAllowed(t *testing.T) {
	state := newCalcState()
	if err := processEntry(state, "3in"); err != nil {
		t.Fatalf("processEntry(3in) error = %v", err)
	}
	if err := processEntry(state, "*2nd"); err != nil {
		t.Errorf("processEntry(*2nd) error = %v, want nil", err)
	}
	want := 6.0
	if got := state.acc[unitIN-1]; math.Abs(got-want) > 1e-9 {
		t.Errorf("inch accumulator = %v, want %v", got, want)
	}
}

func TestClearAndUndo(t *testing.T) {
	state := newCalcState()
	if err := processEntry(state, "3in"); err != nil {
		t.Fatalf("processEntry(3in) error = %v", err)
	}
	before := state.acc

	if err := processEntry(state, "5in"); err != nil {
		t.Fatalf("processEntry(5in) error = %v", err)
	}
	state.undo()
	if state.acc != before {
		t.Errorf("undo() left acc = %v, want %v", state.acc, before)
	}

	state.clear()
	for i, v := range state.acc {
		if v != 0 {
			t.Errorf("clear() left acc[%d] = %v, want 0", i, v)
		}
	}
	if state.exponent != 0 {
		t.Errorf("clear() left exponent = %d, want 0", state.exponent)
	}
}

// TestFractionalFeetInches_KnownConversion enters exactly 3 feet
// 9.78125 inches (3 ft 9 & 25/32 in) as a single inch value and
// checks the mixed-fraction display reduces to the same feet, inches,
// and 32nds independently computed here.
func TestFractionalFeetInches_KnownConversion(t *testing.T) {
	state := newCalcState()
	state.defaultUnit = unitIN
	totalInches := 3*12 + 9 + 25.0/32.0
	if err := processEntry(state, strconv.FormatFloat(totalInches, 'f', -1, 64)); err != nil {
		t.Fatalf("processEntry error = %v", err)
	}

	line, ok := state.fractionalFeetInches()
	if !ok {
		t.Fatal("fractionalFeetInches() ok = false, want true")
	}
	if !strings.Contains(line, "3 ft") || !strings.Contains(line, "9 & 25/32 in") {
		t.Errorf("fractionalFeetInches() = %q, want it to contain \"3 ft\" and \"9 & 25/32 in\"", line)
	}
}
