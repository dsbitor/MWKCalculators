// plate computes a chain-drilling layout for cutting a circular
// metal plate free from flat stock: a ring of holes drilled just
// outside the plate's final diameter, leaving a scalloped edge that
// finish machining cleans up.
//
// Converted from PLATE.C (M. W. Klotz, 6/03), WorkshopUtilities/slug.
package main

import (
	"fmt"
	"os"

	"mwkgo/internal/chaindrill"
	"mwkgo/internal/promptio"
)

// plan wires plate's specific direction, outward (the as-drilled
// piece is larger than the finished plate, since finish machining
// removes the scalloped edge), into the shared chaindrill
// calculation.
func plan(finalDiameter, radialAllowance, drillDiameter, webWanted float64) chaindrill.Plan {
	return chaindrill.Compute(finalDiameter, radialAllowance, drillDiameter, webWanted, true)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "plate:", err)
		os.Exit(1)
	}

	fmt.Println("CHAIN DRILLING A CIRCULAR PLATE")
	fmt.Println()

	finalDiameter := prompter.Float("Diameter of final plate (in)", 3.0)
	radialAllowance := prompter.Float("Radial allowance for finish machining (in)", 0.050)
	drillDiameter := prompter.Float("Drill diameter (in)", 0.25)
	webWanted := prompter.Float("Approximate web thickness (in)", 0.050)

	result := plan(finalDiameter, radialAllowance, drillDiameter, webWanted)

	fmt.Println()
	fmt.Printf("Number of holes = %d\n", result.HoleCount)
	fmt.Printf("Diameter of drilling circle = %.3f in\n", result.DrillingCircleDiameter)
	fmt.Printf("Final web thickness = %.3f in\n", result.WebThickness)
	fmt.Printf("Angle between adjacent holes = %.3f deg\n", result.AngleBetweenHolesDeg)
	fmt.Printf("Chordal distance between adjacent holes = %.3f in\n", result.ChordBetweenHoles)
}
