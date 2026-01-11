// Package fork provides utilities for spawning child processes with
// various privilege levels.
package fork

import (
	"os"
)

// StartOptions contains options for starting a child process.
type StartOptions struct {
	// Path is the path to the executable.
	Path string
	// Args are the command-line arguments.
	Args []string
	// Dir is the working directory. If empty, uses the current directory.
	Dir string
	// Env is the environment variables. If nil, uses the current process's environment.
	Env []string
}

// startProcess starts a new process with the given options and file descriptors.
func startProcess(opts StartOptions) (*os.Process, error) {
	// Prepare attributes
	attr := &os.ProcAttr{
		Dir: opts.Dir,
		Env: opts.Env,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}

	// Prepend executable path to args if not already included
	args := opts.Args
	if len(args) == 0 || args[0] != opts.Path {
		args = append([]string{opts.Path}, args...)
	}

	proc, err := os.StartProcess(opts.Path, args, attr)
	if err != nil {
		return nil, err
	}

	// Release the process so it can run independently
	if err := proc.Release(); err != nil {
		return nil, err
	}

	return proc, nil
}
