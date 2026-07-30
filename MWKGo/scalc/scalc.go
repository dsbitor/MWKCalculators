// scalc is a collection of small specialized calculators: a
// proportion solver, running parallel-resistance/RSS/RMS
// accumulators, running descriptive statistics, degree/DMS
// conversions in both directions, machine screw number to major
// diameter, linear interpolation, and drawing-scale conversion.
//
// Converted from SCALC.C (M. W. Klotz), Math/scalc.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"mwkgo/internal/promptio"
)

// solveRatio solves a/b = c/d for whichever of a, b, c, d is zero,
// given the other three.
func solveRatio(a, b, c, d float64) (float64, float64, float64, float64, error) {
	switch {
	case a != 0 && b != 0 && c != 0:
		d = c * b / a
	case a != 0 && b != 0 && d != 0:
		c = d * a / b
	case a != 0 && c != 0 && d != 0:
		b = a * d / c
	case b != 0 && c != 0 && d != 0:
		a = b * c / d
	default:
		return 0, 0, 0, 0, errors.New("insufficient data")
	}
	return a, b, c, d, nil
}

// parallelResistance returns 1/sum(1/x) for the values in xs (the
// combined resistance of resistors in parallel, or the combined
// value of any other quantity that combines reciprocally).
func parallelResistance(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += 1 / x
	}
	return 1 / sum
}

// rss returns the root-sum-square of the values in xs, preserving
// the original program's own sign-aware squaring (x*|x| rather than
// x*x), which lets a negative entry cancel a previous positive one
// of the same magnitude, matching the original's documented use of
// negative values to "remove a previously entered value".
func rss(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += x * math.Abs(x)
	}
	return math.Sqrt(sum)
}

// rms returns the root-mean-square of the values in xs, with the
// same sign-aware squaring as rss.
func rms(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += x * math.Abs(x)
	}
	return math.Sqrt(sum / float64(len(xs)))
}

// descriptiveStats holds the running statistics of a set of values.
type descriptiveStats struct {
	Max, Min, Median, Mean, StdDev float64
}

// computeStats returns the max, min, median, mean, and sample
// standard deviation of xs. The original program finds the median
// via an in-place quickselect partition; this conversion sorts a
// copy instead, an internal simplification with no effect on the
// result.
func computeStats(xs []float64) descriptiveStats {
	sorted := append([]float64(nil), xs...)
	sort.Float64s(sorted)

	n := len(sorted)
	stats := descriptiveStats{Max: sorted[n-1], Min: sorted[0], Median: sorted[(n+1)/2-1]}

	var sum, sum2 float64
	for _, x := range xs {
		sum += x
		sum2 += x * x
	}
	stats.Mean = sum / float64(n)
	if n > 1 {
		stats.StdDev = math.Sqrt((sum2 - float64(n)*stats.Mean*stats.Mean) / float64(n-1))
	}
	return stats
}

// dmsToDecimalDegrees converts degrees/minutes/seconds to fractional
// degrees, carrying seconds into minutes and minutes into degrees
// whenever they exceed (not merely reach) 60, matching the original
// program's own strict ">" comparison.
func dmsToDecimalDegrees(degrees, minutes, seconds float64) (normDeg, normMin, normSec, decimalDegrees float64) {
	sign := 1.0
	if degrees < 0 {
		degrees = -degrees
		sign = -1
	}
	for seconds > 60 {
		seconds -= 60
		minutes++
	}
	for minutes > 60 {
		minutes -= 60
		degrees++
	}
	return degrees, minutes, seconds, sign * (degrees + minutes/60 + seconds/3600)
}

// decimalDegreesToDMS converts fractional degrees to
// degrees/minutes/seconds, rounding seconds to the nearest whole
// second and carrying into minutes if that rounds up to 60.
func decimalDegreesToDMS(angle float64) (sign float64, degrees, minutes, seconds int) {
	sign = 1
	if angle < 0 {
		sign = -1
		angle = -angle
	}
	degrees = int(math.Floor(angle))
	minutes = int(math.Floor((angle - float64(degrees)) * 60))
	seconds = int(math.Round((angle - float64(degrees) - float64(minutes)/60) * 3600))
	if seconds == 60 {
		minutes++
		seconds = 0
	}
	return sign, degrees, minutes, seconds
}

// screwMajorDiameter returns the major diameter, in inches, of a
// numbered machine screw, per the standard formula.
func screwMajorDiameter(screwNumber float64) float64 {
	return 0.06 + screwNumber*0.013
}

// linearInterpolate returns the y value at x on the line through
// (x1,y1) and (x2,y2).
func linearInterpolate(x1, y1, x2, y2, x float64) float64 {
	slope := (y2 - y1) / (x2 - x1)
	intercept := y1 - slope*x1
	return slope*x + intercept
}

