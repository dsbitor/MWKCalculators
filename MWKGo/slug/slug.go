// slug computes a chain-drilling layout for opening a large hole in
// plate stock by drilling and freeing a "slug": a ring of holes
// drilled just inside the hole's final diameter, leaving a scalloped
// edge that finish machining opens out to size.
//
// Converted from SLUG.C (M. W. Klotz, 10/00), WorkshopUtilities/slug.
package main

import (
	"fmt"
	"os"

	"mwkgo/internal/chaindrill"
	"mwkgo/internal/promptio"
)

// plan wires slug's specific direction, inward (the as-drilled
// circle is smaller than the finished hole, since finish machining
// enlarges it out to size), into the shared chaindrill calculation.
func plan(finalDiameter, radialAllowance, drillDiameter, webWanted float64) chaindrill.Plan {
	return chaindrill.Compute(finalDiameter, radialAllowance, drillDiameter, webWanted, false)
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "slug:", err)
		os.Exit(1)
	}

	fmt.Println("CHAIN DRILLING LARGE HOLES")
	fmt.Println()

	finalDiameter := prompter.Float("Diameter of final hole (in)", 3.0)
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
