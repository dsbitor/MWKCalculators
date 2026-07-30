// fits computes the shaft diameter to machine for a standard
// shaft/hole fit class (shrink, force, drive, push, slide, running,
// and clearance fits), given the nominal (hole) diameter, using a
// small reference table of per-fit-class allowances.
//
// Converted from FITS.C (M. W. Klotz), WorkshopUtilities/fits.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
)

// fit is one row of the shaft/hole fit reference table.
type fit struct {
	Name               string
	ConstantThou       float64 // constant offset, thousandths of an inch
	AllowanceThouPerIn float64 // allowance, thousandths of an inch per inch of diameter
}

// shaftDiameter returns the shaft diameter to machine for f, given
// the nominal (hole) diameter, following the original program's own
// formula.
func shaftDiameter(f fit, nominalDiameter float64) float64 {
	return nominalDiameter + 0.001*(f.AllowanceThouPerIn*nominalDiameter+f.ConstantThou)
}

// loadFits returns every row of the fits reference table, in the
// original data file's own order (list_position), since the fit
// chosen by number depends on that order matching what the user is
// shown.
func loadFits(ctx context.Context, db *sql.DB) ([]fit, error) {
	rows, err := db.QueryContext(ctx, `SELECT name, constant_thou, allowance_thou_per_in FROM fits ORDER BY list_position`)
	if err != nil {
		return nil, fmt.Errorf("query fits: %w", err)
	}
	defer rows.Close()

	var fits []fit
	for rows.Next() {
		var f fit
		if err := rows.Scan(&f.Name, &f.ConstantThou, &f.AllowanceThouPerIn); err != nil {
			return nil, fmt.Errorf("scan fit: %w", err)
		}
		fits = append(fits, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate fits: %w", err)
	}
	return fits, nil
}

func main() {
	ctx := context.Background()

	db, err := refdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fits:", err)
		os.Exit(1)
	}
	defer db.Close()

	fits, err := loadFits(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fits:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fits:", err)
		os.Exit(1)
	}

	fmt.Println("SHAFT/HOLE FIT COMPUTATIONS")
	fmt.Println()
	fmt.Printf("Number of data items read = %d\n", len(fits))
	for i, f := range fits {
		fmt.Printf("%2d  %s\n", i+1, f.Name)
	}
	fmt.Println()

	chosen := promptForFit(prompter, fits)
	diameter := prompter.Float("Nominal diameter of shaft (in)", 1.0)
	dc := shaftDiameter(chosen, diameter)
	fmt.Printf("\nDiameter of shaft for %s Fit = %.5f in\n", chosen.Name, dc)
}

// promptForFit repeatedly prompts for a fit number until a valid one
// (1-based, within range) is entered.
func promptForFit(p *promptio.Prompter, fits []fit) fit {
	for {
		n := p.Int("Fit desired", 4)
		if n < 1 || n > len(fits) {
			fmt.Println("OPTION SELECTION ERROR")
			continue
		}
		return fits[n-1]
	}
}
