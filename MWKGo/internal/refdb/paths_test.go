package refdb

import (
	"path/filepath"
	"strings"
	"testing"
)

// withIsolatedHome points the user's home directory at a temporary
// directory for the duration of the test, so ReferenceDBPath and
// UserDBPath never touch the real machine's configuration directory.
func withIsolatedHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	return home
}

func TestReferenceDBPath_ReturnsPathUnderConfigDirName(t *testing.T) {
	withIsolatedHome(t)

	path, err := ReferenceDBPath()
	if err != nil {
		t.Fatalf("ReferenceDBPath() error = %v", err)
	}
	if filepath.Base(path) != "reference.db" {
		t.Errorf("ReferenceDBPath() = %q, want a path ending in reference.db", path)
	}
	if !strings.Contains(path, configDirName) {
		t.Errorf("ReferenceDBPath() = %q, want it to contain %q", path, configDirName)
	}
}

func TestUserDBPath_ReturnsPathUnderConfigDirName(t *testing.T) {
	withIsolatedHome(t)

	path, err := UserDBPath()
	if err != nil {
		t.Fatalf("UserDBPath() error = %v", err)
	}
	if filepath.Base(path) != "user.db" {
		t.Errorf("UserDBPath() = %q, want a path ending in user.db", path)
	}
}

func TestReferenceDBPath_AndUserDBPath_AreDistinct(t *testing.T) {
	withIsolatedHome(t)

	referencePath, err := ReferenceDBPath()
	if err != nil {
		t.Fatalf("ReferenceDBPath() error = %v", err)
	}
	userPath, err := UserDBPath()
	if err != nil {
		t.Fatalf("UserDBPath() error = %v", err)
	}

	if referencePath == userPath {
		t.Errorf("ReferenceDBPath() and UserDBPath() both returned %q, want distinct files", referencePath)
	}
}
