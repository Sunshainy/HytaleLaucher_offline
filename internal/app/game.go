// Package app provides game-related methods for the application.
package app

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/browser"

	"hytale-launcher/internal/build"
	"hytale-launcher/internal/buildscan"
	"hytale-launcher/internal/deletex"
	"hytale-launcher/internal/hytale"
	"hytale-launcher/internal/ioutil"
	"hytale-launcher/internal/launch"
	"hytale-launcher/internal/net"
	"hytale-launcher/internal/pkg"
	"hytale-launcher/internal/repair"
	"hytale-launcher/internal/session"
)

// updatingMu protects the updating flag.
var updatingMu sync.RWMutex
var updating bool

// isUpdating returns true if an update is currently in progress.
func (a *App) isUpdating() bool {
	updatingMu.RLock()
	defer updatingMu.RUnlock()
	return updating
}

// markAsUpdating sets the updating flag.
func (a *App) markAsUpdating(value bool) {
	updatingMu.Lock()
	defer updatingMu.Unlock()
	updating = value
}

// cancelUpdating cancels any in-progress update.
func (a *App) cancelUpdating() {
	a.markAsUpdating(false)
}

// IsGameAvailable returns true if the game is installed and ready to launch.
func (a *App) IsGameAvailable() bool {
	if a.State == nil {
		return false
	}

	gameDep := a.State.GetDependency("game")
	jreDep := a.State.GetDependency("jre")

	return gameDep != nil && jreDep != nil
}

// GetGameVersion returns the installed game version for the current channel.
func (a *App) GetGameVersion() string {
	if a.State == nil {
		return ""
	}

	gameDep := a.State.GetDependency("game")
	if gameDep == nil {
		return ""
	}

	return gameDep.Version
}

// GetLauncherVersion returns the launcher version information.
func (a *App) GetLauncherVersion() map[string]interface{} {
	return map[string]interface{}{
		"version":      build.Version,
		"release":      build.Release,
		"build_number": build.BuildNumber,
		"platform":     build.OS(),
		"arch":         build.Arch(),
	}
}

// GetInstalledGames returns a list of all installed game builds.
func (a *App) GetInstalledGames() []buildscan.GameInstall {
	return buildscan.ScanInstalledGames(false)
}

// GetInstalledGameDirSizes returns the sizes of installed game directories.
func (a *App) GetInstalledGameDirSizes() map[string]int64 {
	sizes := make(map[string]int64)

	installs := a.GetInstalledGames()
	for _, install := range installs {
		size, err := ioutil.DirSize(install.Dir)
		if err != nil {
			slog.Warn("failed to calculate directory size",
				"dir", install.Dir,
				"error", err,
			)
			continue
		}
		sizes[install.Channel] = size
	}

	return sizes
}

// LaunchGame launches the game with the current configuration.
func (a *App) LaunchGame() error {
	if net.Current() == net.ModeOffline && !a.HasValidSession() {
		return &launch.AuthError{Err: errors.New("offline mode requires a valid session")}
	}

	if a.State == nil {
		return errors.New("no channel selected")
	}

	gameDep := a.State.GetDependency("game")
	if gameDep == nil {
		return errors.New("game not installed")
	}

	jreDep := a.State.GetDependency("jre")
	if jreDep == nil {
		return errors.New("java not installed")
	}

	// Get the game executable path
	gamePath, err := ioutil.FindExecutable(gameDep.Path, []string{".jar", "-server.jar"})
	if err != nil {
		return err
	}
	if gamePath == "" {
		return errors.New("game executable not found")
	}

	// Get the Java executable path
	javaPath, err := ioutil.FindExecutable(jreDep.Path, []string{"java", "java.exe"})
	if err != nil {
		return err
	}
	if javaPath == "" {
		return errors.New("java executable not found")
	}

	// Get session data
	gameSession := a.getGameSession()

	profile := a.getCurrentProfile()
	profileID := ""
	if profile != nil {
		profileID = profile.UUID
	}

	req := &launch.Request{
		GamePath:      gamePath,
		JavaPath:      javaPath,
		WorkingDir:    gameDep.Path,
		Channel:       a.State.Channel,
		SessionToken:  gameSession.SessionToken,
		IdentityToken: gameSession.IdentityToken,
		ProfileID:     profileID,
	}

	slog.Info("launching game",
		"game_path", gamePath,
		"java_path", javaPath,
		"channel", a.State.Channel,
	)

	ctx := context.Background()
	return launch.Do(ctx, req)
}

