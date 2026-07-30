// build-refdb is the offline, throwaway tool (see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2") that builds the single golden reference.db every
// universal-reference-data calculator embeds via internal/refdata.
// It runs by hand, once per change, reading each contributing
// program's original .DAT file straight from its checked-in source
// tree location and writing its rows into that program's own table.
//
// Adding a new reference-bucket program means adding one entry to
// the tables slice in run() below: a migration creating its table,
// and a seed function that parses its .DAT file with
// internal/legacydat and inserts the rows. Nothing else in this file
// changes.
//
// Usage:
//
//	go run ./tools/build-refdb -out internal/refdata/reference.db
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"mwkgo/internal/legacydat"
	"mwkgo/internal/refdb"
)

func main() {
	out := flag.String("out", "internal/refdata/reference.db", "path to write the golden reference database to")
	fitsDAT := flag.String("fits-dat", "../MWKC/WorkshopUtilities/fits/FITS.DAT", "path to FITS.DAT")
	speedDAT := flag.String("speed-dat", "../MWKC/WorkshopUtilities/speed/SPEED.DAT", "path to SPEED.DAT")
	gageDAT := flag.String("gage-dat", "../MWKC/WorkshopUtilities/gage/GAGE.DAT", "path to GAGE.DAT")
	expandDAT := flag.String("expand-dat", "../MWKC/WorkshopUtilities/expand/EXPAND.DAT", "path to EXPAND.DAT")
	findthrdDAT := flag.String("findthrd-dat", "../MWKC/WorkshopUtilities/findthrd/FINDTHRD.DAT", "path to FINDTHRD.DAT")
	weightDAT := flag.String("weight-dat", "../MWKC/WorkshopUtilities/weight/WEIGHT.DAT", "path to WEIGHT.DAT")
	unitDAT := flag.String("unit-dat", "../MWKC/Misc/unit/UNIT.DAT", "path to UNIT.DAT")
	drillDAT := flag.String("drill-dat", "../MWKC/WorkshopUtilities/drill/DRILL.DAT", "path to DRILL.DAT")
	flag.Parse()

	if err := run(*out, *fitsDAT, *speedDAT, *gageDAT, *expandDAT, *findthrdDAT, *weightDAT, *unitDAT, *drillDAT); err != nil {
		log.Fatal(err)
	}
}

