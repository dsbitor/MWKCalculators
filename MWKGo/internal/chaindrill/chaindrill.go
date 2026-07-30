// Package chaindrill implements the chain-drilling layout
// calculation shared by the plate and slug programs: given a target
// circle diameter, a finish-machining allowance, a drill diameter,
// and an approximate web thickness, it finds the whole number of
// holes, drilled on what diameter, that keeps the web thickness
// close to the one requested.
package chaindrill

import (
	"math"

	"mwkgo/internal/mwktrig"
)

// Plan is the result of a chain-drilling layout.
type Plan struct {
	HoleCount              int
	DrillingCircleDiameter float64
	WebThickness           float64
	AngleBetweenHolesDeg   float64
	ChordBetweenHoles      float64
}

// Compute finds the chain-drilling layout for a target diameter of
// finalDiameter, growing or shrinking outward from it by the given
// radialAllowance and drillDiameter depending on outward: true grows
// the drilling circle outward past the final edge (cutting a plate
// free, where the as-drilled piece must be larger than the finished
// one), false shrinks it inward (opening a hole, where the
// as-drilled hole must be smaller than the finished one, since
// finish machining enlarges it).
func Compute(finalDiameter, radialAllowance, drillDiameter, approxWebThickness float64, outward bool) Plan {
	sign := -1.0
	if outward {
		sign = 1.0
	}

	afterAllowance := finalDiameter + sign*2*radialAllowance
	drillingCircleDiameter := afterAllowance + sign*drillDiameter
	drillingCircleRadius := 0.5 * drillingCircleDiameter
	drillRadius := 0.5 * drillDiameter

	theta := 2 * mwktrig.ClampedAsin(drillRadius/drillingCircleRadius)
	phi := approxWebThickness / drillingCircleRadius
	omega := theta + phi

	holeCount := int(math.Floor(2 * math.Pi / omega))
	omega = 2 * math.Pi / float64(holeCount)
	webThickness := drillingCircleRadius * (omega - theta)
	chord := drillingCircleDiameter * math.Sin(0.5*omega)

	return Plan{
		HoleCount:              holeCount,
		DrillingCircleDiameter: drillingCircleDiameter,
		WebThickness:           webThickness,
		AngleBetweenHolesDeg:   omega * 180 / math.Pi,
		ChordBetweenHoles:      chord,
	}
}
