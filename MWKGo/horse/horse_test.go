package main

import (
	"math"
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name                        string
		torque, rpm, hp             float64
		wantTorque, wantRPM, wantHP float64
	}{
		{name: "torque and rpm known, solves for horsepower", torque: 10, rpm: 2626, hp: 0, wantTorque: 10, wantRPM: 2626, wantHP: 10 * 2626 / horsepowerConstant},
		{name: "horsepower and rpm known, solves for torque", torque: 0, rpm: 2626, hp: 5, wantTorque: 5 * horsepowerConstant / 2626, wantRPM: 2626, wantHP: 5},
		{name: "horsepower and torque known, solves for rpm", torque: 10, rpm: 0, hp: 5, wantTorque: 10, wantRPM: 5 * horsepowerConstant / 10, wantHP: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTorque, gotRPM, gotHP := solve(tt.torque, tt.rpm, tt.hp)
			if diff := math.Abs(gotTorque - tt.wantTorque); diff > 1e-9 {
				t.Errorf("torque = %v, want %v", gotTorque, tt.wantTorque)
			}
			if diff := math.Abs(gotRPM - tt.wantRPM); diff > 1e-9 {
				t.Errorf("rpm = %v, want %v", gotRPM, tt.wantRPM)
			}
			if diff := math.Abs(gotHP - tt.wantHP); diff > 1e-9 {
				t.Errorf("hp = %v, want %v", gotHP, tt.wantHP)
			}
		})
	}
}

func TestSolve_AllThreeKnownLeavesInputsUnchanged(t *testing.T) {
	// With no zero among the three, there is nothing to solve for;
	// the original program's structure never reaches this case in
	// practice (it stops asking once two are known), but the
	// function itself must not silently overwrite a value the
	// caller already had.
	torque, rpm, hp := solve(10.0, 5252.0, 10.0)
	if torque != 10 || rpm != 5252 || hp != 10 {
		t.Errorf("solve(10, 5252, 10) = (%v, %v, %v), want the inputs unchanged", torque, rpm, hp)
	}
}

func TestSolve_FewerThanTwoKnownLeavesInputsUnchanged(t *testing.T) {
	torque, rpm, hp := solve(10.0, 0, 0)
	if torque != 10 || rpm != 0 || hp != 0 {
		t.Errorf("solve(10, 0, 0) = (%v, %v, %v), want the inputs unchanged", torque, rpm, hp)
	}
}
