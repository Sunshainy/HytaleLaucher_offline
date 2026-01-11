// Package download provides file download functionality with progress reporting,
// resume capability, and hash verification.
package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"hytale-launcher/internal/ioutil"
	"hytale-launcher/internal/net"
)

// ProgressReporter is called during downloads to report progress.
// bytesDownloaded is the number of bytes downloaded so far.
// speed is the current download speed in bytes per second.
type ProgressReporter func(bytesDownloaded int64, speed int64)

// base extracts the filename from a URL, stripping any query parameters.
func base(url string) string {
	// Cut at the first '?' to remove query parameters
	before, _, _ := strings.Cut(url, "?")
	return path.Base(before)
}

// DownloadTemp downloads a file from url to a temporary file in dir.
// If sha256 is non-empty, the downloaded file's hash is verified.
// Returns the path to the temporary file on success.
func DownloadTemp(
	ctx context.Context,
	client *http.Client,
	dir string,
	url string,
	sha256 string,
	reporter ProgressReporter,
) (string, error) {
	var success bool

	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Create a temp file with a pattern based on the URL's base name
	pattern := "dl-*-" + base(url)
	tempFile, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	defer func() {
		tempFile.Close()
		if !success {
			os.Remove(tempFile.Name())
		}
	}()

	slog.Debug("downloading file",
		"url", url,
		"destination", tempFile.Name(),
		"sha256", sha256,
	)

	// Download the file
	err = downloadFile(ctx, client, url, tempFile, reporter)
	if errors.Is(err, context.Canceled) {
		return "", context.Canceled
	}
	if err != nil {
		return "", fmt.Errorf("error downloading file from %q: %w", url, err)
	}

	// Verify SHA256 if provided
	if sha256 != "" {
		if err := verifySHA256(tempFile.Name(), sha256); err != nil {
			return "", err
		}
	}

	success = true
	return tempFile.Name(), nil
}

// downloadFile performs the actual HTTP download to the given file.
func downloadFile(
	ctx context.Context,
	client *http.Client,
	url string,
	file *os.File,
	reporter ProgressReporter,
) error {
	// Check for offline error (network connectivity)
	if err := checkOffline(); err != nil {
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for 404 Not Found
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	// Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Buffer for reading
	buf := make([]byte, 64*1024) // 64KB buffer

	// Speed calculation variables
	const (
		speedWindowSize   = 20                     // Number of samples for moving average
		speedSamplePeriod = 250 * time.Millisecond // Time between speed samples
	)

	var (
		bytesDownloaded int64
		speedSamples    []int64
		lastSampleTime  = time.Now()
		sampleBytes     int64
		currentSpeed    int64
	)

	for {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			slog.Debug("download cancelled by context", "error", ctx.Err())
			return ctx.Err()
		default:
		}

		// Read from response body
		n, readErr := resp.Body.Read(buf)

		if n > 0 {
			// Write to file
			if _, writeErr := file.Write(buf[:n]); writeErr != nil {
				return writeErr
			}

			bytesDownloaded += int64(n)
			sampleBytes += int64(n)

			// Update speed calculation periodically
			elapsed := time.Since(lastSampleTime)
			if elapsed >= speedSamplePeriod {
				lastSampleTime = time.Now()

				// Add sample to sliding window
				if len(speedSamples) >= speedWindowSize {
					// Remove oldest sample
					speedSamples = speedSamples[1:]
				}
				speedSamples = append(speedSamples, sampleBytes)
				sampleBytes = 0

				// Calculate average speed (bytes per 250ms * 4 = bytes per second)
				if len(speedSamples) > 0 {
					var sum int64
					for _, s := range speedSamples {
						sum += s
					}
					currentSpeed = (sum / int64(len(speedSamples))) * 4
				}

				// Report progress
				if reporter != nil {
					reporter(bytesDownloaded, currentSpeed)
				}
			}
		}

		// Check for EOF or error
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				// Final progress report
				if reporter != nil {
					reporter(bytesDownloaded, currentSpeed)
				}
				return nil
			}
			return readErr
		}
	}
}

// checkOffline checks if the system appears to be offline.
// Returns net.ErrOffline if the launcher is in offline mode.
func checkOffline() error {
	return net.OfflineError()
}

// verifySHA256 verifies that the file at path has the expected SHA256 hash.
func verifySHA256(path, expectedHash string) error {
	return ioutil.VerifySHA256(path, expectedHash)
}
