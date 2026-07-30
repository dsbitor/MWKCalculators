// flywheel computes the drilling and milling layout for a tapered
// (or straight) spoke flywheel machined from solid stock: holes
// drilled through a reduced-thickness web define the corners of
// cutouts that, once milled away, leave the spokes. The rotary table
// angle needed to mill along a spoke's edge is found by solving a
// transcendental equation (the angle at which the inner and outer
// hole edges align) by iterative search.
//
// Converted from FLYWHEEL.C, WorkshopUtilities/flywheel. The
// original writes its results to FLYWHEEL.OUT via fopen; this
// conversion prints to stdout instead and drops that DOS
// file-save-then-page convenience, the same approach used for loan
// (Tier 1 suitability review, Finding 5).
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

func sind(deg float64) float64 { return math.Sin(deg * math.Pi / 180) }
func asind(x float64) float64  { return math.Asin(x) * 180 / math.Pi }

// maxPhiRefinements bounds the zoom-in grid search for phi, matching
// the same technique used in cseg's angle search. Each refinement
// narrows the search window and its step size tenfold; the original
// program relies on this converging (it has no iteration cap of its
// own here), so a generous bound is added per coding-style.md Rule 2.
const maxPhiRefinements = 20

// searchPhi finds the angle phi minimizing |f(phi)| for
// f(p) = r1*sind(theta1+p) - s1 - (r2*sind(theta2+p) - s2), the
// transcendental equation for the rotary table angle that aligns the
// inner and outer hole edges along a spoke, by repeatedly scanning a
// shrinking window at ever finer resolution.
func searchPhi(r1, s1, theta1, r2, s2, theta2 float64) float64 {
	f := func(p float64) float64 {
		return r1*sind(theta1+p) - s1 - (r2*sind(theta2+p) - s2)
	}

	phiLast, lower, upper, step := 0.0, 0.0, 20.0, 1.0
	var phi float64
	for i := 0; i < maxPhiRefinements; i++ {
		best := math.Inf(1)
		for p := lower; p <= upper; p += step {
			if q := math.Abs(f(p)); q < best {
				best, phi = q, p
			}
		}
		if math.Abs(phi-phiLast) <= 0.01 {
			break
		}
		phiLast, lower, upper, step = phi, phi-step, phi+step, step/10
	}
	return phi
}

// flywheelLayout holds the computed drilling/milling layout for one
// solution (either the exact solution or the nearest-integral-angle
// variant).
type flywheelLayout struct {
	R1, Offset, Theta1, Theta2, Phi float64
	X, Sep, Spoke1, Spoke2, W       float64
}

var errFlywheelNoTaper = errors.New("must specify either the outer hole offset or its offset angle")

// computeFlywheelLayout computes the flywheel drilling layout for
// nspoke spokes, inner/outer hole radii r1/r2, inner/outer hole
// diameters d1/d2, and the spoke taper specified as either
// offsetIn (outer hole offset from spoke centerline) or theta2In
// (outer hole angle from spoke centerline) — exactly one of the two
// must be nonzero. It returns the exact solution and the
// nearest-integral-outer-angle variant.
func computeFlywheelLayout(nspoke int, r1, r2, d1, d2, offsetIn, theta2In float64) (exact, nearestIntegral flywheelLayout, err error) {
	if offsetIn == 0 && theta2In == 0 {
		return flywheelLayout{}, flywheelLayout{}, errFlywheelNoTaper
	}

	s1, s2 := 0.5*d1, 0.5*d2
	theta1 := 180 / float64(nspoke)

	theta2, offset := theta2In, offsetIn
	if theta2 == 0 {
		theta2 = asind(offset / r2)
	} else {
		offset = r2 * sind(theta2)
	}

	phi := searchPhi(r1, s1, theta1, r2, s2, theta2)
	exact = flywheelLayout{
		R1: r1, Offset: offset, Theta1: theta1, Theta2: theta2, Phi: phi,
		X:      r1*sind(theta1+phi) - s1,
		Sep:    2 * (theta1 - theta2),
		Spoke1: 2 * (r1*sind(theta1) - s1),
		Spoke2: 2 * (r2*sind(theta2) - s2),
		W:      r2 - r1 + s1 + s2,
	}

	theta2i := math.Round(theta2)
	phii := math.Round(phi)
	r1i := r2 * (sind(theta2i+phii) + (s1-s2)/r2) / sind(theta1+phii)
	nearestIntegral = flywheelLayout{
		R1: r1i, Offset: r2 * sind(theta2i), Theta1: theta1, Theta2: theta2i, Phi: phii,
		X:      r1i*sind(theta1+phii) - s1,
		Sep:    2 * (theta1 - theta2i),
		Spoke1: 2 * (r1i*sind(theta1) - s1),
		Spoke2: 2 * (r2*sind(theta2i) - s2),
		W:      r2 - r1i + s1 + s2,
	}

	return exact, nearestIntegral, nil
}

