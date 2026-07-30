// dplate computes the diameter of disks needed to build a dividing
// plate: a set of equal cylindrical disks arranged around a central
// disk, each touching its neighbors and the center, giving a
// mechanical division into equal parts without a dividing head.
//
// Converted from DPLATE.C (M. W. Klotz, 1/02), WorkshopUtilities/dplate.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// diskDiameter returns the diameter of the disks needed to arrange
// divisions equal disks around a central disk of mountingDiameter,
// each disk touching both its neighbors and the center.
func diskDiameter(divisions float64, mountingDiameter float64) float64 {
	centerRadius := 0.5 * mountingDiameter
	halfAngleSin := math.Sin(0.5 * (360 / divisions) * math.Pi / 180)
	diskRadius := centerRadius * halfAngleSin / (1 - halfAngleSin)
	return 2 * diskRadius
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dplate:", err)
		os.Exit(1)
	}

	fmt.Println("DIVISION PLATE CALCULATION")
	fmt.Println()

	divisions := prompter.Int("Number of divisions", 14)
	mountingDiameter := prompter.Float("Diameter of mounting circle", 112)

	fmt.Printf("Disk diameter = %f\n", diskDiameter(float64(divisions), mountingDiameter))
}
