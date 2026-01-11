package hytale

import (
	"net/http"

	"hytale-launcher/internal/build"
)

// SetUserAgent sets the Hytale launcher HTTP headers on the given request.
// This includes:
//   - User-Agent: hytale-launcher/{version}
//   - X-Hytale-Launcher-Version: {version}
//   - X-Hytale-Launcher-Branch: {release}
func SetUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", build.UserAgent())
	req.Header.Set("X-Hytale-Launcher-Version", build.Version)
	req.Header.Set("X-Hytale-Launcher-Branch", build.Release)
}
