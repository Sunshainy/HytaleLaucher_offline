// Package build provides build-time information about the application.
package build

// These variables are set at build time via ldflags.
var (
	// Release is the release branch/mode (e.g., "release", "dev").
	Release string

	// Version is the build version string (e.g., "2026-01-10-c3879fa").
	Version string
)

// isDevMode returns true if the application is running in development mode.
func isDevMode() bool {
	return Release == "dev"
}

// UserAgent returns the user agent string for HTTP requests.
// For release builds, it returns "hytale-launcher/{version}".
// For non-release builds, it returns "hytale-launcher/{release}/{version}".
func UserAgent() string {
	if Release == "release" {
		return "hytale-launcher/" + Version
	}
	return "hytale-launcher/" + Release + "/" + Version
}
