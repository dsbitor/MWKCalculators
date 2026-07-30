// cmiter computes the compound table saw angles (miter gauge angle
// and blade tilt) needed to cut pieces for a mitred polygonal form
// with sloped sides, such as a tapered box or lampshade.
//
// Converted from CMITER.C (M. W. Klotz, 3/02),
// WorkshopUtilities/cmiter. Reference:
// http://www.betterwoodworking.com/compound_miter.htm
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// compoundMiterAngles is the compound miter setup for cutting one
// side of a sloped polygonal form.
type compoundMiterAngles struct {
	MiterGaugeAngleDeg float64
	BladeTiltMitredDeg float64
	BladeTiltButtedDeg float64
}

// computeCompoundMiterAngles returns the saw setup for a polygonal
// form of the given number of sides and side slope in degrees (0 for
// a flat-sided, vertical-walled form; 90 for a flat form cut in the
// table's plane).
func computeCompoundMiterAngles(sides int, slopeDeg float64) compoundMiterAngles {
	halfAngleDeg := 180 / float64(sides)

	switch slopeDeg {
	case 90:
		return compoundMiterAngles{MiterGaugeAngleDeg: 90, BladeTiltMitredDeg: halfAngleDeg, BladeTiltButtedDeg: 0}
	case 0:
		miterGauge := math.Atan(1/math.Tan(halfAngleDeg*math.Pi/180)) * 180 / math.Pi
		return compoundMiterAngles{MiterGaugeAngleDeg: miterGauge, BladeTiltMitredDeg: 0, BladeTiltButtedDeg: 90}
	}

	miterGauge := math.Atan(1/(math.Cos(slopeDeg*math.Pi/180)*math.Tan(halfAngleDeg*math.Pi/180))) * 180 / math.Pi
	cosMiterGauge := math.Cos(miterGauge * math.Pi / 180)
	tanSlope := math.Tan(slopeDeg * math.Pi / 180)

	return compoundMiterAngles{
		MiterGaugeAngleDeg: miterGauge,
		BladeTiltMitredDeg: math.Atan(cosMiterGauge*tanSlope) * 180 / math.Pi,
		BladeTiltButtedDeg: math.Atan(cosMiterGauge/tanSlope) * 180 / math.Pi,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cmiter:", err)
		os.Exit(1)
	}

	fmt.Println("TABLE SAW ANGLES FOR POLYGONAL FORMS")
	fmt.Println()

	sides := prompter.Int("Number of sides", 4)
	slope := prompter.Float("Slope (deg)", 30.0)

	fmt.Println()
	fmt.Printf("Number of sides = %d\n", sides)
	fmt.Printf("Slope = %.4f deg\n", slope)

	angles := computeCompoundMiterAngles(sides, slope)

	fmt.Printf("Miter gauge angle = %.4f deg\n", angles.MiterGaugeAngleDeg)
	fmt.Printf("Blade tilt (mitered) = %.4f deg\n", angles.BladeTiltMitredDeg)
	fmt.Printf("Blade tilt (butted) = %.4f deg\n", angles.BladeTiltButtedDeg)
}
