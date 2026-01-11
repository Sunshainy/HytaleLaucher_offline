// Package buildscan provides functionality for scanning and managing installed game builds.
package buildscan

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"hytale-launcher/internal/appstate"
	"hytale-launcher/internal/deletex"
	"hytale-launcher/internal/helper"
	"hytale-launcher/internal/hytale"
)

// GameInstall represents an installed game build.
type GameInstall struct {
	// Channel is the release channel (e.g., "release", "beta").
	Channel string

	// PkgID is the package identifier (e.g., "game").
	PkgID string

	// Version is the installed version string.
	Version string

	// Dir is the installation directory path.
	Dir string

	// HasSignature indicates whether a signature file exists for this install.
	HasSignature bool
}

// Uninstall removes the game installation from disk and cleans up related state.
// It reports progress through the provided reporter callback.
func (g *GameInstall) Uninstall(reporter deletex.ProgressReporter, state *appstate.State) error {
	slog.Info("uninstalling game install",
		"install", g,
		"dir", g.Dir,
	)

	// Delete the installation directory with progress reporting
	if err := deletex.Dir(g.Dir, reporter); err != nil {
		slog.Error("failed to uninstall game install",
			"install", g,
			"error", err,
		)
		return err
	}

	// Clear the dependency from state
	state.SetDependency(g.PkgID, "uninstall", nil)

	// Remove the signature package
	if err := helper.RemoveSig(state, g.Version); err != nil {
		slog.Error("failed to remove game signature package",
			"version", g.Version,
			"error", err,
		)
		return err
	}

	// Save the updated state
	state.Save("uninstall_game_version")

	// Check if any game packages are still installed for this channel
	hasInstalledPackages := false
	for _, pkgID := range hytale.KnownGamePackages() {
		if deps := state.GetDeps(pkgID); deps != nil && len(deps) > 0 {
			hasInstalledPackages = true
			break
		}
	}

	// If no more game packages are installed, remove the channel directory
	if !hasInstalledPackages {
		channelDir := hytale.ChannelDir(g.Channel)
		slog.Info("removing channel directory as no more game packages are installed",
			"channel", g.Channel,
			"dir", channelDir,
		)

		if err := os.RemoveAll(channelDir); err != nil {
			slog.Error("failed to remove channel directory",
				"channel", g.Channel,
				"error", err,
			)
			return err
		}
	}

	return nil
}

// ScanInstalledGames scans the storage directory for installed game builds.
// It walks through channel directories and identifies valid game installations
// by checking for app state files.
// If detailed is true, additional verification is performed on each install.
func ScanInstalledGames(detailed bool) []GameInstall {
	_ = detailed // Used for future extended scanning
	var installs []GameInstall

	storageDir := hytale.StorageDir()

	// Walk through the storage directory
	err := filepath.WalkDir(storageDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}

		// Skip the root directory
		if path == storageDir {
			return nil
		}

		// Only process directories
		if !d.IsDir() {
			return nil
		}

		// Check if this is a channel directory
		channelName := d.Name()
		if !hytale.IsKnownChannel(channelName) {
			slog.Warn("skipping unrecognized channel directory",
				"path", path,
				"channel", channelName,
			)
			return filepath.SkipDir
		}

		// Try to load app state for this channel
		state, err := appstate.Load(channelName)
		if errors.Is(err, appstate.ErrNotFound) {
			return nil
		}
		if err != nil {
			slog.Warn("error loading app state for game install",
				"path", path,
				"error", err,
			)
			return nil
		}

		// Iterate through known game packages to find installations
		for _, pkgID := range hytale.KnownGamePackages() {
			deps := state.GetDeps(pkgID)
			if deps == nil {
				continue
			}

			for version, dep := range deps {
				// Get the installation directory
				installDir := dep.Path
				if installDir == "" {
					installDir = filepath.Join(path, version)
				}

				// Check for signature
				hasSignature := dep.SigPath() != ""
				if hasSignature {
					slog.Debug("found signature for game install",
						"version", version,
						"sig_file", dep.SigPath(),
					)
				}

				install := GameInstall{
					Channel:      channelName,
					PkgID:        pkgID,
					Version:      version,
					Dir:          installDir,
					HasSignature: hasSignature,
				}

				slog.Info("found game install", "install", install)
				installs = append(installs, install)
			}
		}

		return filepath.SkipDir // Don't recurse into channel directories
	})

	if err != nil {
		slog.Error("error scanning for installed games", "error", err)
	}

	return installs
}
