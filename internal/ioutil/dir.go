// Package ioutil provides I/O utility functions for the Hytale launcher.
package ioutil

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

// MkdirAll creates a directory and all parent directories with permissions 0755.
// If the directory already exists, it logs a message and returns nil.
func MkdirAll(path string) error {
	err := os.MkdirAll(path, 0o755)
	if errors.Is(err, os.ErrExist) {
		slog.Debug("directory already exists", "path", path)
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// EmptyDir removes a directory and all its contents, then recreates it empty.
func EmptyDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove directory: %w", err)
	}
	if err := MkdirAll(path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// StorageDir returns the application storage directory.
// This is a convenience wrapper around hytale.StorageDir().
func StorageDir() string {
	// Avoid circular import by using a simple implementation here.
	// The actual path is determined by the hytale package.
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("failed to get user home directory", "error", err)
		return ""
	}

	// Default to ~/.local/share/hytale on Linux
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return fmt.Sprintf("%s/hytale", xdg)
	}

	return fmt.Sprintf("%s/.local/share/hytale", home)
}