func run(outPath, fitsDATPath, speedDATPath, gageDATPath, expandDATPath, findthrdDATPath, weightDATPath, unitDATPath, drillDATPath string) error {
	tables := []struct {
		name      string
		migration refdb.Migration
		seed      func(ctx context.Context, db *sql.DB) (int, error)
	}{
		{
			name: "fits",
			migration: refdb.Migration{
				Version:     1,
				Description: "create fits table",
				// list_position preserves FITS.DAT's own row
				// order (SQLite assigns it 1, 2, 3, ... in
				// insertion order as an alias for the rowid):
				// the original program numbers fits by that
				// same order for its menu, and its own default
				// selection (a "push fit") depends on it being
				// entry 4, not on alphabetical order.
				SQL: `CREATE TABLE fits (
					list_position         INTEGER PRIMARY KEY,
					name                  TEXT NOT NULL,
					constant_thou         REAL NOT NULL,
					allowance_thou_per_in REAL NOT NULL
				)`,
			},
			seed: func(ctx context.Context, db *sql.DB) (int, error) {
				return seedFits(ctx, db, fitsDATPath)
			},
		},
		{
			name: "speed",
			migration: refdb.Migration{
				Version:     2,
				Description: "create machining_speeds table",
				// list_position preserves SPEED.DAT's own row
				// order for the same reason as fits.list_position
				// above: SPEED.C numbers materials by file order,
				// and its default selection (material 1) depends
				// on it.
				SQL: `CREATE TABLE machining_speeds (
					list_position INTEGER PRIMARY KEY,
					material      TEXT    NOT NULL,
					low_sfpm      INTEGER NOT NULL,
					high_sfpm     INTEGER NOT NULL
				)`,
			},
			seed: func(ctx context.Context, db *sql.DB) (int, error) {
				return seedSpeed(ctx, db, speedDATPath)
			},
		},
		{
			name: "gage",
			migration: refdb.Migration{
				Version:     3,
				Description: "create wire_gages and sheet_gages tables",
				SQL: `
					CREATE TABLE wire_gages (
						gage    TEXT NOT NULL,
						size_in REAL NOT NULL
					);
					CREATE TABLE sheet_gages (
						gage    TEXT NOT NULL,
						size_in REAL NOT NULL
					);`,
			},
			seed: func(ctx context.Context, db *sql.DB) (int, error) {
				return seedGage(ctx, db, gageDATPath)
			},
		},
		{
			name: "expand",
			migration: refdb.Migration{
				Version:     4,
				Description: "create materials table",
				SQL: `CREATE TABLE materials (
					name             TEXT NOT NULL,
					cte_ppm_per_degf REAL NOT NULL
				)`,
			},
			seed: func(ctx context.Context, db *sql.DB) (int, error) {
				return seedExpand(ctx, db, expandDATPath)
			},
		},
		{
			name: "findthrd",
			migration: refdb.Migration{
				Version:     5,
				Description: "create threads table",
				SQL: `CREATE TABLE threads (
					name      TEXT NOT NULL,
					diam_in   REAL NOT NULL,
					diam_mm   REAL NOT NULL,
					pitch_tpi REAL NOT NULL,
					pitch_mm  REAL NOT NULL
				)`,
			},
			seed: func(ctx context.Context, db *sql.DB) (int, error) {
				return seedFindthrd(ctx, db, findthrdDATPath)
			},
		},
		{
			name: "weight",
			migration: refdb.Migration{
				Version:     6,
				Description: "create weight_materials table",
				SQL: `CREATE TABLE weight_materials (
					name              TEXT NOT NULL,
					density_lb_per_in3 REAL NOT NULL
				)`,
			},
			seed: func(ctx context.Context, db *sql.DB) (int, error) {
				return seedWeight(ctx, db, weightDATPath)
			},
		},
		{
			name: "unit",
			migration: refdb.Migration{
				Version:     7,
				Description: "create unit_prefixes and unit_definitions tables",
				SQL: `
					CREATE TABLE unit_prefixes (
						name  TEXT NOT NULL,
						value REAL NOT NULL
					);
					CREATE TABLE unit_definitions (
						name             TEXT    NOT NULL,
						factor           REAL    NOT NULL,
						dim_length       INTEGER NOT NULL,
						dim_mass         INTEGER NOT NULL,
						dim_time         INTEGER NOT NULL,
						dim_angle        INTEGER NOT NULL,
						dim_solidangle   INTEGER NOT NULL,
						dim_charge       INTEGER NOT NULL,
						dim_amount       INTEGER NOT NULL
					);`,
			},
			seed: func(ctx context.Context, db *sql.DB) (int, error) {
				prefixes, units, err := seedUnit(ctx, db, unitDATPath)
				return prefixes + units, err
			},
		},
		{
			name: "drill",
			migration: refdb.Migration{
				Version:     8,
				Description: "create drills table",
				SQL: `CREATE TABLE drills (
					name        TEXT NOT NULL,
					diameter_in REAL NOT NULL
				)`,
			},
			seed: func(ctx context.Context, db *sql.DB) (int, error) {
				return seedDrill(ctx, db, drillDATPath)
			},
		},
	}

	migrations := make([]refdb.Migration, len(tables))
	for i, t := range tables {
		migrations[i] = t.migration
	}

	ctx := context.Background()
	os.Remove(outPath) // a re-run should produce a clean file, not append to a stale one
	db, err := refdb.Open(ctx, outPath, migrations)
	if err != nil {
		return fmt.Errorf("open %s: %w", outPath, err)
	}
	defer db.Close()

	for _, t := range tables {
		n, err := t.seed(ctx, db)
		if err != nil {
			return fmt.Errorf("seed table for %s: %w", t.name, err)
		}
		fmt.Printf("%s: inserted %d rows\n", t.name, n)
	}

	fmt.Printf("wrote %s\n", outPath)
	return nil
}

