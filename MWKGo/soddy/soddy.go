// soddy computes the diameter of the outer and inner Soddy circles
// for three mutually tangent circles of given diameters: the fourth
// circle tangent to all three, either enclosing them (outer) or
// nestled in the gap between them (inner).
//
// Converted from SODDY.C (M. W. Klotz, 3/02, 12/04),
// WorkshopUtilities/plug.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// soddyDiameter returns the diameter of the fourth circle tangent to
// three mutually tangent circles of diameters d1, d2, and d3. outer
// selects the circle that encloses all three; the non-outer solution
// is the circle nestled in the gap between them. This is Descartes'
// Circle Theorem, expressed in terms of diameters rather than the
// more usual curvatures.
func soddyDiameter(d1, d2, d3 float64, outer bool) float64 {
	sign := 1.0
	if !outer {
		sign = -1.0
	}
	radical := sign * 2 * math.Sqrt(d1*d2*d3*(d1+d2+d3))
	return math.Abs(d1 * d2 * d3 / (d2*d3 + d1*(d2+d3) - radical))
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "soddy:", err)
		os.Exit(1)
	}

	d1 := prompter.Float("Diameter of first circle", 0.245)
	d2 := prompter.Float("Diameter of second circle", 0.249)
	d3 := prompter.Float("Diameter of third circle", 0.253)

	fmt.Printf("Diameter of outer Soddy circle = %.7f\n", soddyDiameter(d1, d2, d3, true))
	fmt.Printf("Diameter of inner Soddy circle = %.7f\n", soddyDiameter(d1, d2, d3, false))
}
