package ioutil

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// VerifySHA256 computes the SHA256 hash of a file and compares it to the expected hash.
// Returns nil if the hashes match, or an error describing the mismatch.
func VerifySHA256(path string, expectedHash string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file for hashing: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("error hashing file: %w", err)
	}

	actualHash := hex.EncodeToString(h.Sum(nil))
	if actualHash != expectedHash {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	return nil
}

// MakeExecutable adds execute permissions (0111) to a file.
// It preserves the existing file mode and adds the execute bits.
func MakeExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat binary: %w", err)
	}

	newMode := info.Mode() | 0o111
	if err := os.Chmod(path, newMode); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	return nil
}

// FindExecutable walks a directory tree looking for a file whose name ends with
// one of the provided suffixes. Returns the path to the first matching file found.
// If no matching file is found, returns an empty string and nil error.
func FindExecutable(dir string, suffixes []string) (string, error) {
	var result string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		for _, suffix := range suffixes {
			if strings.HasSuffix(path, suffix) {
				result = path
				return filepath.SkipAll
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return result, nil
}