// seedFits parses FITS.DAT (fit name, constant in thousandths of an
// inch, allowance in thousandths of an inch per inch of diameter,
// comma-separated, bracketed by STARTOFDATA/ENDOFDATA) and inserts
// each row into the fits table.
func seedFits(ctx context.Context, db *sql.DB, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, fmt.Errorf("scan %s: %w", path, err)
	}

	const insertSQL = `INSERT INTO fits (name, constant_thou, allowance_thou_per_in) VALUES (?, ?, ?)`
	for i, row := range rows {
		fields := legacydat.Fields(row, ",;")
		if len(fields) != 3 {
			return 0, fmt.Errorf("%s line %d: %q does not split into 3 comma-separated fields", path, i+1, row)
		}
		constant, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return 0, fmt.Errorf("%s line %d: constant %q: %w", path, i+1, fields[1], err)
		}
		allowance, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return 0, fmt.Errorf("%s line %d: allowance %q: %w", path, i+1, fields[2], err)
		}
		if _, err := db.ExecContext(ctx, insertSQL, fields[0], constant, allowance); err != nil {
			return 0, fmt.Errorf("insert fit %q: %w", fields[0], err)
		}
	}
	return len(rows), nil
}

// seedSpeed parses SPEED.DAT (material name, low sfpm, high sfpm,
// tab-separated, with no STARTOFDATA marker) and inserts each row
// into the machining_speeds table.
func seedSpeed(ctx context.Context, db *sql.DB, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	// SPEED.DAT has no STARTOFDATA marker: its data region starts at
	// the first non-comment line and runs to ENDOFDATA, matching
	// SPEED.C's own rdata(), which never waits for a start marker.
	rows, err := legacydat.Rows(f, legacydat.Options{EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, fmt.Errorf("scan %s: %w", path, err)
	}

	const insertSQL = `INSERT INTO machining_speeds (material, low_sfpm, high_sfpm) VALUES (?, ?, ?)`
	for i, row := range rows {
		fields := legacydat.Fields(row, "\t;")
		if len(fields) != 3 {
			return 0, fmt.Errorf("%s line %d: %q does not split into 3 tab-separated fields", path, i+1, row)
		}
		low, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, fmt.Errorf("%s line %d: low sfpm %q: %w", path, i+1, fields[1], err)
		}
		high, err := strconv.Atoi(fields[2])
		if err != nil {
			return 0, fmt.Errorf("%s line %d: high sfpm %q: %w", path, i+1, fields[2], err)
		}
		if _, err := db.ExecContext(ctx, insertSQL, fields[0], low, high); err != nil {
			return 0, fmt.Errorf("insert material %q: %w", fields[0], err)
		}
	}
	return len(rows), nil
}

// seedGage parses GAGE.DAT (gage designation, size in inches,
// comma-separated, with no STARTOFDATA marker) into the wire_gages
// and sheet_gages tables. The two tables share one file, separated
// by a sentinel row whose gage designation is literally "xx",
// matching GAGE.C's own rdata(), which switches from wire to sheet
// data the same way.
func seedGage(ctx context.Context, db *sql.DB, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, fmt.Errorf("scan %s: %w", path, err)
	}

	const insertWireSQL = `INSERT INTO wire_gages (gage, size_in) VALUES (?, ?)`
	const insertSheetSQL = `INSERT INTO sheet_gages (gage, size_in) VALUES (?, ?)`

	section := 0 // 0 = wire, 1 = sheet; advances past each "xx" sentinel row
	n := 0
	for i, row := range rows {
		fields := legacydat.Fields(row, ",")
		if len(fields) < 2 {
			return 0, fmt.Errorf("%s line %d: %q does not split into gage and size", path, i+1, row)
		}
		if fields[0] == "xx" {
			section++
			continue
		}
		size, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return 0, fmt.Errorf("%s line %d: size %q: %w", path, i+1, fields[1], err)
		}
		insertSQL := insertWireSQL
		if section == 1 {
			insertSQL = insertSheetSQL
		}
		if _, err := db.ExecContext(ctx, insertSQL, fields[0], size); err != nil {
			return 0, fmt.Errorf("insert gage %q: %w", fields[0], err)
		}
		n++
	}
	return n, nil
}

