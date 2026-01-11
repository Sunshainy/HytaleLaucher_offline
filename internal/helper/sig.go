package helper

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"hytale-launcher/internal/appstate"
	"hytale-launcher/internal/hytale"
)

// sigFilePath returns the path to the signature file for a given version.
func sigFilePath(state *appstate.State, version string) string {
	return filepath.Join(hytale.PackageDir("game", state.Channel, ""), fmt.Sprintf("%s.sig", version))
}

// RemoveSig removes the signature file for a given game version.
func RemoveSig(state *appstate.State, version string) error {
	sigPath := sigFilePath(state, version)

	slog.Debug("removing signature file", "path", sigPath)

	if err := os.Remove(sigPath); err != nil {
		if os.IsNotExist(err) {
			slog.Debug("signature file does not exist", "path", sigPath)
			return nil
		}
		return fmt.Errorf("failed to remove signature file: %w", err)
	}

	slog.Debug("signature file removed", "path", sigPath)
	return nil
}

// SaveSig saves a signature to the signature file for a given game version.
func SaveSig(state *appstate.State, version string, signature []byte) error {
	sigPath := sigFilePath(state, version)

	// Ensure the directory exists
	dir := filepath.Dir(sigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create signature directory: %w", err)
	}

	slog.Debug("saving signature file", "path", sigPath)

	if err := os.WriteFile(sigPath, signature, 0644); err != nil {
		return fmt.Errorf("failed to write signature file: %w", err)
	}

	slog.Debug("signature file saved", "path", sigPath)
	return nil
}

// LoadSig loads the signature from the signature file for a given game version.
func LoadSig(state *appstate.State, version string) ([]byte, error) {
	sigPath := sigFilePath(state, version)

	slog.Debug("loading signature file", "path", sigPath)

	data, err := os.ReadFile(sigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read signature file: %w", err)
	}

	return data, nil
}
