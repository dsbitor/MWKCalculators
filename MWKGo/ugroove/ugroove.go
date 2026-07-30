// ugroove computes a schedule of incremental plunge cuts for milling
// a U-shaped groove (such as a giant O-ring groove for a pipe die)
// on a lathe using a round-nosed tool: successive plunge cuts, each
// offset from the groove's center, leave small cusps that are then
// filed smooth. Cuts are spaced either by a fixed angular increment
// around the groove's profile or by a fixed linear increment.
//
// Converted from UGROOVE.C, WorkshopUtilities/ugroove. The original
// writes its results to UGROOVE.OUT via fopen; this conversion
// prints to stdout instead and drops that DOS file-save-then-page
// convenience, the same approach used for loan (Tier 1 suitability
// review, Finding 5).
package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
)

// grooveCut is one plunge cut: its offset from the groove's center
// (X) and the depth of cut to make there (DOC).
type grooveCut struct {
	X, DOC float64
}

// maxGrooveCuts bounds both cut schedules, per coding-style.md
// Rule 2, against a zero or negative increment.
const maxGrooveCuts = 100_000

// effectiveRadius returns the radius of the ball-nosed tool
// center's path for a groove of radius grooveRadius cut with a tool
// of diameter toolDiameter.
func effectiveRadius(grooveRadius, toolDiameter float64) float64 {
	return grooveRadius - 0.5*toolDiameter
}

// grooveCutsByAngle returns the plunge cut schedule for a U-groove
// of radius grooveRadius cut with a ball-nosed tool of diameter
// toolDiameter, stepping by angleStepDeg degrees from the groove's
// centerline out to its edge.
func grooveCutsByAngle(grooveRadius, toolDiameter, angleStepDeg float64) []grooveCut {
	r := effectiveRadius(grooveRadius, toolDiameter)
	toolRadius := 0.5 * toolDiameter

	var cuts []grooveCut
	for i, theta := 0, 0.0; i < maxGrooveCuts && theta <= 90.0; i, theta = i+1, theta+angleStepDeg {
		x := r * math.Cos(theta*math.Pi/180)
		depth := math.Sqrt(r*r-x*x) + toolRadius
		cuts = append(cuts, grooveCut{X: x, DOC: depth})
	}
	return cuts
}

// grooveCutsByLinearStep returns the plunge cut schedule for the
// same groove and tool as grooveCutsByAngle, but stepping by a fixed
// linear distance linearStep from the groove's edge in toward its
// centerline.
func grooveCutsByLinearStep(grooveRadius, toolDiameter, linearStep float64) []grooveCut {
	r := effectiveRadius(grooveRadius, toolDiameter)
	toolRadius := 0.5 * toolDiameter

	var cuts []grooveCut
	for i, x := 0, r; i < maxGrooveCuts && x >= 0; i, x = i+1, x-linearStep {
		depth := math.Sqrt(r*r-x*x) + toolRadius
		cuts = append(cuts, grooveCut{X: x, DOC: depth})
	}
	return cuts
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ugroove:", err)
		os.Exit(1)
	}

	fmt.Println("INCREMENTAL CUTTING OF A U-SHAPED GROOVE")
	fmt.Println()

	grooveRadius := prompter.Float("Groove radius (in)", 1.0)
	toolDiameter := prompter.Float("Lathe tool diameter (in)", 0.25)
	byLinearStep := strings.ToLower(prompter.Line("[A]ngular or (L)inear increment ? ")) == "l"

	var cuts []grooveCut
	if byLinearStep {
		linearStep := prompter.Float("Linear increment (in)", 0.0625)
		cuts = grooveCutsByLinearStep(grooveRadius, toolDiameter, linearStep)
	} else {
		angleStep := prompter.Float("Angular increment (deg)", 5.0)
		cuts = grooveCutsByAngle(grooveRadius, toolDiameter, angleStep)
	}

	fmt.Printf("Groove radius = %.4f in\n", grooveRadius)
	fmt.Printf("Lathe tool diameter = %.4f in\n", toolDiameter)
	fmt.Println()
	fmt.Printf("%6s  %6s\n", "X", "DOC")
	for _, c := range cuts {
		fmt.Printf("%6.3f  %6.3f\n", c.X, c.DOC)
	}
}
