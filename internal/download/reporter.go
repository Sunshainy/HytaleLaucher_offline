package download

import (
	"context"
	"net/http"
	"os"

	"hytale-launcher/internal/hytale"
)

// ProgressReport contains information about download progress.
// This is used to send progress updates to a status callback.
type ProgressReport struct {
	// StatusKey is the identifier for this download operation
	StatusKey string

	// Data contains additional key-value data about the download
	Data map[string]any

	// Progress is the overall progress (0.0 to 1.0)
	Progress float64

	// StatusType indicates the type of status update
	StatusType string

	// BytesDownloaded is the number of bytes downloaded so far
	BytesDownloaded int64

	// TotalBytes is the expected total size (-1 if unknown)
	TotalBytes int64

	// Speed is the current download speed in bytes per second
	Speed int64
}

// Reporter creates a ProgressReporter that reports download progress
// through a callback function. It throttles updates to avoid overwhelming
// the UI with too many updates.
//
// Parameters:
//   - statusKey: identifier for this download operation (e.g., "downloading_patch")
//   - data: additional metadata to include in progress reports
//   - progressScale: multiplier for progress (for sub-operations)
//   - progressOffset: offset to add to progress (for sequential operations)
//   - callback: function to receive progress updates
//
// The reported progress will be: progressOffset + (actualProgress * progressScale)
func Reporter(
	statusKey string,
	data map[string]any,
	progressScale float64,
	progressOffset float64,
	callback func(ProgressReport),
) ProgressReporter {
	var (
		lastProgress float64
		lastSpeed    int64
	)

	return func(bytesDownloaded int64, speed int64) {
		// Calculate progress (0.0 to 1.0) within the scale
		var progress float64
		if bytesDownloaded > 0 && speed > 0 {
			// We don't have total size here, so we can't calculate true progress
			// The original code seems to calculate this differently based on content-length
			progress = 0.0
		}

		// Apply scaling
		if progressScale > 0 {
			progress = progress * progressScale
		}

		// Calculate final progress with offset
		finalProgress := progressOffset + progress

		// Throttle updates - only report if progress changed significantly
		// or if speed changed
		shouldReport := shouldReportProgress(lastProgress, finalProgress)
		if !shouldReport && speed == lastSpeed {
			return
		}

		lastProgress = finalProgress
		lastSpeed = speed

		// Send the progress report
		report := ProgressReport{
			StatusKey:       statusKey,
			Data:            data,
			Progress:        finalProgress,
			StatusType:      "update_status",
			BytesDownloaded: bytesDownloaded,
			TotalBytes:      -1, // Unknown
			Speed:           speed,
		}

		callback(report)
	}
}

// shouldReportProgress determines if a progress update should be sent
// based on the change in progress value.
// Updates are throttled to roughly 1% increments, except near 0% and 100%.
func shouldReportProgress(lastProgress, currentProgress float64) bool {
	// Always report at boundaries
	if currentProgress < 0.01 {
		return true
	}
	if currentProgress >= 0.99 {
		return true
	}

	// Report if progress changed by at least 1%
	return currentProgress-lastProgress >= 0.01
}

// NewReporter creates a reporter adapter for update status reporting.
// This adapts between download progress and the pkg.ProgressReporter callback.
// Parameters:
//   - baseStatus: the base status structure to use for reports
//   - baseProgress: the starting progress value (0.0 to 1.0)
//   - weight: the fraction of total progress this download represents
//   - callback: the pkg.ProgressReporter to call with status updates
//
// Note: Without knowing the total download size, this reporter cannot calculate
// accurate progress percentages. Progress remains at baseProgress until complete.
func NewReporter(baseStatus interface{}, baseProgress, weight float64, callback interface{}) ProgressReporter {
	var (
		lastProgress float64
		lastSpeed    int64
	)

	return func(bytesDownloaded int64, speed int64) {
		// Calculate progress - without total bytes, we stay at baseProgress
		var progress float64
		if weight > 0 {
			progress = 0.0 // No way to calculate without total bytes
		}

		// Apply scaling
		finalProgress := baseProgress + (progress * weight)

		// Throttle updates using the same logic as Reporter
		shouldReport := shouldReportProgress(lastProgress, finalProgress)
		if !shouldReport && speed == lastSpeed {
			return
		}

		lastProgress = finalProgress
		lastSpeed = speed

		if callback == nil {
			return
		}

		// Try to call the callback using type assertion
		switch cb := callback.(type) {
		case func(ProgressReport):
			cb(ProgressReport{
				Progress:        finalProgress,
				BytesDownloaded: bytesDownloaded,
				TotalBytes:      -1, // Unknown
				Speed:           speed,
			})
		}
	}
}

