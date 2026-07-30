// Package refdata embeds the single golden reference.db every
// Tier 2 "universal reference data" calculator reads from (drill
// sizes, shaft/hole fits, cutting speeds, and similar tables that are
// the same for every user, as opposed to one user's own machine
// configuration — see ai/plans/c-to-go-conversion-plan.md, "Data-file
// strategy for Tier 2").
//
// reference.db is built once, offline, by MWKGo/tools/build-refdb
// from the original programs' own .DAT files, and its bytes are
// committed here alongside this package, following the same
// write-temp-then-rename seeding internal/refdb already implements:
// every calculator that needs a reference table embeds the same
// bytes and calls Open, which is a no-op past the first run on a
// given machine.
package refdata

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"mwkgo/internal/refdb"
)

//go:embed reference.db
var golden []byte

// Open seeds the user's reference database from the embedded golden
// copy if it is not already present, then opens it. The returned
// *sql.DB has the project's standard pragmas applied (see
// internal/refdb) but no migrations are run against it: the schema
// and its rows are already baked into the embedded copy by
// MWKGo/tools/build-refdb.
func Open(ctx context.Context) (*sql.DB, error) {
	path, err := refdb.ReferenceDBPath()
	if err != nil {
		return nil, fmt.Errorf("locate reference database path: %w", err)
	}
	if err := refdb.SeedReferenceDB(path, golden); err != nil {
		return nil, fmt.Errorf("seed reference database: %w", err)
	}
	db, err := refdb.Open(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("open reference database: %w", err)
	}
	return db, nil
}
