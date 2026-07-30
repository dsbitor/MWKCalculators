// wire recommends an AWG wire gage for a given current and current
// density, and reports that gage's diameter, area, resistance, and
// weight per 1000 feet.
//
// Converted from WIRE.C (M. W. Klotz, 11/98), WorkshopUtilities/wire.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// gageStepRatio is the diameter ratio between adjacent AWG gage
// numbers: each step down in gage number multiplies the diameter by
// this factor. It is the standard AWG geometric progression ratio,
// 92^(1/39), matching the original program's constant.
const gageStepRatio = 1.12294049

// diameterMils returns the wire diameter, in thousandths of an inch,
// for the given AWG gage number. gage is not restricted to whole
// numbers; a fractional gage interpolates between the two nearest
// standard sizes.
func diameterMils(gage float64) float64 {
	return 324.87 / math.Pow(gageStepRatio, gage)
}

// gageForDiameter returns the AWG gage number, not necessarily a
// whole number, whose diameter is diameterMils. It is the inverse of
// diameterMils.
func gageForDiameter(diameterMils float64) float64 {
	return (math.Log10(0.32487) - math.Log10(1e-3*diameterMils)) / math.Log10(gageStepRatio)
}

// recommendedGage returns the nearest whole AWG gage number for a
// wire that must carry current amps at no more than currentDensity
// amps per circular mil, rounding to the nearest whole gage (ties
// round up, away from zero for a positive gage).
func recommendedGage(current, currentDensity float64) int {
	requiredDiameter := math.Sqrt(current / currentDensity)
	exact := gageForDiameter(requiredDiameter)
	return int(math.Floor(exact + 0.5))
}

// wireProperties is what a chosen gage tells you about the wire.
type wireProperties struct {
	DiameterMils            float64
	AreaCircularMils        float64
	ResistanceOhmsPer1000Ft float64
	WeightLbsPer1000Ft      float64
}

// propertiesForGage returns the physical properties of the given AWG
// gage of copper wire.
func propertiesForGage(gage int) wireProperties {
	d := diameterMils(float64(gage))
	area := d * d
	return wireProperties{
		DiameterMils:            d,
		AreaCircularMils:        area,
		ResistanceOhmsPer1000Ft: 10370 / area,
		WeightLbsPer1000Ft:      3.02675e-3 * area,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wire:", err)
		os.Exit(1)
	}

	fmt.Println("ELECTRICAL WIRE COMPUTATIONS")
	fmt.Println()

	current := prompter.Float("Current wire must carry (amps)", 12.0)
	density := prompter.Float("Desired current density (amps/cmil recommended)", 0.0025)

	gage := recommendedGage(current, density)
	props := propertiesForGage(gage)

	fmt.Printf("\nRecommended AWG gage = %d\n\n", gage)
	fmt.Printf("Diameter = %.1f mils\n", props.DiameterMils)
	fmt.Printf("Area = %.1f circular mils\n", props.AreaCircularMils)
	fmt.Printf("Resistance = %.3f ohms/1000 ft\n", props.ResistanceOhmsPer1000Ft)
	fmt.Printf("Weight = %.1f lbs/1000 ft\n", props.WeightLbsPer1000Ft)
}
