// tubewall computes a tube's wall thickness from measurements made
// with an ordinary flat-anvil outside micrometer, correcting for the
// anvil bridging across the tube's curved surface rather than
// touching the true diameter.
//
// Converted from TUBEWALL.C (M. W. Klotz, 2/02),
// WorkshopUtilities/tubewall.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// wallThickness returns a tube's wall thickness, given the
// micrometer anvil's diameter, the tube's outside diameter, and the
// micrometer reading obtained measuring across the tube's bore with
// that anvil.
func wallThickness(anvilDiameter, tubeOutsideDiameter, micrometerReading float64) float64 {
	anvilRadius := 0.5 * anvilDiameter
	tubeOutsideRadius := 0.5 * tubeOutsideDiameter

	b := -tubeOutsideDiameter
	c := -anvilRadius*anvilRadius + 2*tubeOutsideRadius*micrometerReading - micrometerReading*micrometerReading

	return -0.5 * (b + math.Sqrt(b*b-4*c))
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tubewall:", err)
		os.Exit(1)
	}

	fmt.Println("MEASURING TUBE WALL THICKNESS WITH OUTSIDE MICROMETER")
	fmt.Println()

	anvilDiameter := prompter.Float("Micrometer anvil diameter", 0.249)
	tubeOutsideDiameter := prompter.Float("Tube outside diameter", 0.879)
	micrometerReading := prompter.Float("Micrometer measurement", 0.0625)

	thickness := wallThickness(anvilDiameter, tubeOutsideDiameter, micrometerReading)

	fmt.Printf("\nTube wall thickness = %.4f\n", thickness)
	fmt.Printf("Tube inside diameter = %.4f\n", tubeOutsideDiameter-2*thickness)
}
