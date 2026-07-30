// rfrac approximates a decimal number as a rational fraction, using
// a continued-fraction expansion carried out only as far as needed
// to reach a requested approximation accuracy.
//
// Converted from RFRAC.C, Math/rfrac.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// maxContinuedFractionTerms bounds the continued-fraction expansion,
// matching the original program's own NCF constant.
const maxContinuedFractionTerms = 100

// continuedFractionConvergent evaluates the continued fraction
// [terms[0]; terms[1], terms[2], ...] and returns its value as a
// numerator and denominator, matching the original program's cfrac
// function.
func continuedFractionConvergent(terms []float64) (num, den int64) {
	k := len(terms) - 1
	lc, lb, la := 1.0, 0.0, terms[k]
	for k > 0 {
		k--
		lc, lb, la = lc*la+lb, lc, terms[k]
	}
	return int64(lc), int64(lc*la + lb)
}

// approximateFractionalPart returns the numerator and denominator of
// the simplest rational fraction approximating fractionalPart (which
// must be in [0, 1)) to within relativeError, built up one continued
// fraction term at a time until that accuracy is reached or
// maxContinuedFractionTerms is exhausted, matching the original
// program's ffrac function.
func approximateFractionalPart(fractionalPart, relativeError float64) (num, den int64) {
	if fractionalPart == 0 {
		return 0, 1
	}

	var terms []float64
	z := fractionalPart
	error := math.Inf(1)
	for len(terms) < maxContinuedFractionTerms && math.Abs(error) > relativeError {
		z = 1 / z
		term := math.Floor(z)
		terms = append(terms, term)
		z -= term
		num, den = continuedFractionConvergent(terms)
		error = (float64(num)/float64(den) - fractionalPart) / fractionalPart
	}
	return num, den
}

// gcd returns the greatest common divisor of x and y.
func gcd(x, y int64) int64 {
	for y != 0 {
		x, y = y, x%y
	}
	if x < 0 {
		return -x
	}
	return x
}

// rationalApproximation is a rational approximation to a number,
// split into its whole part and a reduced fractional part.
type rationalApproximation struct {
	Whole    int64
	Num, Den int64
	Value    float64
	ErrorPct float64
}

// errWholeNumberHasNoFraction reports that f has no fractional part
// to approximate. The original program has no such guard, and would
// divide by zero (in its own relative-error calculation) for a whole
// number input; this conversion reports the condition explicitly
// instead.
var errWholeNumberHasNoFraction = errors.New("number is already a whole number; nothing to approximate")

// approximateFraction returns a rational approximation to f, accurate
// to within accuracyPct percent, split into a whole part plus a
// reduced numerator/denominator for the fractional part.
func approximateFraction(f, accuracyPct float64) (rationalApproximation, error) {
	f = math.Abs(f)
	whole, fractionalPart := math.Modf(f)
	if fractionalPart == 0 {
		return rationalApproximation{}, errWholeNumberHasNoFraction
	}

	num, den := approximateFractionalPart(fractionalPart, accuracyPct/100)
	if g := gcd(num, den); g > 1 {
		num, den = num/g, den/g
	}

	value := float64(num) / float64(den)
	return rationalApproximation{
		Whole:    int64(whole),
		Num:      num,
		Den:      den,
		Value:    whole + value,
		ErrorPct: 100 * (value - fractionalPart) / fractionalPart,
	}, nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rfrac:", err)
		os.Exit(1)
	}

	fmt.Println("RATIONAL FRACTION APPROXIMATION")
	fmt.Println()

	f := math.Abs(prompter.Float("Number to approximate", 3.14159))

	var accuracyPct float64
	for {
		accuracyPct = prompter.Float("Desired approximation accuracy (%)", 0.01)
		if accuracyPct >= 0 && accuracyPct <= 100 {
			break
		}
	}

	result, err := approximateFraction(f, accuracyPct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rfrac:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("Number = %g\n", f)
	fmt.Printf("%d & %d/%d = %d/%d = %g\n", result.Whole, result.Num, result.Den, result.Den*result.Whole+result.Num, result.Den, result.Value)
	fmt.Printf("Approximation Error = %g %%\n", result.ErrorPct)
}
