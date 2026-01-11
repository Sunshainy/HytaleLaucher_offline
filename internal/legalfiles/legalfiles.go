// Package legalfiles handles extraction and management of legal document files
// such as EULA, Terms of Service, and Privacy Policy.
package legalfiles

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/getsentry/sentry-go"
)

// alreadyExists checks if a legal file already exists at the given path with the expected size.
// Returns true if the file exists with matching size, false otherwise.
// Logs warnings or errors for unexpected conditions (permission errors, size mismatches).
func alreadyExists(path string, expectedSize int64) bool {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err != nil {
		sentry.CaptureException(err)
		slog.Error("error checking if legal file exists", "path", path, "error", err)
		return false
	}

	actualSize := info.Size()
	if actualSize != expectedSize {
		slog.Warn("legal file exists but size does not match expected",
			"path", path,
			"expected_size", expectedSize,
			"actual_size", actualSize,
		)
		return false
	}

	return true
}

// Extract writes a legal file to the specified path if it does not already exist
// with the expected content size. If the file already exists with the correct size,
// no action is taken.
func Extract(path string, data []byte) error {
	if alreadyExists(path, int64(len(data))) {
		slog.Warn("legal file already exists with expected size", "path", path, "size", len(data))
		return nil
	}

	slog.Info("creating legal file", "path", path)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("could not create legal file %q: %w", path, err)
	}

	return nil
}
