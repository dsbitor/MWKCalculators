// chord computes the chord length for stepping off a circle into a
// given number of equal divisions.
//
// Converted from CHORD.C (M. W. Klotz, 2/00), WorkshopUtilities/chord.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// chordLength returns the straight-line distance between two adjacent
// points when a circle of the given diameter is divided into
// divisions equal parts. divisions and diameter are taken as given;
// a divisions of zero produces NaN through ordinary floating-point
// division, matching the unguarded arithmetic of the original
// program rather than introducing new validation for a case the
// original never checked either.
func chordLength(divisions, diameter float64) float64 {
	angleDegrees := 360.0 / divisions
	return diameter * math.Sin((angleDegrees/2)*math.Pi/180)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "chord:", err)
		os.Exit(1)
	}

	fmt.Println("CHORD LENGTH CALCULATION")
	fmt.Println()

	divisions := prompter.Float("Number of divisions", 5)
	diameter := prompter.Float("Diameter of circle", 1)

	fmt.Printf("\nChord length = %.4f\n", chordLength(divisions, diameter))
}
