// catenary computes the droop at the center of a hanging cable and
// its total length, given the cable's tension, its weight per unit
// length, and the straight-line distance between its supports.
//
// Converted from CATENARY.C (M. W. Klotz, 6/09), Math/catenary.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// catenaryShape returns the droop at the center of the span and the
// total cable length. density must not be zero; the original program
// does not guard against it either, since a weightless cable is
// outside the physical scope of this calculation and division by
// zero surfaces immediately as +Inf or NaN in the result rather than
// silently producing a plausible-looking wrong answer.
func catenaryShape(tension, density, distance float64) (droop, length float64) {
	param := tension / density
	halfSpan := 0.5 * distance
	droop = param * (math.Cosh(halfSpan/param) - 1)
	length = 2 * param * math.Sinh(0.5*distance/param)
	return droop, length
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "catenary:", err)
		os.Exit(1)
	}

	tension := prompter.Float("Tension in cable (lbf)", 0.674427074)
	density := prompter.Float("Cable weight per unit length (lb/ft)", 0.00044091)
	distance := prompter.Float("Straightline distance between supports (ft)", 1.640419948)

	droop, length := catenaryShape(tension, density, distance)

	fmt.Println()
	fmt.Printf("Cable droop at center of span = %.8g ft\n", droop)
	fmt.Printf("Cable length = %.10g ft\n", length)
}
