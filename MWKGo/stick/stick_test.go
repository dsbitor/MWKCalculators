package main

import (
	"math"
	"testing"
)

func TestNearInteger(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		tolerance float64
		want      bool
	}{
		{name: "exact integer", x: 5.0, tolerance: 0.000001, want: true},
		{name: "within tolerance above an integer", x: 5.0000001, tolerance: 0.000001, want: true},
		{name: "outside tolerance above an integer", x: 5.00001, tolerance: 0.000001, want: false},
		{name: "within tolerance below the next integer", x: 5.9999995, tolerance: 0.000001, want: true},
		{name: "clearly not near an integer", x: 5.5, tolerance: 0.000001, want: false},
		{name: "zero is an integer", x: 0, tolerance: 0.000001, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nearInteger(tt.x, tt.tolerance); got != tt.want {
				t.Errorf("nearInteger(%v, %v) = %v, want %v", tt.x, tt.tolerance, got, tt.want)
			}
		})
	}
}

func TestStickLength_ExactRatioConvergesImmediately(t *testing.T) {
	// A thread pitch exactly double the leadscrew pitch means the
	// very first candidate (one leadscrew pitch) already lands on a
	// whole number of thread pitches too.
	got, err := stickLength(4, 8, 100)
	if err != nil {
		t.Fatalf("stickLength() error = %v", err)
	}
	want := 1.0 / 4
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("stickLength(4, 8, 100) = %v, want %v", got, want)
	}
}

func TestStickLength_RequiresASecondStep(t *testing.T) {
	// leadscrewFactor=4 (step 0.25), threadFactor=6: the first
	// candidate (0.25) gives 6*0.25=1.5, not near an integer; the
	// second candidate (0.5) gives 6*0.5=3.0 exactly.
	got, err := stickLength(4, 6, 100)
	if err != nil {
		t.Fatalf("stickLength() error = %v", err)
	}
	want := 0.5
	if diff := math.Abs(got - want); diff > 1e-9 {
		t.Errorf("stickLength(4, 6, 100) = %v, want %v", got, want)
	}
}

func TestStickLength_UnreachableWithinLimitReturnsError(t *testing.T) {
	// An irrational thread factor essentially never lands near a
	// whole-number multiple for small step counts, so a small
	// iteration limit reliably exhausts without converging: the
	// contract here is a returned error, not an infinite loop or a
	// panic, matching the original program's own manual bailout.
	_, err := stickLength(1, math.Sqrt2, 10)
	if err == nil {
		t.Fatalf("stickLength() error = nil, want an error when the search limit is exhausted")
	}
}
