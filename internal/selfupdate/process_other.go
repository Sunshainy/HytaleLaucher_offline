//go:build !windows

package selfupdate

import (
	"os"
	"syscall"
)

// processExists checks if a process with the given PID exists.
// On Unix systems, this is done by sending signal 0 to the process,
// which checks for existence without actually sending a signal.
func processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Signal 0 checks if process exists without actually sending a signal
	err = process.Signal(syscall.Signal(0))
	return err == nil
}
