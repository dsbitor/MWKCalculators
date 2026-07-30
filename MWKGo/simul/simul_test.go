package main

import (
	"math"
	"strings"
	"testing"
)

func TestGaussJordan_KnownSolution(t *testing.T) {
	// x + y = 3, x - y = 1  =>  x=2, y=1
	sys := linearSystem{
		a: [][]float64{{1, 1}, {1, -1}},
		b: []float64{3, 1},
	}
	if err := gaussJordan(sys); err != nil {
		t.Fatalf("gaussJordan() error = %v", err)
	}
	if math.Abs(sys.b[0]-2) > 1e-9 || math.Abs(sys.b[1]-1) > 1e-9 {
		t.Errorf("solution = %v, want [2, 1]", sys.b)
	}
}

// TestGaussJordan_TXTWorkedExample reproduces SIMUL.TXT's own worked
// example (4 equations, 4 unknowns) and is self-verifying: rather
// than asserting specific solution values, it substitutes the
// solution back into each original equation and checks the residual
// is negligible — the actual defining property of a correct
// solution.
func TestGaussJordan_TXTWorkedExample(t *testing.T) {
	original := [][]float64{
		{0, 3, -2, 7.5},
		{4.25, 5, 6, 3},
		{1, 3.75, 5, 7},
		{4, -3, 2.5, 5},
	}
	constants := []float64{34.25, 53.375, 59.875, 29.75}

	a := make([][]float64, 4)
	for i, row := range original {
		a[i] = append([]float64(nil), row...)
	}
	sys := linearSystem{a: a, b: append([]float64(nil), constants...)}

	if err := gaussJordan(sys); err != nil {
		t.Fatalf("gaussJordan() error = %v", err)
	}
	x := sys.b

	for i, row := range original {
		got := 0.0
		for j, coeff := range row {
			got += coeff * x[j]
		}
		if math.Abs(got-constants[i]) > 1e-6 {
			t.Errorf("equation %d: substituting solution gives %v, want %v", i, got, constants[i])
		}
	}
}

func TestGaussJordan_SingularMatrix(t *testing.T) {
	// x + y = 3, 2x + 2y = 6 (the second is not independent of the
	// first, exactly SIMUL.TXT's own example of a non-independent
	// system).
	sys := linearSystem{
		a: [][]float64{{1, 1}, {2, 2}},
		b: []float64{3, 6},
	}
	if err := gaussJordan(sys); err != errSingular {
		t.Errorf("gaussJordan() error = %v, want errSingular", err)
	}
}

func TestLoadSystem_TXTWorkedExample(t *testing.T) {
	data := `STARTOFDATA
4			;number of unknowns (n)

0    , 3    ,-2    , 7.5 , 34.25
4.25 , 5    , 6    , 3   , 53.375
1    , 3.75 , 5    , 7   , 59.875
4    ,-3    , 2.5  , 5   , 29.75
ENDOFDATA
`
	sys, err := loadSystem(strings.NewReader(data))
	if err != nil {
		t.Fatalf("loadSystem() error = %v", err)
	}
	if len(sys.a) != 4 {
		t.Fatalf("loadSystem() has %d equations, want 4", len(sys.a))
	}
	if sys.a[0][1] != 3 || sys.a[1][0] != 4.25 || sys.b[3] != 29.75 {
		t.Errorf("loadSystem() = %+v, does not match the expected worked example values", sys)
	}
}

func TestLoadSystem_WrongEquationCount(t *testing.T) {
	data := "STARTOFDATA\n2\n1,1,3\nENDOFDATA\n"
	if _, err := loadSystem(strings.NewReader(data)); err == nil {
		t.Error("loadSystem() error = nil, want an error when fewer equations are given than the declared unknown count")
	}
}
