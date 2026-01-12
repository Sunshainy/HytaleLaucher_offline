// Package app provides update-related methods for the application.
package app

import (
	"context"
	"log/slog"
	"sync"

	"github.com/getsentry/sentry-go"

	"hytale-launcher/internal/pkg"
	"hytale-launcher/internal/update"
)

// cancelContext holds the current update cancellation context.
var cancelContext context.CancelFunc
var cancelMu sync.Mutex

// PendingUpdates returns information about pending updates.
func (a *App) PendingUpdates() []update.Item {
	if a.Updater == nil {
		return nil
	}

	packages := a.Updater.GetPackages()
	var pending []update.Item

	for _, p := range packages {
		if p.AvailableUpdate != nil {
			pending = append(pending, *p.AvailableUpdate)
		}
	}

	return pending
}

// ApplyUpdates applies all pending updates.
func (a *App) ApplyUpdates() error {
	if a.Updater == nil || a.State == nil {
		return nil
	}

	if a.isUpdating() {
		slog.Warn("update already in progress")
		return nil
	}

	a.markAsUpdating(true)
	defer a.markAsUpdating(false)

	ctx, cancel := context.WithCancel(context.Background())

	cancelMu.Lock()
	cancelContext = cancel
	cancelMu.Unlock()

	defer func() {
		cancelMu.Lock()
		cancelContext = nil
		cancelMu.Unlock()
	}()

	slog.Info("applying updates")

	// Apply updates through the updater
	if err := a.Updater.ApplyUpdates(a.State); err != nil {
		sentry.CaptureException(err)
		slog.Error("failed to apply updates", "error", err)
		a.Emit("update:error", err.Error())
		return err
	}

	// Check if context was cancelled
	select {
	case <-ctx.Done():
		slog.Info("update cancelled")
		a.Emit("update:cancelled")
		return ctx.Err()
	default:
	}

	slog.Info("updates applied successfully")
	a.Emit("update:complete")
	return nil
}

// CancelUpdates cancels any in-progress updates.
func (a *App) CancelUpdates() error {
	slog.Info("cancelling updates")

	cancelMu.Lock()
	if cancelContext != nil {
		cancelContext()
	}
	cancelMu.Unlock()

	a.cancelUpdating()
	a.Emit("update:cancelled")
	return nil
}

// CheckForFreestandingLauncherUpdate checks for launcher updates outside of the normal flow.
func (a *App) CheckForFreestandingLauncherUpdate() (bool, error) {
	slog.Debug("checking for freestanding launcher update")

	ctx := context.Background()
	launcherUpdate, err := pkg.CheckForLauncherUpdate(ctx)
	if err != nil {
		slog.Warn("error checking for launcher update", "error", err)
		return false, err
	}

	if launcherUpdate != nil {
		slog.Info("launcher update available")
		return true, nil
	}

	return false, nil
}

// InvalidateVersionManifests clears cached version manifests.
func (a *App) InvalidateVersionManifests() {
	slog.Debug("invalidating version manifests")
	pkg.InvalidateVersionManifests()
}
