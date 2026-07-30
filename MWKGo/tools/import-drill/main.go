// import-drill is a throwaway, run-once tool (see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for Tier
// 2") that reads the legacy DRILL.DAT reference file and writes its
// rows into a SQLite "drills" table, producing the golden
// reference.db content the eventual drill calculator embeds and
// seeds on first run. It carries none of the production
// error-handling or test requirements ops-standards.md asks of
// shipped code: it runs by hand, once, and its output is committed.
//
// Usage:
//
//	go run ./tools/import-drill -in DRILL.DAT -out reference.db
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"mwkgo/internal/legacydat"
	"mwkgo/internal/refdb"
)

var migrations = []refdb.Migration{
	{
		Version:     1,
		Description: "create drills table",
		SQL: `CREATE TABLE drills (
			name        TEXT NOT NULL,
			diameter_in REAL NOT NULL
		)`,
	},
}

func main() {
	in := flag.String("in", "DRILL.DAT", "path to the legacy DRILL.DAT file")
	out := flag.String("out", "reference.db", "path to write the SQLite database to")
	flag.Parse()

	if err := run(*in, *out); err != nil {
		log.Fatal(err)
	}
}

func run(inPath, outPath string) error {
	f, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", inPath, err)
	}
	defer f.Close()

	rows, err := legacydat.Rows(f, legacydat.Options{StartMarker: "STARTOFDATA", EndMarker: "ENDOFDATA"})
	if err != nil {
		return fmt.Errorf("scan %s: %w", inPath, err)
	}

	type drill struct {
		name string
		size float64
	}
	drills := make([]drill, 0, len(rows))
	for i, row := range rows {
		fields := legacydat.Fields(row, "\t")
		if len(fields) != 2 {
			return fmt.Errorf("%s line %d: %q does not split into name and diameter on a tab", inPath, i+1, row)
		}
		size, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return fmt.Errorf("%s line %d: diameter %q: %w", inPath, i+1, fields[1], err)
		}
		drills = append(drills, drill{name: fields[0], size: size})
	}

	ctx := context.Background()
	os.Remove(outPath) // a re-run should produce a clean file, not append to a stale one
	db, err := refdb.Open(ctx, outPath, migrations)
	if err != nil {
		return fmt.Errorf("open %s: %w", outPath, err)
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	const insertSQL = `INSERT INTO drills (name, diameter_in) VALUES (?, ?)`
	for _, d := range drills {
		if _, err := tx.ExecContext(ctx, insertSQL, d.name, d.size); err != nil {
			return fmt.Errorf("insert drill %q: %w", d.name, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	fmt.Printf("imported %d drills from %s into %s\n", len(drills), inPath, outPath)
	return nil
}
