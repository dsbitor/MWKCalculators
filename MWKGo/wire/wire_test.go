package main

import (
	"math"
	"testing"
)

func TestDiameterMils_MatchesStandardAWGFormula(t *testing.T) {
	// The standard published AWG geometric progression formula,
	// diameter_mils = 5 * 92^((36-gage)/39), is a widely documented
	// reference independent of this code and this program's own
	// constant. The two should agree closely, though not to full
	// floating-point precision, since gageStepRatio is a finite
	// decimal approximation of the true 92^(1/39).
	tests := []float64{0, 10, 20, 36}

	for _, gage := range tests {
		got := diameterMils(gage)
		want := 5 * math.Pow(92, (36-gage)/39)
		if diff := math.Abs(got - want); diff > 0.05 {
			t.Errorf("diameterMils(%v) = %v, want approximately %v (standard AWG formula)", gage, got, want)
		}
	}
}

func TestGageForDiameter_IsTheInverseOfDiameterMils(t *testing.T) {
	tests := []float64{0, 10, 20, 36}

	for _, gage := range tests {
		roundTripped := gageForDiameter(diameterMils(gage))
		if diff := math.Abs(roundTripped - gage); diff > 1e-9 {
			t.Errorf("gageForDiameter(diameterMils(%v)) = %v, want %v", gage, roundTripped, gage)
		}
	}
}

func TestRecommendedGage_ExactMatchRoundsToThatGage(t *testing.T) {
	// Constructing current and density so the required diameter
	// exactly matches AWG 10's diameter removes any rounding
	// ambiguity from the test.
	density := 1.0
	exactDiameter := diameterMils(10)
	current := exactDiameter * exactDiameter * density

	got := recommendedGage(current, density)
	if got != 10 {
		t.Errorf("recommendedGage(%v, %v) = %v, want 10", current, density, got)
	}
}

func TestRecommendedGage_DocumentedDefaultInput(t *testing.T) {
	// The documented default input, evaluated against the ported
	// formula directly.
	current, density := 12.0, 0.0025
	requiredDiameter := math.Sqrt(current / density)
	want := int(math.Floor(gageForDiameter(requiredDiameter) + 0.5))

	got := recommendedGage(current, density)
	if got != want {
		t.Errorf("recommendedGage(%v, %v) = %v, want %v", current, density, got, want)
	}
}

func TestPropertiesForGage_AreConsistentWithDiameter(t *testing.T) {
	props := propertiesForGage(12)
	d := props.DiameterMils

	if diff := math.Abs(props.AreaCircularMils - d*d); diff > 1e-9 {
		t.Errorf("AreaCircularMils = %v, want %v", props.AreaCircularMils, d*d)
	}
	if diff := math.Abs(props.ResistanceOhmsPer1000Ft - 10370/(d*d)); diff > 1e-9 {
		t.Errorf("ResistanceOhmsPer1000Ft = %v, want %v", props.ResistanceOhmsPer1000Ft, 10370/(d*d))
	}
	if diff := math.Abs(props.WeightLbsPer1000Ft - 3.02675e-3*d*d); diff > 1e-9 {
		t.Errorf("WeightLbsPer1000Ft = %v, want %v", props.WeightLbsPer1000Ft, 3.02675e-3*d*d)
	}
}
