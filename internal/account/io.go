package account

import (
	"encoding/json"
	"fmt"

	"hytale-launcher/internal/crypto"
)

// ReadFile reads and decrypts an account file from the given path.
// The file is expected to be encrypted with the account encryption key.
// Returns the deserialized Account and any error encountered.
func ReadFile(filePath string) (*Account, error) {
	data, err := crypto.ReadFile(filePath, keyName)
	if err != nil {
		return nil, err
	}

	acct := newAccount(filePath)

	if err := json.Unmarshal(data, acct); err != nil {
		return nil, fmt.Errorf("could not unmarshal account data: %w", err)
	}

	return acct, nil
}

// Write serializes and encrypts the account data to the given path.
// The data is encrypted with the account encryption key.
func (a *Account) Write(filePath string) error {
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("could not marshal account data: %w", err)
	}

	return crypto.WriteFile(filePath, keyName, data)
}

// SaveFile saves the account data to its original file path.
// This is a convenience method that calls Write with the stored file path.
func (a *Account) SaveFile() error {
	if err := a.Write(a.filePath); err != nil {
		return fmt.Errorf("error writing account file: %w", err)
	}

	return nil
}
