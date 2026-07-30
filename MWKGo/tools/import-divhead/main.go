// import-divhead is a throwaway, run-once tool (see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for Tier
// 2") demonstrating the machine-specific-data half of that split,
// using the legacy DIVHEAD.DAT file as a worked example.
//
// Unlike import-drill's universal reference data, a dividing head's
// worm-gear ratio and hole-circle plates are specific to one
// person's machine, so the original author's own DIVHEAD.DAT is not
// fit to ship as every user's default. This tool instead:
//
//  1. Creates the user database's schema (the divhead_settings and
//     divhead_hole_circles tables) with no rows in it, which is what
//     actually ships.
//  2. Parses DIVHEAD.DAT anyway, and writes its values out as example
//     CSVs (via the same internal/csvtable.Export a user's own
//     eventual "export my current data" command would use), so a new
//     user has a concrete, correctly-shaped template to edit and
//     import with internal/csvtable.Import rather than an empty
//     table and a blank page.
//
// Usage:
//
//	go run ./tools/import-divhead -in DIVHEAD.DAT -userdb userdata.db -examples-dir examples
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"mwkgo/internal/csvtable"
	"mwkgo/internal/legacydat"
	"mwkgo/internal/refdb"
)

var migrations = []refdb.Migration{
	{
		Version:     1,
		Description: "create divhead_settings and divhead_hole_circles tables",
		SQL: `
			CREATE TABLE divhead_settings (
				id                INTEGER PRIMARY KEY CHECK (id = 1),
				worm_gear_ratio   INTEGER NOT NULL,
				rapid_index_holes INTEGER NOT NULL
			);
			CREATE TABLE divhead_hole_circles (
				holes INTEGER NOT NULL
			);`,
	},
}

func main() {
	in := flag.String("in", "DIVHEAD.DAT", "path to the legacy DIVHEAD.DAT file")
	userDB := flag.String("userdb", "userdata.db", "path to create the empty, schema-only user database at")
	examplesDir := flag.String("examples-dir", "examples", "directory to write the example CSVs into")
	flag.Parse()

	if err := run(*in, *userDB, *examplesDir); err != nil {
		log.Fatal(err)
	}
}

func run(inPath, userDBPath, examplesDir string) error {
	ratio, rapidIndexHoles, holeCircles, err := parseDivhead(inPath)
	if err != nil {
		return err
	}

	ctx := context.Background()

	os.Remove(userDBPath)
	userDB, err := refdb.Open(ctx, userDBPath, migrations)
	if err != nil {
		return fmt.Errorf("create schema-only user database %s: %w", userDBPath, err)
	}
	defer userDB.Close()
	fmt.Printf("created empty user database schema at %s (no rows seeded)\n", userDBPath)

	exampleDB, err := refdb.Open(ctx, ":memory:", migrations)
	if err != nil {
		return fmt.Errorf("open in-memory example database: %w", err)
	}
	defer exampleDB.Close()

	const insertSettings = `INSERT INTO divhead_settings (id, worm_gear_ratio, rapid_index_holes) VALUES (1, ?, ?)`
	if _, err := exampleDB.ExecContext(ctx, insertSettings, ratio, rapidIndexHoles); err != nil {
		return fmt.Errorf("seed example settings: %w", err)
	}
	const insertHoles = `INSERT INTO divhead_hole_circles (holes) VALUES (?)`
	for _, holes := range holeCircles {
		if _, err := exampleDB.ExecContext(ctx, insertHoles, holes); err != nil {
			return fmt.Errorf("seed example hole circle %d: %w", holes, err)
		}
	}

	if err := os.MkdirAll(examplesDir, 0o755); err != nil {
		return fmt.Errorf("create examples directory %s: %w", examplesDir, err)
	}
	for _, table := range []string{"divhead_settings", "divhead_hole_circles"} {
		if err := exportTable(ctx, exampleDB, table, filepath.Join(examplesDir, table+".example.csv")); err != nil {
			return err
		}
	}

	fmt.Printf("wrote example CSVs for %d hole circles (ratio %d:1, rapid index %d) from %s into %s\n",
		len(holeCircles), ratio, rapidIndexHoles, inPath, examplesDir)
	return nil
}

func exportTable(ctx context.Context, db *sql.DB, table, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	if err := csvtable.Export(ctx, db, table, f); err != nil {
		return fmt.Errorf("export table %s to %s: %w", table, path, err)
	}
	return nil
}

// parseDivhead reads the worm gear ratio and rapid-indexing-plate
// hole count (the file's first two data lines, in that order) and
// the flat list of hole-circle counts that follows, matching
// DIVHEAD.C's own rdata(): the first data line is always the ratio,
// the second is always the rapid-index plate's hole count, and every
// remaining line is one hole circle's hole count.
func parseDivhead(path string) (ratio, rapidIndexHoles int, holeCircles []int, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, 0, nil, fmt.Errorf("scan %s: %w", path, err)
	}
	if len(rows) < 2 {
		return 0, 0, nil, fmt.Errorf("%s: data region has %d lines, want at least 2 (ratio, rapid index holes)", path, len(rows))
	}

	values := make([]int, len(rows))
	for i, row := range rows {
		fields := legacydat.Fields(row, ",\t;")
		if len(fields) == 0 {
			return 0, 0, nil, fmt.Errorf("%s line %d: %q has no value", path, i+1, row)
		}
		v, err := strconv.Atoi(fields[0])
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%s line %d: %q: %w", path, i+1, fields[0], err)
		}
		values[i] = v
	}

	return values[0], values[1], values[2:], nil
}
