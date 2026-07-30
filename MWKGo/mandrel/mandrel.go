// mandrel computes the recommended mandrel diameter for winding a
// coil spring to a given inside diameter, following Kozo Hiraoka's
// empirical correction factor (Home Shop Machinist, July/August
// 1987) for the spring's springback after winding.
//
// Converted from MANDREL.C (M. W. Klotz, 2/05),
// WorkshopUtilities/mandrel.
package main

import (
	"fmt"
	"os"

	"mwkgo/internal/promptio"
)

// wireType selects which of Hiraoka's two empirical coefficient
// sets to use.
type wireType int

const (
	musicWire        wireType = 0
	phosphorusBronze wireType = 1
)

// coefficients holds Hiraoka's empirical constant and first-order
// coefficients, indexed by wireType.
var (
	constantCoefficient = [2]float64{0.980364, 0.012436}
	linearCoefficient   = [2]float64{-0.012436, -0.11018}
)

// mandrelDiameter returns the recommended mandrel diameter for
// winding a spring of the given wire diameter to the given inside
// diameter, using Hiraoka's empirical correction for the given wire
// type.
func mandrelDiameter(wt wireType, wireDiameter, springInsideDiameter float64) float64 {
	averageSpringDiameter := springInsideDiameter + wireDiameter
	factor := constantCoefficient[wt] + linearCoefficient[wt]*averageSpringDiameter/wireDiameter
	return factor*averageSpringDiameter - wireDiameter
}

// normalizeWireType treats any input other than 1 as music wire (0),
// matching the original program's own validation.
func normalizeWireType(input int) wireType {
	if input == int(phosphorusBronze) {
		return phosphorusBronze
	}
	return musicWire
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mandrel:", err)
		os.Exit(1)
	}

	fmt.Println("Kozo Hiraoka's SPRING WINDING MANDREL DIAMETER CALCULATION")
	fmt.Println()

	wt := normalizeWireType(prompter.Int("Wire type: music wire [0] or phosphorus bronze (1)", 0))

	wireDiameter := prompter.Float("Wire diameter (in)", 0.040)
	insideDiameter := prompter.Float("Spring inside diameter (in)", 0.203)

	fmt.Printf("\nRecommended mandrel diameter = %.3f in\n", mandrelDiameter(wt, wireDiameter, insideDiameter))
}
