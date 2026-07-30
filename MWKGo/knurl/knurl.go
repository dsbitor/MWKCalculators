// knurl computes the workpiece diameter that produces a perfect,
// evenly spaced knurl pattern with a given knurl wheel.
//
// Converted from KNURL.C (M. W. Klotz, 11/98), WorkshopUtilities/knurl.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// knurlResult is the outcome of fitting a whole number of knurl
// crests around a workpiece close to its nominal diameter.
type knurlResult struct {
	ToothSpacing      float64 // circumferential spacing between knurl teeth
	CrestCount        int     // whole number of crests that fit
	Circumference     float64 // circumference needed for that many crests
	WorkpieceDiameter float64 // diameter corresponding to that circumference
}

// perfectKnurlFit returns the workpiece diameter closest to
// nominalDiameter that fits a whole number of teeth from a knurl
// wheel of the given diameter and tooth count. toothCount must be at
// least 1.
func perfectKnurlFit(wheelDiameter float64, toothCount int, nominalDiameter float64) knurlResult {
	toothSpacing := math.Pi * wheelDiameter / float64(toothCount)
	crestCount := int(math.Floor(math.Pi * nominalDiameter / toothSpacing))
	circumference := float64(crestCount) * toothSpacing

	return knurlResult{
		ToothSpacing:      toothSpacing,
		CrestCount:        crestCount,
		Circumference:     circumference,
		WorkpieceDiameter: circumference / math.Pi,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "knurl:", err)
		os.Exit(1)
	}

	fmt.Println("WORKPIECE DIAMETER FOR PERFECT KNURLING")
	fmt.Println()

	wheelDiameter := prompter.Float("Diameter of knurl wheel (in)", 0.625)
	toothCount := prompter.Int("Number of teeth on knurl wheel", 40)
	nominalDiameter := prompter.Float("Nominal diameter of workpiece (in)", 0.87)

	result := perfectKnurlFit(wheelDiameter, toothCount, nominalDiameter)

	fmt.Printf("Spacing between knurl teeth = %.3f in\n", result.ToothSpacing)
	fmt.Printf("Integer number of crests to make on workpiece = %d\n", result.CrestCount)
	fmt.Printf("Required workpiece circumference = %.3f in\n", result.Circumference)
	fmt.Printf("\nWORKPIECE DIAMETER FOR PERFECT KNURLING = %.3f in\n", result.WorkpieceDiameter)
}
