//go:build darwin

package keyring

import (
	"encoding/base64"
	"errors"

	gokeyring "github.com/zalando/go-keyring"
)

// darwinKeyStore implements keyStore for macOS using the system Keychain.
type darwinKeyStore struct{}

// newKeyStore creates a new macOS keyring implementation.
func newKeyStore() keyStore {
	return &darwinKeyStore{}
}

// get retrieves a value from the macOS Keychain.
func (k *darwinKeyStore) get(service, key string) ([]byte, error) {
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

// set stores a value in the macOS Keychain.
func (k *darwinKeyStore) set(service, key string, value []byte) error {
	// Encode to base64 for storage
	encoded := base64.StdEncoding.EncodeToString(value)
	return gokeyring.Set(service, key, encoded)
}
