// osborne demonstrates the "Osborne Maneuver" (from Guy Lautard's
// Home Machinist's Bedside Reader #2) for centering round stock in a
// milling machine using only an edge finder: alternately finding the
// edge on the x and y axes and moving in by half the diameter each
// time converges rapidly on the true center.
//
// Converted from OSBORNE.C (M. W. Klotz, 11/99),
// WorkshopUtilities/osborne.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// osborneStep runs one iteration of the maneuver: given the current
// offset on one axis and the workpiece radius, it returns the
// resulting offset on the other axis and the combined radial error
// (the root-sum-square of both offsets).
func osborneStep(offset, radius float64) (nextOffset, radialError float64) {
	thetaDeg := mwktrig.ClampedAsinDeg(offset / radius)
	nextOffset = radius * (1 - math.Cos(thetaDeg*math.Pi/180))
	radialError = math.Hypot(offset, nextOffset)
	return nextOffset, radialError
}

const iterationCount = 6

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "osborne:", err)
		os.Exit(1)
	}

	fmt.Println("OSBORNE MANEUVER")
	fmt.Println()

	diameter := prompter.Float("Workpiece diameter", 2.0)
	offset := prompter.Float("Initial offset", 0.1)
	radius := 0.5 * diameter

	fmt.Println()
	for i := 1; i <= iterationCount; i++ {
		nextOffset, radialError := osborneStep(offset, radius)
		fmt.Printf("iteration: del1,del2,error= %d: %.8f, %.8f, %.8f\n", i, offset, nextOffset, radialError)
		offset = nextOffset
	}
}
