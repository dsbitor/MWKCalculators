// sine computes the second cylinder diameter needed to build a
// custom sine bar from two cylinders connected by a pair of drilled
// links a fixed center-to-center distance apart, along with how
// sensitive the resulting angle is to small errors in either
// cylinder's diameter or the link spacing.
//
// Converted from SINE.C (M. W. Klotz, 11/98), WorkshopUtilities/sine.
// The later SINENL.C, converted separately as sinenl, removes the
// link entirely by using two cylinders in direct contact instead.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// secondCylinderDiameter returns the diameter needed for the second
// cylinder so that, held centerDistance apart from the first
// cylinder (of firstDiameter) by a pair of links, a flat plate laid
// across both cylinders sits at angleDegrees.
func secondCylinderDiameter(firstDiameter, centerDistance, angleDegrees float64) float64 {
	return 2*centerDistance*math.Sin(0.5*angleDegrees*math.Pi/180) + firstDiameter
}

// fits reports whether the two cylinders are small enough to be
// connected by links centerDistance apart: their diameters must not
// sum to more than twice the center distance, or they would collide
// with each other or the links.
func fits(firstDiameter, secondDiameter, centerDistance float64) bool {
	return 0.5*(firstDiameter+secondDiameter) <= centerDistance
}

// angleSensitivity returns how many degrees the resulting angle
// changes per inch of error in the first cylinder's diameter, and
// per inch of error in the center distance, for a sine bar already
// built to the given dimensions.
func angleSensitivity(firstDiameter, secondDiameter, centerDistance float64) (perDiameterError, perDistanceError float64) {
	halfDiameterDifferenceRatio := 0.5 * (secondDiameter - firstDiameter) / centerDistance
	slopeFactor := 1 / math.Sqrt(1-halfDiameterDifferenceRatio*halfDiameterDifferenceRatio)

	perDiameterError = 0.5 * slopeFactor / centerDistance
	perDistanceError = slopeFactor * math.Abs(halfDiameterDifferenceRatio) / centerDistance
	return perDiameterError, perDistanceError
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sine:", err)
		os.Exit(1)
	}

	fmt.Println("SINEBARS MADE WITH ACCURATELY SPACED CYLINDERS")
	fmt.Println()

	firstDiameter := prompter.Float("Diameter of first cylinder (d1) (in)", 0.375)
	centerDistance := prompter.Float("Distance between cylinder centers (r) (in)", 3.0)
	angle := prompter.Float("Desired angle (deg)", 10.0)

	secondDiameter := secondCylinderDiameter(firstDiameter, centerDistance, angle)

	if !fits(firstDiameter, secondDiameter, centerDistance) {
		fmt.Printf("\nSecond cylinder too large to fit with r = %.4f in\n", centerDistance)
		fmt.Println("Either use a smaller 'first' cylinder or larger r")
		return
	}

	perDiameterError, perDistanceError := angleSensitivity(firstDiameter, secondDiameter, centerDistance)
	degreesPerRadian := 180 / math.Pi

	fmt.Printf("\nDIAMETER OF SECOND CYLINDER (d2) = %.4f in\n\n", secondDiameter)
	fmt.Printf("Angle error for d1 error of 0.001 in = %.4f deg\n", 0.001*degreesPerRadian*perDiameterError)
	fmt.Printf("Angle error for  r error of 0.001 in = %.4f deg\n", 0.001*degreesPerRadian*perDistanceError)
}
