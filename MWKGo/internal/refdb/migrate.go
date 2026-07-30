package refdb

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
)

// Migration is one versioned, idempotent schema change. Version
// numbers must be unique; migrations are applied in ascending
// version order regardless of the order given to Open.
type Migration struct {
	Version     int
	Description string
	SQL         string
}

const schemaMigrationsTableDDL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version     INTEGER PRIMARY KEY,
    description TEXT    NOT NULL,
    applied_at  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
)`

// applyMigrations brings db up to date, applying every migration in
// migrations whose version is not already recorded in
// schema_migrations. Each migration runs in its own transaction: a
// migration that fails leaves the schema at the previous version,
// and later migrations are not attempted.
func applyMigrations(ctx context.Context, db *sql.DB, migrations []Migration) error {
	if _, err := db.ExecContext(ctx, schemaMigrationsTableDDL); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return err
	}

	ordered := make([]Migration, len(migrations))
	copy(ordered, migrations)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Version < ordered[j].Version })

	for _, migration := range ordered {
		if applied[migration.Version] {
			continue
		}
		if err := applyOne(ctx, db, migration); err != nil {
			return err
		}
	}
	return nil
}

func appliedVersions(ctx context.Context, db *sql.DB) (map[int]bool, error) {
	rows, err := db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied migration version: %w", err)
		}
		applied[version] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}
	return applied, nil
}

func applyOne(ctx context.Context, db *sql.DB, migration Migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction for migration %d: %w", migration.Version, err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return fmt.Errorf("apply migration %d (%s): %w", migration.Version, migration.Description, err)
	}

	const recordSQL = `INSERT INTO schema_migrations (version, description) VALUES (?, ?)`
	if _, err := tx.ExecContext(ctx, recordSQL, migration.Version, migration.Description); err != nil {
		return fmt.Errorf("record migration %d: %w", migration.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %d: %w", migration.Version, err)
	}
	return nil
}
