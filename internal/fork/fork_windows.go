//go:build windows

package fork

import (
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

// IsElevated returns true if the current process is running with elevated privileges.
// On Windows, this checks if the process has administrator privileges.
func IsElevated() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}

// RunElevated starts a process with elevated privileges using Windows UAC.
// This will trigger a UAC prompt for the user to approve.
func RunElevated(path string, args []string) (*os.Process, error) {
	verb := "runas"
	cwd, _ := os.Getwd()

	// Join args for ShellExecute
	argString := ""
	for i, arg := range args {
		if i > 0 {
			argString += " "
		}
		// Quote arguments that contain spaces
		if containsSpace(arg) {
			argString += "\"" + arg + "\""
		} else {
			argString += arg
		}
	}

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	pathPtr, _ := syscall.UTF16PtrFromString(path)
	argsPtr, _ := syscall.UTF16PtrFromString(argString)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)

	err := windows.ShellExecute(0, verbPtr, pathPtr, argsPtr, cwdPtr, windows.SW_SHOWNORMAL)
	if err != nil {
		return nil, err
	}

	// ShellExecute doesn't return a process handle, so we return nil
	return nil, nil
}

// RunAsUser starts a process as the current user.
func RunAsUser(path string) (*os.Process, error) {
	return startProcess(StartOptions{
		Path: path,
		Args: []string{path},
	})
}

// containsSpace returns true if the string contains a space.
func containsSpace(s string) bool {
	for _, c := range s {
		if c == ' ' {
			return true
		}
	}
	return false
}
