package refdb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSeedReferenceDB_WritesGoldenContentWhenAbsent(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "reference.db")
	golden := []byte("golden database bytes")

	if err := SeedReferenceDB(dest, golden); err != nil {
		t.Fatalf("SeedReferenceDB() error = %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read seeded file: %v", err)
	}
	if string(got) != string(golden) {
		t.Errorf("seeded content = %q, want %q", got, golden)
	}
}

func TestSeedReferenceDB_ExistingFileIsNotOverwritten(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "reference.db")
	original := []byte("a user's existing database, not a golden copy")
	if err := os.WriteFile(dest, original, 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	if err := SeedReferenceDB(dest, []byte("golden database bytes")); err != nil {
		t.Fatalf("SeedReferenceDB() error = %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read file after seed: %v", err)
	}
	if string(got) != string(original) {
		t.Errorf("existing content was overwritten: got %q, want %q", got, original)
	}
}

func TestSeedReferenceDB_MissingParentDirectory_ReturnsError(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "does-not-exist", "reference.db")

	err := SeedReferenceDB(dest, []byte("golden database bytes"))
	if err == nil {
		t.Fatalf("SeedReferenceDB() error = nil, want an error when the parent directory is missing")
	}
}

func TestSeedReferenceDB_LeavesNoTempFileBehindOnSuccess(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "reference.db")

	if err := SeedReferenceDB(dest, []byte("golden database bytes")); err != nil {
		t.Fatalf("SeedReferenceDB() error = %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read directory: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "reference.db" {
		t.Errorf("directory contents = %v, want only reference.db", entries)
	}
}
