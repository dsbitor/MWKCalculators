// pulley solves for the diameter of the second pulley in a two-pulley
// belt system, given the first (known) pulley's diameter, the
// separation between the pulley centers, and the desired belt length.
//
// Converted from PULLEY.C (M. W. Klotz), WorkshopUtilities/belt. See
// MWKGo/belt/belt.go and docs/calculators/belt.md for the rest of
// this archive's programs.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// maxSearchPasses bounds the coarse-scan-and-bracket-refine search,
// replacing the original's interactive keypress abort (there is no
// keyboard to poll here) per coding-style.md Rule 2. Each pass
// narrows the bracket by 10x, so convergence in practice takes only a
// handful of passes; this is far larger than ever needed.
const maxSearchPasses = 100

// twoPulleyBeltLength is the closed-form two-pulley external-tangent
// belt length for known radii and separation, the same formula
// QBELT.C's own twoPulleyWrap computes (see MWKGo/qbelt/qbelt.go);
// duplicated here since PULLEY.C's own search calls it once per trial
// diameter rather than as a separate step.
func twoPulleyBeltLength(r1, r2, sep float64) (beta1, beta2, wrap1, wrap2, span, length float64) {
	eps := r2 - r1
	theta := 2 * mwktrig.ClampedAsin(eps/sep)
	beta1 = math.Pi - theta
	beta2 = math.Pi + theta
	wrap1 = r1 * beta1
	wrap2 = r2 * beta2
	span = math.Sqrt(sep*sep - eps*eps)
	length = 2*span + wrap1 + wrap2
	return beta1, beta2, wrap1, wrap2, span, length
}

// result is one converged solution: the found diameter, both
// pulleys' wrap angles and wrap lengths, and the belt span.
type result struct {
	D2           float64
	Beta1, Beta2 float64
	Wrap1, Wrap2 float64
	Span         float64
}

// findUnknownDiameter searches for the second pulley's diameter that
// makes the two-pulley belt length match targetLength within
// accuracy, matching PULLEY.C's own search exactly: a coarse 10-step
// scan of d2 between a tenth and ten times the known diameter, then
// each pass narrows to the bracket where the length-vs-target
// difference changes sign, repeating until within accuracy or the
// pass limit is reached.
func findUnknownDiameter(knownDiam, sep, targetLength, accuracy float64) (result, error) {
	r1 := 0.5 * knownDiam
	ds, df := 0.1*knownDiam, 10*knownDiam
	diff := -1.0

	for pass := 0; pass < maxSearchPasses; pass++ {
		step := 0.1 * (df - ds)
		var dl, dh float64
		bracketed := false

		for d2 := ds; d2 <= df; d2 += step {
			difflast := diff
			r2 := 0.5 * d2
			beta1, beta2, wrap1, wrap2, span, length := twoPulleyBeltLength(r1, r2, sep)
			diff = length - targetLength
			if math.Abs(diff) < accuracy {
				return result{D2: d2, Beta1: beta1, Beta2: beta2, Wrap1: wrap1, Wrap2: wrap2, Span: span}, nil
			}
			if difflast < 0 && diff > 0 {
				dh = d2
				bracketed = true
				break
			}
			dl = d2
		}
		if !bracketed {
			return result{}, fmt.Errorf("search did not bracket a solution in range [%.6g, %.6g]", ds, df)
		}
		ds, df = dl, dh
	}
	return result{}, fmt.Errorf("search did not converge within %d passes", maxSearchPasses)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pulley:", err)
		os.Exit(1)
	}

	fmt.Println("TWO PULLEY BELT LENGTH CALCULATION")
	fmt.Println()

	d1 := prompter.Float("Diameter of known pulley", 1.4)
	sep := prompter.Float("Separation between pulley centers", 2.5)
	targetLength := prompter.Float("Belt length", 8.21)
	accuracy := prompter.Float("Calculation accuracy", 0.001)

	r, err := findUnknownDiameter(d1, sep, targetLength, accuracy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pulley:", err)
		os.Exit(1)
	}

	fmt.Println("\nFor known pulley:")
	fmt.Printf("Diameter = %.3f\n", d1)
	fmt.Printf("Wrap angle = %.2f deg\n", r.Beta1*180/math.Pi)
	fmt.Printf("Wrap length = %.3f\n", r.Wrap1)
	fmt.Println("\nFor calculated pulley:")
	fmt.Printf("Diameter = %.3f\n", r.D2)
	fmt.Printf("Wrap angle = %.2f deg\n", r.Beta2*180/math.Pi)
	fmt.Printf("Wrap length = %.3f\n", r.Wrap2)
	fmt.Printf("\nBelt span between pulleys = %.3f\n", r.Span)
	fmt.Printf("Belt length = %.3f\n", 2*r.Span+r.Wrap1+r.Wrap2)
}
