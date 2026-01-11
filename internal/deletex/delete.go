// Package deletex provides utilities for deleting directories with progress reporting.
package deletex

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// ProgressReporter is a callback function that reports deletion progress.
type ProgressReporter func()

// Dir deletes all files in a directory with progress reporting,
// then removes the directory tree.
func Dir(dir string, reporter ProgressReporter) error {
	slog.Info("scanning directory for files to delete", "dir", dir)

	// Collect all files to delete
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		slog.Error("failed to walk directory", "dir", dir, "error", err)
		return err
	}

	slog.Info("found files to delete", "total", len(files))

	// Delete each file with progress reporting
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			slog.Error("failed to delete file", "path", file, "error", err)
			return err
		}
		if reporter != nil {
			reporter()
		}
	}

	slog.Info("deleted all files, removing directory tree", "dir", dir)

	// Remove the directory tree
	if err := os.RemoveAll(dir); err != nil {
		slog.Error("failed to remove directory", "dir", dir, "error", err)
		return err
	}

	slog.Info("directory deletion complete", "dir", dir)
	return nil
}
