package osadapter_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/Piya-Boy/devsync/osadapter"
)

func TestCurrent(t *testing.T) {
	os := osadapter.Current()
	if os == osadapter.Unknown {
		t.Errorf("unexpected unknown OS for %s", runtime.GOOS)
	}
	// Must match runtime.GOOS
	if string(os) != runtime.GOOS {
		t.Errorf("osadapter.Current() = %q, want %q", os, runtime.GOOS)
	}
}

func TestNormalizePath_Unix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix path test skipped on Windows")
	}
	result := osadapter.NormalizePath("/some/path/to/file")
	if !strings.HasPrefix(result, "/") {
		t.Errorf("expected Unix path, got %q", result)
	}
}

func TestNormalizePath_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows path test skipped on non-Windows")
	}
	result := osadapter.NormalizePath("C:/Users/test")
	if !strings.Contains(result, `\`) {
		t.Errorf("expected backslashes on Windows path, got %q", result)
	}
}

func TestToForwardSlash(t *testing.T) {
	input := `C:\Users\deploy\app`
	result := osadapter.ToForwardSlash(input)
	if strings.Contains(result, `\`) {
		t.Errorf("expected no backslashes, got %q", result)
	}
	if !strings.Contains(result, "/") {
		t.Errorf("expected forward slashes, got %q", result)
	}
}

func TestIsWindows(t *testing.T) {
	want := runtime.GOOS == "windows"
	got := osadapter.IsWindows()
	if got != want {
		t.Errorf("IsWindows() = %v, want %v", got, want)
	}
}

func TestCommand_NotNil(t *testing.T) {
	cmd := osadapter.Command("echo hello")
	if cmd == nil {
		t.Error("Command() returned nil")
	}
}

func TestValidatePlatformSupport_SMB_NonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("SMB is valid on Windows")
	}
	err := osadapter.ValidatePlatformSupport("smb")
	if err == nil {
		t.Error("expected error for SMB on non-Windows")
	}
}

func TestValidatePlatformSupport_Unknown(t *testing.T) {
	// Unknown transport should not error
	err := osadapter.ValidatePlatformSupport("ftp")
	if err != nil {
		t.Errorf("unexpected error for unknown transport: %v", err)
	}
}
