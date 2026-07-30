// mixture solves for the amounts and concentration of a mixture of
// two solutions, given the concentration of each solution and any
// two of: the amount of solution a, the amount of solution b, the
// total amount of mixture, or the mixture's concentration. Diluting
// a solution with a pure diluent (e.g. water) is treated as mixing
// with a solution of 0% concentration.
//
// Converted from MIXTURE.C, WorkshopUtilities/mixture.
package main

import (
	"errors"
	"fmt"
	"os"

	"mwkgo/internal/promptio"
)

// mixtureResult holds the fully solved amounts and concentrations of
// a two-solution mixture.
type mixtureResult struct {
	AmountA, AmountB, AmountMixture float64
	ConcentrationMixturePct         float64
}

// errInsufficientMixtureData reports that fewer than two of the four
// solvable quantities (amount of a, amount of b, amount of mixture,
// concentration of mixture) were given.
var errInsufficientMixtureData = errors.New("insufficient data for solution")

// solveMixture solves for the amounts of solution a, solution b, and
// the resulting mixture, given the two solutions' concentrations
// (as percentages) and any two of amountA, amountB, amountMixture,
// or concentrationMixturePct (zero meaning unknown). It recognizes
// the same six combinations as the original program.
func solveMixture(concentrationAPct, concentrationBPct, amountA, amountB, amountMixture, concentrationMixturePct float64) (mixtureResult, error) {
	pa := 0.01 * concentrationAPct
	pb := 0.01 * concentrationBPct
	pm := 0.01 * concentrationMixturePct

	switch {
	case amountA != 0 && amountB != 0:
		amountMixture = amountA + amountB
		pm = (amountA*pa + amountB*pb) / amountMixture
	case amountA != 0 && concentrationMixturePct != 0:
		amountB = amountA * (pa - pm) / (pm - pb)
		amountMixture = amountA + amountB
	case amountA != 0 && amountMixture != 0:
		amountB = amountMixture - amountA
		pm = (amountA*pa + amountB*pb) / amountMixture
	case amountB != 0 && concentrationMixturePct != 0:
		amountA = amountB * (pm - pb) / (pa - pm)
		amountMixture = amountA + amountB
	case amountB != 0 && amountMixture != 0:
		amountA = amountMixture - amountB
		pm = (amountA*pa + amountB*pb) / amountMixture
	case amountMixture != 0 && concentrationMixturePct != 0:
		amountA = amountMixture * (pm - pb) / (pa - pb)
		amountB = amountMixture - amountA
	default:
		return mixtureResult{}, errInsufficientMixtureData
	}

	return mixtureResult{
		AmountA:                 amountA,
		AmountB:                 amountB,
		AmountMixture:           amountMixture,
		ConcentrationMixturePct: 100 * pm,
	}, nil
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mixture:", err)
		os.Exit(1)
	}

	fmt.Println("MIXTURE AND DILUTION CALCULATOR")
	fmt.Println()
	fmt.Println("Concentrations of both solutions, Ca and Cb, must BOTH be specified.")
	fmt.Println("If one solution is a dilutant (e.g., pure water), enter its concentration")
	fmt.Println("as zero.")
	fmt.Println()

	concentrationA := prompter.Float("Ca = concentration of solution a (%)", 10.0)
	concentrationB := prompter.Float("Cb = concentration of solution b (%)", 0.0)

	fmt.Println("\nInput whatever data you know; press return if not known.")
	fmt.Println("You must enter two data items to obtain a solution.")
	fmt.Println()

	amountA := prompter.Float("amount of solution a", 0)
	amountB := prompter.Float("amount of solution b", 0)

	known := 0
	for _, v := range []float64{amountA, amountB} {
		if v != 0 {
			known++
		}
	}

	var amountMixture, concentrationMixture float64
	if known < 2 {
		amountMixture = prompter.Float("amount of mixture", 0)
		if amountMixture != 0 {
			known++
		}
	}
	if known < 2 {
		concentrationMixture = prompter.Float("concentration of mixture (%)", 0)
		if concentrationMixture != 0 {
			known++
		}
	}
	if known < 2 {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	result, err := solveMixture(concentrationA, concentrationB, amountA, amountB, amountMixture, concentrationMixture)
	if err != nil {
		fmt.Println("\nINSUFFICIENT DATA FOR SOLUTION")
		return
	}

	fmt.Printf("\namount of solution a = %.5g\n", result.AmountA)
	fmt.Printf("concentration of solution a = %.5g %%\n", concentrationA)
	fmt.Printf("amount of solution b = %.5g\n", result.AmountB)
	fmt.Printf("concentration of solution b = %.5g %%\n", concentrationB)
	fmt.Printf("amount of mixture = %.5g\n", result.AmountMixture)
	fmt.Printf("concentration of mixture = %.5g %%\n", result.ConcentrationMixturePct)
}
