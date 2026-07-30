package csvtable

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	const createTable = `CREATE TABLE drills (name TEXT NOT NULL, diameter_in REAL NOT NULL)`
	if _, err := db.Exec(createTable); err != nil {
		t.Fatalf("create test table: %v", err)
	}
	return db
}

func TestExport_WritesHeaderAndRows(t *testing.T) {
	db := openTestDB(t)
	if _, err := db.Exec(`INSERT INTO drills (name, diameter_in) VALUES (?, ?), (?, ?)`,
		"1/4 in", 0.25, "M6", 0.2362); err != nil {
		t.Fatalf("seed rows: %v", err)
	}

	var out strings.Builder
	if err := Export(context.Background(), db, "drills", &out); err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	got := out.String()
	if !strings.HasPrefix(got, "name,diameter_in\n") {
		t.Errorf("Export() output %q does not start with the expected header", got)
	}
	if !strings.Contains(got, "1/4 in,0.25\n") {
		t.Errorf("Export() output %q does not contain the first row", got)
	}
}

func TestExport_EmptyTable_WritesHeaderOnly(t *testing.T) {
	db := openTestDB(t)

	var out strings.Builder
	if err := Export(context.Background(), db, "drills", &out); err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	if out.String() != "name,diameter_in\n" {
		t.Errorf("Export() = %q, want just the header row", out.String())
	}
}

func TestExport_UnknownTable_ReturnsError(t *testing.T) {
	db := openTestDB(t)

	var out strings.Builder
	err := Export(context.Background(), db, "does_not_exist", &out)
	if err == nil {
		t.Fatalf("Export() error = nil, want an error for a table that does not exist")
	}
}

func TestImport_ReplacesTableContents(t *testing.T) {
	db := openTestDB(t)
	if _, err := db.Exec(`INSERT INTO drills (name, diameter_in) VALUES (?, ?)`, "stale row", 1.0); err != nil {
		t.Fatalf("seed stale row: %v", err)
	}

	csvData := "name,diameter_in\n1/4 in,0.25\nM6,0.2362\n"
	if err := Import(context.Background(), db, "drills", strings.NewReader(csvData)); err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	rows, err := db.Query(`SELECT name, diameter_in FROM drills ORDER BY name`)
	if err != nil {
		t.Fatalf("query after import: %v", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		var diameter float64
		if err := rows.Scan(&name, &diameter); err != nil {
			t.Fatalf("scan row: %v", err)
		}
		names = append(names, name)
	}
	if len(names) != 2 || names[0] != "1/4 in" || names[1] != "M6" {
		t.Errorf("rows after import = %v, want exactly the two imported names", names)
	}
}

func TestImport_UnknownColumn_LeavesTableUnchanged(t *testing.T) {
	db := openTestDB(t)
	if _, err := db.Exec(`INSERT INTO drills (name, diameter_in) VALUES (?, ?)`, "original row", 1.0); err != nil {
		t.Fatalf("seed original row: %v", err)
	}

	csvData := "name,not_a_real_column\n1/4 in,0.25\n"
	err := Import(context.Background(), db, "drills", strings.NewReader(csvData))
	if err == nil {
		t.Fatalf("Import() error = nil, want an error for an unrecognised column")
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM drills`).Scan(&count); err != nil {
		t.Fatalf("count rows after failed import: %v", err)
	}
	if count != 1 {
		t.Errorf("row count after a failed import = %d, want 1 (the original row, untouched)", count)
	}
}

func TestImport_RaggedRows_ReturnsError(t *testing.T) {
	db := openTestDB(t)

	csvData := "name,diameter_in\n1/4 in,0.25,extra\n"
	err := Import(context.Background(), db, "drills", strings.NewReader(csvData))
	if err == nil {
		t.Fatalf("Import() error = nil, want an error for a row with the wrong number of fields")
	}
}

func TestImport_ClosedDatabase_ReturnsError(t *testing.T) {
	db := openTestDB(t)
	db.Close()

	err := Import(context.Background(), db, "drills", strings.NewReader("name,diameter_in\n"))
	if err == nil {
		t.Fatalf("Import() error = nil, want an error when the database connection is closed")
	}
}
