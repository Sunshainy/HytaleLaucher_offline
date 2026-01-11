// Package session provides types for managing game session data.
package session

import (
	"time"
)

// GameSession represents an authenticated game session with its tokens
// and expiration information.
type GameSession struct {
	// SessionToken is the session token used for game authentication.
	SessionToken string `json:"sessionToken"`

	// IdentityToken is the identity token for the player.
	IdentityToken string `json:"identityToken"`

	// ExpiresAt is the time when the session tokens expire.
	ExpiresAt time.Time `json:"expiresAt"`
}

// IsValid returns true if the session is non-empty and not expired.
func (s *GameSession) IsValid() bool {
	if s == nil {
		return false
	}
	if s.SessionToken == "" || s.IdentityToken == "" {
		return false
	}
	return time.Now().Before(s.ExpiresAt)
}

// IsExpired returns true if the session has expired.
func (s *GameSession) IsExpired() bool {
	if s == nil {
		return true
	}
	return time.Now().After(s.ExpiresAt)
}

// TimeUntilExpiry returns the duration until the session expires.
// Returns a negative duration if already expired.
func (s *GameSession) TimeUntilExpiry() time.Duration {
	if s == nil {
		return 0
	}
	return time.Until(s.ExpiresAt)
}

// NeedsRefresh returns true if the session should be refreshed.
// This returns true if the session will expire within the given margin.
func (s *GameSession) NeedsRefresh(margin time.Duration) bool {
	if s == nil {
		return true
	}
	return time.Now().Add(margin).After(s.ExpiresAt)
}
