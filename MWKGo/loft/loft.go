// loft computes the minimum length of thread engagement needed to
// guarantee that a screw shears in tension before its threads strip
// out of the hole, given the screw's basic diameter and thread
// pitch, assuming the screw and the tapped hole are the same
// material.
//
// Converted from LOFT.C (M. W. Klotz, 1/05), WorkshopUtilities/loft.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

// threadEngagement holds the computed screw thread engagement
// figures: PitchCircleDiameter and the two areas that TensileArea
// and ShearArea are compared through to reach EngagementLength.
type threadEngagement struct {
	PitchCircleDiameter float64
	TensileArea         float64
	ShearArea           float64
	EngagementLength    float64
}

// computeThreadEngagement returns the minimum thread engagement
// length for a screw of basic diameter diameter and thread pitch
// pitch (both in the same unit), per the standard approximation for
// 60 degree thread forms: the length at which the screw's tensile
// area and the internal thread's shear area balance.
func computeThreadEngagement(diameter, pitch float64) threadEngagement {
	pitchCircleDiameter := diameter - 0.64952*pitch
	tensileArea := 0.25 * math.Pi * (diameter - 0.938194*pitch) * (diameter - 0.938194*pitch)
	shearArea := 0.5 * math.Pi * pitchCircleDiameter
	engagementLength := 2 * tensileArea / shearArea

	return threadEngagement{
		PitchCircleDiameter: pitchCircleDiameter,
		TensileArea:         tensileArea,
		ShearArea:           shearArea,
		EngagementLength:    engagementLength,
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "loft:", err)
		os.Exit(1)
	}

	fmt.Println("MINIMUM LENGTH OF THREAD ENGAGEMENT")
	fmt.Println("TO ENSURE SCREW BREAKS BEFORE THREADS STRIP")
	fmt.Println()

	metric := strings.ToLower(prompter.Line("(m)etric or [i]mperial? ")) == "m"

	unitLength, unitArea := "in", "in^2"
	defaultDiameter := 0.25
	if metric {
		unitLength, unitArea = "mm", "mm^2"
		defaultDiameter = 4.0
	}

	diameter := prompter.Float(fmt.Sprintf("Basic diameter of screw (%s)", unitLength), defaultDiameter)

	var pitch float64
	if metric {
		pitch = prompter.Float("Pitch of screw (mm)", 0.7)
	} else {
		threadsPerInch := prompter.Float("Pitch of screw (tpi)", 20.0)
		pitch = 1 / threadsPerInch
	}

	e := computeThreadEngagement(diameter, pitch)

	fmt.Println()
	fmt.Printf("Pitch circle diameter of thread = %.4f %s\n", e.PitchCircleDiameter, unitArea)
	fmt.Printf("Screw thread tensile area = %.4f %s\n", e.TensileArea, unitArea)
	fmt.Printf("Thread shear area = %.4f %s\n", e.ShearArea, unitArea)
	fmt.Println()
	fmt.Printf("Length of thread engagement = %.4f %s\n", e.EngagementLength, unitLength)
	fmt.Printf("  or %.4f threads\n", e.EngagementLength/pitch)

	fmt.Println("\nNote:")
	fmt.Println("If J=(tensile strength of screw material)/(tensile strength of hole material)")
	fmt.Println("is greater than unity, multiply engagement length by J.")
	fmt.Println("[Calculations assume screw and hole are made of same material.]")
}
