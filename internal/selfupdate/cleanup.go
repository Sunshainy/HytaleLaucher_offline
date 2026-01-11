package selfupdate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"hytale-launcher/internal/crypto"
	"hytale-launcher/internal/hytale"
)

// cleanupNoteFile returns the path to the cleanup note file.
var cleanupNoteFile = func() string {
	return crypto.DatFile(filepath.Join(hytale.DataDir(), "selfupdate"))
}

// cleanupNote stores information about a pending cleanup operation
// after a self-update has been performed.
type cleanupNote struct {
	// Channel is the channel of the old launcher version that should be cleaned up.
	Channel string `json:"channel"`
	// Version is the version string of the old launcher.
	Version string `json:"version"`
}

// WriteFile writes the cleanup note to the filesystem.
func (n *cleanupNote) WriteFile() error {
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return crypto.WriteFile(cleanupNoteFile(), cleanupNoteKeyName, data)
}

// consumeCleanupNote reads and deletes the cleanup note file.
// It returns the cleanup note if one exists, or nil if not.
func consumeCleanupNote() (*cleanupNote, error) {
	defer func() {
		os.Remove(cleanupNoteFile())
	}()

	data, err := crypto.ReadFile(cleanupNoteFile(), cleanupNoteKeyName)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	note := &cleanupNote{}
	if err := json.Unmarshal(data, note); err != nil {
		return nil, err
	}

	return note, nil
}

// CleanupOldLauncher reads the cleanup note and removes the old launcher directory
// if a cleanup is pending from a previous self-update operation.
func CleanupOldLauncher() error {
	note, err := consumeCleanupNote()
	if err != nil {
		return err
	}
	if note == nil {
		return nil
	}

	dir := hytale.PackageDir("launcher", note.Channel, note.Version)
	slog.Debug("cleaning up old launcher", "dir", dir)

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("error removing old launcher directory: %w", err)
	}

	return nil
}
