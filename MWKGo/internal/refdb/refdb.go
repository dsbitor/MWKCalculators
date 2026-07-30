// Package refdb opens and migrates the two SQLite databases the
// data-backed calculators share: a shipped reference database of
// universal lookup tables, and a user database of machine-specific
// data such as a lathe's change gears, which is never seeded with a
// fabricated default and is populated only through CSV import.
//
// Every connection applies the pragma policy from this project's
// sqlite.md standard (WAL journalling, full synchronous durability,
// foreign key enforcement, and a busy timeout) before any migration
// runs, using the modernc.org/sqlite driver so no CGO toolchain is
// required to build a calculator that uses one of these databases.
package refdb

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// pragmas are applied to every connection immediately after opening,
// in this order, before any migration or query runs.
var pragmas = []string{
	"PRAGMA journal_mode=WAL",
	"PRAGMA synchronous=FULL",
	"PRAGMA foreign_keys=ON",
	"PRAGMA busy_timeout=5000",
}

// Open opens the SQLite database at path, applying the project's
// standard pragmas and then bringing the schema up to date with
// migrations. path is created if it does not already exist, which is
// the correct behaviour for the user database (it starts empty) and
// a no-op for the reference database (SeedReferenceDB should have
// already put a populated file at path before Open is called).
func Open(ctx context.Context, path string, migrations []Migration) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database %s: %w", path, err)
	}

	// A single connection is more than enough for a calculator one
	// person runs at a time, and it is required for correctness
	// against an in-memory database (":memory:" in tests): each
	// pooled connection to ":memory:" would otherwise see its own
	// separate, empty database.
	db.SetMaxOpenConns(1)

	for _, pragma := range pragmas {
		if _, err := db.ExecContext(ctx, pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("apply pragma %q to %s: %w", pragma, path, err)
		}
	}

	if err := applyMigrations(ctx, db, migrations); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate database %s: %w", path, err)
	}

	return db, nil
}
