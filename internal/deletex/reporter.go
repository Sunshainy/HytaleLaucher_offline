package deletex

import (
	"hytale-launcher/internal/throttle"
)

// Reporter provides progress reporting for file deletion operations.
type Reporter struct {
	eventName string
	metadata  map[string]any
	scale     float64
	offset    float64
	gate      *throttle.ProgressGate
	callback  ProgressCallback
	total     int
}

// ProgressCallback is a function that receives progress updates.
// It receives the event name, metadata, and progress value (0.0 to 1.0).
type ProgressCallback func(eventName string, metadata map[string]any, progress float64)

// NewReporter creates a new deletion progress reporter.
func NewReporter(eventName string, metadata map[string]any, scale, offset float64, callback ProgressCallback) *Reporter {
	return &Reporter{
		eventName: eventName,
		metadata:  metadata,
		scale:     scale,
		offset:    offset,
		gate:      &throttle.ProgressGate{},
		callback:  callback,
	}
}

// SetTotal sets the total number of files to delete.
func (r *Reporter) SetTotal(total int) {
	r.total = total
}

// Report reports progress for a file deletion.
// current is the number of files deleted so far.
func (r *Reporter) Report(current, total int) {
	if total < 1 {
		return
	}

	// Clamp current to total
	if current > total {
		current = total
	}

	// Calculate progress as a fraction
	progress := float64(current) / float64(total)

	// Apply scaling if configured
	if r.scale > 0 {
		progress = progress * r.scale
	}

	// Apply throttling - only report if gate allows
	if r.gate.Update(progress) {
		r.callback(r.eventName, r.metadata, progress+r.offset)
	}
}

// ReportDone reports 100% completion.
func (r *Reporter) ReportDone() {
	r.callback(r.eventName, r.metadata, r.scale+r.offset)
}
