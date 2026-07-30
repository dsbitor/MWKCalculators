package main

import (
	"math"
	"testing"
)

func TestComputeHelicalGear_ZeroHelixAngleIsAStandardSpurGear(t *testing.T) {
	// A zero degree helix angle is a plain spur gear, whose pitch
	// diameter is teeth/diametralPitch exactly: the standard spur
	// gear formula, independent of this code's helix-specific
	// trigonometry.
	got := computeHelicalGear(20, 10, 0, 1, 0.125)
	want := 20.0 / 10.0
	if diff := math.Abs(got.PitchDiameter - want); diff > 1e-9 {
		t.Errorf("PitchDiameter = %v, want %v", got.PitchDiameter, want)
	}
}

func TestComputeHelicalGear_WholeDepthTierBoundary(t *testing.T) {
	tests := []struct {
		name           string
		diametralPitch float64
		want           float64
	}{
		{name: "coarse pitch uses the 2.2/P+0.002 formula", diametralPitch: 10, want: 2.2/10 + 0.002},
		{name: "fine pitch uses the 2.157/P formula", diametralPitch: 40, want: 2.157 / 40},
		{name: "boundary at exactly 20 uses the coarse formula", diametralPitch: 20, want: 2.2/20 + 0.002},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeHelicalGear(6, tt.diametralPitch, 80, 1, 0.125)
			if diff := math.Abs(got.WholeDepth - tt.want); diff > 1e-9 {
				t.Errorf("WholeDepth = %v, want %v", got.WholeDepth, tt.want)
			}
		})
	}
}

func TestComputeHelicalGear_DocumentedDefaultInput(t *testing.T) {
	// The documented default input, evaluated against the ported
	// formula directly.
	teeth, diametralPitch, helixAngle, mandrelHub, templateThickness := 6.0, 40.0, 80.0, 1.0, 0.125

	cosHelix := math.Cos(helixAngle * math.Pi / 180)
	tanHelix := math.Tan(helixAngle * math.Pi / 180)
	wantPD := teeth / (diametralPitch * cosHelix)
	wantLead := math.Pi * wantPD / tanHelix
	wantBlank := wantPD + 2/diametralPitch
	wantTemplateAngle := math.Atan(wantLead/(math.Pi*mandrelHub+templateThickness)) * 180 / math.Pi

	got := computeHelicalGear(teeth, diametralPitch, helixAngle, mandrelHub, templateThickness)
	if diff := math.Abs(got.PitchDiameter - wantPD); diff > 1e-9 {
		t.Errorf("PitchDiameter = %v, want %v", got.PitchDiameter, wantPD)
	}
	if diff := math.Abs(got.HelixLead - wantLead); diff > 1e-9 {
		t.Errorf("HelixLead = %v, want %v", got.HelixLead, wantLead)
	}
	if diff := math.Abs(got.BlankDiameter - wantBlank); diff > 1e-9 {
		t.Errorf("BlankDiameter = %v, want %v", got.BlankDiameter, wantBlank)
	}
	if diff := math.Abs(got.TemplateAngleDeg - wantTemplateAngle); diff > 1e-9 {
		t.Errorf("TemplateAngleDeg = %v, want %v", got.TemplateAngleDeg, wantTemplateAngle)
	}
}
