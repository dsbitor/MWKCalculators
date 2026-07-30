// expand computes thermal expansion: given a material (or a
// user-supplied coefficient), a nominal length, and either a
// temperature change or a length change, it computes the other one.
//
// Converted from EXPAND.C (M. W. Klotz), WorkshopUtilities/expand.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
)

// material is one row of the coefficient-of-thermal-expansion
// reference table.
type material struct {
	Name          string
	CTEPPMPerDegF float64
}

// lengthChange returns the change in length for an object of the
// given nominal length and coefficient of thermal expansion (as a
// fraction, not ppm) subjected to a temperature change.
func lengthChange(cte, nominalLength, deltaT float64) float64 {
	return cte * nominalLength * deltaT
}

// temperatureChange returns the temperature change that would
// produce the given length change, the inverse of lengthChange.
func temperatureChange(cte, nominalLength, deltaLength float64) float64 {
	return deltaLength / (cte * nominalLength)
}

// loadMaterials returns every row of the materials reference table,
// alphabetically (case-insensitively), matching EXPAND.C's own
// sorted display order.
func loadMaterials(ctx context.Context, db *sql.DB) ([]material, error) {
	rows, err := db.QueryContext(ctx, `SELECT name, cte_ppm_per_degf FROM materials ORDER BY name COLLATE NOCASE`)
	if err != nil {
		return nil, fmt.Errorf("query materials: %w", err)
	}
	defer rows.Close()

	var materials []material
	for rows.Next() {
		var m material
		if err := rows.Scan(&m.Name, &m.CTEPPMPerDegF); err != nil {
			return nil, fmt.Errorf("scan material: %w", err)
		}
		materials = append(materials, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate materials: %w", err)
	}
	return materials, nil
}

func main() {
	ctx := context.Background()

	db, err := refdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "expand:", err)
		os.Exit(1)
	}
	defer db.Close()

	materials, err := loadMaterials(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "expand:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "expand:", err)
		os.Exit(1)
	}

	fmt.Println("MATERIAL EXPANSION CALCULATIONS")
	fmt.Println()
	fmt.Printf("Number of data items read = %d\n", len(materials))

	chosen := promptForMaterial(prompter, materials)
	length := prompter.Float("Nominal length of object (in)", 1.0)

	fmt.Println()
	fmt.Println("A. Find length change given temperature change")
	fmt.Println("B. Find temperature change given length change")

	cte := chosen.CTEPPMPerDegF * 1e-6
	var deltaT, deltaLength float64
	if strings.ToLower(prompter.Line("([A],B) ? ")) == "b" {
		deltaLength = prompter.Float("Length change (in)", 0.0001)
		deltaT = temperatureChange(cte, length, deltaLength)
	} else {
		deltaT = prompter.Float("Temperature change (degF)", 100.0)
		deltaLength = lengthChange(cte, length, deltaT)
	}

	fmt.Printf("\nCTE of %s = %.4f ppm/degF\n", chosen.Name, cte*1e6)
	fmt.Printf("Nominal length = %.4f in\n", length)
	fmt.Printf("Length change = %.6f in\n", deltaLength)
	fmt.Printf("Temperature change = %.4f degF\n", deltaT)
}

// promptForMaterial lists every material by number and prompts for a
// choice; a number outside the list (including the offered "one past
// the end" option) instead asks for a custom coefficient, matching
// EXPAND.C's own mselect(): the material's name is then reported as
// "??" as the original does, since it has no name of its own.
func promptForMaterial(p *promptio.Prompter, materials []material) material {
	for i, m := range materials {
		fmt.Printf("%d  %s\n", i+1, m.Name)
	}
	customOption := len(materials) + 1
	fmt.Printf("%d  User input\n", customOption)

	n := p.Int("\nMaterial number", customOption)
	if n < 1 || n > len(materials) {
		cte := p.Float("Material coefficient of thermal expansion (ppm/degF)", 10.0)
		return material{Name: "??", CTEPPMPerDegF: cte}
	}
	return materials[n-1]
}
