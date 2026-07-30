// ungula computes the volume of an ungula: the shape water takes up
// in a cylindrical glass tilted just enough that the water doesn't
// completely cover the bottom. Named for its resemblance to a cow's
// hoof (ungula in Latin).
//
// Converted from UNGULA.C (M. W. Klotz, 6/04), Math/ungula.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// ungulaVolume returns the volume of fluid in a cylinder of radius
// cylinderRadius, tilted so the wet part of the base is a circular
// segment of the given sagitta (the distance from the wall on the
// wet side to the dry/wet boundary chord, measured along a
// diameter), with fluid depth rising linearly from zero at that
// chord to height at the wall on the wet side.
//
// sagitta must be greater than zero and at most 2*cylinderRadius (the
// full diameter); the original program enforces this at the prompt
// rather than inside the calculation, and this conversion keeps that
// split. At sagitta == 2*cylinderRadius (the entire base is wet), the
// intermediate value a is zero, making (sagitta-cylinderRadius)/a
// evaluate to +Inf; math.Atan(+Inf) correctly returns pi/2 in that
// limit under IEEE 754, so the formula still gives the right answer
// without a special-cased branch. This was verified independently by
// numerical integration of the physical shape, not just by re-running
// this formula: see ungula_test.go and docs/calculators/ungula.md,
// which also corrects an error in the original UNGULA.TXT's own
// stated special case.
func ungulaVolume(height, sagitta, cylinderRadius float64) float64 {
	a := math.Sqrt(2*sagitta*cylinderRadius - sagitta*sagitta)
	phi := 0.5*math.Pi + math.Atan((sagitta-cylinderRadius)/a)
	sinPhi := math.Sin(phi)
	cosPhi := math.Cos(phi)

	numerator := 3*sinPhi - 3*phi*cosPhi - sinPhi*sinPhi*sinPhi
	return height * cylinderRadius * cylinderRadius * (numerator / (1 - cosPhi)) / 3
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ungula:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("COMPUTE VOLUME OF AN UNGULA")
	fmt.Println()

	cylinderDiameter := prompter.Float("Cylinder diameter", 2.0)
	cylinderRadius := 0.5 * cylinderDiameter
	height := prompter.Float("Height of ungula", 10.0)

	for {
		sagitta := prompter.Float("Sagitta of ungula base", 1.0)
		if sagitta <= 0 || sagitta > cylinderDiameter {
			fmt.Println()
			fmt.Println("0 < sagitta <= diameter of cylinder is required")
			continue
		}

		fmt.Printf("\nVolume of ungula = %.4f\n", ungulaVolume(height, sagitta, cylinderRadius))
		return
	}
}
