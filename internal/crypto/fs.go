package crypto

import (
	"fmt"
	"os"

	"hytale-launcher/internal/keyring"
)

// ReadFile reads a file and decrypts it if necessary.
// The keyName is used to retrieve the encryption key from the keyring.
func ReadFile(path string, keyName string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, err := keyring.GetOrGenKey(keyName)
	if err != nil {
		return nil, fmt.Errorf("could not get encryption key %q: %w", keyName, err)
	}

	decrypted, err := Decrypt(data, key)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// WriteFile encrypts data and writes it to a file.
// The keyName is used to retrieve the encryption key from the keyring.
// The file is written with 0644 permissions.
func WriteFile(path string, keyName string, data []byte) error {
	key, err := keyring.GetOrGenKey(keyName)
	if err != nil {
		return fmt.Errorf("could not get encryption key %q: %w", keyName, err)
	}

	encrypted, err := Encrypt(data, key)
	if err != nil {
		return fmt.Errorf("could not encrypt data for %q: %w", path, err)
	}

	if err := os.WriteFile(path, encrypted, 0644); err != nil {
		return err
	}

	return nil
}
