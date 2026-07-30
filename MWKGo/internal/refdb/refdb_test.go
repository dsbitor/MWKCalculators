package refdb

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
)

func TestOpen_NoMigrations_CreatesOnlyTheMigrationsTable(t *testing.T) {
	db, err := Open(context.Background(), ":memory:", nil)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	assertTableExists(t, db, "schema_migrations")
}

func TestOpen_AppliesMigrationsInVersionOrder(t *testing.T) {
	migrations := []Migration{
		{Version: 2, Description: "add widgets", SQL: `CREATE TABLE widgets (id INTEGER PRIMARY KEY)`},
		{Version: 1, Description: "add drills", SQL: `CREATE TABLE drills (id INTEGER PRIMARY KEY)`},
	}

	db, err := Open(context.Background(), ":memory:", migrations)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	assertTableExists(t, db, "drills")
	assertTableExists(t, db, "widgets")

	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			t.Fatalf("scan version: %v", err)
		}
		versions = append(versions, v)
	}
	if len(versions) != 2 || versions[0] != 1 || versions[1] != 2 {
		t.Errorf("recorded versions = %v, want [1 2]", versions)
	}
}

func TestOpen_ReopeningSkipsAlreadyAppliedMigrations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	migration := Migration{Version: 1, Description: "add drills", SQL: `CREATE TABLE drills (id INTEGER PRIMARY KEY)`}

	first, err := Open(context.Background(), path, []Migration{migration})
	if err != nil {
		t.Fatalf("first Open() error = %v", err)
	}
	first.Close()

	second, err := Open(context.Background(), path, []Migration{migration})
	if err != nil {
		t.Fatalf("second Open() error = %v, want no error re-applying an already-recorded migration", err)
	}
	defer second.Close()

	assertTableExists(t, second, "drills")
}

func TestOpen_FailedMigrationLeavesEarlierMigrationsApplied(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	migrations := []Migration{
		{Version: 1, Description: "add drills", SQL: `CREATE TABLE drills (id INTEGER PRIMARY KEY)`},
		{Version: 2, Description: "broken migration", SQL: `NOT VALID SQL`},
	}

	_, err := Open(context.Background(), path, migrations)
	if err == nil {
		t.Fatalf("Open() error = nil, want an error from the invalid migration SQL")
	}

	recovered, err := Open(context.Background(), path, migrations[:1])
	if err != nil {
		t.Fatalf("reopening after a failed migration: %v", err)
	}
	defer recovered.Close()

	assertTableExists(t, recovered, "drills")
}

func TestOpen_UnwritableDirectory_ReturnsWrappedError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist", "test.db")

	_, err := Open(context.Background(), path, nil)
	if err == nil {
		t.Fatalf("Open() error = nil, want an error for a database path in a missing directory")
	}
}

func assertTableExists(t *testing.T, db *sql.DB, name string) {
	t.Helper()
	var found string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", name).Scan(&found)
	if err != nil {
		t.Errorf("table %q not found: %v", name, err)
	}
}
