package hytale

import (
	"path/filepath"
)

// Known channels for game releases.
var knownChannels = map[string]bool{
	"release": true,
	"beta":    true,
	"alpha":   true,
	"dev":     true,
}

// Known game package identifiers.
var knownGamePackages = []string{
	"game",
}

// ChannelDir returns the directory path for a given channel.
func ChannelDir(channel string) string {
	return filepath.Join(StorageDir(), channel)
}

// PackageDir returns the directory path for a specific package version.
// The path follows the pattern: StorageDir/channel/package/pkgID/version
func PackageDir(pkgID, channel, version string) string {
	return filepath.Join(ChannelDir(channel), "package", pkgID, version)
}

// IsKnownChannel returns true if the channel name is a recognized release channel.
func IsKnownChannel(channel string) bool {
	return knownChannels[channel]
}

// KnownGamePackages returns a slice of all known game package identifiers.
func KnownGamePackages() []string {
	result := make([]string, len(knownGamePackages))
	copy(result, knownGamePackages)
	return result
}
