// lvern designs a linear vernier scale: given the spacing of the
// main scale's major divisions, how many subdivisions the main scale
// already has, and how fine a reading the vernier should resolve,
// it computes how many divisions the vernier scale needs and how
// long each of those divisions is.
//
// Converted from LVERN.C (M. W. Klotz, 11/98),
// WorkshopUtilities/lvern.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/promptio"
)

// parseDecimalOrFraction parses s as either a plain decimal number
// ("1.5") or a mixed-number fraction ("1/2" or "1-1/2"), matching
// the original program's dfinput function. An empty string is an
// error; callers substitute their own default for that case.
func parseDecimalOrFraction(s string) (float64, error) {
	if s == "" {
		return 0, errors.New("empty input")
	}

	slashIndex := strings.IndexByte(s, '/')
	if slashIndex == -1 {
		return strconv.ParseFloat(s, 64)
	}

	fractionPart := s
	var whole float64
	if dashIndex := strings.IndexByte(s, '-'); dashIndex != -1 {
		var err error
		if whole, err = strconv.ParseFloat(s[:dashIndex], 64); err != nil {
			return 0, err
		}
		fractionPart = s[dashIndex+1:]
	}

	numDen := strings.SplitN(fractionPart, "/", 2)
	if len(numDen) != 2 {
		return 0, fmt.Errorf("parse fraction %q: missing '/'", fractionPart)
	}
	num, err := strconv.ParseFloat(strings.TrimSpace(numDen[0]), 64)
	if err != nil {
		return 0, err
	}
	den, err := strconv.ParseFloat(strings.TrimSpace(numDen[1]), 64)
	if err != nil {
		return 0, err
	}
	return whole + num/den, nil
}

// maxContinuedFractionSteps bounds the continued-fraction search in
// nearestRationalFraction. A finite-precision float64 always
// terminates in well under this many steps; the bound exists only to
// satisfy coding-style.md Rule 2 (bound every loop) against a
// pathological input.
const maxContinuedFractionSteps = 100

// nearestRationalFraction returns the nearest simple rational
// fraction to q, as a whole part plus numerator/denominator, using
// the same continued-fraction convergent algorithm as the original
// program's frac function: An = Xn*An-1 + An-2, where X = floor(1/F)
// and F is the fractional remainder at each step.
func nearestRationalFraction(q float64) (whole, num, den int64) {
	f := q
	a1, a2 := 1.0, 0.0
	for i := 0; math.Abs(f) > 1e-6 && i < maxContinuedFractionSteps; i++ {
		z := 1 / f
		x := math.Floor(z)
		f = z - x
		an := x*a1 + a2
		a2, a1 = a1, an
	}

	denominator := a1
	remainder := q * denominator
	var wholePart float64
	for remainder > denominator {
		remainder -= denominator
		wholePart++
	}

	return int64(wholePart), int64(math.Round(remainder)), int64(denominator)
}

// formatNearestFraction renders nearestRationalFraction's result the
// way the original program's frac function prints it: as
// "whole-num/den" if there's a whole part, or just "num/den"
// otherwise.
func formatNearestFraction(q float64) string {
	whole, num, den := nearestRationalFraction(q)
	if whole != 0 {
		return fmt.Sprintf("%d-%d/%d", whole, num, den)
	}
	return fmt.Sprintf("%d/%d", num, den)
}

// vernierMainScale returns the distance spanned by one division of
// the main scale (mainDivisionSpacing / mainSubdivisions) and the
// number of vernier subdivisions needed to resolve down to
// vernierResolution, along with whether that count came out to an
// exact whole number.
func vernierMainScale(mainDivisionSpacing float64, mainSubdivisions int, vernierResolution float64) (mainDivisionLength float64, vernierSubdivisions int, exact bool) {
	mainDivisionLength = mainDivisionSpacing / float64(mainSubdivisions)
	exactCount := mainDivisionLength / vernierResolution
	vernierSubdivisions = int(math.Floor(exactCount))
	exact = float64(vernierSubdivisions) == exactCount
	return mainDivisionLength, vernierSubdivisions, exact
}

// vernierDivisionLength returns the length of one division on a
// vernier scale of total span vernierSpan divided into
// vernierSubdivisions equal parts.
func vernierDivisionLength(vernierSpan float64, vernierSubdivisions int) float64 {
	return vernierSpan / float64(vernierSubdivisions)
}

// promptDecimalOrFraction prints label and def in the original
// program's own decimal-or-fraction prompt format, then returns the
// parsed value entered, or def if the line is blank or unparseable.
func promptDecimalOrFraction(prompter *promptio.Prompter, label string, def float64) float64 {
	line := prompter.Line(fmt.Sprintf("\ndefault = %.6f\n{[d].dd or [d-]d/d (e.g. 1.5 or 1-1/2)}\n%s ? ", def, label))
	if line == "" {
		return def
	}
	value, err := parseDecimalOrFraction(line)
	if err != nil {
		fmt.Printf("%q is not a valid number or fraction; using default %g\n", line, def)
		return def
	}
	return value
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "lvern:", err)
		os.Exit(1)
	}

	fmt.Println("VERNIER DESIGN")

	mainDivisionSpacing := promptDecimalOrFraction(prompter, "distance spanned by one major division on main scale", 1.0)
	mainSubdivisions := prompter.Int("number of subdivisions of one major division on main scale", 8)
	vernierResolution := promptDecimalOrFraction(prompter, "distance measured by one division on vernier scale", 1.0/32.0)
	fmt.Printf("Nearest Rational Fraction = %s\n", formatNearestFraction(vernierResolution))

	mainDivisionLength, vernierSubdivisions, exact := vernierMainScale(mainDivisionSpacing, mainSubdivisions, vernierResolution)
	if !exact {
		fmt.Println("\nnumber of vernier subdivisions not an integer")
		return
	}

	fmt.Printf("\ndistance spanned by one division on main scale = %.4f\n", mainDivisionLength)
	fmt.Printf("Nearest Rational Fraction = %s\n", formatNearestFraction(mainDivisionLength))
	fmt.Printf("number of subdivisions required on vernier scale = %d\n", vernierSubdivisions)

	typicalVernierSpan := mainDivisionSpacing - mainDivisionLength
	fmt.Printf("\na typical vernier scale would span a distance of %.4f\n", typicalVernierSpan)
	vernierSpan := prompter.Float("distance spanned by vernier scale", typicalVernierSpan)

	vernierDivision := vernierDivisionLength(vernierSpan, vernierSubdivisions)
	fmt.Printf("\ndistance spanned by one division on vernier scale = %.4f\n", vernierDivision)
	fmt.Printf("Nearest Rational Fraction = %s\n", formatNearestFraction(vernierDivision))
}
