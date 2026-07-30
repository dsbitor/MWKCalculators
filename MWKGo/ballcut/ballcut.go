// ballcut computes a schedule of incremental plunge cuts for
// roughing a spherical (or hemispherical) shape on a lathe with a
// squared-off cutoff tool: each cut is a step in a "staircase"
// approximation of the sphere, later filed smooth. Cuts are stepped
// either by a fixed angular increment around the quarter-circle
// profile or by a fixed increment along the lathe bed axis.
//
// Converted from BALLCUT.C, WorkshopUtilities/ballcut. The original
// writes its results to BALLCUT.OUT via fopen; this conversion
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

// ballCutStep is one plunge cut in the schedule: its axial position
// along the lathe bed (XF), the increment from the previous cut
// (DX), the depth of cut (YF), the increment in depth from the
// previous cut (DY), and the resulting work diameter (WD).
type ballCutStep struct {
	N      int
	XF, DX float64
	YF, DY float64
	WD     float64
}

// maxBallCutSteps bounds both cut schedules. The quarter-circle
// profile is fully covered well within this many steps for any
// realistic angular or linear increment; the bound exists to satisfy
// coding-style.md Rule 2 (bound every loop) against a zero or
// negative increment.
const maxBallCutSteps = 100_000

// ballCutScheduleByAngle returns the plunge cut schedule for a
// sphere of diameter sphereDiameter cut from stock of diameter
// stockDiameter, stepping by angleStepDeg degrees around the
// quarter-circle profile. Cuts stop once the required depth would
// exceed the available stock.
func ballCutScheduleByAngle(sphereDiameter, stockDiameter, angleStepDeg float64) []ballCutStep {
	r := 0.5 * sphereDiameter
	rs := 0.5 * stockDiameter

	var steps []ballCutStep
	var lastAxial, lastDepth float64
	for i, theta := 0, 0.0; i < maxBallCutSteps && theta <= 90.0; i, theta = i+1, theta+angleStepDeg {
		axial := r - r*math.Cos(theta*math.Pi/180)
		radial := r * math.Sin(theta*math.Pi/180)
		depth := rs - radial
		if depth < 0 {
			break
		}

		var dAxial, dDepth float64
		if i > 0 {
			dAxial = axial - lastAxial
			dDepth = depth - lastDepth
		}
		steps = append(steps, ballCutStep{N: i, XF: axial, DX: dAxial, YF: depth, DY: dDepth, WD: 2 * radial})
		lastAxial, lastDepth = axial, depth
	}
	return steps
}

// ballCutScheduleByAxialStep returns the plunge cut schedule for the
// same sphere and stock as ballCutScheduleByAngle, but stepping by a
// fixed axial distance axialStep along the lathe bed instead of a
// fixed angle.
func ballCutScheduleByAxialStep(sphereDiameter, stockDiameter, axialStep float64) []ballCutStep {
	r := 0.5 * sphereDiameter
	rs := 0.5 * stockDiameter

	var steps []ballCutStep
	var lastAxial, lastDepth float64
	for i, axial := 0, 0.0; i < maxBallCutSteps && axial <= r; i, axial = i+1, axial+axialStep {
		radial := math.Sqrt(r*r - (axial-r)*(axial-r))
		depth := rs - radial
		if depth < 0 {
			break
		}

		var dAxial, dDepth float64
		if i > 0 {
			dAxial = axial - lastAxial
			dDepth = depth - lastDepth
		}
		steps = append(steps, ballCutStep{N: i, XF: axial, DX: dAxial, YF: depth, DY: dDepth, WD: 2 * radial})
		lastAxial, lastDepth = axial, depth
	}
	return steps
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ballcut:", err)
		os.Exit(1)
	}

	fmt.Println("INCREMENTAL SPHERE SHAPING ON LATHE")
	fmt.Println()

	sphereDiameter := prompter.Float("Sphere diameter (in)", 1.0)
	stockDiameter := prompter.Float("Stock diameter (in)", sphereDiameter)
	byAxialStep := strings.ToLower(prompter.Line("Constant [A]ngle step or constant (X) step ? ")) == "x"

	var steps []ballCutStep
	if byAxialStep {
		axialStep := prompter.Float("X increment (in)", 0.02)
		steps = ballCutScheduleByAxialStep(sphereDiameter, stockDiameter, axialStep)
	} else {
		angleStep := prompter.Float("Angular increment (deg)", 5.0)
		steps = ballCutScheduleByAngle(sphereDiameter, stockDiameter, angleStep)
	}

	fmt.Println("\nIncremental Sphere Turning Data")
	fmt.Printf("Sphere diameter = %.4f in\n", sphereDiameter)
	fmt.Printf("Stock diameter = %.4f in\n", stockDiameter)
	fmt.Println()
	fmt.Println("N = cut number")
	fmt.Println("XF = axial (along lathe bed) position of tool")
	fmt.Println("DX = increment in x from last cut")
	fmt.Println("YF = depth of cut")
	fmt.Println("DY = increment in y from last cut")
	fmt.Println("WD = work diameter resulting from depth of cut YF")
	fmt.Println()
	fmt.Printf("%3s%9s%9s%9s%9s%9s\n\n", "N", "XF", "DX", "YF", "DY", "WD")
	for _, s := range steps {
		fmt.Printf("%3d%9.3f%+9.3f%9.3f%+9.3f%9.3f\n", s.N, s.XF, s.DX, s.YF, s.DY, s.WD)
	}
}
