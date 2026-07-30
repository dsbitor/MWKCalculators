// Package numdisplay implements the variable-precision number display
// used by several of the original MWK.LIB-based calculators (their
// own vp and dplaces functions): fixed-point display with optional
// self-adjusting decimal places and comma grouping, plus engineering
// and scientific notation, with an automatic escalation to
// engineering notation for magnitudes fixed-point display cannot show
// sensibly.
package numdisplay

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"mwkgo/internal/mwkfmt"
)

// Notations.
const (
	Fixed = iota
	Engineering
	Scientific
)

// selfAdjustEpsilon is how closely a formatted-and-reparsed value
// must match the original for the self-adjusting branch to accept
// that number of decimal places, matching the original's own 1e-15
// tolerance (double precision has about 15 significant decimal
// digits).
const selfAdjustEpsilon = 1e-15

// maxDecimalSearch bounds the search for the number of decimal
// places needed to show a non-zero fraction; double precision cannot
// usefully distinguish more than about 15 decimal digits.
const maxDecimalSearch = 15

// maxEngineeringShift bounds the loop that walks an engineering-
// notation exponent down to the nearest multiple of three with a
// mantissa magnitude >= 1. The exponent moves by exactly one per
// step and only ever needs to move by at most two steps (the
// distance to the nearest lower multiple of three), so this ceiling
// is never reached in practice; it exists only so the loop has an
// enforced maximum.
const maxEngineeringShift = 10

// Format formats x for display in the given notation, with decimals
// places (or, in fixed-point self-adjusting mode, at least enough
// places to round-trip x exactly). It automatically escalates to
// engineering notation when x is far enough from 1 (log10 magnitude
// beyond 15) that fixed-point display would be unreadable, matching
// the original program's own automatic override.
func Format(x float64, decimals int, selfAdjust bool, notation int) string {
	magnitude := 1.0
	if x != 0 {
		magnitude = math.Log10(math.Abs(x))
	}
	effective := notation
	if notation == Fixed && math.Abs(magnitude) > maxDecimalSearch {
		effective = Engineering
	}
	places := decimalPlacesNeeded(x, decimals, selfAdjust, effective)

	switch effective {
	case Engineering:
		return formatEngineering(x, places)
	case Scientific:
		return formatScientific(x, places)
	default:
		return formatFixed(x, places)
	}
}

// decimalPlacesNeeded returns the number of decimal places to
// display x with. In fixed-point mode with self-adjusting decimals,
// it finds the fewest places that let x round-trip through
// formatting to within selfAdjustEpsilon; in fixed-point mode with a
// set decimals, it uses that count unless it would hide the whole
// fractional part as "0", in which case it grows just enough to show
// something non-zero. Engineering and scientific notation always use
// the configured decimals directly.
func decimalPlacesNeeded(x float64, decimals int, selfAdjust bool, notation int) int {
	if notation != Fixed {
		return decimals
	}
	_, frac := math.Modf(math.Abs(x))
	minPlaces := decimals
	if selfAdjust {
		minPlaces = 0
	}
	if frac == 0 {
		return minPlaces
	}
	for d := minPlaces; d <= maxDecimalSearch; d++ {
		wd, _ := strconv.ParseFloat(strconv.FormatFloat(x, 'f', d, 64), 64)
		if selfAdjust {
			if math.Abs(x-wd) < selfAdjustEpsilon {
				return d
			}
		} else if wd != 0 {
			return d
		}
	}
	return maxDecimalSearch
}

// formatFixed formats x with the given number of decimal places,
// inserting comma separators every three digits of the integer part
// once |x| reaches 1000.
func formatFixed(x float64, places int) string {
	s := strconv.FormatFloat(x, 'f', places, 64)
	if math.Abs(x) < 1000 {
		return s
	}

	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}
	intPart, fracPart := s, ""
	if dot := strings.IndexByte(s, '.'); dot != -1 {
		intPart, fracPart = s[:dot], s[dot:]
	}
	n, err := strconv.ParseUint(intPart, 10, 64)
	if err != nil {
		return sign + s
	}
	return sign + mwkfmt.GroupedUint(n) + fracPart
}

// formatEngineering formats x with a mantissa magnitude in [1,1000)
// and an exponent that is a multiple of three, the convention used
// by digital multimeters and similar instruments so the exponent
// lines up with metric prefixes (k, M, etc.).
func formatEngineering(x float64, places int) string {
	xd := math.Abs(x)
	mantissa, exponent := x, 0
	if xd != 0 {
		exponent = int(math.Floor(math.Log10(xd)))
		mantissa = xd * math.Pow(10, float64(-exponent))
		for shift := 0; shift < maxEngineeringShift && (exponent%3 != 0 || mantissa < 1); shift++ {
			mantissa *= 10
			exponent--
		}
		if x < 0 {
			mantissa = -mantissa
		}
	}
	return fmt.Sprintf("%.*f E%+04d", places, mantissa, exponent)
}

func formatScientific(x float64, places int) string {
	s := strconv.FormatFloat(x, 'E', places, 64)
	idx := strings.IndexByte(s, 'E')
	return s[:idx] + " " + s[idx:]
}
