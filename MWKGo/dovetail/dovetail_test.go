package main

import (
	"math"
	"testing"
)

func TestMaleAndFemaleDovetail_DocumentedDefaultsAreAMatchedPair(t *testing.T) {
	// The original program's own default inputs form a matched pair:
	// a male dovetail measured across its default pin separation
	// (2.5) and a female dovetail measured across its default pin
	// separation (1.0283) describe the same physical dovetail, so
	// the male's top/bottom must equal the female's bottom/top.
	// This is a stronger check than re-deriving the formula, since
	// it also confirms the two formulas are consistent with each
	// other, not just each internally self-consistent.
	angle, pinDiameter, height := 60.0, 0.375, 0.5

	maleBottom, maleTop := maleDovetail(angle, pinDiameter, height, 2.5)
	femaleBottom, femaleTop := femaleDovetail(angle, pinDiameter, height, 1.0283)

	if math.Abs(maleTop-femaleBottom) > 1e-3 {
		t.Errorf("male top (%v) should match female bottom (%v)", maleTop, femaleBottom)
	}
	if math.Abs(maleBottom-femaleTop) > 1e-3 {
		t.Errorf("male bottom (%v) should match female top (%v)", maleBottom, femaleTop)
	}
}

func TestMaleDovetail_DocumentedDefaultInput(t *testing.T) {
	bottom, top := maleDovetail(60.0, 0.375, 0.5, 2.5)
	wantBottom := 1.475480947161671
	wantTop := 2.052831216351297
	if math.Abs(bottom-wantBottom) > 1e-9 {
		t.Errorf("bottom = %v, want %v", bottom, wantBottom)
	}
	if math.Abs(top-wantTop) > 1e-9 {
		t.Errorf("top = %v, want %v", top, wantTop)
	}
}

func TestMaleDovetail_ZeroPinDiameterMeasuresBottomDirectly(t *testing.T) {
	// A pin of zero diameter contributes no correction at all: the
	// "across pins" measurement is then, trivially, the bottom
	// measurement itself.
	bottom, top := maleDovetail(60.0, 0.0, 0.5, 2.5)
	if math.Abs(bottom-2.5) > 1e-9 {
		t.Errorf("bottom = %v, want 2.5 (measurement unchanged with zero pin diameter)", bottom)
	}
	if !(top > bottom) {
		t.Errorf("top (%v) should be greater than bottom (%v) for a positive dovetail height", top, bottom)
	}
}
