package refdb

import (
	"fmt"
	"os"
	"path/filepath"
)

// configDirName is the subdirectory created under the user's
// standard per-OS configuration directory (for example
// ~/Library/Application Support on macOS, ~/.config on Linux) to
// hold both SQLite databases.
const configDirName = "mwkgo"

// ReferenceDBPath returns the path to the shipped reference
// database, creating its parent directory if necessary. This
// database holds universal lookup data (see SeedReferenceDB) and is
// safe to replace wholesale on an upgrade.
func ReferenceDBPath() (string, error) {
	return dbPath("reference.db")
}

// UserDBPath returns the path to the user's own database, creating
// its parent directory if necessary. This database holds
// machine-specific data that only the user can supply, such as a
// particular lathe's change gears, and an upgrade must never
// overwrite it.
func UserDBPath() (string, error) {
	return dbPath("user.db")
}

func dbPath(filename string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locate user configuration directory: %w", err)
	}

	dir := filepath.Join(configDir, configDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create configuration directory %s: %w", dir, err)
	}

	return filepath.Join(dir, filename), nil
}
