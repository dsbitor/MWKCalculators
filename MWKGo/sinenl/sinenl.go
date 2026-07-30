// sinenl computes the diameter of the smaller of two touching
// cylinders needed to set a given angle on a sine bar, using only
// two cylinders and no connecting link.
//
// Converted from SINENL.C (M. W. Klotz, 5/01), WorkshopUtilities/sine.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// smallerCylinderDiameter returns the diameter of the smaller
// cylinder that, together with a larger cylinder of largerDiameter
// touching it, sets angleDegrees on a sine bar with no connecting
// link between the two cylinders.
func smallerCylinderDiameter(largerDiameter, angleDegrees float64) float64 {
	halfAngleSin := math.Sin(0.5 * angleDegrees * math.Pi / 180)
	return largerDiameter * (1 - halfAngleSin) / (1 + halfAngleSin)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sinenl:", err)
		os.Exit(1)
	}

	fmt.Println("SINEBARS MADE WITH TWO CYLINDERS")
	fmt.Println()

	largerDiameter := prompter.Float("Diameter of larger cylinder (in)", 0.75)
	angle := prompter.Float("Desired angle (deg)", 1.5)

	fmt.Printf("\nDIAMETER OF SMALLER CYLINDER = %.4f in\n", smallerCylinderDiameter(largerDiameter, angle))
}
