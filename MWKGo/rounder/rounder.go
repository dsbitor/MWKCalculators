// rounder generates a table of ball-end mill positions for rounding
// over an edge on the milling machine: successive scallop cuts, each
// tangent to the desired radius profile, are made by stepping the
// mill through an angle theta and positioning its center at the
// corresponding table coordinates.
//
// Converted from ROUNDER.C (M. W. Klotz), WorkshopUtilities/rounder.
// The original writes its results to ROUNDER.OUT via fopen; this
// conversion prints to stdout instead and drops that DOS
// file-save-then-page convenience, the same approach used for loan
// (Tier 1 suitability review, Finding 5).
package main

import (
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
)

// rounderRow is one angular step's mill-center offsets, using the
// same A/B/C/D labels as the original program's table header.
type rounderRow struct {
	AngleDeg   float64
	A, B, C, D float64
}

// rounderTable returns the mill-center offset table from theta=0 to
// theta=90 degrees in steps of stepDeg, for a workpiece radius plus
// ball mill radius of combinedRadius.
func rounderTable(combinedRadius, stepDeg float64) []rounderRow {
	var rows []rounderRow
	for theta := 0.0; theta <= 90.0; theta += stepDeg {
		sinT := math.Sin(theta * math.Pi / 180)
		cosT := math.Cos(theta * math.Pi / 180)
		rows = append(rows, rounderRow{
			AngleDeg: theta,
			A:        combinedRadius * sinT,
			B:        combinedRadius * cosT,
			C:        combinedRadius * (1 - sinT),
			D:        combinedRadius * (1 - cosT),
		})
	}
	return rows
}

func main() {
	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rounder:", err)
		os.Exit(1)
	}

	fmt.Println("ROUNDING OVER CALCULATIONS")
	fmt.Println()

	workpieceRadius := prompter.Float("Workpiece radius (in)", 1.0)
	millDiameter := prompter.Float("Ball-end mill diameter (in)", 0.25)
	stepDeg := prompter.Float("Angular increment (deg)", 5.0)
	combinedRadius := workpieceRadius + 0.5*millDiameter

	fmt.Printf("Workpiece radius = %.3f in\n", workpieceRadius)
	fmt.Printf("Ball mill diameter = %.3f in\n", millDiameter)
	fmt.Println("A = (R+r)*sin(theta)")
	fmt.Println("B = (R+r)*cos(theta)")
	fmt.Println("C = (R+r)*(1 - sin(theta))")
	fmt.Println("D = (R+r)*(1 - cos(theta))")
	fmt.Println()
	fmt.Printf("%6s  %6s  %6s  %6s  %6s\n", "ANGLE", "A", "B", "C", "D")
	for _, row := range rounderTable(combinedRadius, stepDeg) {
		fmt.Printf("%6.1f  %6.3f  %6.3f  %6.3f  %6.3f\n", row.AngleDeg, row.A, row.B, row.C, row.D)
	}
}
