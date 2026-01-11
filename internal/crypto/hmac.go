package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HMAC computes an HMAC-SHA256 of the data using the provided key,
// and returns the result as a hexadecimal string.
func HMAC(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
