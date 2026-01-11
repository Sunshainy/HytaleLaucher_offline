// Package launch handles launching the Hytale game process.
// It manages process creation, command line arguments, and Java/game execution.
package launch

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// AuthError represents an authentication error that occurred during launch.
type AuthError struct {
	Err error
}

// Error returns the error message for AuthError.
func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("authentication failed: %v", e.Err)
	}
	return "authentication failed"
}

// Unwrap returns the underlying error.
func (e *AuthError) Unwrap() error {
	return e.Err
}

// Request contains all the parameters needed to launch the game.
type Request struct {
	// GamePath is the path to the game executable or JAR file.
	GamePath string

	// JavaPath is the path to the Java executable.
	JavaPath string

	// WorkingDir is the working directory for the game process.
	WorkingDir string

	// Channel is the game channel (e.g., "release", "beta").
	Channel string

	// SessionToken is the authentication session token.
	SessionToken string

	// IdentityToken is the identity token for authentication.
	IdentityToken string

	// ProfileID is the user's profile identifier.
	ProfileID string

	// ExtraArgs are additional command line arguments.
	ExtraArgs []string

	// Env contains additional environment variables.
	Env []string
}

// appendSessionArgs appends session-related arguments to the command line.
func (r *Request) appendSessionArgs(args []string) []string {
	if r.SessionToken != "" {
		args = append(args, "--sessionToken", r.SessionToken)
	}
	if r.IdentityToken != "" {
		args = append(args, "--identityToken", r.IdentityToken)
	}
	if r.ProfileID != "" {
		args = append(args, "--profileId", r.ProfileID)
	}
	return args
}

// appendAuthArgs appends authentication-related arguments to the command line.
func (r *Request) appendAuthArgs(args []string) []string {
	if r.Channel != "" {
		args = append(args, "--channel", r.Channel)
	}
	return args
}

// waitResult holds the result of waiting for a process to complete.
type waitResult struct {
	state *os.ProcessState
	err   error
}

// launchEnv returns the environment variables for the game process.
// It combines the current environment with any additional variables.
func launchEnv(extra []string) []string {
	env := os.Environ()
	env = append(env, extra...)
	return env
}

// startAndWait starts the command and waits for it to complete.
// It returns an error if the process fails to start or exits with a non-zero code.
func startAndWait(ctx context.Context, cmd *exec.Cmd) error {
	slog.Info("starting game process",
		"path", cmd.Path,
		"args", cmd.Args,
		"dir", cmd.Dir,
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start game process: %w", err)
	}

	// Create a channel to receive the wait result
	done := make(chan waitResult, 1)

	go func() {
		state, err := cmd.Process.Wait()
		done <- waitResult{state: state, err: err}
	}()

	// Wait for either context cancellation or process completion
	select {
	case <-ctx.Done():
		// Context was cancelled, kill the process
		slog.Info("context cancelled, killing game process")
		if err := cmd.Process.Kill(); err != nil {
			slog.Error("failed to kill game process", "error", err)
		}
		return ctx.Err()

	case result := <-done:
		if result.err != nil {
			return fmt.Errorf("game process error: %w", result.err)
		}

		if !result.state.Success() {
			exitCode := result.state.ExitCode()
			slog.Warn("game process exited with non-zero code", "exitCode", exitCode)
			return fmt.Errorf("game exited with code %d", exitCode)
		}

		slog.Info("game process completed successfully")
		return nil
	}
}

// Do launches the game with the given request parameters.
// It constructs the command line arguments, sets up the environment,
// and waits for the game process to complete.
func Do(ctx context.Context, req *Request) error {
	if req.GamePath == "" {
		return errors.New("game path is required")
	}

	if req.JavaPath == "" {
		return errors.New("java path is required")
	}

	slog.Info("launching game",
		"gamePath", req.GamePath,
		"javaPath", req.JavaPath,
		"channel", req.Channel,
		"workingDir", req.WorkingDir,
	)

	// Build command line arguments
	args := []string{}

	// Add the game JAR as the first argument after java
	args = append(args, "-jar", req.GamePath)

	// Add session arguments
	args = req.appendSessionArgs(args)

	// Add auth arguments
	args = req.appendAuthArgs(args)

	// Add any extra arguments
	args = append(args, req.ExtraArgs...)

	// Create the command
	cmd := exec.CommandContext(ctx, req.JavaPath, args...)

	// Set working directory
	if req.WorkingDir != "" {
		cmd.Dir = req.WorkingDir
	}

	// Set environment
	cmd.Env = launchEnv(req.Env)

	// Connect stdout and stderr to the current process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start and wait for the game
	if err := startAndWait(ctx, cmd); err != nil {
		// Check if this is an authentication error
		var authErr *AuthError
		if errors.As(err, &authErr) {
			return authErr
		}
		return err
	}

	return nil
}
