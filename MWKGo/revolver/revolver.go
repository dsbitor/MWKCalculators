// revolver computes the dimensions needed to make a revolver
// cylinder style tool holder: a rotating cylinder with holes
// arranged around its axis, each holding a tool with its tip
// visible below.
//
// Converted from REVOLVER.C (M. W. Klotz, 8/00),
// WorkshopUtilities/revolver.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// cylinderDimensions returns the radius at which the tool holes
// should be centered, and the required outer diameter of the
// cylinder stock, for holeCount holes of holeDiameter arranged
// around the cylinder with edgeSpacing between adjacent holes and
// wallThickness of material required beyond the outer edge of each
// hole.
func cylinderDimensions(holeCount float64, holeDiameter, edgeSpacing, wallThickness float64) (holeRadius, cylinderDiameter float64) {
	anglePerHole := 360 / holeCount
	holeRadius = (edgeSpacing + holeDiameter) / (2 * math.Sin(0.5*anglePerHole*math.Pi/180))

	cylinderRadius := holeRadius + 0.5*holeDiameter + wallThickness
	cylinderDiameter = 2 * cylinderRadius
	return holeRadius, cylinderDiameter
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "revolver:", err)
		os.Exit(1)
	}

	fmt.Println("REVOLVER CYLINDER TOOL HOLDER CALCULATIONS")
	fmt.Println()

	holeCount := prompter.Int("Number of holes", 6)
	fmt.Println()
	fmt.Println("Hole size to fit tool body")
	holeDiameter := prompter.Float("Diameter of holes", 0.25)
	fmt.Println()
	fmt.Println("Remember to allow for diameter of tool top if larger than tool body")
	edgeSpacing := prompter.Float("Spacing between hole edges", 0.5)
	fmt.Println()
	fmt.Println("Holes cannot be tangent to cylinder circumference")
	wallThickness := prompter.Float("Thickness required at outer edge of holes", 0.25)

	holeRadius, cylinderDiameter := cylinderDimensions(float64(holeCount), holeDiameter, edgeSpacing, wallThickness)

	fmt.Printf("\nRadius for hole placement = %.3f in\n", holeRadius)
	fmt.Printf("Required stock diameter for cylinder = %.3f in\n", cylinderDiameter)
}
