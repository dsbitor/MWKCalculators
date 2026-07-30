// Package mwktrig implements the domain-clamped inverse trig
// functions used throughout the original programs' shared macros
// (ASN, ACS, ASND, ACSD in MWK.LIB). A value passed to arcsine or
// arccosine that has drifted fractionally outside [-1,1], typically
// from floating-point round-off in an otherwise valid calculation,
// is clamped to the nearest edge of the domain instead of producing
// NaN.
package mwktrig

import "math"

// sign returns -1, 0, or 1 according to the sign of x, matching the
// original programs' SGN macro.
func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

// ClampedAsin returns the arcsine of x in radians, clamped the same
// way as ClampedAsinDeg: if x is outside [-1,1], it is replaced with
// its own sign before the arcsine is taken. This matches the
// original programs' ASN macro (the radians counterpart of ASND).
func ClampedAsin(x float64) float64 {
	if math.Abs(x) >= 1 {
		x = sign(x)
	}
	return math.Asin(x)
}

// ClampedAsinDeg returns the arcsine of x in degrees. If x is
// outside [-1,1], it is replaced with its own sign (1 or -1) before
// the arcsine is taken, giving a clamped result of -90 or 90 degrees
// instead of NaN.
func ClampedAsinDeg(x float64) float64 {
	if math.Abs(x) >= 1 {
		x = sign(x)
	}
	return math.Asin(x) * 180 / math.Pi
}

// ClampedAcosDeg returns the arccosine of x in degrees, clamped the
// same way as ClampedAsinDeg.
func ClampedAcosDeg(x float64) float64 {
	if math.Abs(x) >= 1 {
		x = sign(x)
	}
	return math.Acos(x) * 180 / math.Pi
}
