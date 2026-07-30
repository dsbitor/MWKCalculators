package main

import "testing"

func TestSolveDiophantine_DocumentedDefaultInput(t *testing.T) {
	// A linear Diophantine equation's own definition is the
	// strongest possible test: whatever x and y the solver returns,
	// a*x + b*y must equal c exactly. This is checked directly
	// rather than comparing against a specific hand-picked (x, y)
	// pair, since any of infinitely many solutions is equally
	// correct.
	a, b, c := int64(172), int64(20), int64(1000)
	x, y, ok := solveDiophantine(a, b, c)
	if !ok {
		t.Fatal("solveDiophantine() ok = false, want true")
	}
	if got := a*x + b*y; got != c {
		t.Errorf("%d*%d + %d*%d = %d, want %d", a, x, b, y, got, c)
	}
}

func TestSolveDiophantine_UnsolvableWhenNotDivisibleByGCD(t *testing.T) {
	// gcd(4, 6) = 2, so no integer x, y can satisfy 4x + 6y = 7 (an
	// odd number can never result from a sum of even multiples).
	_, _, ok := solveDiophantine(4, 6, 7)
	if ok {
		t.Error("solveDiophantine(4, 6, 7) ok = true, want false")
	}
}

func TestSolveDiophantine_EverySampleSolutionSatisfiesEquation(t *testing.T) {
	// The general solution family (x + k*b/g, y - k*a/g) must
	// satisfy the original equation for every integer k, not just
	// k=0: the same self-verifying check applied across the whole
	// family the program prints as "sample solutions".
	a, b, c := int64(172), int64(20), int64(1000)
	x, y, ok := solveDiophantine(a, b, c)
	if !ok {
		t.Fatal("solveDiophantine() ok = false, want true")
	}
	g := gcd(a, b)
	for k := int64(-4); k <= 4; k++ {
		xp := x + k*b/g
		yp := y - k*a/g
		if got := a*xp + b*yp; got != c {
			t.Errorf("k=%d: %d*%d + %d*%d = %d, want %d", k, a, xp, b, yp, got, c)
		}
	}
}

func TestGCD(t *testing.T) {
	cases := []struct{ a, b, want int64 }{
		{172, 20, 4},
		{17, 5, 1},
		{0, 0, 1},
		{0, 5, 5},
	}
	for _, c := range cases {
		if got := gcd(c.a, c.b); got != c.want {
			t.Errorf("gcd(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}
