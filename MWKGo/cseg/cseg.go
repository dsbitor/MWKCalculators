// cseg solves for the properties of a circular segment (the region
// between a chord and the arc it cuts off) given any two of: radius,
// included angle, chord length, height (sagitta), or arc length.
//
// Converted from CSEG.C, Math/cseg.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// circularSegment holds the five related dimensions of a circular
// segment, plus its area.
type circularSegment struct {
	Radius, AngleDeg, Chord, Height, Arc, Area float64
}

// errInsufficientSegmentData reports that fewer than two of the five
// solvable quantities were given.
var errInsufficientSegmentData = errors.New("insufficient data for solution")

// maxRefinementSteps bounds the zoom-in grid search used when
// solving from chord+arc or height+arc, which have no closed form,
// matching the original program's own fixed 6-step refinement.
const maxRefinementSteps = 6

// refinementGridSteps is the number of angle samples per refinement
// step, matching the original program's own 1-degree-then-finer
// search.
const refinementGridSteps = 180

// searchAngleFor finds the angle (in degrees, 0-180) at which f
// (given that angle) most closely matches target, by repeatedly
// scanning a shrinking window at ever finer resolution, matching the
// original program's own search used for the chord+arc and
// height+arc cases.
func searchAngleFor(target float64, f func(angleDeg float64) float64) float64 {
	lower, upper, step := 1.0, 180.0, 1.0
	best := math.Inf(1)
	var bestAngle float64
	for i := 0; i < maxRefinementSteps; i++ {
		for angle := lower; angle <= upper; angle += step {
			if diff := math.Abs(f(angle) - target); diff < best {
				best = diff
				bestAngle = angle
			}
		}
		lower, upper = bestAngle-step, bestAngle+step
		step *= 0.1
	}
	return bestAngle
}

// computeFromRadiusAngle fills in the remaining segment dimensions
// given the radius and included angle.
func computeFromRadiusAngle(radius, angleDeg float64) circularSegment {
	halfAngleRad := 0.5 * angleDeg * math.Pi / 180
	height := radius - radius*math.Cos(halfAngleRad)
	chord := 2 * radius * math.Sin(halfAngleRad)
	arc := radius * angleDeg * math.Pi / 180
	return finishSegment(radius, angleDeg, chord, height, arc)
}

func finishSegment(radius, angleDeg, chord, height, arc float64) circularSegment {
	area := math.Pi*radius*radius*(angleDeg/360) - 0.5*(radius-height)*chord
	return circularSegment{Radius: radius, AngleDeg: angleDeg, Chord: chord, Height: height, Arc: arc, Area: area}
}

// solveCircularSegment solves for every circular segment dimension
// given any two of radius, angleDeg, chord, height, or arc (zero
// meaning unknown). It recognizes the same ten combinations as the
// original program.
func solveCircularSegment(radius, angleDeg, chord, height, arc float64) (circularSegment, error) {
	switch {
	case radius != 0 && angleDeg != 0:
		return computeFromRadiusAngle(radius, angleDeg), nil

	case radius != 0 && chord != 0:
		angleDeg = 2 * math.Asin(0.5*chord/radius) * 180 / math.Pi
		return computeFromRadiusAngle(radius, angleDeg), nil

	case radius != 0 && height != 0:
		angleDeg = 2 * math.Acos((radius-height)/radius) * 180 / math.Pi
		return computeFromRadiusAngle(radius, angleDeg), nil

	case radius != 0 && arc != 0:
		angleDeg = 180 / math.Pi * arc / radius
		return computeFromRadiusAngle(radius, angleDeg), nil

	case angleDeg != 0 && chord != 0:
		radius = 0.5 * chord / math.Sin(0.5*angleDeg*math.Pi/180)
		return computeFromRadiusAngle(radius, angleDeg), nil

	case angleDeg != 0 && height != 0:
		radius = height / (1 - math.Cos(0.5*angleDeg*math.Pi/180))
		return computeFromRadiusAngle(radius, angleDeg), nil

	case angleDeg != 0 && arc != 0:
		radius = 180 / math.Pi * arc / angleDeg
		return computeFromRadiusAngle(radius, angleDeg), nil

	case chord != 0 && height != 0:
		radius = (4*height*height + chord*chord) / (8 * height)
		angleDeg = 2 * math.Acos((radius-height)/radius) * 180 / math.Pi
		return computeFromRadiusAngle(radius, angleDeg), nil

	case chord != 0 && arc != 0:
		angleDeg = searchAngleFor(chord, func(a float64) float64 {
			r := 180 / math.Pi * arc / a
			return 2 * r * math.Sin(0.5*a*math.Pi/180)
		})
		radius = 180 / math.Pi * arc / angleDeg
		return computeFromRadiusAngle(radius, angleDeg), nil

	case height != 0 && arc != 0:
		angleDeg = searchAngleFor(height, func(a float64) float64 {
			r := 180 / math.Pi * arc / a
			return r * (1 - math.Cos(0.5*a*math.Pi/180))
		})
		radius = 180 / math.Pi * arc / angleDeg
		return computeFromRadiusAngle(radius, angleDeg), nil

	default:
		return circularSegment{}, errInsufficientSegmentData
	}
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cseg:", err)
		os.Exit(1)
	}

	fmt.Println("CIRCULAR SEGMENT CALCULATIONS")
	fmt.Println()
	fmt.Println("Input whatever data you know; press return if not known.")
	fmt.Println("You must enter two data items to obtain a solution.")
	fmt.Println()

	radius := prompter.Float("Radius of segment", 0)
	angleDeg := prompter.Float("Segment included angle (deg)", 0)

	known := 0
	for _, v := range []float64{radius, angleDeg} {
		if v != 0 {
			known++
		}
	}

	var chord, height, arc float64
	if known < 2 {
		chord = prompter.Float("Chord of segment", 0)
		if chord != 0 {
			known++
		}
	}
	if known < 2 {
		height = prompter.Float("Height of segment (sagitta)", 0)
		if height != 0 {
			known++
		}
	}
	if known < 2 {
		arc = prompter.Float("Arc length", 0)
		if arc != 0 {
			known++
		}
	}
	if known < 2 {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	seg, err := solveCircularSegment(radius, angleDeg, chord, height, arc)
	if err != nil {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	fmt.Printf("\nRadius           = %.4f\n", seg.Radius)
	fmt.Printf("Angle            = %.4f deg\n", seg.AngleDeg)
	fmt.Printf("Chord            = %.4f\n", seg.Chord)
	fmt.Printf("Height (sagitta) = %.4f\n", seg.Height)
	fmt.Printf("Arc Length       = %.4f\n", seg.Arc)
	fmt.Printf("Area             = %.4f\n", seg.Area)
}
