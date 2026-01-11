// Package updater provides functionality for checking and managing updates
// across multiple update sources.
package updater

import (
	"fmt"
	"log/slog"
	"sync"

	"hytale-launcher/internal/appstate"
	"hytale-launcher/internal/auth"
	"hytale-launcher/internal/update"
)

// Package represents an updatable component with its update source.
type Package struct {
	// Name is the package identifier.
	Name string

	// Pkg is the update package implementation.
	Pkg update.Package

	// AvailableUpdate holds the pending update info, if any.
	AvailableUpdate *update.Item
}

// Updater manages a collection of updatable packages.
type Updater struct {
	// packages is the list of registered update packages.
	packages []*Package

	// listener receives update events and notifications.
	listener update.Listener

	// mu protects access to packages and their state.
	mu sync.RWMutex
}

// New creates a new Updater instance with the given listener and packages.
func New(listener update.Listener, pkgs ...Package) *Updater {
	u := &Updater{
		packages: make([]*Package, 0, len(pkgs)),
		listener: listener,
	}

	for i := range pkgs {
		p := pkgs[i]
		u.packages = append(u.packages, &Package{
			Name: p.Name,
			Pkg:  p.Pkg,
		})
	}

	return u
}

// CheckForUpdates checks all registered packages for available updates.
// It returns the number of updates found and any error encountered.
func (u *Updater) CheckForUpdates(state *appstate.State, auth *auth.Controller) (int, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Clear previous update info.
	for _, pkg := range u.packages {
		pkg.AvailableUpdate = nil
	}

	channel := ""
	if state != nil {
		channel = state.Channel
	}

	updateCount := 0
	for _, pkg := range u.packages {
		slog.Debug("checking for update",
			"package", pkg.Name,
			"channel", channel,
		)

		// Emit checking event.
		if u.listener != nil {
			u.listener.Event(update.Event{
				Name:    "checking",
				Package: pkg.Name,
			})
		}

		// TODO: Actually check for updates using the package implementation.
		// For now, we'll just log that we checked.
		slog.Debug("update check complete for package",
			"package", pkg.Name,
			"has_update", pkg.AvailableUpdate != nil,
		)

		if pkg.AvailableUpdate != nil {
			updateCount++
		}
	}

	return updateCount, nil
}

// Register adds a package to the updater.
func (u *Updater) Register(pkg *Package) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.packages = append(u.packages, pkg)
}

// GetPackages returns a copy of the registered packages.
func (u *Updater) GetPackages() []*Package {
	u.mu.RLock()
	defer u.mu.RUnlock()

	result := make([]*Package, len(u.packages))
	copy(result, u.packages)
	return result
}

// GetPackage returns a package by name, or nil if not found.
func (u *Updater) GetPackage(name string) *Package {
	u.mu.RLock()
	defer u.mu.RUnlock()

	for _, pkg := range u.packages {
		if pkg.Name == name {
			return pkg
		}
	}
	return nil
}

// HasPendingUpdates returns true if any package has an available update.
func (u *Updater) HasPendingUpdates() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()

	for _, pkg := range u.packages {
		if pkg.AvailableUpdate != nil {
			return true
		}
	}
	return false
}

// HasBlockingUpdates returns true if any pending update is blocking.
func (u *Updater) HasBlockingUpdates() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()

	for _, pkg := range u.packages {
		if pkg.AvailableUpdate != nil && pkg.AvailableUpdate.IsBlocking {
			return true
		}
	}
	return false
}

// ClearPendingUpdates clears all pending update info.
func (u *Updater) ClearPendingUpdates() {
	u.mu.Lock()
	defer u.mu.Unlock()

	for _, pkg := range u.packages {
		pkg.AvailableUpdate = nil
	}
}

// ApplyUpdates applies all pending updates.
// It returns an error if any update fails.
func (u *Updater) ApplyUpdates(state *appstate.State) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	for _, pkg := range u.packages {
		if pkg.AvailableUpdate == nil {
			continue
		}

		slog.Info("applying update",
			"package", pkg.Name,
			"version", pkg.AvailableUpdate.Version,
		)

		// Emit applying event.
		if u.listener != nil {
			u.listener.Event(update.Event{
				Name:    "applying",
				Package: pkg.Name,
				Version: pkg.AvailableUpdate.Version,
			})
		}

		// TODO: Actually apply the update.
		// For now, just mark as complete.

		// Emit complete event.
		if u.listener != nil {
			u.listener.Event(update.Event{
				Name:    "complete",
				Package: pkg.Name,
				Version: pkg.AvailableUpdate.Version,
			})
		}

		pkg.AvailableUpdate = nil
	}

	return nil
}

// Verify verifies the integrity of all installed packages.
func (u *Updater) Verify(state *appstate.State) error {
	u.mu.RLock()
	defer u.mu.RUnlock()

	for _, pkg := range u.packages {
		slog.Debug("verifying package",
			"package", pkg.Name,
		)

		// TODO: Actually verify the package.
	}

	return nil
}

// reportError sends an error event to the listener.
func (u *Updater) reportError(pkg string, err error) {
	if u.listener != nil {
		u.listener.Event(update.Event{
			Name:    "error",
			Package: pkg,
			Error:   err.Error(),
		})
	}
}

// reportProgress sends a progress notification to the listener.
func (u *Updater) reportProgress(pkg string, downloaded, total int64, progress float64) {
	if u.listener != nil {
		u.listener.Notify(update.Notification{
			Package:         pkg,
			BytesDownloaded: downloaded,
			BytesTotal:      total,
			Progress:        progress,
		})
	}
}

// ensure reportError and reportProgress are used (avoid unused warnings)
var _ = (*Updater)(nil).reportError
var _ = (*Updater)(nil).reportProgress
var _ = fmt.Errorf
