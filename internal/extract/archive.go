// Package extract provides functionality for extracting archive files.
// It supports zip and tar.gz formats with progress reporting and path transformation.
package extract

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// ProgressFunc is called during extraction to report progress.
// It receives the current file index and total file count.
type ProgressFunc func(current, total int)

// NameTransformerFunc transforms file names during extraction.
// It receives the original name and returns the transformed name.
type NameTransformerFunc func(name string) string

// StripRootDir removes the first path component from a file path.
// For example, "root/dir/file.txt" becomes "dir/file.txt".
// If no separator is found, an empty string is returned.
func StripRootDir(path string) string {
	idx := strings.Index(path, "/")
	if idx == -1 {
		return ""
	}
	return path[idx+1:]
}

// safePath validates and constructs a safe file path within the destination directory.
// It prevents path traversal attacks by ensuring the resulting path is within destDir.
func safePath(destDir, name string) (string, error) {
	fullPath := filepath.Join(destDir, name)
	cleanDest := filepath.Clean(destDir) + string(filepath.Separator)
	cleanFull := filepath.Clean(fullPath)

	// Check that the full path starts with the destination directory
	if strings.HasPrefix(cleanFull, cleanDest) || cleanFull == filepath.Clean(destDir) {
		return fullPath, nil
	}

	return "", fmt.Errorf("illegal file path in archive: %s", name)
}

// ArchiveWithoutCleanup extracts an archive file to the destination directory without removing it first.
// It supports .zip, .tar.gz, and .tgz formats.
// The progress function is called with (current, total) file counts during extraction.
// The nameTransformer function can be used to modify file names during extraction (e.g., StripRootDir).
func ArchiveWithoutCleanup(archivePath, destDir string, progress ProgressFunc, nameTransformer NameTransformerFunc) error {
	slog.Debug("extracting archive without cleanup", "archive_path", archivePath, "dest_dir", destDir)

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Determine archive type by extension
	ext := strings.ToLower(filepath.Ext(archivePath))

	switch ext {
	case ".zip":
		return extractZip(archivePath, destDir, progress, nameTransformer)
	case ".gz", ".tgz":
		return extractTarGz(archivePath, destDir, progress, nameTransformer)
	default:
		return fmt.Errorf("unsupported archive format: %s", ext)
	}
}

// Archive extracts an archive file to the destination directory.
// It supports .zip, .tar.gz, and .tgz formats.
// The progress function is called with (current, total) file counts during extraction.
// The nameTransformer function can be used to modify file names during extraction (e.g., StripRootDir).
func Archive(archivePath, destDir string, progress ProgressFunc, nameTransformer NameTransformerFunc) error {
	slog.Debug("extracting archive", "archive_path", archivePath, "dest_dir", destDir)

	// Remove existing destination directory
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to remove existing directory: %w", err)
	}

	// Determine archive type by extension
	ext := strings.ToLower(filepath.Ext(archivePath))

	switch ext {
	case ".zip":
		return extractZip(archivePath, destDir, progress, nameTransformer)
	case ".gz", ".tgz":
		return extractTarGz(archivePath, destDir, progress, nameTransformer)
	default:
		return fmt.Errorf("unsupported archive format: %s", ext)
	}
}

// extractZipFile extracts a single file from a zip archive.
func extractZipFile(f *zip.File, destDir string, nameTransformer NameTransformerFunc) error {
	name := f.Name

	// Apply name transformer if provided
	if nameTransformer != nil {
		name = nameTransformer(name)
	}

	// Skip if name is empty after transformation
	if name == "" {
		return nil
	}

	// Validate the path
	destPath, err := safePath(destDir, name)
	if err != nil {
		return err
	}

	// Handle directories
	if f.Mode().IsDir() {
		return os.MkdirAll(destPath, f.Mode().Perm())
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}

	// Open the file in the archive
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create the destination file
	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode().Perm())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy the contents
	_, err = io.Copy(outFile, rc)
	return err
}

// extractZip extracts a zip archive to the destination directory.
func extractZip(archivePath, destDir string, progress ProgressFunc, nameTransformer NameTransformerFunc) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create the destination directory
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	total := len(reader.File)

	for i, f := range reader.File {
		if err := extractZipFile(f, destDir, nameTransformer); err != nil {
			return err
		}

		if progress != nil {
			progress(i+1, total)
		}
	}

	return nil
}

// countTarGzFiles counts the number of regular files in a tar.gz archive.
func countTarGzFiles(r io.Reader) (int, error) {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return 0, err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	count := 0

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return count, nil
		}
		if err != nil {
			return 0, err
		}

		// Count only regular files
		if header.Typeflag == tar.TypeReg {
			count++
		}
	}
}

// extractTarGz extracts a tar.gz archive to the destination directory.
func extractTarGz(archivePath, destDir string, progress ProgressFunc, nameTransformer NameTransformerFunc) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create the destination directory
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	// Count files if progress reporting is needed
	var total int
	if progress != nil {
		total, err = countTarGzFiles(file)
		if err != nil {
			return err
		}
		// Seek back to the beginning
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return err
		}
	}

	// Create gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	fileIndex := 0

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		name := header.Name

		// Apply name transformer if provided
		if nameTransformer != nil {
			name = nameTransformer(name)
		}

		// Skip if name is empty after transformation
		if name == "" {
			continue
		}

		// Validate the path
		destPath, err := safePath(destDir, name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, 0o755); err != nil {
				return err
			}

		case tar.TypeReg:
			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
				return err
			}

			// Create the destination file
			outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode().Perm())
			if err != nil {
				return err
			}

			// Copy the contents
			_, copyErr := io.Copy(outFile, tarReader)
			closeErr := outFile.Close()

			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return fmt.Errorf("error closing extracted file %s: %w", destPath, closeErr)
			}

			fileIndex++
			if progress != nil {
				progress(fileIndex, total)
			}
		}
	}
}
