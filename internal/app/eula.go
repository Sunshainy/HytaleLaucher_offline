// Package app provides EULA-related methods for the application.
package app

import (
	"embed"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/getsentry/sentry-go"

	"hytale-launcher/internal/hytale"
	"hytale-launcher/internal/legalfiles"
)

//go:embed eula.txt third-party-licenses.txt
var legalContent embed.FS

const (
	eulaFilename     = "eula.txt"
	licensesFilename = "third-party-licenses.txt"
)

// GetEULA returns the contents of the EULA file.
func (a *App) GetEULA() string {
	data, err := legalContent.ReadFile(eulaFilename)
	if err != nil {
		sentry.CaptureException(err)
		slog.Error("failed to read EULA file", "error", err)
		return ""
	}
	return string(data)
}

// HasAcceptedEULA returns true if the user has accepted the EULA.
func (a *App) HasAcceptedEULA() bool {
	acct := a.Auth.GetAccount()
	if acct == nil {
		return false
	}
	return acct.EULAAcceptedAt != nil
}

// AcceptEULA marks the EULA as accepted by the user.
func (a *App) AcceptEULA() error {
	acct := a.Auth.GetAccount()
	if acct == nil {
		slog.Error("cannot accept EULA: no user logged in")
		return nil
	}

	slog.Info("EULA accepted by user")

	now := time.Now()
	acct.EULAAcceptedAt = &now
	a.Auth.SaveAccount("eula_accepted")

	// Extract legal files to storage directory
	storageDir := hytale.StorageDir()

	eulaData, err := legalContent.ReadFile(eulaFilename)
	if err != nil {
		sentry.CaptureException(err)
		slog.Error("failed to read EULA for extraction", "error", err)
	} else {
		eulaPath := filepath.Join(storageDir, eulaFilename)
		if err := legalfiles.Extract(eulaPath, eulaData); err != nil {
			sentry.CaptureException(err)
			slog.Error("failed to extract EULA file", "error", err)
		}
	}

	licensesData, err := legalContent.ReadFile(licensesFilename)
	if err != nil {
		sentry.CaptureException(err)
		slog.Error("failed to read licenses for extraction", "error", err)
	} else {
		licensesPath := filepath.Join(storageDir, licensesFilename)
		if err := legalfiles.Extract(licensesPath, licensesData); err != nil {
			sentry.CaptureException(err)
			slog.Error("failed to extract licenses file", "error", err)
		}
	}

	a.Emit("eula_accepted")
	return nil
}

// DeclineEULA indicates the user has declined the EULA.
// This logs the user out.
func (a *App) DeclineEULA() error {
	slog.Info("EULA declined by user")
	return a.Logout()
}

// writeLegalFiles writes the legal files (EULA, licenses) to the storage directory.
// This is called during initialization to ensure the files are present.
func writeLegalFiles() {
	storageDir := hytale.StorageDir()

	eulaData, err := legalContent.ReadFile(eulaFilename)
	if err != nil {
		slog.Warn("failed to read embedded EULA", "error", err)
	} else {
		eulaPath := filepath.Join(storageDir, eulaFilename)
		if err := os.WriteFile(eulaPath, eulaData, 0o644); err != nil {
			slog.Warn("failed to write EULA file", "error", err)
		}
	}

	licensesData, err := legalContent.ReadFile(licensesFilename)
	if err != nil {
		slog.Warn("failed to read embedded licenses", "error", err)
	} else {
		licensesPath := filepath.Join(storageDir, licensesFilename)
		if err := os.WriteFile(licensesPath, licensesData, 0o644); err != nil {
			slog.Warn("failed to write licenses file", "error", err)
		}
	}
}
