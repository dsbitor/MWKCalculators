// dovetail computes the measurements across the top and bottom of a
// dovetail from a pin-based measurement: cylindrical pins are set
// into the angles of the dovetail and the distance across (male) or
// between (female) the pins is measured with calipers, giving an
// indirect but accurate way to check a dovetail angle and width
// without directly measuring the sloped surfaces.
//
// Converted from DOVETAIL.C (M. W. Klotz), WorkshopUtilities/dovetail.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

// dovetailMeasurements is the pin geometry shared by both the male
// and female formulas: the vertical offset each pin's contact point
// contributes to the top-to-bottom measurement (verticalOffset) and
// the horizontal correction the pin's radius contributes at the
// contact angle (radiusFactor).
func dovetailMeasurements(angleDeg, height float64) (verticalOffset, radiusFactor float64) {
	verticalOffset = height / math.Tan(angleDeg*math.Pi/180)
	halfAngleTan := math.Tan(0.5 * angleDeg * math.Pi / 180)
	radiusFactor = 1 + 1/halfAngleTan
	return verticalOffset, radiusFactor
}

// maleDovetail returns the bottom and top measurements of a male
// (external) dovetail, given the angle, pin diameter, dovetail
// height, and the measurement across the pins.
func maleDovetail(angleDeg, pinDiameter, height, acrossPins float64) (bottom, top float64) {
	verticalOffset, radiusFactor := dovetailMeasurements(angleDeg, height)
	bottom = acrossPins - pinDiameter*radiusFactor
	top = bottom + 2*verticalOffset
	return bottom, top
}

// femaleDovetail returns the bottom and top measurements of a
// female (internal) dovetail, given the angle, pin diameter,
// dovetail depth, and the measurement between the pins.
func femaleDovetail(angleDeg, pinDiameter, depth, betweenPins float64) (bottom, top float64) {
	verticalOffset, radiusFactor := dovetailMeasurements(angleDeg, depth)
	bottom = betweenPins + pinDiameter*radiusFactor
	top = bottom - 2*verticalOffset
	return bottom, top
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dovetail:", err)
		os.Exit(1)
	}

	fmt.Println("DOVETAIL MEASUREMENTS")
	fmt.Println()

	var male bool
	for {
		choice := strings.ToLower(prompter.Line("(M)ale or (F)emale dovetail? "))
		if choice == "m" {
			male = true
			break
		}
		if choice == "f" {
			male = false
			break
		}
		fmt.Println("error, try again")
	}

	angle := prompter.Float("Dovetail angle (deg)", 60.0)
	pinDiameter := prompter.Float("Pin diameter", 0.375)

	var height, measurement, bottom, top float64
	if male {
		height = prompter.Float("Height of dovetail", 0.5)
		measurement = prompter.Float("Measurement across pins", 2.5)
		bottom, top = maleDovetail(angle, pinDiameter, height, measurement)
	} else {
		height = prompter.Float("Depth of dovetail", 0.5)
		measurement = prompter.Float("Measurement between pins", 1.0283)
		bottom, top = femaleDovetail(angle, pinDiameter, height, measurement)
	}

	fmt.Println()
	fmt.Printf("Dovetail angle = %.4f deg\n", angle)
	fmt.Printf("Pin diameter = %.4f\n", pinDiameter)
	fmt.Printf("Measurement across/between pins = %.4f\n", measurement)
	fmt.Printf("Height/depth of dovetail = %.4f\n", height)
	fmt.Printf("Measurement across top of dovetail = %.4f\n", top)
	fmt.Printf("Measurement across bottom of dovetail = %.4f\n", bottom)
}
