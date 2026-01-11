package pkg

import (
	"log/slog"
)

// InvalidateVersionManifests invalidates all cached version manifests,
// forcing them to be re-fetched on next access.
func InvalidateVersionManifests() {
	slog.Debug("invalidating cached version manifests")

	// Invalidate game manifest cache
	if gameManifest != nil {
		gameManifest.Invalidate()
	}

	// Invalidate Java manifest cache
	if javaManifest != nil {
		javaManifest.Invalidate()
	}

	// Invalidate launcher manifest cache
	if launcherManifest != nil {
		launcherManifest.Invalidate()
	}
}

// GetGameManifest returns the game version manifest getter.
func GetGameManifest() interface{} {
	return gameManifest
}

// GetJavaManifest returns the Java version manifest getter.
func GetJavaManifest() interface{} {
	return javaManifest
}

// GetLauncherManifest returns the launcher version manifest getter.
func GetLauncherManifest() interface{} {
	return launcherManifest
}
