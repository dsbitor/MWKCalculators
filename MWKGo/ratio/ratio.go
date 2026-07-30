// ratio generates a "preferred numbers" series (American National
// Standard, also known as a Renard series): a geometric progression
// from 1 to 10 in a chosen number of steps, each consecutive pair a
// fixed percentage apart, optionally scaled and rounded to the
// nearest integer — a standard tool for choosing a sensible spread of
// sizes (gears, resistor values, packaging sizes, and the like)
// without gaps that are too large or too fine to be useful.
//
// Converted from RATIO.C (M. W. Klotz), WorkshopUtilities/ratio. The
// Tier 2 preparatory scan found this program misclassified: its own
// dfile variable ("DATA.DAT") is dead code, never actually passed to
// fopen — only an *output* report file is opened — so it belongs
// with the Tier 1 no-data-file calculators, not Tier 2; see
// ai/plans/c-to-go-conversion-plan.md's clean-up phase note. (The
// repository also had a byte-for-byte duplicate extraction of this
// same archive under a "ratio 2" directory name, not a second,
// distinct program.) Like several other Tier 1 programs, this drops
// the original's file-save feature and prints straight to stdout.
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// stepCounts matches RATIO.C's own mm[] array: the preset series
// lengths for its five standard preferred-number series (R5, R10,
// R20, R40, R80 in the standard's own naming).
var stepCounts = [5]int{5, 10, 20, 40, 80}

// preferredNumber is one entry in the series: its position, its raw
// value (always starting at 1), the value after applying the scale
// factor, and the scaled value rounded to the nearest integer.
type preferredNumber struct {
	Index   int
	Value   float64
	Scaled  float64
	Nearest uint64
}

// consecutiveRatioPercent is the fixed percentage step between
// consecutive numbers in an m-step series, matching RATIO.C's own
// formula: each step multiplies by 10^(1/m), so the step itself is
// that factor minus 1, as a percentage.
func consecutiveRatioPercent(m int) float64 {
	return 100 * (math.Pow(10, 1/float64(m)) - 1)
}

// roundFrac rounds s to the nearest integer, matching RATIO.C's own
// rounding rule exactly: only a fractional part strictly greater than
// 0.5 rounds up; a fractional part of exactly 0.5 (or less) rounds
// down. This is round-half-down, not the more familiar round-half-up.
func roundFrac(s float64) uint64 {
	intPart, frac := math.Modf(s)
	if frac > 0.5 {
		return uint64(intPart + 1)
	}
	return uint64(intPart)
}

// computeSeries returns the m-step preferred-number series (a
// geometric progression from 1 to just short of 10), each entry
// scaled by scale, matching RATIO.C's own wdata() computation.
func computeSeries(m int, scale float64) []preferredNumber {
	series := make([]preferredNumber, m)
	for i := 0; i < m; i++ {
		value := math.Pow(10, float64(i)/float64(m))
		scaled := scale * value
		series[i] = preferredNumber{Index: i, Value: value, Scaled: scaled, Nearest: roundFrac(scaled)}
	}
	return series
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ratio:", err)
		os.Exit(1)
	}

	fmt.Println("AMERICAN NATIONAL STANDARD FOR PREFERRED NUMBERS")
	fmt.Println()
	fmt.Println("1 -  5 NUMBERS ABOUT 58% APART")
	fmt.Println("2 - 10 NUMBERS ABOUT 26% APART")
	fmt.Println("3 - 20 NUMBERS ABOUT 12% APART")
	fmt.Println("4 - 40 NUMBERS ABOUT  6% APART")
	fmt.Println("5 - 80 NUMBERS ABOUT  3% APART")
	fmt.Println("6 - DESIGN YOUR OWN SET OF NUMBERS")

	op := prompter.Int("Option", 2)
	var m int
	switch {
	case op == 6:
		m = prompter.Int("How many numbers in series", 15)
	case op >= 1 && op <= 5:
		m = stepCounts[op-1]
	default:
		fmt.Println("invalid option")
		os.Exit(1)
	}
	if m < 1 {
		fmt.Println("number of numbers in series must be at least 1")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("The first number in each series is always 1.  You can scale the series")
	fmt.Println("by applying a scale factor (typically a power of ten) to each number")
	fmt.Println("in the series.")
	fmt.Println()
	scale := prompter.Float("Scale factor", 10)

	fmt.Println("\nAMERICAN NATIONAL STANDARD FOR PREFERRED NUMBERS")
	fmt.Println()
	fmt.Printf("RATIO OF CONSECUTIVE NUMBERS = %.2f %%\n\n", consecutiveRatioPercent(m))
	fmt.Println("N     NUMBER      SCALED       NEAREST INTEGER")
	for _, p := range computeSeries(m, scale) {
		fmt.Printf("%-3d   %-10.4f  %-10.4f   %-10d\n", p.Index, p.Value, p.Scaled, p.Nearest)
	}
}
