// gear computes standard 20-degree full-depth involute spur gear
// dimensions for a mating pair of gears, given their tooth counts,
// diametral pitch, and pressure angle, per John A. Cooper's
// "Spur Gears and Pinions" (Machinist's Workshop, 4/99).
//
// Converted from GEAR.C, WorkshopUtilities/gear. The original writes
// its results to GEAR.OUT via fopen; this conversion prints to
// stdout instead and drops that DOS file-save-then-page convenience,
// the same approach used for loan (Tier 1 suitability review,
// Finding 5).
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// gearDimensions holds the standard dimensions of a single gear in a
// mating pair.
type gearDimensions struct {
	Teeth              float64
	OutsideDiameter    float64
	Addendum           float64
	Dedendum           float64
	WholeDepth         float64
	CircularPitch      float64
	ToothThickness     float64
	PitchDiameter      float64
	BaseCircleRadius   float64
	ToothProfileRadius float64
}

// gearPair holds the dimensions shared by both gears in a mating
// pair, plus each gear's own dimensions.
type gearPair struct {
	GearRatio        float64
	DiametralPitch   float64
	PressureAngleDeg float64
	CenterDistance   float64
	Gears            [2]gearDimensions
}

// computeGearPair returns the standard spur gear dimensions for two
// mating gears with teeth1 and teeth2 teeth, a diametral pitch of
// diametralPitch, and a pressure angle of pressureAngleDeg degrees.
func computeGearPair(teeth1, teeth2, diametralPitch, pressureAngleDeg float64) gearPair {
	teeth := [2]float64{teeth1, teeth2}
	pair := gearPair{
		GearRatio:        math.Max(teeth1, teeth2) / math.Min(teeth1, teeth2),
		DiametralPitch:   diametralPitch,
		PressureAngleDeg: pressureAngleDeg,
	}

	for i, n := range teeth {
		circularPitch := math.Pi / diametralPitch
		pitchDiameter := n / diametralPitch
		pair.Gears[i] = gearDimensions{
			Teeth:              n,
			OutsideDiameter:    (n + 2) / diametralPitch,
			Addendum:           1.0 / diametralPitch,
			Dedendum:           1.200 / diametralPitch,
			WholeDepth:         2.200 / diametralPitch,
			CircularPitch:      circularPitch,
			ToothThickness:     0.48 * circularPitch,
			PitchDiameter:      pitchDiameter,
			BaseCircleRadius:   0.5 * pitchDiameter * math.Cos(pressureAngleDeg*math.Pi/180),
			ToothProfileRadius: 0.5 * pitchDiameter * math.Sin(pressureAngleDeg*math.Pi/180),
		}
	}
	pair.CenterDistance = 0.5 * (pair.Gears[0].PitchDiameter + pair.Gears[1].PitchDiameter)

	return pair
}

const mmPerInch = 25.4

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gear:", err)
		os.Exit(1)
	}

	fmt.Println("GEAR CALCULATIONS")
	fmt.Println()

	teeth1 := prompter.Float("Number of teeth on gear 1", 45)
	teeth2 := prompter.Float("Number of teeth on gear 2", 20)
	diametralPitch := prompter.Float("Diametral Pitch (25.4/mod)", 20)
	pressureAngle := prompter.Float("Pressure Angle (deg)", 20)

	pair := computeGearPair(teeth1, teeth2, diametralPitch, pressureAngle)

	fmt.Println("\nFOR BOTH GEARS:")
	fmt.Printf("Gear Ratio = %.6f:1\n", pair.GearRatio)
	fmt.Printf("Diametral Pitch = %.4f\n", pair.DiametralPitch)
	fmt.Printf("Pressure Angle = %.4f deg\n", pair.PressureAngleDeg)
	fmt.Printf("Center Distance = %.4f in = %.4f mm\n", pair.CenterDistance, pair.CenterDistance*mmPerInch)

	for i, g := range pair.Gears {
		fmt.Printf("\nFOR GEAR %d:\n", i+1)
		fmt.Printf("Number of teeth = %.0f\n", g.Teeth)
		fmt.Printf("Outside Diameter = %.4f in = %.4f mm\n", g.OutsideDiameter, g.OutsideDiameter*mmPerInch)
		fmt.Printf("Addendum = %.4f in = %.4f mm\n", g.Addendum, g.Addendum*mmPerInch)
		fmt.Printf("Dedendum = %.4f in = %.4f mm\n", g.Dedendum, g.Dedendum*mmPerInch)
		fmt.Printf("Whole Depth = %.4f in = %.4f mm\n", g.WholeDepth, g.WholeDepth*mmPerInch)
		fmt.Printf("Circular Pitch = %.4f in = %.4f mm\n", g.CircularPitch, g.CircularPitch*mmPerInch)
		fmt.Printf("Tooth Thickness = %.4f in = %.4f mm\n", g.ToothThickness, g.ToothThickness*mmPerInch)
		fmt.Printf("Pitch Diameter = %.4f in = %.4f mm\n", g.PitchDiameter, g.PitchDiameter*mmPerInch)
		fmt.Printf("Base Circle Radius = %.4f in = %.4f mm\n", g.BaseCircleRadius, g.BaseCircleRadius*mmPerInch)
		fmt.Printf("Tooth Profile Radius = %.4f in = %.4f mm\n", g.ToothProfileRadius, g.ToothProfileRadius*mmPerInch)
	}
}
