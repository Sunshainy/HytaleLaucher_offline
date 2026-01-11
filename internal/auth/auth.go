// Package auth provides authentication management for the Hytale launcher.
// It handles OAuth token storage, restoration, and session management.
package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"golang.org/x/oauth2"

	"hytale-launcher/internal/account"
	"hytale-launcher/internal/crypto"
)

// storageDir is a function that returns the application storage directory.
// This should be set by the application during initialization via SetStorageDir.
var storageDir func() string

// oauthConfig holds the OAuth2 configuration.
// This should be set by the application during initialization via SetOAuthConfig.
var oauthConfig *oauth2.Config

// SetStorageDir sets the function used to get the storage directory.
// This must be called before Init() is used.
func SetStorageDir(fn func() string) {
	storageDir = fn
}

// SetOAuthConfig sets the OAuth2 configuration used for token refresh.
func SetOAuthConfig(cfg *oauth2.Config) {
	oauthConfig = cfg
}

// getAccountFilePath returns the path to the account data file.
// Returns empty string if storageDir is not set.
func getAccountFilePath() string {
	if storageDir == nil {
		return ""
	}
	return crypto.DatFile(filepath.Join(storageDir(), "account"))
}

// Controller manages authentication state and OAuth token lifecycle.
type Controller struct {
	// Account holds the current user account data, including tokens and profiles.
	Account *account.Account

	// client is the HTTP client configured with OAuth token source.
	client *http.Client

	mu sync.RWMutex
}

// Init initializes the auth controller by loading the account from disk.
// If the account file exists and is valid, it restores the OAuth session.
// If the file is corrupted or invalid, it is removed and a fresh state is used.
// Returns nil on success; errors are logged but do not cause Init to fail.
func (c *Controller) Init() error {
	filePath := getAccountFilePath()
	if filePath == "" {
		return nil
	}

	acct, err := account.ReadFile(filePath)

	// If file doesn't exist, that's fine - user hasn't logged in yet
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		// Log and report the error, but continue with fresh state
		sentry.CaptureException(err)
		slog.Error("unable to read account file",
			"error", err,
			"file", filePath,
		)

		// Try to remove the corrupted file
		if removeErr := os.Remove(filePath); removeErr != nil {
			sentry.CaptureException(removeErr)
			slog.Error("failed to remove invalid account file",
				"file", filePath,
				"error", removeErr,
			)
		}

		return nil
	}

	// Account file loaded successfully - restore the OAuth session
	if acct != nil {
		c.restore(acct)
	}

	return nil
}

// restore rebuilds the OAuth client from a previously saved account.
// It creates a token source that monitors for token refreshes and persists changes.
func (c *Controller) restore(acct *account.Account) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get the current profile's token, or fall back to account-level token if available
	var profileToken account.Token
	if acct.CurrentProfile != nil {
		profileToken = acct.CurrentProfile.Token
	} else if len(acct.Profiles) > 0 {
		profileToken = acct.Profiles[0].Token
	}

	// Convert stored token data to oauth2.Token
	token := &oauth2.Token{
		AccessToken:  profileToken.AccessToken,
		RefreshToken: profileToken.RefreshToken,
		Expiry:       profileToken.Expiry,
	}

	// Create an HTTP client with token watching capability
	// The callback will be invoked when tokens are refreshed
	c.client = newWatchClient(
		context.Background(),
		token,
		c.tokenChanged,
	)

	c.Account = acct
}

// tokenChanged is called when the OAuth token is refreshed.
// It updates the current profile with the new token values and persists to disk.
func (c *Controller) tokenChanged(newToken *oauth2.Token) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Account == nil || c.Account.CurrentProfile == nil {
		return
	}

	// Check if token actually changed (compare access and refresh tokens)
	oldToken := c.Account.CurrentProfile.Token
	if oldToken.AccessToken == newToken.AccessToken &&
		oldToken.RefreshToken == newToken.RefreshToken {
		return
	}

	slog.Debug("oauth token(s) changed")

	// Update current profile with new token values
	c.Account.CurrentProfile.Token = account.Token{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		Expiry:    newToken.Expiry,
	}

	c.saveAccountLocked("token_changed")
}

