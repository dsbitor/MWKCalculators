// horse solves for whichever of torque, rotational speed, or
// horsepower is unknown, given the other two.
//
// Converted from HORSE.C (M. W. Klotz, 4/01), WorkshopUtilities/horse.
package main

import (
	"fmt"
	"os"

	"mwkgo/internal/promptio"
)

// horsepowerConstant is the standard conversion constant relating
// torque in ft-lb, speed in rpm, and horsepower:
// horsepower = torque * rpm / horsepowerConstant.
const horsepowerConstant = 5252.0

// solve fills in whichever of torque, rpm, or hp is zero, given the
// other two nonzero values. If fewer than two of the three are
// nonzero, or all three already are, the inputs are returned
// unchanged: solving requires exactly one unknown.
func solve(torque, rpm, hp float64) (newTorque, newRPM, newHP float64) {
	switch {
	case hp != 0 && rpm != 0 && torque == 0:
		torque = hp * horsepowerConstant / rpm
	case hp != 0 && torque != 0 && rpm == 0:
		rpm = hp * horsepowerConstant / torque
	case torque != 0 && rpm != 0 && hp == 0:
		hp = torque * rpm / horsepowerConstant
	}
	return torque, rpm, hp
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "horse:", err)
		os.Exit(1)
	}

	fmt.Println("TORQUE-SPEED-HORSEPOWER CALCULATIONS")
	fmt.Println()

	fmt.Println("Input whatever data you know; press return if not known.")
	fmt.Println("You must enter two data items to obtain a solution.")
	fmt.Println()
	fmt.Println("ft-lb = in-lb/12")

	torque := prompter.Float("Torque (ft-lb)", 0)
	rpm := prompter.Float("Rotational speed (rpm)", 0)

	known := 0
	if torque != 0 {
		known++
	}
	if rpm != 0 {
		known++
	}

	var hp float64
	if known < 2 {
		hp = prompter.Float("Horsepower (hp)", 0)
		if hp != 0 {
			known++
		}
	}

	if known < 2 {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	torque, rpm, hp = solve(torque, rpm, hp)

	fmt.Println()
	fmt.Printf("Torque     = %.4f ft-lb = %.4f in-lb\n", torque, torque*12)
	fmt.Printf("Speed      = %.4f rpm\n", rpm)
	fmt.Printf("Horsepower = %.4f hp\n", hp)
}
