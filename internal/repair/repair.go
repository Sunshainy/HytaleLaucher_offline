// Package repair provides functionality for verifying and repairing game installations.
package repair

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"hytale-launcher/internal/ioutil"
)

// FileStatus represents the verification status of a file.
type FileStatus int

const (
	// FileStatusOK indicates the file passed verification.
	FileStatusOK FileStatus = iota

	// FileStatusMissing indicates the file is missing.
	FileStatusMissing

	// FileStatusCorrupted indicates the file exists but has an invalid checksum.
	FileStatusCorrupted
)

// FileResult contains the verification result for a single file.
type FileResult struct {
	// Path is the relative path of the file.
	Path string

	// Status is the verification status.
	Status FileStatus

	// ExpectedHash is the expected SHA256 hash.
	ExpectedHash string

	// ActualHash is the computed hash, if the file exists.
	ActualHash string

	// Error contains any error encountered during verification.
	Error error
}

// Result contains the complete verification result for an installation.
type Result struct {
	// TotalFiles is the number of files checked.
	TotalFiles int

	// OKFiles is the number of files that passed verification.
	OKFiles int

	// MissingFiles contains files that are missing.
	MissingFiles []FileResult

	// CorruptedFiles contains files with checksum mismatches.
	CorruptedFiles []FileResult

	// Errors contains any errors encountered during verification.
	Errors []FileResult
}

// IsHealthy returns true if all files passed verification.
func (r *Result) IsHealthy() bool {
	return len(r.MissingFiles) == 0 && len(r.CorruptedFiles) == 0 && len(r.Errors) == 0
}

// NeedsRepair returns true if any files need to be repaired.
func (r *Result) NeedsRepair() bool {
	return len(r.MissingFiles) > 0 || len(r.CorruptedFiles) > 0
}

// ProgressReporter is called during verification to report progress.
type ProgressReporter func(current, total int, path string)

// Verify checks the integrity of an installation against expected checksums.
// The checksums map contains relative paths to their expected SHA256 hashes.
func Verify(installDir string, checksums map[string]string, reporter ProgressReporter) (*Result, error) {
	if checksums == nil || len(checksums) == 0 {
		return nil, errors.New("no checksums provided for verification")
	}

	result := &Result{
		TotalFiles: len(checksums),
	}

	current := 0
	for relativePath, expectedHash := range checksums {
		current++

		fullPath := filepath.Join(installDir, relativePath)

		if reporter != nil {
			reporter(current, result.TotalFiles, relativePath)
		}

		fileResult := verifyFile(fullPath, relativePath, expectedHash)

		switch fileResult.Status {
		case FileStatusOK:
			result.OKFiles++
		case FileStatusMissing:
			result.MissingFiles = append(result.MissingFiles, fileResult)
		case FileStatusCorrupted:
			result.CorruptedFiles = append(result.CorruptedFiles, fileResult)
		default:
			if fileResult.Error != nil {
				result.Errors = append(result.Errors, fileResult)
			}
		}
	}

	return result, nil
}

// verifyFile checks a single file against its expected hash.
func verifyFile(fullPath, relativePath, expectedHash string) FileResult {
	result := FileResult{
		Path:         relativePath,
		ExpectedHash: expectedHash,
	}

	// Check if file exists
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		result.Status = FileStatusMissing
		return result
	}
	if err != nil {
		result.Status = FileStatusCorrupted
		result.Error = err
		return result
	}

	// Skip directories
	if info.IsDir() {
		result.Status = FileStatusOK
		return result
	}

	// Verify checksum
	if err := ioutil.VerifySHA256(fullPath, expectedHash); err != nil {
		result.Status = FileStatusCorrupted
		result.Error = err
		return result
	}

	result.Status = FileStatusOK
	return result
}

// VerifyDirectory verifies all files in a directory against a checksums map.
// This is useful when verifying an installation directory where checksums are
// provided externally (e.g., from a manifest file).
func VerifyDirectory(dir string, checksums map[string]string, reporter ProgressReporter) (*Result, error) {
	return Verify(dir, checksums, reporter)
}

// RepairFile attempts to repair a single corrupted or missing file by re-downloading it.
// It downloads the file to a temporary location, verifies its checksum, and then
// atomically replaces the target file.
func RepairFile(installDir, relativePath, downloadURL, expectedHash string) error {
	slog.Info("repairing file",
		"path", relativePath,
		"download_url", downloadURL,
	)

	fullPath := filepath.Join(installDir, relativePath)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Download to a temp file in the same directory (for atomic rename)
	tempFile, err := os.CreateTemp(filepath.Dir(fullPath), ".repair-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		tempFile.Close()
		os.Remove(tempPath) // Clean up on failure
	}()

	// Download the file
	ctx, cancel := context.WithTimeout(context.Background(), 300*http.DefaultClient.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Write to temp file
	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	tempFile.Close()

	// Verify checksum if provided
	if expectedHash != "" {
		if err := ioutil.VerifySHA256(tempPath, expectedHash); err != nil {
			return fmt.Errorf("checksum verification failed: %w", err)
		}
	}

	// Atomic replace
	if err := os.Rename(tempPath, fullPath); err != nil {
		return fmt.Errorf("failed to replace file: %w", err)
	}

	slog.Info("file repaired successfully", "path", relativePath)
	return nil
}

// CleanupOrphanedFiles removes files that are not in the expected file list.
// This can be useful for cleaning up after failed updates.
func CleanupOrphanedFiles(installDir string, expectedFiles map[string]bool) ([]string, error) {
	var removed []string

	err := filepath.WalkDir(installDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}

		if d.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(installDir, path)
		if err != nil {
			return nil
		}

		if !expectedFiles[relativePath] {
			slog.Debug("removing orphaned file", "path", relativePath)
			if err := os.Remove(path); err != nil {
				slog.Warn("failed to remove orphaned file",
					"path", relativePath,
					"error", err,
				)
			} else {
				removed = append(removed, relativePath)
			}
		}

		return nil
	})

	return removed, err
}