// Client returns the OAuth-configured HTTP client.
// Returns nil if no authenticated session exists.
func (c *Controller) Client() *http.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client
}

// GetAccount returns the current account, or nil if not logged in.
func (c *Controller) GetAccount() *account.Account {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Account
}

// SaveAccount persists the current account state to disk.
// The cause parameter is logged for debugging purposes.
func (c *Controller) SaveAccount(cause string) {
	c.mu.RLock()
	acct := c.Account
	c.mu.RUnlock()

	if acct == nil {
		return
	}

	slog.Debug("requesting account save", "cause", cause)

	if err := acct.SaveFile(); err != nil {
		sentry.CaptureException(err)
		slog.Error("unable to save account file",
			"error", err,
			"file", getAccountFilePath(),
		)
	}
}

// saveAccountLocked saves the account without acquiring the lock.
// Caller must hold c.mu.
func (c *Controller) saveAccountLocked(cause string) {
	if c.Account == nil {
		return
	}

	slog.Debug("requesting account save", "cause", cause)

	if err := c.Account.SaveFile(); err != nil {
		sentry.CaptureException(err)
		slog.Error("unable to save account file",
			"error", err,
			"file", getAccountFilePath(),
		)
	}
}

// SetAccount updates the controller with a new account and persists it.
// This is typically called after successful OAuth login flow.
func (c *Controller) SetAccount(acct *account.Account, client *http.Client) {
	c.mu.Lock()
	c.Account = acct
	c.client = client
	c.mu.Unlock()

	c.SaveAccount("account_set")
}

// Logout clears the current session and removes the account file.
func (c *Controller) Logout() error {
	c.mu.Lock()
	c.Account = nil
	c.client = nil
	c.mu.Unlock()

	filePath := getAccountFilePath()
	if filePath == "" {
		return nil
	}

	if err := os.Remove(filePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	slog.Info("logged out successfully")
	return nil
}

// IsLoggedIn returns true if there is an active authenticated session.
func (c *Controller) IsLoggedIn() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Account != nil && c.client != nil
}

// newWatchClient creates an HTTP client with an OAuth token source that
// monitors for token changes and invokes the callback when tokens are refreshed.
func newWatchClient(ctx context.Context, token *oauth2.Token, onChange func(*oauth2.Token)) *http.Client {
	var tokenSource oauth2.TokenSource

	// If we have an OAuth config, use it for token refresh capability
	if oauthConfig != nil {
		tokenSource = oauthConfig.TokenSource(ctx, token)
	} else {
		// Fall back to static token source (no refresh capability)
		tokenSource = oauth2.StaticTokenSource(token)
	}

	// Wrap with watch capability
	src := &watchTokenSource{
		src:      tokenSource,
		onChange: onChange,
		prev:     token,
	}

	// Create client with 10 second timeout
	client := oauth2.NewClient(ctx, src)
	client.Timeout = 10 * time.Second

	return client
}

// watchTokenSource wraps an oauth2.TokenSource and calls onChange
// when a new token is obtained that differs from the previous one.
type watchTokenSource struct {
	mu       sync.Mutex
	src      oauth2.TokenSource
	onChange func(*oauth2.Token)
	prev     *oauth2.Token
}

// Token implements oauth2.TokenSource.
func (s *watchTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	token, err := s.src.Token()
	if err != nil {
		return nil, err
	}

	// Check if token changed
	if !tokenEqual(s.prev, token) {
		s.prev = token
		if s.onChange != nil {
			s.onChange(token)
		}
	}

	return token, nil
}

// tokenEqual checks if two tokens are equivalent.
func tokenEqual(a, b *oauth2.Token) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.AccessToken == b.AccessToken &&
		a.RefreshToken == b.RefreshToken &&
		a.Expiry.Equal(b.Expiry)
}
