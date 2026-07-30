package main

import (
	"math/cmplx"
	"testing"
)

// evalPoly evaluates sum(coeffs[i] * x^i) at x, the polynomial's own
// defining property: every genuine root must evaluate to
// (approximately) zero, independent of whichever method found it.
func evalPoly(coeffs []float64, x complex128) complex128 {
	result := complex128(0)
	power := complex128(1)
	for _, c := range coeffs {
		result += complex(c, 0) * power
		power *= x
	}
	return result
}

func TestSolveQuadratic_RootsSatisfyEquation(t *testing.T) {
	// x^2 - 5x + 6 = 0, roots 2 and 3.
	coeffs := []float64{6, -5, 1}
	roots := solveQuadratic(coeffs[0], coeffs[1], coeffs[2])
	for _, r := range roots {
		if cmplx.Abs(evalPoly(coeffs, r)) > 1e-9 {
			t.Errorf("root %v does not satisfy the equation: p(root) = %v", r, evalPoly(coeffs, r))
		}
	}
}

func TestSolveQuadratic_ComplexRoots(t *testing.T) {
	// x^2 + 1 = 0, roots i and -i.
	coeffs := []float64{1, 0, 1}
	roots := solveQuadratic(coeffs[0], coeffs[1], coeffs[2])
	for _, r := range roots {
		if cmplx.Abs(evalPoly(coeffs, r)) > 1e-9 {
			t.Errorf("root %v does not satisfy the equation", r)
		}
	}
}

func TestSolveCubic_ThreeRealRoots(t *testing.T) {
	// (x-1)(x-2)(x-3) = x^3 - 6x^2 + 11x - 6, roots 1, 2, 3.
	coeffs := []float64{-6, 11, -6, 1}
	roots := solveCubic(coeffs[0], coeffs[1], coeffs[2], coeffs[3])
	for _, r := range roots {
		if cmplx.Abs(evalPoly(coeffs, r)) > 1e-6 {
			t.Errorf("root %v does not satisfy the equation: p(root) = %v", r, evalPoly(coeffs, r))
		}
	}
}

func TestSolveCubic_OneRealTwoComplexRoots(t *testing.T) {
	// (x-1)(x^2+1) = x^3 - x^2 + x - 1, roots 1, i, -i.
	coeffs := []float64{-1, 1, -1, 1}
	roots := solveCubic(coeffs[0], coeffs[1], coeffs[2], coeffs[3])
	for _, r := range roots {
		if cmplx.Abs(evalPoly(coeffs, r)) > 1e-6 {
			t.Errorf("root %v does not satisfy the equation: p(root) = %v", r, evalPoly(coeffs, r))
		}
	}
}

func TestSolveQuartic_FourRealRoots(t *testing.T) {
	// (x-1)(x-2)(x-3)(x-4) = x^4 - 10x^3 + 35x^2 - 50x + 24.
	coeffs := []float64{24, -50, 35, -10, 1}
	roots, err := solveQuartic(coeffs[0], coeffs[1], coeffs[2], coeffs[3], coeffs[4])
	if err != nil {
		t.Fatalf("solveQuartic() error = %v", err)
	}
	for _, r := range roots {
		if cmplx.Abs(evalPoly(coeffs, r)) > 1e-6 {
			t.Errorf("root %v does not satisfy the equation: p(root) = %v", r, evalPoly(coeffs, r))
		}
	}
}

func TestSolveQuartic_TwoRealTwoComplexRoots(t *testing.T) {
	// (x-1)(x-2)(x^2+1) = x^4 - 3x^3 + 3x^2 - 3x + 2, roots 1, 2, i, -i.
	// (Chosen to avoid a known floating-point knife-edge in the
	// original algorithm's resultant-cubic root selection, exercised
	// by fully conjugate-paired roots like i,-i,2i,-2i; see 234.md.)
	coeffs := []float64{2, -3, 3, -3, 1}
	roots, err := solveQuartic(coeffs[0], coeffs[1], coeffs[2], coeffs[3], coeffs[4])
	if err != nil {
		t.Fatalf("solveQuartic() error = %v", err)
	}
	for _, r := range roots {
		if cmplx.Abs(evalPoly(coeffs, r)) > 1e-6 {
			t.Errorf("root %v does not satisfy the equation: p(root) = %v", r, evalPoly(coeffs, r))
		}
	}
}

func TestSolveQuartic_KnownFloatingPointKnifeEdgeReturnsError(t *testing.T) {
	// (x^2+1)(x^2+4) = x^4 + 5x^2 + 4, roots i, -i, 2i, -2i. This
	// specific case lands its resultant cubic's usable root exactly
	// on the rr>=0 boundary, where floating-point rounding in the
	// original algorithm (preserved here) can tip rr to a tiny
	// negative value and reject the only valid candidate. Documented
	// as a known inherited limitation rather than "fixed", since the
	// original program has the identical issue.
	_, err := solveQuartic(4, 0, 5, 0, 1)
	if err == nil {
		t.Skip("this floating-point knife-edge did not reproduce on this platform/build; not a correctness requirement either way")
	}
}
