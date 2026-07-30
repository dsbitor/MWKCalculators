// gage looks up wire and sheet metal gage sizes: given a gage
// designation (e.g. "12" or "000000"), it reports the corresponding
// diameter or thickness in inches; given a diameter or thickness, it
// reports the closest matching gage designation.
//
// Converted from GAGE.C (M. W. Klotz), WorkshopUtilities/gage.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"
	"strings"

	"mwkgo/internal/promptio"
	"mwkgo/internal/refdata"
)

// gageEntry is one row of a wire or sheet gage table.
type gageEntry struct {
	Gage   string
	SizeIn float64
}

func loadGages(ctx context.Context, db *sql.DB, table string) ([]gageEntry, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf(`SELECT gage, size_in FROM %s`, table))
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", table, err)
	}
	defer rows.Close()

	var entries []gageEntry
	for rows.Next() {
		var e gageEntry
		if err := rows.Scan(&e.Gage, &e.SizeIn); err != nil {
			return nil, fmt.Errorf("scan %s row: %w", table, err)
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", table, err)
	}
	return entries, nil
}

// findByDesignation returns the entry with the exact gage
// designation, if present.
func findByDesignation(entries []gageEntry, gage string) (gageEntry, bool) {
	for _, e := range entries {
		if e.Gage == gage {
			return e, true
		}
	}
	return gageEntry{}, false
}

// findClosest returns the entry whose size is closest to size.
func findClosest(entries []gageEntry, size float64) gageEntry {
	best := entries[0]
	bestDiff := math.Abs(size - best.SizeIn)
	for _, e := range entries[1:] {
		if diff := math.Abs(size - e.SizeIn); diff < bestDiff {
			best, bestDiff = e, diff
		}
	}
	return best
}

func main() {
	ctx := context.Background()

	db, err := refdata.Open(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gage:", err)
		os.Exit(1)
	}
	defer db.Close()

	wire, err := loadGages(ctx, db, "wire_gages")
	if err != nil {
		fmt.Fprintln(os.Stderr, "gage:", err)
		os.Exit(1)
	}
	sheet, err := loadGages(ctx, db, "sheet_gages")
	if err != nil {
		fmt.Fprintln(os.Stderr, "gage:", err)
		os.Exit(1)
	}

	prompter, err := promptio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gage:", err)
		os.Exit(1)
	}

	fmt.Println("WIRE AND SHEET GAGE UTILITY")
	fmt.Println()
	fmt.Printf("Number of (wire,sheet) data items read = %d,%d\n", len(wire), len(sheet))

	kind := "wire"
	entries := wire
	if strings.ToLower(prompter.Line("[W]ire or (S)heet ? ")) == "s" {
		kind, entries = "sheet", sheet
	}

	if strings.ToLower(prompter.Line("Find (G)age or [S]ize ? ")) == "g" {
		size := prompter.Float("Thickness (in)", 0.1)
		closest := findClosest(entries, size)
		fmt.Printf("CLOSEST %s GAGE = %s WITH SIZE %.4f in\n", strings.ToUpper(kind), closest.Gage, closest.SizeIn)
		return
	}

	designation := strings.TrimSpace(prompter.Line("Gage designation ? "))
	entry, ok := findByDesignation(entries, designation)
	if !ok {
		fmt.Printf("GAGE %s NOT FOUND\n", designation)
		return
	}
	fmt.Printf("SIZE OF %s GAGE %s = %.4f in\n", strings.ToUpper(kind), entry.Gage, entry.SizeIn)
}