// NewReporterWithSize creates a reporter adapter that knows the expected file size.
// This allows accurate progress calculation based on bytes downloaded.
// Parameters:
//   - statusKey: identifier for this download operation
//   - data: additional metadata to include in progress reports
//   - totalBytes: expected total file size in bytes
//   - progressScale: multiplier for progress (for sub-operations)
//   - progressOffset: offset to add to progress (for sequential operations)
//   - callback: function to receive progress updates
func NewReporterWithSize(
	statusKey string,
	data map[string]any,
	totalBytes int64,
	progressScale float64,
	progressOffset float64,
	callback func(ProgressReport),
) ProgressReporter {
	var (
		lastProgress float64
		lastSpeed    int64
	)

	return func(bytesDownloaded int64, speed int64) {
		// Calculate progress (0.0 to 1.0)
		var progress float64
		if totalBytes > 0 {
			downloaded := bytesDownloaded
			if downloaded > totalBytes {
				downloaded = totalBytes
			}
			progress = float64(downloaded) / float64(totalBytes)
		}

		// Apply scaling
		if progressScale > 0 {
			progress = progress * progressScale
		}

		// Calculate final progress with offset
		finalProgress := progressOffset + progress

		// Throttle updates
		shouldReport := shouldReportProgress(lastProgress, finalProgress)
		if !shouldReport && speed == lastSpeed {
			return
		}

		lastProgress = finalProgress
		lastSpeed = speed

		// Send the progress report
		report := ProgressReport{
			StatusKey:       statusKey,
			Data:            data,
			Progress:        finalProgress,
			StatusType:      "update_status",
			BytesDownloaded: bytesDownloaded,
			TotalBytes:      totalBytes,
			Speed:           speed,
		}

		callback(report)
	}
}

// DownloadTempSimple downloads a file to a temp directory and returns the path.
// This is a simplified version that uses default settings.
func DownloadTempSimple(ctx context.Context, url string, reporter ProgressReporter) (string, error) {
	client := http.DefaultClient
	cacheDir := hytale.InStorageDir("cache")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	return DownloadTemp(ctx, client, cacheDir, url, "", reporter)
}

// ReporterWithTotal creates a ProgressReporter that knows the expected total size.
// This allows for accurate progress percentage calculation.
func ReporterWithTotal(
	statusKey string,
	data map[string]any,
	totalBytes int64,
	progressScale float64,
	progressOffset float64,
	callback func(ProgressReport),
) ProgressReporter {
	var (
		lastProgress float64
		lastSpeed    int64
	)

	return func(bytesDownloaded int64, speed int64) {
		// Calculate progress (0.0 to 1.0)
		var progress float64
		if totalBytes > 0 {
			downloaded := bytesDownloaded
			if downloaded > totalBytes {
				downloaded = totalBytes
			}
			progress = float64(downloaded) / float64(totalBytes)
		}

		// Apply scaling
		if progressScale > 0 {
			progress = progress * progressScale
		}

		// Calculate final progress with offset
		finalProgress := progressOffset + progress

		// Throttle updates
		shouldReport := shouldReportProgress(lastProgress, finalProgress)
		if !shouldReport && speed == lastSpeed {
			return
		}

		lastProgress = finalProgress
		lastSpeed = speed

		// Send the progress report
		report := ProgressReport{
			StatusKey:       statusKey,
			Data:            data,
			Progress:        finalProgress,
			StatusType:      "update_status",
			BytesDownloaded: bytesDownloaded,
			TotalBytes:      totalBytes,
			Speed:           speed,
		}

		callback(report)
	}
}
