// mradius computes the radius of curvature of a part's edge, using
// two rollers of known diameter, gage blocks, and a surface plate.
//
// Converted from MRADIUS.C (M. W. Klotz, 3/02), WorkshopUtilities/mradius.
package main

import (
	"fmt"
	"os"

	"mwkgo/internal/promptio"
)

// radiusOfCurvature returns the radius of curvature of a part
// measured with two rollers of diameter rollerDiameter, resting
// against the part with their outer edges measurementAcross apart,
// separated at the bottom by a gage block stack of thickness gap.
func radiusOfCurvature(rollerDiameter, measurementAcross, gap float64) float64 {
	rollerRadius := 0.5 * rollerDiameter
	halfWidth := 0.5*measurementAcross - rollerRadius
	adjustedRadius := rollerRadius - gap

	return (adjustedRadius*adjustedRadius + halfWidth*halfWidth - rollerRadius*rollerRadius) / (2 * (rollerDiameter - gap))
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mradius:", err)
		os.Exit(1)
	}

	rollerDiameter := prompter.Float("Diameter of rollers (d)", 0.25)
	measurementAcross := prompter.Float("Measurement across rollers (m)", 2.0)
	gap := prompter.Float("Gap (g)", 0.1)

	fmt.Printf("\nRadius of curvature = %.4f\n", radiusOfCurvature(rollerDiameter, measurementAcross, gap))
}
