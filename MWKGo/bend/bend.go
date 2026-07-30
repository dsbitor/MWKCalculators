// bend computes sheet metal bend allowances: the length of material
// consumed forming a bend of a given angle and radius in material of
// a given thickness.
//
// Converted from BEND.C (M. W. Klotz, 11/98), WorkshopUtilities/bend.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// neutralAxisOffset returns the k-factor offset of the bend's
// neutral axis from its inner radius, as a fraction of thickness.
// The original program uses a three-tier lookup by radius-to-
// thickness ratio rather than a continuous k-factor curve: 1/3 for
// a tight bend (radius under twice the thickness), 1/2 for a gentle
// bend (radius over four times the thickness), and 0.4 in between.
func neutralAxisOffset(thickness, radius float64) float64 {
	switch {
	case radius < 2*thickness:
		return thickness / 3
	case radius > 4*thickness:
		return thickness / 2
	default:
		return 0.4 * thickness
	}
}

// bendLengths returns the exterior length (measured to the outer
// surface of the bend), the interior length (measured to the inner
// surface), and the bend allowance (the length of flat material
// consumed forming the bend, measured at the neutral axis) for a
// bend of angleDegrees in material of the given thickness and
// radius.
func bendLengths(thickness, radius, angleDegrees float64) (exterior, interior, allowance float64) {
	angleRadians := angleDegrees * math.Pi / 180
	offset := neutralAxisOffset(thickness, radius)

	exterior = angleRadians * (radius + thickness)
	interior = angleRadians * radius
	allowance = angleRadians * (radius + offset)
	return exterior, interior, allowance
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bend:", err)
		os.Exit(1)
	}

	fmt.Println("BEND ALLOWANCE COMPUTATION")
	fmt.Println()

	thickness := prompter.Float("Thickness of material (in)", 0.125)
	radius := prompter.Float("Radius of bend (in)", 3.0)
	angle := prompter.Float("Angle of bend (deg)", 180.0)

	exterior, interior, allowance := bendLengths(thickness, radius, angle)

	fmt.Printf("\nLength of bend exterior = %.4f in\n", exterior)
	fmt.Printf("Length of bend interior = %.4f in\n", interior)
	fmt.Printf("\nLength of material required to form bend = %.4f in\n", allowance)
}
