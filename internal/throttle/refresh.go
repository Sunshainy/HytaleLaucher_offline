package throttle

import (
	"context"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
)

// RefreshFunc is a function that performs a refresh operation.
// It returns an error if the refresh fails.
type RefreshFunc func() error

// Refresher periodically calls a refresh function at a specified interval.
// It handles errors by logging them and reporting to Sentry.
type Refresher struct {
	cancel context.CancelFunc
	ctx    context.Context
	fn     RefreshFunc
}

// NewRefresher creates a new Refresher with the given refresh function.
func NewRefresher(fn RefreshFunc) *Refresher {
	return &Refresher{
		fn: fn,
	}
}

// Start begins the periodic refresh loop in a background goroutine.
// The loop will continue until Stop is called.
// The interval parameter specifies how often to call the refresh function.
func (r *Refresher) Start(interval time.Duration) {
	r.ctx, r.cancel = context.WithCancel(context.Background())
	go r.loop(interval)
}

// Stop halts the refresh loop.
// It cancels the context, which causes the loop goroutine to exit.
func (r *Refresher) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
}

// loop runs the periodic refresh operation.
// It creates a ticker and repeatedly calls the refresh function,
// handling any errors that occur.
func (r *Refresher) loop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		// Get the done channel from the context
		done := r.ctx.Done()

		select {
		case <-done:
			return
		case <-ticker.C:
			if err := r.fn(); err != nil {
				slog.Error("error refreshing application state", "error", err)
				sentry.CaptureException(err)
			}
		}
	}
}
