// links computes the taper angles and shim height needed to mill a
// tapered, radiused-end link: a flat strip holding two holes apart
// at a fixed distance, with each end radiused and the sides tapered
// to blend into the end radii.
//
// Converted from LINKS.C (M. W. Klotz, 10/00), WorkshopUtilities/links.
package main

import (
	"fmt"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// linkAngles returns the taper half angle (the angle each tapered
// side makes with the line joining the hole centers), the included
// angle at the small end, and the included angle at the big end, for
// a link whose ends have radii smallEndRadius and bigEndRadius,
// bigEndRadius the larger, holeCenterDistance apart.
func linkAngles(smallEndRadius, bigEndRadius, holeCenterDistance float64) (halfAngle, smallEndAngle, bigEndAngle float64) {
	halfAngle = mwktrig.ClampedAsinDeg((bigEndRadius - smallEndRadius) / holeCenterDistance)
	smallEndAngle = 180 - 2*halfAngle
	bigEndAngle = 180 + 2*halfAngle
	return halfAngle, smallEndAngle, bigEndAngle
}

// shimHeight returns the shim height needed under the small end pin
// to produce the correct taper angle when milling, from each end's
// radius and hole diameter.
func shimHeight(smallEndRadius, bigEndRadius, smallHoleDiameter, bigHoleDiameter float64) float64 {
	return bigEndRadius - smallEndRadius + 0.5*(bigHoleDiameter-smallHoleDiameter)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "links:", err)
		os.Exit(1)
	}

	fmt.Println("CALCULATIONS FOR TAPERED, RADIUSED-END LINKS")
	fmt.Println()

	smallEndRadius := prompter.Float("Small end radius (in)", 1.0/16)
	smallHoleDiameter := prompter.Float("Small end hole diameter (in)", 1.0/16)
	bigEndRadius := prompter.Float("Big end radius (in)", 3.0/32)
	bigHoleDiameter := prompter.Float("Big end hole diameter (in)", 3.0/16)
	centerDistance := prompter.Float("Distance between hole centers (in)", 1.0)

	halfAngle, smallEndAngle, bigEndAngle := linkAngles(smallEndRadius, bigEndRadius, centerDistance)
	height := shimHeight(smallEndRadius, bigEndRadius, smallHoleDiameter, bigHoleDiameter)

	fmt.Println()
	fmt.Printf("Angle of tapered side = %.2f deg\n", halfAngle)
	fmt.Printf("Included angle of small end = %.2f deg\n", smallEndAngle)
	fmt.Printf("Included angle of big end = %.2f\n", bigEndAngle)
	fmt.Printf("Small end shim height = %.4f in\n", height)
}
