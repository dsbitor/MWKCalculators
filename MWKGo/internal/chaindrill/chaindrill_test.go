package chaindrill

import (
	"math"
	"testing"
)

// These worked examples are taken directly from SLUG.TXT, the
// author's own annotated sample runs for chain-drilling a large
// hole (the inward, non-outward case).
func TestCompute_WorkedExamples(t *testing.T) {
	tests := []struct {
		name                                                               string
		finalDiameter, radialAllowance, drillDiameter, webWanted           float64
		wantHoleCount                                                      int
		wantDrillingCircleDiameter, wantWebThickness, wantAngle, wantChord float64
	}{
		{
			name:          "1/4in drill, from SLUG.TXT's first example",
			finalDiameter: 3, radialAllowance: 0.05, drillDiameter: 0.25, webWanted: 0.05,
			wantHoleCount: 27, wantDrillingCircleDiameter: 2.650, wantWebThickness: 0.058,
			wantAngle: 13.333, wantChord: 0.308,
		},
		{
			name:          "3/8in drill, from SLUG.TXT's second example",
			finalDiameter: 3, radialAllowance: 0.05, drillDiameter: 0.375, webWanted: 0.05,
			wantHoleCount: 18, wantDrillingCircleDiameter: 2.525, wantWebThickness: 0.064,
			wantAngle: 20.000, wantChord: 0.438,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compute(tt.finalDiameter, tt.radialAllowance, tt.drillDiameter, tt.webWanted, false)
			if got.HoleCount != tt.wantHoleCount {
				t.Errorf("HoleCount = %v, want %v", got.HoleCount, tt.wantHoleCount)
			}
			if diff := math.Abs(got.DrillingCircleDiameter - tt.wantDrillingCircleDiameter); diff > 0.001 {
				t.Errorf("DrillingCircleDiameter = %v, want %v", got.DrillingCircleDiameter, tt.wantDrillingCircleDiameter)
			}
			if diff := math.Abs(got.WebThickness - tt.wantWebThickness); diff > 0.001 {
				t.Errorf("WebThickness = %v, want %v", got.WebThickness, tt.wantWebThickness)
			}
			if diff := math.Abs(got.AngleBetweenHolesDeg - tt.wantAngle); diff > 0.001 {
				t.Errorf("AngleBetweenHolesDeg = %v, want %v", got.AngleBetweenHolesDeg, tt.wantAngle)
			}
			if diff := math.Abs(got.ChordBetweenHoles - tt.wantChord); diff > 0.001 {
				t.Errorf("ChordBetweenHoles = %v, want %v", got.ChordBetweenHoles, tt.wantChord)
			}
		})
	}
}

func TestCompute_OutwardAndInwardAreMirrorImages(t *testing.T) {
	// Growing outward by an allowance and drill diameter, then
	// shrinking back inward by the same amounts, must return to the
	// same final diameter: an identity independent of this code's
	// specific hole-count arithmetic, checked here on the
	// intermediate drilling circle diameter rather than the full
	// plan, since hole count and web thickness legitimately differ
	// between the two directions.
	finalDiameter, allowance, drill := 3.0, 0.05, 0.25
	outward := Compute(finalDiameter, allowance, drill, 0.05, true)
	inward := Compute(finalDiameter, allowance, drill, 0.05, false)

	outwardCircleFromFinal := finalDiameter + 2*allowance + drill
	inwardCircleFromFinal := finalDiameter - 2*allowance - drill

	if diff := math.Abs(outward.DrillingCircleDiameter - outwardCircleFromFinal); diff > 1e-9 {
		t.Errorf("outward DrillingCircleDiameter = %v, want %v", outward.DrillingCircleDiameter, outwardCircleFromFinal)
	}
	if diff := math.Abs(inward.DrillingCircleDiameter - inwardCircleFromFinal); diff > 1e-9 {
		t.Errorf("inward DrillingCircleDiameter = %v, want %v", inward.DrillingCircleDiameter, inwardCircleFromFinal)
	}
}

func TestCompute_HoleCountNeverExceedsWhatFitsWithoutOverlap(t *testing.T) {
	// The chosen omega must always be at least theta (the angle the
	// drill itself subtends), or adjacent holes would overlap more
	// than the geometry allows: a sanity property independent of
	// this code's specific web-thickness arithmetic.
	plan := Compute(3, 0.05, 0.25, 0.05, false)
	angleRadians := plan.AngleBetweenHolesDeg * math.Pi / 180
	drillRadius := 0.5 * 0.25
	drillingCircleRadius := 0.5 * plan.DrillingCircleDiameter
	theta := 2 * math.Asin(drillRadius/drillingCircleRadius)

	if angleRadians < theta {
		t.Errorf("angle between holes (%v rad) is less than the drill's own subtended angle (%v rad)", angleRadians, theta)
	}
}
