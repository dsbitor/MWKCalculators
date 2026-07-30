// pcd (Pulley Center Distance) solves for the separation between two
// pulley centers, given both pulley diameters and the desired belt
// length — the complement of pulley, which instead solves for an
// unknown diameter.
//
// Converted from PCD.C (M. W. Klotz), WorkshopUtilities/belt. See
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
// duplicated here since PCD.C's own search calls it once per trial
// separation rather than as a separate step.
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

// result is one converged solution: the found separation, both
// pulleys' wrap angles and wrap lengths, and the belt span.
type result struct {
	Sep          float64
	Beta1, Beta2 float64
	Wrap1, Wrap2 float64
	Span         float64
}

// findCenterDistance searches for the pulley separation that makes
// the two-pulley belt length match targetLength within accuracy,
// matching PCD.C's own search exactly (the same coarse-scan-and-
// bracket-refine structure as PULLEY.C's own search for an unknown
// diameter, but searching separation instead): a coarse 10-step scan
// of separation between a tenth and ten times the target belt
// length, each pass narrowing to the bracket where the length-vs-
// target difference changes sign, repeating until within accuracy or
// the pass limit is reached.
func findCenterDistance(d1, d2, targetLength, accuracy float64) (result, error) {
	r1, r2 := 0.5*d1, 0.5*d2
	ss, sf := 0.1*targetLength, 10*targetLength
	diff := -1.0

	for pass := 0; pass < maxSearchPasses; pass++ {
		step := 0.1 * (sf - ss)
		var sl, sh float64
		bracketed := false

		for sep := ss; sep <= sf; sep += step {
			difflast := diff
			beta1, beta2, wrap1, wrap2, span, length := twoPulleyBeltLength(r1, r2, sep)
			diff = length - targetLength
			if math.Abs(diff) < accuracy {
				return result{Sep: sep, Beta1: beta1, Beta2: beta2, Wrap1: wrap1, Wrap2: wrap2, Span: span}, nil
			}
			if difflast < 0 && diff > 0 {
				sh = sep
				bracketed = true
				break
			}
			sl = sep
		}
		if !bracketed {
			return result{}, fmt.Errorf("search did not bracket a solution in range [%.6g, %.6g]", ss, sf)
		}
		ss, sf = sl, sh
	}
	return result{}, fmt.Errorf("search did not converge within %d passes", maxSearchPasses)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pcd:", err)
		os.Exit(1)
	}

	fmt.Println("TWO PULLEY CENTER DISTANCE CALCULATION")
	fmt.Println()

	d1 := prompter.Float("Diameter of driver pulley", 1.4)
	d2 := prompter.Float("Diameter of driven pulley", 0.603)
	targetLength := prompter.Float("Belt length", 8.21)
	accuracy := prompter.Float("Calculation accuracy", 0.0001)

	r, err := findCenterDistance(d1, d2, targetLength, accuracy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pcd:", err)
		os.Exit(1)
	}

	fmt.Println("\nFor known pulley:")
	fmt.Printf("Diameter = %.3f\n", d1)
	fmt.Printf("Wrap angle = %.2f deg\n", r.Beta1*180/math.Pi)
	fmt.Printf("Wrap length = %.3f\n", r.Wrap1)
	fmt.Println("\nFor calculated pulley:")
	fmt.Printf("Diameter = %.3f\n", d2)
	fmt.Printf("Wrap angle = %.2f deg\n", r.Beta2*180/math.Pi)
	fmt.Printf("Wrap length = %.3f\n", r.Wrap2)
	fmt.Printf("\nBelt span between pulleys = %.3f\n", r.Span)
	fmt.Printf("Belt length = %.3f\n", 2*r.Span+r.Wrap1+r.Wrap2)
	fmt.Printf("Pulley center distance = %.3f\n", r.Sep)
}
