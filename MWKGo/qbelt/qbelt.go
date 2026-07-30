// qbelt computes the flat belt length for a simple two-pulley system,
// given both pulley diameters and the separation between their
// centers — a faster, no-data-file alternative to belt for the common
// two-pulley case.
//
// Converted from QBELT.C (M. W. Klotz), WorkshopUtilities/belt. That
// archive's own BELT.TXT describes this as belt's simpler,
// interactive-only companion for exactly this case; see
// MWKGo/belt/belt.go and docs/calculators/belt.md.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// twoPulleyWrap is the closed-form two-pulley external-tangent belt
// calculation QBELT.C, PULLEY.C, and PCD.C all share: given both
// pulley diameters and their center separation, it returns each
// pulley's own wrap angle (radians) and wrap length, the straight
// belt span between them, and the total belt length. The caller
// decides which pulley is "1" and which is "2" (QBELT itself sorts
// them smaller-first before calling; PULLEY and PCD assign the roles
// by which pulley is the known one, not by size).
func twoPulleyWrap(d1, d2, sep float64) (beta1, beta2, wrap1, wrap2, span, totalLength float64) {
	r1, r2 := 0.5*d1, 0.5*d2
	eps := r2 - r1
	theta := 2 * mwktrig.ClampedAsin(eps/sep)
	beta1 = math.Pi - theta
	beta2 = math.Pi + theta
	wrap1 = r1 * beta1
	wrap2 = r2 * beta2
	span = math.Sqrt(sep*sep - eps*eps)
	totalLength = 2*span + wrap1 + wrap2
	return beta1, beta2, wrap1, wrap2, span, totalLength
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "qbelt:", err)
		os.Exit(1)
	}

	fmt.Println("TWO PULLEY BELT LENGTH CALCULATION")
	fmt.Println()

	d1 := prompter.Float("Diameter of smaller pulley", 2)
	d2 := prompter.Float("Diameter of larger pulley", 6)
	sep := prompter.Float("Separation between pulley centers", 6)
	if d1 > d2 {
		d1, d2 = d2, d1
	}

	beta1, beta2, wrap1, wrap2, span, totalLength := twoPulleyWrap(d1, d2, sep)

	fmt.Println("\nFor smaller pulley:")
	fmt.Printf("Diameter = %.3f\n", d1)
	fmt.Printf("Wrap angle = %.2f deg\n", beta1*180/math.Pi)
	fmt.Printf("Wrap length = %.3f\n", wrap1)
	fmt.Println("\nFor larger pulley:")
	fmt.Printf("Diameter = %.3f\n", d2)
	fmt.Printf("Wrap angle = %.2f deg\n", beta2*180/math.Pi)
	fmt.Printf("Wrap length = %.3f\n", wrap2)
	fmt.Printf("\nBelt span between pulleys = %.3f\n", span)
	fmt.Printf("Belt length = %.3f\n", totalLength)
}
