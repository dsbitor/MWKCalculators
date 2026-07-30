// protrac implements a sinebar-like protractor: two bars hinged
// together with two precision pins mounted at a fixed radius from
// the hinge. Measuring or setting the distance between the pins
// gives the angle between the bars with sinebar-like accuracy. Given
// the pin radius, pin diameter, and the pin separation when the
// protractor is fully closed, this program converts between pin
// separation and included angle in either direction.
//
// Converted from PROTRAC.C (M. W. Klotz), WorkshopUtilities/protrac.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// closedAngleDeg returns the angle between the pins when the
// protractor is fully closed: nonzero because the pins, being of
// finite diameter, cannot fully close the gap between the bars.
func closedAngleDeg(radius, pinDiameter, closedGap float64) float64 {
	return 2 * mwktrig.ClampedAsinDeg(0.5*(pinDiameter+closedGap)/radius)
}

// angleFromSeparation returns the included angle for a measured pin
// separation.
func angleFromSeparation(radius, pinDiameter, separation, closedAngle float64) float64 {
	return 2*mwktrig.ClampedAsinDeg(0.5*(pinDiameter+separation)/radius) - closedAngle
}

// separationFromAngle returns the pin separation for a given
// included angle.
func separationFromAngle(radius, pinDiameter, angle, closedAngle float64) float64 {
	halfAngleRad := 0.5 * (angle + closedAngle) * math.Pi / 180
	return 2*radius*math.Sin(halfAngleRad) - pinDiameter
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "protrac:", err)
		os.Exit(1)
	}

	fmt.Println("SINEBAR-LIKE PROTRACTOR")
	fmt.Println()

	radius := prompter.Float("Radius on which pins are mounted (in)", 3.0)
	pinDiameter := prompter.Float("Pin diameter (in)", 0.25)
	closedGap := prompter.Float("Pin separation when protractor is closed (in)", 0.22)

	closedAngle := closedAngleDeg(radius, pinDiameter, closedGap)
	fmt.Printf("\nAngle between pins when closed = %.4f deg\n", closedAngle)

	for {
		fmt.Println("\nA  find angle given pin separation")
		fmt.Println("D  find pin separation given angle")
		fmt.Println("Q  quit")
		choice := strings.ToLower(prompter.Line("(A,D,Q) ? "))

		var separation, angle float64
		switch choice {
		case "a":
			separation = prompter.Float("Pin separation (in)", 1.0)
			angle = angleFromSeparation(radius, pinDiameter, separation, closedAngle)
		case "d":
			angle = prompter.Float("Angle (deg)", 15.0)
			separation = separationFromAngle(radius, pinDiameter, angle, closedAngle)
		case "q":
			return
		default:
			fmt.Println("NOT A VALID OPTION")
			continue
		}
		fmt.Printf("Pin separation = %.4f in\n", separation)
		fmt.Printf("Angle = %.4f deg\n", angle)
	}
}