// rotarySettings returns the sequence of rotary table angle settings
// (degrees) needed to drill the inner holes, matching the original
// program's own "while (t < 360)" table-generation loop, bounded
// here by the fact that theta1 is always positive and the loop
// strictly increases t.
func innerHoleSettings(theta1 float64) []float64 {
	var settings []float64
	for t := theta1; t < 360; t += 2 * theta1 {
		settings = append(settings, t)
	}
	return settings
}

// outerHoleSettings returns the sequence of rotary table angle
// settings (degrees) needed to drill the outer holes, alternating
// between two step sizes (sep, then 2*theta2) per web space.
func outerHoleSettings(theta2, sep float64) []float64 {
	var settings []float64
	useSep := false
	for t := theta2; t < 360; {
		settings = append(settings, t)
		if useSep {
			t += sep
		} else {
			t += 2 * theta2
		}
		useSep = !useSep
	}
	return settings
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "flywheel:", err)
		os.Exit(1)
	}

	fmt.Println("TAPERED (OR STRAIGHT) SPOKE FLYWHEEL CALCULATIONS")
	fmt.Println()

	nspoke := prompter.Int("Number of spokes", 6)
	r1 := prompter.Float("Radius on which inner holes are located", 0.440)
	r2 := prompter.Float("Radius on which outer holes are located", 1.367)
	d1 := prompter.Float("Diameter of inner holes", 0.1875)
	d2 := prompter.Float("Diameter of outer holes", 0.1875)

	var exact, nearestIntegral flywheelLayout
	for {
		fmt.Println("YOU MUST SPECIFY EITHER THE OUTER HOLE OFFSET OR ITS OFFSET ANGLE.")
		fmt.Println("Note: Set offset = 0 and theta2 = 7 deg for demo.")
		offset := prompter.Float("Offset from spoke CL to outer hole center", 0)
		var theta2 float64
		if offset == 0 {
			theta2 = prompter.Float("Angle from spoke CL to outer hole center", 0)
		}

		exact, nearestIntegral, err = computeFlywheelLayout(nspoke, r1, r2, d1, d2, offset, theta2)
		if err == nil {
			break
		}
		fmt.Println("Failed to specify one or the other.  Try again.")
	}

	printFlywheelLayout(exact, r2, d1, d2)
	fmt.Println("\n***** NEAREST INTEGRAL ANGLE SOLUTION *****")
	printFlywheelLayout(nearestIntegral, r2, d1, d2)
}

func printFlywheelLayout(l flywheelLayout, r2, d1, d2 float64) {
	fmt.Printf("Radius on which inner holes are located (R1) = %.3f in\n", l.R1)
	fmt.Printf("Radius on which outer holes are located (R2) = %.3f in\n", r2)
	fmt.Printf("Diameter of inner holes (2*r1) = %.3f in\n", d1)
	fmt.Printf("Diameter of outer holes (2*r2) = %.3f in\n", d2)
	fmt.Printf("Distance from spoke CL to outer hole center (d2) = %.3f in\n", l.Offset)
	fmt.Println()
	fmt.Printf("Angle from spoke CL to inner hole center (theta1) = %.3f deg\n", l.Theta1)
	fmt.Printf("Angle from spoke CL to outer hole center (theta2) = %.3f deg\n", l.Theta2)
	fmt.Printf("Angle btw two outer holes in same web space = %.3f deg\n", l.Sep)
	fmt.Printf("Inner spoke width = %.3f in\n", l.Spoke1)
	fmt.Printf("Outer spoke width = %.3f in\n", l.Spoke2)
	fmt.Printf("Minimum radial web length (W) = %.3f in\n", l.W)
	fmt.Println()
	fmt.Println("Assuming that a spoke CL is initially aligned with the mill y axis,")
	fmt.Printf("the rotary table must be rotated (phi=) %.3f deg to bring the spoke edge\n", l.Phi)
	fmt.Println("parallel to the mill y axis.")
	fmt.Printf("The table must then be offset by (x=) %.3f in plus half the cutter diameter.\n", l.X)
	fmt.Println()
	fmt.Println("The rotary table settings (deg) for the inner holes are:")
	for _, t := range innerHoleSettings(l.Theta1) {
		fmt.Printf("%.3f\n", t)
	}
	fmt.Println()
	fmt.Println("The rotary table settings (deg) for the outer holes are:")
	for _, t := range outerHoleSettings(l.Theta2, l.Sep) {
		fmt.Printf("%.3f\n", t)
	}
	fmt.Println()
}
