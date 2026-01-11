package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"hytale-launcher/internal/build"
)

// Decrypt decrypts data that was encrypted with AES-GCM.
// Encrypted data is expected to start with 'E' byte marker, followed by nonce, then ciphertext.
func Decrypt(data, key []byte) ([]byte, error) {
	if len(data) == 0 || data[0] != 'E' {
		return data, nil
	}

	data = data[1:] // Strip the 'E' marker

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Encrypt encrypts data using AES-GCM.
// The encrypted output is prefixed with 'E' byte marker, followed by nonce, then ciphertext.
// In dev mode, data is returned unencrypted.
func Encrypt(data, key []byte) ([]byte, error) {
	if build.Release == "dev" {
		return data, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Prepend 'E' marker to indicate encrypted data
	result := make([]byte, 1+len(ciphertext))
	result[0] = 'E'
	copy(result[1:], ciphertext)

	return result, nil
}
