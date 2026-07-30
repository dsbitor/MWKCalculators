package userdata

import (
	"context"
	"testing"
)

func TestOpen_CreatesEmptySchema(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	db, err := Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	for _, table := range []string{
		"divhead_settings", "divhead_hole_circles",
		"diffthrd_pitches", "diam_available_speeds", "spaceblk_blocks",
	} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Errorf("query %s: %v", table, err)
			continue
		}
		if count != 0 {
			t.Errorf("table %s has %d rows, want 0 (no fabricated default)", table, count)
		}
	}
}

func TestOpen_SecondCallReusesExistingFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	first, err := Open(context.Background())
	if err != nil {
		t.Fatalf("first Open() error = %v", err)
	}
	if _, err := first.Exec(`INSERT INTO diffthrd_pitches (pitch_tpi) VALUES (20)`); err != nil {
		t.Fatalf("seed a row: %v", err)
	}
	first.Close()

	second, err := Open(context.Background())
	if err != nil {
		t.Fatalf("second Open() error = %v", err)
	}
	defer second.Close()

	var count int
	if err := second.QueryRow("SELECT COUNT(*) FROM diffthrd_pitches").Scan(&count); err != nil {
		t.Fatalf("query diffthrd_pitches: %v", err)
	}
	if count != 1 {
		t.Errorf("diffthrd_pitches has %d rows after reopening, want 1 (the row inserted before closing)", count)
	}
}
