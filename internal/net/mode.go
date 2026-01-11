package net

import (
	"errors"
	"sync"
	"time"

	"hytale-launcher/internal/build"
)

// Mode represents the current network mode of the launcher.
type Mode string

const (
	// ModeOnline indicates the launcher is operating in online mode.
	ModeOnline Mode = "online"
	// ModeOffline indicates the launcher is operating in offline mode.
	ModeOffline Mode = "offline"
)

var (
	// modeMu protects access to the current mode.
	modeMu sync.RWMutex
	// currentMode holds the current network mode.
	currentMode Mode = ModeOnline
)

// Current returns the current network mode.
func Current() Mode {
	modeMu.RLock()
	defer modeMu.RUnlock()
	return currentMode
}

// SetMode updates the current network mode.
func SetMode(mode Mode) {
	modeMu.Lock()
	defer modeMu.Unlock()
	currentMode = mode
}

// ErrOffline is returned when an operation cannot be performed because
// the launcher is in offline mode.
var ErrOffline = errors.New("launcher is in offline mode")

// OfflineError returns ErrOffline if the launcher is currently in offline mode,
// otherwise returns nil. In development builds, this can be overridden using
// the HYTALE_LAUNCHER_OFFLINE_MODE environment variable.
func OfflineError() error {
	// In development mode, check for environment variable override
	if build.OfflineMode() {
		// Sleep briefly when using the override
		time.Sleep(2 * time.Second)
		return nil
	}

	if Current() == ModeOffline {
		return ErrOffline
	}

	return nil
}
