//go:build linux

package keyring

import (
	"encoding/base64"
	"errors"
	"os"

	gokeyring "github.com/zalando/go-keyring"
)

// linuxKeyStore implements keyStore for Linux using the system keyring.
type linuxKeyStore struct {
	enabled bool
}

// newKeyStore creates a new Linux keyring implementation.
func newKeyStore() keyStore {
	// Check if keyring is enabled via environment variable
	_, enabled := os.LookupEnv("HYTALE_LAUNCHER_ENABLE_KEYRING")
	return &linuxKeyStore{enabled: enabled}
}

// get retrieves a value from the Linux keyring.
func (k *linuxKeyStore) get(service, key string) ([]byte, error) {
	if !k.enabled {
		return nil, nil
	}

	secret, err := gokeyring.Get(service, key)
	if err != nil {
		if errors.Is(err, gokeyring.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

// set stores a value in the Linux keyring.
func (k *linuxKeyStore) set(service, key string, value []byte) error {
	if !k.enabled {
		return nil
	}

	// Encode to base64 for storage
	encoded := base64.StdEncoding.EncodeToString(value)
	return gokeyring.Set(service, key, encoded)
}
