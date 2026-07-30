package refdb

import (
	"fmt"
	"os"
	"path/filepath"
)

// SeedReferenceDB writes golden to dest if dest does not already
// exist, using a write-temp-file-then-rename sequence so a process
// interrupted midway through the write never leaves a partially
// written database at dest. It is a no-op, not an error, when dest
// already exists, which makes it safe to call unconditionally every
// time a calculator starts.
func SeedReferenceDB(dest string, golden []byte) error {
	if _, err := os.Stat(dest); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check for existing reference database %s: %w", dest, err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(dest), ".reference-db-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file for %s: %w", dest, err)
	}
	tmpName := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpName) // no-op if the rename below already succeeded
	}()

	if _, err := tmp.Write(golden); err != nil {
		return fmt.Errorf("write temp file %s: %w", tmpName, err)
	}
	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("sync temp file %s: %w", tmpName, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file %s: %w", tmpName, err)
	}
	if err := os.Rename(tmpName, dest); err != nil {
		return fmt.Errorf("rename %s to %s: %w", tmpName, dest, err)
	}
	return nil
}
