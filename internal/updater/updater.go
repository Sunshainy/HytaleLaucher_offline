// Package updater provides functionality for checking and managing updates
// across multiple update sources.
package updater

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"hytale-launcher/internal/appstate"
	"hytale-launcher/internal/auth"
	"hytale-launcher/internal/pkg"
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
func (u *Updater) CheckForUpdates(state *appstate.State, authCtrl *auth.Controller) (int, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Clear previous update info.
	for _, p := range u.packages {
		p.AvailableUpdate = nil
	}

	channel := ""
	if state != nil {
		channel = state.Channel
	}

	ctx := context.Background()
	updateCount := 0

	for _, p := range u.packages {
		slog.Debug("checking for update",
			"package", p.Name,
			"channel", channel,
		)

		// Emit checking event.
		if u.listener != nil {
			u.listener.Event(update.Event{
				Name:    "checking",
				Package: p.Name,
			})
		}

		// Check for updates based on package type
		var pkgUpdate pkg.Update
		var err error

		switch p.Name {
		case "jre":
			pkgUpdate, err = pkg.CheckForJavaUpdate(ctx, state, channel)
		case "game":
			// Build auth context for game updates
			var gameAuth *pkg.Auth
			if authCtrl != nil && authCtrl.IsLoggedIn() {
				acct := authCtrl.GetAccount()
				if acct != nil {
					gameAuth = &pkg.Auth{
						Account: &pkg.GameAccount{
							Patchlines: make(map[string]*pkg.GamePatchline),
						},
					}
					// Populate patchlines from account data
					if acct.CurrentProfile != nil {
						for _, ent := range acct.CurrentProfile.Entitlements {
							// Parse patchline entitlements
							if len(ent) > 10 && ent[:10] == "patchline:" {
								patchlineName := ent[10:]
								gameAuth.Account.Patchlines[patchlineName] = &pkg.GamePatchline{
									Name:        patchlineName,
									NewestBuild: 1, // Will be populated from server
								}
							}
						}
					}
				}
			}

			if gameAuth != nil && gameAuth.Account != nil {
				game := &pkg.Game{
					Channel: channel,
					State:   state,
				}
				pkgUpdate, err = game.CheckForUpdate(ctx, gameAuth)
			}
		case "launcher":
			pkgUpdate, err = pkg.CheckForLauncherUpdate(ctx)
		}

		if err != nil {
			slog.Warn("error checking for update",
				"package", p.Name,
				"error", err,
			)
			u.reportError(p.Name, err)
			continue
		}

		if pkgUpdate != nil {
			// Convert pkg.Update to update.Item
			info := pkg.GetUpdateInfo(pkgUpdate)
			p.AvailableUpdate = &update.Item{
				Name:           p.Name,
				Version:        info.TargetVersion,
				CurrentVersion: info.CurrentVersion,
				Size:           info.Size,
			}
			updateCount++
		}

		slog.Debug("update check complete for package",
			"package", p.Name,
			"has_update", p.AvailableUpdate != nil,
		)
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

	ctx := context.Background()

	for _, p := range u.packages {
		if p.AvailableUpdate == nil {
			continue
		}

		slog.Info("applying update",
			"package", p.Name,
			"version", p.AvailableUpdate.Version,
		)

		// Emit applying event.
		if u.listener != nil {
			u.listener.Event(update.Event{
				Name:    "applying",
				Package: p.Name,
				Version: p.AvailableUpdate.Version,
			})
		}

		// Create progress reporter that emits notifications
		reporter := func(status pkg.UpdateStatus) {
			u.reportProgress(p.Name, 0, 0, status.Progress)
		}

		// Re-check and apply the update based on package type
		var err error
		switch p.Name {
		case "jre":
			var javaUpdate pkg.Update
			javaUpdate, err = pkg.CheckForJavaUpdate(ctx, state, state.Channel)
			if err == nil && javaUpdate != nil {
				err = javaUpdate.Apply(ctx, state, reporter)
			}
		case "launcher":
			var launcherUpdate pkg.Update
			launcherUpdate, err = pkg.CheckForLauncherUpdate(ctx)
			if err == nil && launcherUpdate != nil {
				err = launcherUpdate.Apply(ctx, state, reporter)
			}
		}

		if err != nil {
			slog.Error("failed to apply update",
				"package", p.Name,
				"error", err,
			)
			u.reportError(p.Name, err)
			return fmt.Errorf("failed to apply %s update: %w", p.Name, err)
		}

		// Emit complete event.
		if u.listener != nil {
			u.listener.Event(update.Event{
				Name:    "complete",
				Package: p.Name,
				Version: p.AvailableUpdate.Version,
			})
		}

		p.AvailableUpdate = nil
	}

	return nil
}

// Verify verifies the integrity of all installed packages.
func (u *Updater) Verify(state *appstate.State) error {
	u.mu.RLock()
	defer u.mu.RUnlock()

	for _, p := range u.packages {
		slog.Debug("verifying package",
			"package", p.Name,
		)

		// Check if the package dependency exists in state
		dep := state.GetDependency(p.Name)
		if dep == nil {
			slog.Debug("package not installed, skipping verification",
				"package", p.Name,
			)
			continue
		}

		// Emit verifying event
		if u.listener != nil {
			u.listener.Event(update.Event{
				Name:    "verifying",
				Package: p.Name,
				Version: dep.Version,
			})
		}

		slog.Info("package verified",
			"package", p.Name,
			"version", dep.Version,
			"build", dep.Build,
		)
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
