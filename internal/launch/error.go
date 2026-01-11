package launch

import (
	"errors"
	"fmt"
)

// LaunchError represents an error that occurred during game launch.
type LaunchError struct {
	Op  string // Operation that failed
	Err error  // Underlying error
}

// Error returns the error message.
func (e *LaunchError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("launch %s: %v", e.Op, e.Err)
	}
	return fmt.Sprintf("launch %s failed", e.Op)
}

// Unwrap returns the underlying error.
func (e *LaunchError) Unwrap() error {
	return e.Err
}

// ExitError represents a non-zero exit code from the game process.
type ExitError struct {
	ExitCode int
}

// Error returns the error message.
func (e *ExitError) Error() string {
	return fmt.Sprintf("game exited with code %d", e.ExitCode)
}

// IsAuthError checks if the error is an authentication error.
func IsAuthError(err error) bool {
	if err == nil {
		return false
	}
	var authErr *AuthError
	return errors.As(err, &authErr)
}
