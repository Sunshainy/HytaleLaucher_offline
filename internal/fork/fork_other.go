//go:build !windows

package fork

import (
	"os"
)

// IsElevated returns true if the current process is running with elevated privileges (root).
// On Unix-like systems, this checks if the effective user ID is 0.
func IsElevated() bool {
	return os.Geteuid() == 0
}

// RunElevated starts a process with elevated privileges.
// On Unix-like systems, this simply starts the process as-is since privilege
// escalation typically requires user interaction (e.g., sudo prompt) which
// should be handled externally.
func RunElevated(path string, args []string) (*os.Process, error) {
	// On Linux/macOS, we just start the process normally
	// The launcher should already be running with appropriate permissions
	// or use pkexec/sudo externally
	allArgs := append([]string{path}, args...)

	return startProcess(StartOptions{
		Path: path,
		Args: allArgs,
	})
}

// RunAsUser starts a process as the current user.
// On Unix-like systems, this simply starts the process normally.
func RunAsUser(path string) (*os.Process, error) {
	return startProcess(StartOptions{
		Path: path,
		Args: []string{path},
	})
}