// seedExpand parses EXPAND.DAT (material name, coefficient of linear
// thermal expansion in ppm/degF, comma-separated, bracketed by
// STARTOFDATA/ENDOFDATA) and inserts each row into the materials
// table.
func seedExpand(ctx context.Context, db *sql.DB, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, fmt.Errorf("scan %s: %w", path, err)
	}

	const insertSQL = `INSERT INTO materials (name, cte_ppm_per_degf) VALUES (?, ?)`
	for i, row := range rows {
		fields := legacydat.Fields(row, ",;")
		if len(fields) != 2 {
			return 0, fmt.Errorf("%s line %d: %q does not split into name and cte", path, i+1, row)
		}
		cte, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return 0, fmt.Errorf("%s line %d: cte %q: %w", path, i+1, fields[1], err)
		}
		if _, err := db.ExecContext(ctx, insertSQL, fields[0], cte); err != nil {
			return 0, fmt.Errorf("insert material %q: %w", fields[0], err)
		}
	}
	return len(rows), nil
}

// seedWeight parses WEIGHT.DAT (material name, density in lb/in^3,
// comma-separated, bracketed by STARTOFDATA/ENDOFDATA) and inserts
// each row into the weight_materials table. WEIGHT.DAT's own rows are
// already alphabetical, matching WEIGHT.C's own rdata(), which sorts
// them before use; this conversion queries with ORDER BY name rather
// than relying on insertion order, so the file's own order is not
// load-bearing.
func seedWeight(ctx context.Context, db *sql.DB, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, fmt.Errorf("scan %s: %w", path, err)
	}

	const insertSQL = `INSERT INTO weight_materials (name, density_lb_per_in3) VALUES (?, ?)`
	for i, row := range rows {
		fields := legacydat.Fields(row, ",;")
		if len(fields) != 2 {
			return 0, fmt.Errorf("%s line %d: %q does not split into name and density", path, i+1, row)
		}
		density, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return 0, fmt.Errorf("%s line %d: density %q: %w", path, i+1, fields[1], err)
		}
		if _, err := db.ExecContext(ctx, insertSQL, fields[0], density); err != nil {
			return 0, fmt.Errorf("insert material %q: %w", fields[0], err)
		}
	}
	return len(rows), nil
}

// fixedField returns the substring of s from start, length runes
// long, trimmed of surrounding whitespace, or "" if start is at or
// past the end of s. Unlike FINDTHRD.C's own strncpy-based column
// extraction (which, for a line shorter than a given column's
// offset, copies whatever leftover bytes happened to still be in its
// reused line buffer from the previous, longer line — undefined
// behavior no Go conversion should try to reproduce), a short line's
// missing trailing columns are treated as genuinely absent here, not
// as garbage.
func fixedField(s string, start, length int) string {
	if start >= len(s) {
		return ""
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	return strings.TrimSpace(s[start:end])
}

// parseFixedFieldFloat parses a fixed-width column already extracted
// by fixedField, returning 0 for a blank column (matching atof("")
// == 0) rather than treating it as an error, since FINDTHRD.DAT's
// shorter entries genuinely omit trailing columns (no metric pitch
// recorded for some thread standards, for instance).
func parseFixedFieldFloat(field string) (float64, error) {
	if field == "" {
		return 0, nil
	}
	return strconv.ParseFloat(field, 64)
}

// seedFindthrd parses FINDTHRD.DAT's fixed-width columns (name in
// [0,12), diameter in inches in [12,18), diameter in mm in [20,27),
// pitch in tpi in [28,34), pitch in mm in [35,41); two more columns
// follow for core diameter and thread depth, which FINDTHRD.C itself
// reads into its th[] array but never uses, so they are not carried
// over here either) into the threads table.
func seedFindthrd(ctx context.Context, db *sql.DB, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, fmt.Errorf("scan %s: %w", path, err)
	}

	const insertSQL = `INSERT INTO threads (name, diam_in, diam_mm, pitch_tpi, pitch_mm) VALUES (?, ?, ?, ?, ?)`
	for i, row := range rows {
		name := fixedField(row, 0, 12)
		diamIn, err := parseFixedFieldFloat(fixedField(row, 12, 6))
		if err != nil {
			return 0, fmt.Errorf("%s line %d: diameter (in) in %q: %w", path, i+1, row, err)
		}
		diamMM, err := parseFixedFieldFloat(fixedField(row, 20, 7))
		if err != nil {
			return 0, fmt.Errorf("%s line %d: diameter (mm) in %q: %w", path, i+1, row, err)
		}
		pitchTPI, err := parseFixedFieldFloat(fixedField(row, 28, 6))
		if err != nil {
			return 0, fmt.Errorf("%s line %d: pitch (tpi) in %q: %w", path, i+1, row, err)
		}
		pitchMM, err := parseFixedFieldFloat(fixedField(row, 35, 6))
		if err != nil {
			return 0, fmt.Errorf("%s line %d: pitch (mm) in %q: %w", path, i+1, row, err)
		}
		if _, err := db.ExecContext(ctx, insertSQL, name, diamIn, diamMM, pitchTPI, pitchMM); err != nil {
			return 0, fmt.Errorf("insert thread %q: %w", name, err)
		}
	}
	return len(rows), nil
}

