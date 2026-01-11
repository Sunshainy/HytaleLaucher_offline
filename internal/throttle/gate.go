// Package throttle provides utilities for throttling and rate limiting operations.
package throttle

// ProgressGate is a gate that throttles progress updates.
// It stores the last reported progress and only allows updates through when:
// - Progress is less than 1% (0.01) - always report start
// - Progress is 99% or more (0.99) - always report near-completion
// - Progress has changed by at least 1% since last update
//
// This prevents excessive UI updates for rapidly changing progress values
// while ensuring the start and end of progress are always reported.
type ProgressGate struct {
	last float64
}

// Update checks if a progress update should be emitted.
// Returns true if the progress should be reported, false if it should be suppressed.
// The progress parameter should be a value between 0.0 and 1.0.
func (g *ProgressGate) Update(progress float64) bool {
	// Always report progress at the boundaries (start and end)
	if progress < 0.01 || progress >= 0.99 {
		g.last = progress
		return true
	}

	// Suppress update if change is less than 1%
	if progress-g.last < 0.01 {
		return false
	}

	// Update passed the gate, store and emit
	g.last = progress
	return true
}
