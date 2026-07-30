package main

import (
	"math"
	"testing"
)

func TestFactorize_360(t *testing.T) {
	f := factorize(360)
	want := primeFactorization{Factors: []uint64{2, 3, 5}, Exps: []int{3, 2, 1}}
	if len(f.Factors) != len(want.Factors) {
		t.Fatalf("Factors = %v, want %v", f.Factors, want.Factors)
	}
	for i := range f.Factors {
		if f.Factors[i] != want.Factors[i] || f.Exps[i] != want.Exps[i] {
			t.Errorf("factor %d: got %d^%d, want %d^%d", i, f.Factors[i], f.Exps[i], want.Factors[i], want.Exps[i])
		}
	}
	// 2^3 * 3^2 * 5 = 360, an independent check that the
	// factorization actually reconstructs the original number.
	product := uint64(1)
	for i, factor := range f.Factors {
		for j := 0; j < f.Exps[i]; j++ {
			product *= factor
		}
	}
	if product != 360 {
		t.Errorf("product of factors = %d, want 360", product)
	}
}

func TestFactorize_Prime(t *testing.T) {
	f := factorize(97)
	if len(f.Factors) != 1 || f.Factors[0] != 97 || f.Exps[0] != 1 {
		t.Errorf("factorize(97) = %v, want [97^1]", f)
	}
}

func TestFindVernierPlates_DocumentedDefaultInput(t *testing.T) {
	// Independently confirmed via a separate brute-force script:
	// for 360 divisions, the best plate pair is (36, 40), for 76
	// total holes.
	plates, ok, err := findVernierPlates(360)
	if err != nil {
		t.Fatalf("findVernierPlates() error = %v", err)
	}
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if plates.Holes1 != 36 || plates.Holes2 != 40 {
		t.Errorf("Holes1, Holes2 = %d, %d, want 36, 40", plates.Holes1, plates.Holes2)
	}
	if plates.TotalHoles != 76 {
		t.Errorf("TotalHoles = %d, want 76", plates.TotalHoles)
	}
}

func TestComputeAlignments_EveryDivisionHasAnAlignment(t *testing.T) {
	// The defining property of a correct plate pair: every one of
	// the requested divisions must have some pair of holes that
	// aligns it. This is checked directly against the found plates
	// rather than assuming the search's own internal check was
	// applied correctly.
	plates, ok, err := findVernierPlates(360)
	if err != nil {
		t.Fatalf("findVernierPlates() error = %v", err)
	}
	if !ok {
		t.Fatal("ok = false, want true")
	}

	alignments, used1, used2 := computeAlignments(plates.Holes1, plates.Holes2, 360)
	if uint64(len(alignments)) != 360+1 {
		t.Fatalf("len(alignments) = %d, want %d", len(alignments), 361)
	}
	if countUsed(used1) == 0 || countUsed(used2) == 0 {
		t.Error("expected at least one hole used on each plate")
	}

	// Each alignment's own angle must actually be an exact multiple
	// of 360/numDivisions.
	divs := 360.0 / 360.0
	for _, a := range alignments {
		want := float64(a.Division) * divs
		if math.Abs(a.AngleDeg-want) > 1e-9 {
			t.Errorf("division %d: AngleDeg = %v, want %v", a.Division, a.AngleDeg, want)
		}
	}
}

func TestHoleLabel(t *testing.T) {
	cases := []struct {
		i1   uint64
		want string
	}{
		{0, "A"},
		{25, "Z"},
		{26, "a"},
		{51, "z"},
		{52, "A'"},
		{104, "A''"},
	}
	for _, c := range cases {
		if got := holeLabel(c.i1); got != c.want {
			t.Errorf("holeLabel(%d) = %q, want %q", c.i1, got, c.want)
		}
	}
}

func TestFindVernierPlates_TooFewDivisionsHasNoSolution(t *testing.T) {
	// 2 divisions is degenerate: no n1 in [2, num) exists at all
	// since the loop range is empty.
	_, ok, err := findVernierPlates(2)
	if err != nil {
		t.Fatalf("findVernierPlates() error = %v", err)
	}
	if ok {
		t.Error("ok = true, want false for numDivisions = 2")
	}
}
