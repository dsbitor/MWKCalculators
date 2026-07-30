// speed recommends a spindle RPM range for a given material and
// cutting-tool (or workpiece) diameter, from a reference table of
// recommended cutting speeds (surface feet per minute) by material.
//
// Converted from SPEED.C (M. W. Klotz), WorkshopUtilities/speed.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
)

// material is one row of the machining-speed reference table.
type material struct {
	Name              string
	LowSFPM, HighSFPM int // recommended cutting speed range, surface feet per minute
}

// rpmRange converts a recommended surface-speed range (in surface
// feet per minute) to a recommended spindle RPM range for the given
// tool or workpiece diameter, via the standard cutting-speed
// relation: surface speed = pi * diameter * rpm / 12 (diameter in
// inches, surface speed in feet per minute), solved for rpm.
func rpmRange(m material, diameterIn float64) (low, high float64) {
	c := 12 / math.Pi
	return c * float64(m.LowSFPM) / diameterIn, c * float64(m.HighSFPM) / diameterIn
}

// loadMaterials returns every row of the machining_speeds reference
// table, in the original data file's own order (list_position),
// since the material chosen by number depends on that order matching
// what the user is shown.
func loadMaterials(ctx context.Context, db *sql.DB) ([]material, error) {
	rows, err := db.QueryContext(ctx, `SELECT material, low_sfpm, high_sfpm FROM machining_speeds ORDER BY list_position`)
	if err != nil {
		return nil, fmt.Errorf("query machining_speeds: %w", err)
	}
	defer rows.Close()

	var materials []material
	for rows.Next() {
		var m material
		if err := rows.Scan(&m.Name, &m.LowSFPM, &m.HighSFPM); err != nil {
			return nil, fmt.Errorf("scan material: %w", err)
		}
		materials = append(materials, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate machining_speeds: %w", err)
	}
	return materials, nil
}

func main() {
	ctx := context.Background()

	db, err := refdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "speed:", err)
		os.Exit(1)
	}
	defer db.Close()

	materials, err := loadMaterials(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "speed:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "speed:", err)
		os.Exit(1)
	}

	fmt.Println("MACHINING SPEED UTILITY")
	fmt.Println()
	fmt.Printf("number of data entries read = %d\n\n", len(materials))
	for i, m := range materials {
		fmt.Printf("%4d  %s\n", i+1, m.Name)
	}
	fmt.Println()

	chosen := promptForMaterial(prompter, materials)

	fmt.Println()
	fmt.Println("Input relevant diameter.")
	fmt.Println("Diameter of work if lathe, diameter of cutter/drill if mill/drillpress.")
	diameter := prompter.Float("Diameter (in)", 1.0)

	low, high := rpmRange(chosen, diameter)
	fmt.Printf("\nFor %s and diameter %.4f in:\n", chosen.Name, diameter)
	fmt.Printf("Recommended rpm between %.0f and %.0f\n", low, high)
}

// promptForMaterial repeatedly prompts for a material number until a
// valid one (1-based, within range) is entered.
func promptForMaterial(p *promptio.Prompter, materials []material) material {
	for {
		n := p.Int("Material Number", 1)
		if n < 1 || n > len(materials) {
			fmt.Println("INVALID NUMBER")
			continue
		}
		return materials[n-1]
	}
}
