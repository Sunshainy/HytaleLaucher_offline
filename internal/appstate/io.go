package appstate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/getsentry/sentry-go"

	"hytale-launcher/internal/build"
	"hytale-launcher/internal/crypto"
	"hytale-launcher/internal/hytale"
)

// encryptionKeyName is the keyring key name used for state file encryption.
const encryptionKeyName = "B7F94324-4365-4EB7-A3FC-7FADAA2EEA2F"

// writeFile marshals the state to JSON and writes it to the encrypted env file.
func (s *State) writeFile() error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("error marshaling launcher state for write: %w", err)
	}
	return crypto.WriteFile(s.envFile(), encryptionKeyName, data)
}

// Save persists the state to disk. It logs the operation and captures
// any errors to Sentry.
func (s *State) Save(cause string) {
	slog.Debug("requesting launcher state save", "channel", s.Channel, "cause", cause)

	if err := s.writeFile(); err != nil {
		slog.Error("error saving launcher state", "channel", s.Channel, "error", err)
		sentry.CaptureException(err)
	}
}

// envFile returns the path to the state file for this state's channel.
func (s *State) envFile() string {
	return crypto.DatFile(filepath.Join(hytale.ChannelDir(s.Channel), "env"))
}

// validatePlatform checks if the saved platform matches the current platform.
// Returns an error if there is a mismatch.
func validatePlatform(s *State) error {
	if s.Platform == nil {
		return nil
	}

	current := build.GetPlatform()
	if *s.Platform == *current {
		return nil
	}

	return fmt.Errorf(
		"platform mismatch: saved platform %s does not match current platform %s",
		s.Platform.String(),
		current.String(),
	)
}

// New creates a new State for the given channel with the current platform.
// The IsNew flag is set to true to indicate a fresh state.
func New(channel string) *State {
	return &State{
		Channel:  channel,
		IsNew:    true,
		Platform: build.GetPlatform(),
	}
}

// ErrNotFound is returned when a state file does not exist.
var ErrNotFound = errors.New("state not found")

// Load attempts to load an existing state from disk for the given channel.
// If the state file doesn't exist, it returns ErrNotFound.
// Returns the state and any error that occurred during loading or validation.
func Load(channel string) (*State, error) {
	s := &State{
		Channel: channel,
	}

	data, err := crypto.ReadFile(s.envFile(), encryptionKeyName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	if err := validatePlatform(s); err != nil {
		return nil, err
	}

	return s, nil
}
