//go:build windows

package keyring

import (
	"encoding/base64"
	"errors"

	gokeyring "github.com/zalando/go-keyring"
)

// windowsKeyStore implements keyStore for Windows using the Credential Manager.
type windowsKeyStore struct{}

// newKeyStore creates a new Windows keyring implementation.
func newKeyStore() keyStore {
	return &windowsKeyStore{}
}

// get retrieves a value from Windows Credential Manager.
func (k *windowsKeyStore) get(service, key string) ([]byte, error) {
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

// set stores a value in Windows Credential Manager.
func (k *windowsKeyStore) set(service, key string, value []byte) error {
	// Encode to base64 for storage
	encoded := base64.StdEncoding.EncodeToString(value)
	return gokeyring.Set(service, key, encoded)
}
