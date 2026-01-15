package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/browser"

	"hytale-launcher/internal/build"
	"hytale-launcher/internal/buildscan"
	"hytale-launcher/internal/deletex"
	"hytale-launcher/internal/extract"
	"hytale-launcher/internal/hytale"
	"hytale-launcher/internal/ioutil"
	"hytale-launcher/internal/net"
	"hytale-launcher/internal/pkg"
	"hytale-launcher/internal/playerprofile"
	"hytale-launcher/internal/repair"
	"hytale-launcher/internal/session"
)

// updatingMu protects the updating flag.
var updatingMu sync.RWMutex
var updating bool

// serverProcess holds the running server process
var serverProcess *os.Process
var serverMu sync.RWMutex

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
	return errors.New("use LaunchGame with player name instead")
}

// launchPackage launches a specific package version.
func (a *App) launchPackage(pkgID, version string) error {
	slog.Info("launching package",
		"package", pkgID,
		"version", version,
	)
	return errors.New("use LaunchGame with player name instead")
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

// LaunchGameRequest contains parameters for launching the game.
type LaunchGameRequest struct {
	PlayerName string `json:"playerName"`
}

// LaunchGame launches the Hytale game with offline mode.
func (a *App) LaunchGame(req LaunchGameRequest) error {
	// Validate player name
	playerName := req.PlayerName
	if playerName == "" {
		return errors.New("player name is required")
	}

	// Save player name for next time
	if err := a.savePlayerName(playerName); err != nil {
		slog.Warn("failed to save player name", "error", err)
	}

	// Define paths
	gameExe := hytale.InStorageDir("package/game/latest/Client/HytaleClient.exe")
	appDir := hytale.InStorageDir("package/game/latest")
	userDir := hytale.InStorageDir("UserData")
	javaExe := hytale.InStorageDir("package/jre/latest/bin/java.exe")

	// Create UserData folder if missing
	if err := ioutil.MkdirAll(userDir); err != nil {
		return err
	}

	// Generate or retrieve unique UUID for this player
	profileManager := playerprofile.New(hytale.InStorageDir("player_profiles.json"))
	if err := profileManager.Load(); err != nil {
		slog.Warn("failed to load player profiles", "error", err)
	}

	// Get or create UUID for this player (deterministic based on player name)
	playerUUID, err := profileManager.GetUUID(playerName)
	if err != nil {
		return fmt.Errorf("failed to generate player UUID: %w", err)
	}

	// Build command arguments
	args := []string{
		"--app-dir", appDir,
		"--user-dir", userDir,
		"--java-exec", javaExe,
		"--auth-mode", "offline",
		"--uuid", playerUUID,
		"--name", playerName,
	}

	slog.Info("launching Hytale",
		"exe", gameExe,
		"playerName", playerName,
		"userDir", userDir,
	)

	// Create the command
	cmd := exec.Command(gameExe, args...)
	cmd.Dir = appDir

	// Hide console window on Windows
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	// Start the game process (don't wait for it)
	if err := cmd.Start(); err != nil {
		return err
	}

	slog.Info("game process started successfully", "pid", cmd.Process.Pid)

	// Emit event to frontend that game has launched
	a.Emit("game:launched")

	// Minimize launcher window after game starts
	a.minimizeLauncher()

	return nil
}

// GetPlayerName returns the saved player name.
func (a *App) GetPlayerName() string {
	name, _ := a.loadPlayerName()
	return name
}

// savePlayerName saves the player name to a file.
func (a *App) savePlayerName(name string) error {
	configPath := hytale.InStorageDir("player.txt")
	return os.WriteFile(configPath, []byte(name), 0644)
}

// loadPlayerName loads the player name from a file.
func (a *App) loadPlayerName() (string, error) {
	configPath := hytale.InStorageDir("player.txt")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsGameInstalled checks if the game is installed.
func (a *App) IsGameInstalled() bool {
	gameExe := hytale.InStorageDir("package/game/latest/Client/HytaleClient.exe")
	_, err := os.Stat(gameExe)
	return err == nil
}

// findGameArchive finds the game ZIP archive next to the launcher executable.
func (a *App) findGameArchive() (string, error) {
	// Get launcher executable path
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Get directory containing the launcher
	launcherDir := filepath.Dir(exePath)

	// Look for .zip file in the launcher directory
	files, err := os.ReadDir(launcherDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".zip" {
			return filepath.Join(launcherDir, file.Name()), nil
		}
	}

	return "", errors.New("game archive not found")
}

// InstallGame installs the game from a ZIP archive.
func (a *App) InstallGame() error {
	// Find the game archive
	archivePath, err := a.findGameArchive()
	if err != nil {
		return fmt.Errorf("failed to find game archive: %w", err)
	}

	slog.Info("installing game from archive", "archive", archivePath)

	// Get destination directory
	destDir := hytale.StorageDir()

	// Extract the archive with progress reporting
	progressFunc := func(current, total int) {
		progress := float64(current) / float64(total) * 100
		a.Emit("install:progress", map[string]interface{}{
			"current":  current,
			"total":    total,
			"progress": progress,
		})
	}

	// Extract without removing existing directory (to preserve log files)
	if err := extract.ArchiveWithoutCleanup(archivePath, destDir, progressFunc, nil); err != nil {
		return fmt.Errorf("failed to extract game archive: %w", err)
	}

	slog.Info("game extraction completed")

	// Delete the archive after successful extraction
	if err := os.Remove(archivePath); err != nil {
		slog.Warn("failed to delete archive after installation", "error", err)
	}

	a.Emit("install:complete")
	return nil
}

// StartServer starts the Hytale server process.
func (a *App) StartServer() error {
	serverMu.Lock()
	defer serverMu.Unlock()

	// Check if server is already running
	if serverProcess != nil {
		return errors.New("server is already running")
	}

	// Server paths
	serverJar := hytale.InStorageDir("package/game/latest/Server/HytaleServer.jar")
	serverDir := hytale.InStorageDir("package/game/latest/Server")
	assetsZip := "../Assets.zip"
	javaExe := hytale.InStorageDir("package/jre/latest/bin/java.exe")
	logFilePath := hytale.InStorageDir("server.log")

	// Check if server exists
	if _, err := os.Stat(serverJar); os.IsNotExist(err) {
		return errors.New("server not found")
	}

	// Check if Java exists
	if _, err := os.Stat(javaExe); os.IsNotExist(err) {
		return errors.New("Java runtime not found")
	}

	// Build command arguments
	args := []string{
		"-jar", serverJar,
		"--assets", assetsZip,
		"--auth-mode", "offline",
	}

	slog.Info("starting Hytale server",
		"jar", serverJar,
		"java", javaExe,
		"workDir", serverDir,
	)

	// Create the command
	cmd := exec.Command(javaExe, args...)
	cmd.Dir = serverDir

	// Hide console window on Windows
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	// Get stdout and stderr pipes for real-time log monitoring
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the server process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	serverProcess = cmd.Process
	slog.Info("server process started", "pid", serverProcess.Pid)

	// Emit "starting" event
	a.Emit("server:starting")

	// Monitor server output in background
	go a.monitorServerOutput(stdout, stderr, logFilePath, cmd)

	return nil
}

// monitorServerOutput monitors server output for boot completion message
func (a *App) monitorServerOutput(stdout, stderr io.ReadCloser, logFilePath string, cmd *exec.Cmd) {
	// Create log file
	logFile, err := os.Create(logFilePath)
	if err != nil {
		slog.Warn("failed to create server log file", "error", err)
	}
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	// Create multi-writer to write to both log file and scan for boot message
	serverBooted := false
	bootCheckDone := make(chan bool, 1)

	// Monitor stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			
			// Write to log file
			if logFile != nil {
				logFile.WriteString(line + "\n")
			}

			// Check if server has booted
			if !serverBooted && strings.Contains(line, "Hytale Server Booted!") {
				serverBooted = true
				slog.Info("server has fully booted")
				a.Emit("server:ready")
				bootCheckDone <- true
			}
		}
	}()

	// Monitor stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			
			// Write to log file
			if logFile != nil {
				logFile.WriteString("[ERROR] " + line + "\n")
			}
		}
	}()

	// Wait for process to complete with timeout for boot check
	go func() {
		select {
		case <-bootCheckDone:
			// Server booted successfully
		case <-time.After(60 * time.Second):
			// Timeout - assume server failed to boot properly
			if !serverBooted {
				slog.Warn("server boot timeout - server may not have started properly")
				a.Emit("server:boot_timeout")
			}
		}
	}()

	// Wait for process to exit
	state, err := cmd.Process.Wait()

	serverMu.Lock()
	serverProcess = nil
	serverMu.Unlock()

	if err != nil {
		slog.Error("server process error", "error", err)
		a.Emit("server:stopped", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		exitCode := 0
		if state != nil {
			exitCode = state.ExitCode()
		}
		slog.Info("server process stopped", "exitCode", exitCode)
		a.Emit("server:stopped", map[string]interface{}{
			"exitCode": exitCode,
		})
	}
}

// StopServer stops the running Hytale server process.
func (a *App) StopServer() error {
	serverMu.Lock()
	defer serverMu.Unlock()

	if serverProcess == nil {
		return errors.New("server is not running")
	}

	slog.Info("stopping server process", "pid", serverProcess.Pid)

	// Kill the server process
	if err := serverProcess.Kill(); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	serverProcess = nil
	slog.Info("server process stopped")

	a.Emit("server:stopped", map[string]interface{}{})
	return nil
}

// IsServerRunning returns true if the server is currently running.
func (a *App) IsServerRunning() bool {
	serverMu.RLock()
	defer serverMu.RUnlock()
	return serverProcess != nil
}