// unitFact decodes UNIT.DAT's own conversion-factor syntax: either a
// plain number, or "1/x" (UNIT.C's own fact() only ever inverts the
// text after the slash — a numerator other than 1 is never actually
// used in the shipped data, so this doesn't generalize to arbitrary
// a/b fractions, matching the original exactly).
func unitFact(s string) (float64, error) {
	if idx := strings.IndexByte(s, '/'); idx >= 0 {
		denom, err := strconv.ParseFloat(strings.TrimSpace(s[idx+1:]), 64)
		if err != nil {
			return 0, fmt.Errorf("denominator %q: %w", s[idx+1:], err)
		}
		return 1 / denom, nil
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("factor %q: %w", s, err)
	}
	return v, nil
}

// seedDrill parses DRILL.DAT (drill designation, size in inches,
// tab-separated, bracketed by STARTOFDATA/ENDOFDATA) and inserts each
// row into the drills table. DRILL.C itself sorts by size before use;
// this conversion queries with ORDER BY diameter_in instead, so the
// file's own order is not load-bearing.
func seedDrill(ctx context.Context, db *sql.DB, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, fmt.Errorf("scan %s: %w", path, err)
	}

	const insertSQL = `INSERT INTO drills (name, diameter_in) VALUES (?, ?)`
	for i, row := range rows {
		fields := legacydat.Fields(row, "\t")
		if len(fields) != 2 {
			return 0, fmt.Errorf("%s line %d: %q does not split into name and size", path, i+1, row)
		}
		size, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return 0, fmt.Errorf("%s line %d: size %q: %w", path, i+1, fields[1], err)
		}
		if _, err := db.ExecContext(ctx, insertSQL, strings.TrimSpace(fields[0]), size); err != nil {
			return 0, fmt.Errorf("insert drill %q: %w", fields[0], err)
		}
	}
	return len(rows), nil
}