// drawingScale returns the scale factor (drawing units per real
// unit) for a drawing where a feature of length drawingLength
// represents a real feature of length realLength.
func drawingScale(drawingLength, realLength float64) float64 {
	return drawingLength / realLength
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scalc:", err)
		os.Exit(1)
	}

	fmt.Println("A COLLECTION OF SPECIALIZED CALCULATORS")

	for {
		fmt.Println()
		fmt.Println("A Ratio (a/b = c/d)")
		fmt.Println("B 1/Z = 1/x1 + 1/x2 + ...")
		fmt.Println("C rss = sqrt (x1*x1 + x2*x2 + ...)")
		fmt.Println("D rms = sqrt [(x1*x1 + x2*x2 + ...)/n]")
		fmt.Println("E max, min, median, mean and standard deviation")
		fmt.Println("F deg:min:sec -> fractional degrees")
		fmt.Println("G fractional degrees -> deg:min:sec")
		fmt.Println("H screw number -> major diameter")
		fmt.Println("I interpolation")
		fmt.Println("J scaling")
		fmt.Println("Q/Esc quit")
		fmt.Println()

		switch strings.ToLower(prompter.Line("(A-J,Q) ? ")) {
		case "a":
			runRatio(prompter)
		case "b":
			runAccumulator(prompter, "1/Z = 1/x1 + 1/x2 + ...", "z", parallelResistance)
		case "c":
			runAccumulator(prompter, "rss = sqrt (x1^2 + x2^2 + ...)", "rss", rss)
		case "d":
			runAccumulator(prompter, "rms = sqrt [(x1^2 + x2^2 + ...)/n]", "rms", rms)
		case "e":
			runStats(prompter)
		case "f":
			runDMSToDecimal(prompter)
		case "g":
			runDecimalToDMS(prompter)
		case "h":
			sn := prompter.Float("Screw number", 0)
			fmt.Printf("Screw diameter = %.4f in\n", screwMajorDiameter(sn))
		case "i":
			runInterpolation(prompter)
		case "j":
			runScale(prompter)
		case "q":
			return
		}
	}
}

func runRatio(p *promptio.Prompter) {
	fmt.Println("solves the ratio: a/b = c/d")
	fmt.Println("enter values for any three of a,b,c,d")
	fmt.Println("press return for unknown value")

	for {
		a := p.Float("a", 0)
		b := p.Float("b", 0)
		c := p.Float("c", 0)
		d := p.Float("d", 0)

		solvedA, solvedB, solvedC, solvedD, err := solveRatio(a, b, c, d)
		if err != nil {
			fmt.Println("\nInsufficient data - try again")
			continue
		}
		switch {
		case a == 0:
			fmt.Printf("a = %g\n", solvedA)
		case b == 0:
			fmt.Printf("b = %g\n", solvedB)
		case c == 0:
			fmt.Printf("c = %g\n", solvedC)
		default:
			fmt.Printf("d = %g\n", solvedD)
		}
		return
	}
}

func runAccumulator(p *promptio.Prompter, header, label string, f func([]float64) float64) {
	fmt.Println(header)
	fmt.Println("negative values of x may be used to remove a previously entered value")
	fmt.Println("press return when done")

	var xs []float64
	for {
		x := p.Float("x", 0)
		if x == 0 {
			return
		}
		xs = append(xs, x)
		fmt.Printf("%s = %g\n", label, f(xs))
	}
}

func runStats(p *promptio.Prompter) {
	fmt.Println("press return when done")

	var xs []float64
	for {
		x := p.Float("x", 0)
		if x == 0 {
			return
		}
		xs = append(xs, x)
		s := computeStats(xs)
		fmt.Printf("max value = %g, min value = %g\n", s.Max, s.Min)
		fmt.Printf("median = %g, mean = %g, standard deviation = %g\n", s.Median, s.Mean, s.StdDev)
	}
}

func runDMSToDecimal(p *promptio.Prompter) {
	d := p.Float("degrees (deg)", 0)
	m := p.Float("minutes (min)", 0)
	s := p.Float("seconds (sec)", 0)

	normDeg, normMin, normSec, decimal := dmsToDecimalDegrees(d, m, s)
	fmt.Printf("\n%g:%g:%g d:m:s = %g deg\n", normDeg, normMin, normSec, decimal)
}

func runDecimalToDMS(p *promptio.Prompter) {
	theta := p.Float("angle (deg)", 0)
	sign, degrees, minutes, seconds := decimalDegreesToDMS(theta)
	fmt.Printf("\n%g deg = %d:%d:%d d:m:s\n", theta, int(sign)*degrees, minutes, seconds)
	fmt.Printf("error = %g deg\n", float64(degrees)+float64(minutes)/60+float64(seconds)/3600-math.Abs(theta))
}

func runInterpolation(p *promptio.Prompter) {
	for {
		x1 := p.Float("X1", 2.3)
		y1 := p.Float("Y1", math.Log(x1))
		x2 := p.Float("X2", 2.5)
		y2 := p.Float("Y2", math.Log(x2))
		x := p.Float("x", 2.4)

		fmt.Printf("Interpolated value = %.6f\n", linearInterpolate(x1, y1, x2, y2, x))

		if strings.ToLower(p.Line("Another interpolation (Y/[N])? ")) != "y" {
			return
		}
	}
}

func runScale(p *promptio.Prompter) {
	draw := p.Float("Length of object on drawing (in)", 0.25)
	real := p.Float("Length of real object (in)", 12.0)
	sf := drawingScale(draw, real)
	fmt.Printf("Scale = %.4f in/in = %.4f in/ft\n", sf, sf*12)

	for {
		xd := p.Float("Length of object on drawing (zero to quit) (in)", 0.1)
		if xd == 0 {
			return
		}
		fmt.Printf("Length of real object = %.4f in\n", xd/sf)
	}
}
