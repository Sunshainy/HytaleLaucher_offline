package crypto

import (
	"os"
)

// DatFile returns the appropriate file path with extension.
// It checks if a .json file exists and is not a directory; if so, returns path + ".json".
// Otherwise, returns path + ".dat" (encrypted file format).
func DatFile(path string) string {
	jsonPath := path + ".json"
	info, err := os.Stat(jsonPath)
	if err == nil && !info.IsDir() {
		return jsonPath
	}
	return path + ".dat"
}