// seedUnit parses UNIT.DAT's own section-based format (BEGINPREFIX,
// BEGINPRIMARY, BEGINMIXED, ENDOFDATA, each line tokenized on
// "=;,") into the unit_prefixes and unit_definitions tables, matching
// UNIT.C's own readdata() state machine exactly.
//
// Inside a BEGINPRIMARY block, a unit line gives only its name and
// factor; its 7 dimension exponents are inherited from the most
// recent NEWUNIT=d0,...,d6 directive, not given per line. Inside a
// BEGINMIXED block, every field is explicit on the line itself
// (NAME=factor,d0,...,d6) — except the shipped data's own NEWTONMETER
// entry, which is missing its 7th dimension field; any BEGINMIXED
// line short of 7 dimension fields has its remaining fields treated
// as 0, matching what the original's own static, zero-initialized
// unit[] array would have left there.
func seedUnit(ctx context.Context, db *sql.DB, path string) (prefixCount, unitCount int, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{EndMarker: "ENDOFDATA"})
	if err != nil {
		return 0, 0, fmt.Errorf("scan %s: %w", path, err)
	}

	const insertPrefixSQL = `INSERT INTO unit_prefixes (name, value) VALUES (?, ?)`
	const insertUnitSQL = `INSERT INTO unit_definitions
		(name, factor, dim_length, dim_mass, dim_time, dim_angle, dim_solidangle, dim_charge, dim_amount)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	const dtypeUndefined, dtypePrefix, dtypePrimary, dtypeMixed = -1, 0, 1, 2
	dtype := dtypeUndefined
	var currentDims [7]int

	for i, row := range rows {
		switch {
		case strings.Contains(row, "BEGINPREFIX"):
			dtype = dtypePrefix
			continue
		case strings.Contains(row, "BEGINPRIMARY"):
			dtype = dtypePrimary
			continue
		case strings.Contains(row, "BEGINMIXED"):
			dtype = dtypeMixed
			continue
		}

		fields := legacydat.Fields(row, "=;,")
		for k, field := range fields {
			fields[k] = strings.TrimSpace(field)
		}

		switch dtype {
		case dtypePrefix:
			if !strings.Contains(row, "=") || len(fields) < 2 {
				continue // must be an assignment, matching UNIT.C's own strchr(line,'=') guard
			}
			val, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return 0, 0, fmt.Errorf("%s line %d: prefix value %q: %w", path, i+1, fields[1], err)
			}
			if _, err := db.ExecContext(ctx, insertPrefixSQL, fields[0], val); err != nil {
				return 0, 0, fmt.Errorf("insert prefix %q: %w", fields[0], err)
			}
			prefixCount++

		case dtypePrimary:
			if strings.Contains(row, "NEWUNIT") {
				for k := 0; k < 7; k++ {
					if k+1 >= len(fields) {
						return 0, 0, fmt.Errorf("%s line %d: %q has fewer than 7 dimension values", path, i+1, row)
					}
					v, err := strconv.Atoi(fields[k+1])
					if err != nil {
						return 0, 0, fmt.Errorf("%s line %d: dimension %q: %w", path, i+1, fields[k+1], err)
					}
					currentDims[k] = v
				}
				continue
			}
			if len(fields) < 2 {
				return 0, 0, fmt.Errorf("%s line %d: %q does not split into name and factor", path, i+1, row)
			}
			factor, err := unitFact(fields[1])
			if err != nil {
				return 0, 0, fmt.Errorf("%s line %d: %w", path, i+1, err)
			}
			if _, err := db.ExecContext(ctx, insertUnitSQL, fields[0], factor,
				currentDims[0], currentDims[1], currentDims[2], currentDims[3], currentDims[4], currentDims[5], currentDims[6]); err != nil {
				return 0, 0, fmt.Errorf("insert unit %q: %w", fields[0], err)
			}
			unitCount++

		case dtypeMixed:
			if len(fields) < 2 {
				return 0, 0, fmt.Errorf("%s line %d: %q does not split into name and factor", path, i+1, row)
			}
			factor, err := unitFact(fields[1])
			if err != nil {
				return 0, 0, fmt.Errorf("%s line %d: %w", path, i+1, err)
			}
			var dims [7]int
			for k := 0; k < 7; k++ {
				if k+2 >= len(fields) {
					break // short line: remaining dimensions default to 0
				}
				v, err := strconv.Atoi(fields[k+2])
				if err != nil {
					return 0, 0, fmt.Errorf("%s line %d: dimension %q: %w", path, i+1, fields[k+2], err)
				}
				dims[k] = v
			}
			if _, err := db.ExecContext(ctx, insertUnitSQL, fields[0], factor,
				dims[0], dims[1], dims[2], dims[3], dims[4], dims[5], dims[6]); err != nil {
				return 0, 0, fmt.Errorf("insert unit %q: %w", fields[0], err)
			}
			unitCount++

		default:
			return 0, 0, fmt.Errorf("%s line %d: %q appears before any BEGIN section", path, i+1, row)
		}
	}
	return prefixCount, unitCount, nil
}
