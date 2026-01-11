package oauth

import (
	"golang.org/x/oauth2"
)

// callbackData holds data received from an OAuth callback.
// Based on decompiled structure analysis:
// - Offset 0x00: success (bool)
// - Offset 0x10: port (int64)
type callbackData struct {
	Success bool
	Port    int64
}

// stateData contains OAuth state information for CSRF protection.
// Based on decompiled structure analysis:
// - Offset 0x08: state string length
// - Offset 0x10: state string data pointer
// - Offset 0x18: verifier string length
// - Offset 0x20: verifier string data pointer
type stateData struct {
	State    string
	Verifier string
}

// result represents the outcome of an OAuth flow.
// Based on decompiled structure analysis:
// - Offset 0x00: token pointer
// - Offset 0x08: error interface type
// - Offset 0x10: error interface data
type result struct {
	Token *oauth2.Token
	Err   error
}

// Loopback handles OAuth authentication via a local HTTP server.
// Based on decompiled structure analysis:
// - Offset 0x08: ClientID string
// - Offset 0x18: RedirectURL string
// - Offset 0x20: Port int64
// - Offset 0x28: Config interface type (oauth2.Config)
// - Offset 0x30: Config interface data
type Loopback struct {
	ClientID    string
	RedirectURL string
	Port        int64
	Config      *oauth2.Config
}
