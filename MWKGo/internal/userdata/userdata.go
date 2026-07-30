// Package userdata opens the single userdata.db every Tier 2
// "machine-specific data" calculator reads from and writes to (a
// lathe's available change-gear-cut thread pitches, a dividing
// head's hole-circle plates, a machine's available spindle speeds,
// and similar data that is specific to one person's equipment, as
// opposed to universal reference data — see
// ai/plans/c-to-go-conversion-plan.md, "Data-file strategy for
// Tier 2").
//
// Unlike internal/refdata's reference.db, userdata.db is never
// seeded with a fabricated default: it starts with its schema
// created and every table empty, and a user populates it through
// internal/csvtable.Import, using the example CSV each contributing
// calculator ships. Migrations lists every machine-specific table's
// schema in one place, since every contributing program's table
// lives in this same shared file and version numbers must be unique
// across all of them.
package userdata

import (
	"context"
	"database/sql"
	"fmt"

	"mwkgo/internal/refdb"
)

// Migrations is the combined, ever-growing list of every
// machine-specific table's schema, one entry contributed per program
// as each is converted.
var Migrations = []refdb.Migration{
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
	{
		Version:     2,
		Description: "create diffthrd_pitches table",
		SQL: `CREATE TABLE diffthrd_pitches (
			pitch_tpi REAL NOT NULL
		)`,
	},
	{
		Version:     3,
		Description: "create diam_available_speeds table",
		SQL: `CREATE TABLE diam_available_speeds (
			rpm INTEGER NOT NULL
		)`,
	},
	{
		Version:     4,
		Description: "create spaceblk_blocks table",
		SQL: `CREATE TABLE spaceblk_blocks (
			size REAL NOT NULL
		)`,
	},
	{
		Version:     5,
		Description: "create gearatio_gears table",
		SQL: `CREATE TABLE gearatio_gears (
			teeth INTEGER NOT NULL
		)`,
	},
	{
		Version:     6,
		Description: "create change_settings, change_leadscrew_pitches, and change_gears tables",
		SQL: `
			CREATE TABLE change_settings (
				id        INTEGER PRIMARY KEY CHECK (id = 1),
				max_pairs INTEGER NOT NULL
			);
			CREATE TABLE change_leadscrew_pitches (
				pitch_tpi REAL NOT NULL
			);
			CREATE TABLE change_gears (
				teeth INTEGER NOT NULL
			);`,
	},
	{
		Version:     7,
		Description: "create ddh_settings, ddh_hole_circles, and ddh_gears tables",
		SQL: `
			CREATE TABLE ddh_settings (
				id                INTEGER PRIMARY KEY CHECK (id = 1),
				worm_gear_ratio   INTEGER NOT NULL,
				rapid_index_holes INTEGER NOT NULL
			);
			CREATE TABLE ddh_hole_circles (
				holes INTEGER NOT NULL
			);
			CREATE TABLE ddh_gears (
				teeth INTEGER NOT NULL
			);`,
	},
	{
		Version:     8,
		Description: "create sun_settings table",
		SQL: `CREATE TABLE sun_settings (
			id                        INTEGER PRIMARY KEY CHECK (id = 1),
			latitude_deg              REAL NOT NULL,
			longitude_deg             REAL NOT NULL,
			standard_time_offset_hr   REAL NOT NULL,
			daylight_time_offset_hr   REAL NOT NULL
		)`,
	},
}

// Open opens (creating if necessary) the user's own database and
// brings its schema up to date with Migrations. Every table this
// returns starts, and stays, empty until a user imports their own
// data with internal/csvtable.Import: no row in userdata.db is ever
// a fabricated default.
func Open(ctx context.Context) (*sql.DB, error) {
	path, err := refdb.UserDBPath()
	if err != nil {
		return nil, fmt.Errorf("locate user database path: %w", err)
	}
	db, err := refdb.Open(ctx, path, Migrations)
	if err != nil {
		return nil, fmt.Errorf("open user database: %w", err)
	}
	return db, nil
}
