// ftaper measures the half angle and included angle of a female
// taper (a tapered hole or socket) using two spheres of different
// diameters dropped into the taper and their depths measured.
//
// Converted from FTAPER.C (M. W. Klotz, 12/05, 3/07),
// WorkshopUtilities/ftaper.
package main

import (
	"fmt"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// taperAngle returns the taper's half angle in degrees and as a
// rise-over-run ratio (inches per inch), from the diameters of two
// spheres of different size and the measured depth to the top of
// each when dropped into the taper.
func taperAngle(smallDiameter, smallDepth, largeDiameter, largeDepth float64) (halfAngleDeg, ratio float64) {
	smallRadius := 0.5 * smallDiameter
	largeRadius := 0.5 * largeDiameter

	verticalSeparation := (smallDepth + smallRadius) - (largeDepth + largeRadius)
	ratio = (largeRadius - smallRadius) / verticalSeparation
	halfAngleDeg = mwktrig.ClampedAsinDeg(ratio)
	return halfAngleDeg, ratio
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ftaper:", err)
		os.Exit(1)
	}

	fmt.Println("MEASUREMENT OF FEMALE TAPERS WITH TWO SPHERES")
	fmt.Println()

	smallDiameter := prompter.Float("Diameter of small sphere", 0.5)
	smallDepth := prompter.Float("Depth to top of small sphere", 2.0)
	largeDiameter := prompter.Float("Diameter of large sphere", 0.75)
	largeDepth := prompter.Float("Depth to top of large sphere", 1.0)

	halfAngle, ratio := taperAngle(smallDiameter, smallDepth, largeDiameter, largeDepth)

	fmt.Printf("\nTaper half angle = %.4f deg = %.4f in/in\n", halfAngle, ratio)
	fmt.Printf("Taper included angle = %.4f deg = %.4f in/in\n", 2*halfAngle, 2*ratio)
}
