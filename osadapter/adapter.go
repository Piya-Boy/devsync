// Package osadapter detects the current OS, normalizes paths, and provides
// safe cross-platform command execution helpers.
package osadapter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// OS represents a known operating system.
type OS string

const (
	Windows OS = "windows"
	Linux   OS = "linux"
	Darwin  OS = "darwin"
	Unknown OS = "unknown"
)

// Current returns the OS the binary is running on.
func Current() OS {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "linux":
		return Linux
	case "darwin":
		return Darwin
	default:
		return Unknown
	}
}

// NormalizePath converts a path to the native format for the current OS.
// On Windows forward-slashes are converted to backslashes.
// On Unix backslashes are converted to forward-slashes.
func NormalizePath(path string) string {
	if Current() == Windows {
		return filepath.FromSlash(path)
	}
	return filepath.ToSlash(path)
}

// ToForwardSlash always returns a path with forward slashes, regardless of OS.
// This is useful when building remote (Unix) paths from local paths.
func ToForwardSlash(path string) string {
	return strings.ReplaceAll(path, `\`, "/")
}

// Command returns an exec.Cmd using the platform-appropriate shell.
// On Windows: cmd /C <command>
// On Unix:    sh -c <command>
func Command(command string) *exec.Cmd {
	if Current() == Windows {
		return exec.Command("cmd", "/C", command)
	}
	return exec.Command("sh", "-c", command)
}

// IsWindows returns true when running on Windows.
func IsWindows() bool { return Current() == Windows }

// String returns the OS name as a string.
func (o OS) String() string { return string(o) }

// ValidatePlatformSupport returns an error if the required tool is unavailable
// for the given transport name on the current OS.
func ValidatePlatformSupport(transportName string) error {
	switch transportName {
	case "smb":
		if Current() != Windows {
			return fmt.Errorf("SMB transport requires Windows (current OS: %s)", Current())
		}
	case "ssh":
		if _, err := exec.LookPath("rsync"); err != nil {
			return fmt.Errorf("rsync not found in PATH; install rsync to use SSH transport")
		}
	}
	return nil
}
