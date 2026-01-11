// Package keyring provides secure credential storage using the system keyring.
package keyring

import (
	"crypto/rand"
	"fmt"
)

const (
	// ServiceName is the name used to identify credentials in the system keyring.
	ServiceName = "com.hypixel.hytale-launcher"
)

// keyStore is the interface for platform-specific keyring implementations.
type keyStore interface {
	get(service, key string) ([]byte, error)
	set(service, key string, value []byte) error
}

// store is the platform-specific keyring implementation.
var store keyStore

func init() {
	store = newKeyStore()
}

// Get retrieves a value from the keyring.
func Get(key string) ([]byte, error) {
	return store.get(ServiceName, key)
}

// Set stores a value in the keyring.
func Set(key string, value []byte) error {
	return store.set(ServiceName, key, value)
}

// GetOrGenKey retrieves a key from the keyring, or generates a new one if it doesn't exist.
// The key is 32 bytes (256 bits) suitable for use with AES-256.
func GetOrGenKey(key string) ([]byte, error) {
	// Try to get existing key
	existingKey, err := store.get(ServiceName, key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key '%s': %w", key, err)
	}

	// If key exists, return it
	if existingKey != nil {
		return existingKey, nil
	}

	// Generate a new 32-byte key
	newKey := make([]byte, 32)
	if _, err := rand.Read(newKey); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	// Store the new key
	if err := store.set(ServiceName, key, newKey); err != nil {
		return nil, fmt.Errorf("failed to store key '%s': %w", key, err)
	}

	return newKey, nil
}
