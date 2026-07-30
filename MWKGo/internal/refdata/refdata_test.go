package refdata

import (
	"context"
	"testing"
)

func TestOpen_SeedsAndQueriesEmbeddedTables(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	db, err := Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	var fitsCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM fits").Scan(&fitsCount); err != nil {
		t.Fatalf("query fits: %v", err)
	}
	if fitsCount == 0 {
		t.Error("fits table is empty, want the embedded reference rows")
	}

	var speedCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM machining_speeds").Scan(&speedCount); err != nil {
		t.Fatalf("query machining_speeds: %v", err)
	}
	if speedCount == 0 {
		t.Error("machining_speeds table is empty, want the embedded reference rows")
	}
}

func TestOpen_SecondCallReusesSeededFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	first, err := Open(context.Background())
	if err != nil {
		t.Fatalf("first Open() error = %v", err)
	}
	first.Close()

	second, err := Open(context.Background())
	if err != nil {
		t.Fatalf("second Open() error = %v, want no error re-seeding an already-present file", err)
	}
	defer second.Close()

	var count int
	if err := second.QueryRow("SELECT COUNT(*) FROM fits").Scan(&count); err != nil {
		t.Fatalf("query fits after second Open: %v", err)
	}
	if count == 0 {
		t.Error("fits table is empty after second Open()")
	}
}