// getGameSession returns the current game session or creates a new one.
func (a *App) getGameSession() *session.GameSession {
	// In a real implementation, this would fetch the session from the API
	return &session.GameSession{}
}

// HasValidSession returns true if there is a valid game session.
func (a *App) HasValidSession() bool {
	gameSession := a.getGameSession()
	return gameSession != nil && gameSession.IsValid()
}

// GetLastKnownGoodVersion returns the last known good version of the game.
func (a *App) GetLastKnownGoodVersion() string {
	if a.State == nil {
		return ""
	}

	gameDep := a.State.GetDependency("game")
	if gameDep == nil {
		return ""
	}

	return gameDep.Version
}

// LaunchLastKnownGood launches the last known good version of the game.
func (a *App) LaunchLastKnownGood() error {
	slog.Info("launching last known good version")
	return a.LaunchGame()
}

// launchPackage launches a specific package version.
func (a *App) launchPackage(pkgID, version string) error {
	slog.Info("launching package",
		"package", pkgID,
		"version", version,
	)
	return a.LaunchGame()
}

// UninstallGame uninstalls the game from the specified channel.
func (a *App) UninstallGame(channel string) error {
	slog.Info("uninstalling game", "channel", channel)

	installs := buildscan.ScanInstalledGames(false)

	for _, install := range installs {
		if install.Channel != channel {
			continue
		}

		var current int64
		reporter := func() {
			current++
			a.Emit("uninstall:progress", map[string]interface{}{
				"current": current,
			})
		}

		if err := install.Uninstall(reporter, a.State); err != nil {
			sentry.CaptureException(err)
			return err
		}
	}

	a.Emit("uninstall:complete")
	return nil
}

// ValidateGameFiles validates the integrity of game files.
func (a *App) ValidateGameFiles() error {
	if a.State == nil {
		return errors.New("no channel selected")
	}

	gameDep := a.State.GetDependency("game")
	if gameDep == nil {
		return errors.New("game not installed")
	}

	slog.Info("validating game files",
		"channel", a.State.Channel,
		"version", gameDep.Version,
	)

	reporter := func(current, total int, path string) {
		progress := float64(current) / float64(total)
		a.Emit("validate:progress", map[string]interface{}{
			"current":  current,
			"total":    total,
			"progress": progress,
			"path":     path,
		})
	}

	// Load checksums from signature file if available
	checksums := make(map[string]string)
	// In a real implementation, load from gameDep.SigPath()

	result, err := repair.Verify(gameDep.Path, checksums, reporter)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	if !result.IsHealthy() {
		a.Emit("validate:failed", map[string]interface{}{
			"missing":   len(result.MissingFiles),
			"corrupted": len(result.CorruptedFiles),
		})
	} else {
		a.Emit("validate:success")
	}

	return nil
}

// ResetGameSettings resets game settings to defaults.
func (a *App) ResetGameSettings() error {
	slog.Info("resetting game settings")

	if a.State == nil {
		return errors.New("no channel selected")
	}

	// Clear cached data
	pkg.InvalidateVersionManifests()

	a.Emit("settings:reset")
	return nil
}

// OpenHytaleDir opens the Hytale storage directory in the file explorer.
func (a *App) OpenHytaleDir() error {
	storageDir := hytale.StorageDir()
	slog.Info("opening Hytale directory", "dir", storageDir)
	return browser.OpenFile(storageDir)
}

// CanDeleteUserData returns true if user data can be deleted.
func (a *App) CanDeleteUserData() bool {
	// Check if there are no running game processes
	return !a.isUpdating()
}

// DeleteUserData deletes all user data from the storage directory.
func (a *App) DeleteUserData() error {
	if !a.CanDeleteUserData() {
		return errors.New("cannot delete user data while updating")
	}

	slog.Warn("deleting all user data")

	// Logout first
	if err := a.Logout(); err != nil {
		slog.Warn("error during logout before delete", "error", err)
	}

	storageDir := hytale.StorageDir()

	var current int64
	reporter := func() {
		current++
		a.Emit("delete:progress", map[string]interface{}{
			"current": current,
		})
	}

	if err := deletex.Dir(storageDir, reporter); err != nil {
		sentry.CaptureException(err)
		return err
	}

	a.Emit("delete:complete")
	return nil
}

// GetLaunchAuthMode returns the authentication mode for launching.
func (a *App) GetLaunchAuthMode() string {
	if net.Current() == net.ModeOffline {
		return "offline"
	}
	return "online"
}
