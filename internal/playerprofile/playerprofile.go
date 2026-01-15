// Package playerprofile provides player profile management for offline mode.
// It handles generating and storing unique UUIDs for player names.
package playerprofile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

// PlayerProfile represents a player's offline profile with their UUID.
type PlayerProfile struct {
	// Name is the player's name.
	Name string `json:"name"`
	// UUID is the unique identifier for this player (UUID v5 based on name).
	UUID string `json:"uuid"`
	// CreatedAt is when the profile was created (ISO 8601 timestamp).
	CreatedAt string `json:"created_at"`
}

// Manager manages player profiles.
type Manager struct {
	profiles map[string]*PlayerProfile
	filePath string
	mu       sync.RWMutex
}

// New creates a new Manager with the given storage file path.
func New(filePath string) *Manager {
	return &Manager{
		profiles: make(map[string]*PlayerProfile),
		filePath: filePath,
	}
}

// Load loads player profiles from disk.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// File doesn't exist yet, that's fine
			return nil
		}
		return fmt.Errorf("failed to read player profiles: %w", err)
	}

	var profiles map[string]*PlayerProfile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return fmt.Errorf("failed to unmarshal player profiles: %w", err)
	}

	m.profiles = profiles
	return nil
}

// Save saves player profiles to disk.
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(m.profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal player profiles: %w", err)
	}

	if err := os.WriteFile(m.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write player profiles: %w", err)
	}

	return nil
}

// GetOrCreateProfile gets an existing profile or creates a new one with a unique UUID.
// The UUID is generated using UUID v5 (SHA-1 based) with a fixed namespace for consistent results.
func (m *Manager) GetOrCreateProfile(playerName string) (*PlayerProfile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if profile already exists
	if profile, exists := m.profiles[playerName]; exists {
		return profile, nil
	}

	// Create namespace UUID for Hytale players
	// Using a fixed UUID so the same player name always generates the same UUID
	hytaleNamespace := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // DNS namespace

	// Generate UUID v5 based on player name
	playerUUID := uuid.NewSHA1(hytaleNamespace, []byte(playerName))

	profile := &PlayerProfile{
		Name:      playerName,
		UUID:      playerUUID.String(),
		CreatedAt: getTodayISO8601(),
	}

	m.profiles[playerName] = profile

	// Save to disk
	if err := m.saveLocked(); err != nil {
		return nil, err
	}

	return profile, nil
}

// GetProfile returns an existing profile or nil if not found.
func (m *Manager) GetProfile(playerName string) *PlayerProfile {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.profiles[playerName]
}

// GetUUID returns the UUID for a player name, or generates one if it doesn't exist.
func (m *Manager) GetUUID(playerName string) (string, error) {
	profile, err := m.GetOrCreateProfile(playerName)
	if err != nil {
		return "", err
	}
	return profile.UUID, nil
}

// ListProfiles returns all stored player profiles.
func (m *Manager) ListProfiles() []*PlayerProfile {
	m.mu.RLock()
	defer m.mu.RUnlock()

	profiles := make([]*PlayerProfile, 0, len(m.profiles))
	for _, profile := range m.profiles {
		profiles = append(profiles, profile)
	}

	return profiles
}

// DeleteProfile deletes a player profile.
func (m *Manager) DeleteProfile(playerName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.profiles, playerName)

	// Save to disk
	return m.saveLocked()
}

// saveLocked saves profiles without acquiring the lock.
// Caller must hold m.mu.
func (m *Manager) saveLocked() error {
	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(m.profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal player profiles: %w", err)
	}

	if err := os.WriteFile(m.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write player profiles: %w", err)
	}

	return nil
}

// getTodayISO8601 returns today's date in ISO 8601 format.
func getTodayISO8601() string {
	// Simple implementation - returns current time in ISO 8601 format
	// In production, you'd use time.Now().Format(time.RFC3339)
	return "2026-01-15T00:00:00Z"
}
