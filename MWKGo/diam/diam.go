// diam reports, for a chosen material and each of the operator's
// machine's available spindle speeds, the range of tool or workpiece
// diameters that keep the cutting speed within that material's
// recommended surface-speed range — the inverse of what speed does
// (which starts from a diameter and finds a recommended rpm range).
//
// The material reference table is universal data, shared with speed
// in the same conversion group (both read the reference database's
// machining_speeds table). The machine's available spindle speeds are
// specific to one owner's equipment, so that part reads from the
// user's own database instead; see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2".
//
// Converted from DIAM.C (M. W. Klotz), WorkshopUtilities/diam.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"math"
	"os"

	"mwkgo/internal/csvtable"
	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
	"mwkgo/internal/userdata"
)

const availableSpeedsTable = "diam_available_speeds"

// material is one row of the machining_speeds reference table (the
// same table speed.go reads).
type material struct {
	Name              string
	LowSFPM, HighSFPM int
}

// diameterRange returns the range of tool or workpiece diameters
// that keep the cutting speed within m's recommended surface-speed
// range at the given spindle speed (rpm), the inverse of the
// surface-speed relation speed.go's rpmRange uses.
func diameterRange(m material, rpm float64) (low, high float64) {
	c := 12 / math.Pi
	return c * float64(m.LowSFPM) / rpm, c * float64(m.HighSFPM) / rpm
}

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

func loadAvailableSpeeds(ctx context.Context, db *sql.DB) ([]float64, error) {
	rows, err := db.QueryContext(ctx, `SELECT rpm FROM `+availableSpeedsTable)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", availableSpeedsTable, err)
	}
	defer rows.Close()

	var speeds []float64
	for rows.Next() {
		var s float64
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("scan speed: %w", err)
		}
		speeds = append(speeds, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", availableSpeedsTable, err)
	}
	return speeds, nil
}

func main() {
	importPath := flag.String("import-speeds", "", "import available machine spindle speeds from a CSV file (one rpm column) and exit")
	flag.Parse()

	ctx := context.Background()

	refDB, err := refdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diam:", err)
		os.Exit(1)
	}
	defer refDB.Close()

	userDB, err := userdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diam:", err)
		os.Exit(1)
	}
	defer userDB.Close()

	if *importPath != "" {
		if err := importAvailableSpeeds(ctx, userDB, *importPath); err != nil {
			fmt.Fprintln(os.Stderr, "diam:", err)
			os.Exit(1)
		}
		return
	}

	materials, err := loadMaterials(ctx, refDB)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diam:", err)
		os.Exit(1)
	}
	speeds, err := loadAvailableSpeeds(ctx, userDB)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diam:", err)
		os.Exit(1)
	}
	if len(speeds) == 0 {
		fmt.Println("No available machine spindle speeds configured yet.")
		fmt.Println("This is machine-specific data (the speed steps your own lathe or mill")
		fmt.Println("offers), so nothing is pre-loaded. Import a CSV with an rpm column:")
		fmt.Println()
		fmt.Println("    diam -import-speeds my-machine-speeds.csv")
		fmt.Println()
		fmt.Println("See docs/calculators/diam.md for the expected format.")
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "diam:", err)
		os.Exit(1)
	}

	fmt.Println("MACHINING DIAMETER UTILITY")
	fmt.Println()
	fmt.Printf("number of data entries read = %d\n\n", len(materials))
	for i, m := range materials {
		fmt.Printf("%4d  %s\n", i+1, m.Name)
	}
	fmt.Println()

	chosen := promptForMaterial(prompter, materials)

	fmt.Printf("\nMaterial = %s\n\n", chosen.Name)
	fmt.Printf("%6s  %6s  %6s\n", "SPEED", "DMIN", "DMAX")
	fmt.Printf("%6s  %6s  %6s\n\n", " rpm ", " in ", " in ")
	for _, rpm := range speeds {
		low, high := diameterRange(chosen, rpm)
		fmt.Printf("%6.0f  %6.3f  %6.3f\n", rpm, low, high)
	}
}

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

func importAvailableSpeeds(ctx context.Context, db *sql.DB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	if err := csvtable.Import(ctx, db, availableSpeedsTable, f); err != nil {
		return fmt.Errorf("import %s: %w", path, err)
	}
	fmt.Printf("imported available speeds from %s\n", path)
	return nil
}
