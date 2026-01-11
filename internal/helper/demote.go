// Package helper provides utility functions for managing game installations.
package helper

import (
	"fmt"
	"log/slog"
	"os"

	"hytale-launcher/internal/appstate"
	"hytale-launcher/internal/hytale"
)

// buildDir constructs the directory path for a specific game build.
func buildDir(state *appstate.State, buildID int) string {
	return hytale.PackageDir("game", state.Channel, fmt.Sprintf("build-%d", buildID))
}

// DemoteLatestGame moves the current game version to a numbered build directory,
// effectively demoting it to a "last known good" (lkg) state.
func DemoteLatestGame(state *appstate.State) error {
	// Get the current game dependency
	currentDep := state.GetDependency("game")

	if currentDep == nil {
		slog.Info("no current game version found, nothing to demote")
		state.SetDependency("lkg", "", nil)
		return nil
	}

	// Get the build directory path for the current version
	destDir := buildDir(state, currentDep.BuildID)
	srcDir := currentDep.Path

	slog.Info("demoting current game version to numbered build directory",
		"build_id", currentDep.BuildID,
		"from", srcDir,
		"to", destDir,
	)

	slog.Info("removing destination directory if it exists", "dir", destDir)
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to remove directory: %w", err)
	}

	slog.Info("renaming source directory to destination directory",
		"from", srcDir,
		"to", destDir,
	)
	if err := os.Rename(srcDir, destDir); err != nil {
		return fmt.Errorf("failed to rename directory: %w", err)
	}

	// Update the dependency to point to the new location
	currentDep.Path = destDir

	// Clear the game dependency and set it as LKG
	state.SetDependency("game", "", nil)
	state.SetDependency("lkg", "", currentDep)

	return nil
}

// DeleteInstall removes a game installation and its associated signature.
func DeleteInstall(state *appstate.State, depName string, dep *appstate.Dep) error {
	if dep == nil {
		return nil
	}

	slog.Info("deleting game installation",
		"dependency", depName,
		"path", dep.Path,
		"version", dep.Version,
	)

	// Remove the installation directory
	if err := os.RemoveAll(dep.Path); err != nil {
		slog.Error("failed to remove installation directory",
			"path", dep.Path,
			"error", err,
		)
		return fmt.Errorf("failed to remove installation: %w", err)
	}

	// Remove the signature file
	if err := RemoveSig(state, dep.Version); err != nil {
		slog.Warn("failed to remove signature file",
			"version", dep.Version,
			"error", err,
		)
		// Continue even if signature removal fails
	}

	// Clear the dependency
	state.SetDependency(depName, "", nil)

	slog.Info("installation deleted successfully", "dependency", depName)
	return nil
}
