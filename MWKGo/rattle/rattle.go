// rattle implements Guy Lautard's technique (Home Machinist's
// Bedside Reader #1, pg. 11) for measuring a bore too large to span
// with available calipers: insert a stick slightly shorter than the
// bore diameter and measure how far it "rattles" back and forth,
// then compute the actual diameter from the stick length and rattle
// distance.
//
// Converted from RATTLE.C (M. W. Klotz, 11/99),
// WorkshopUtilities/rattle.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/mwktrig"
	"mwkgo/internal/promptio"
)

// boreDiameter returns the bore diameter and its difference from the
// stick length, given the stick length and the peak-to-peak rattle
// distance measured with the stick inserted in the bore.
func boreDiameter(stick, rattle float64) (diameter, diff float64) {
	theta := mwktrig.ClampedAsinDeg(0.5 * rattle / stick)
	beta := mwktrig.ClampedAsinDeg(rattle / stick)
	diameter = stick * math.Cos(theta*math.Pi/180) / (1 - 0.5*(1-math.Cos(beta*math.Pi/180)))
	diff = diameter - stick
	return diameter, diff
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rattle:", err)
		os.Exit(1)
	}

	fmt.Println("DIAMETER MEASUREMENT")
	fmt.Println()

	stick := prompter.Float("Measured stick/caliper distance", 4.0)
	rattle := prompter.Float("Rattle distance", 0.2)
	fmt.Println()

	diameter, diff := boreDiameter(stick, rattle)
	fmt.Printf("Diameter = %.6f\n", diameter)
	fmt.Printf("Diameter - stick = %.6f\n", diff)
}
