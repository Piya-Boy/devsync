package transport

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Piya-Boy/devsync/config"
)

// SMBTransport syncs using net use + robocopy (Windows only).
type SMBTransport struct{}

func (SMBTransport) Name() string { return "smb" }

func (SMBTransport) Sync(server *config.Server, localPath, remotePath string, opts Options) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("SMB transport is only supported on Windows")
	}

	share := buildUNCPath(server.Host, remotePath)

	if err := netUse(share, server.User, server.Password, opts.DryRun); err != nil {
		return fmt.Errorf("net use failed: %w", err)
	}

	driveLetter := "Z:"
	if err := robocopy(localPath, driveLetter+`\`+toWindowsPath(remotePath), opts.DryRun); err != nil {
		return fmt.Errorf("robocopy failed: %w", err)
	}

	return netUseDelete(driveLetter, opts.DryRun)
}

func buildUNCPath(host, remotePath string) string {
	// Convert /path/to/share → \\host\path\to\share
	parts := strings.Split(strings.TrimLeft(remotePath, "/"), "/")
	return `\\` + host + `\` + strings.Join(parts, `\`)
}

func toWindowsPath(p string) string {
	return strings.ReplaceAll(p, "/", `\`)
}

func netUse(share, user, password string, dryRun bool) error {
	args := []string{"use", "Z:", share}
	if user != "" {
		args = append(args, "/user:"+user)
	}
	if password != "" {
		args = append(args, password)
	}
	if dryRun {
		fmt.Printf("[dry-run] net %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.Command("net", args...)
	cmd.Stdout = newPrefixWriter("[net] ", os.Stdout)
	cmd.Stderr = newPrefixWriter("[net] ", os.Stderr)
	return cmd.Run()
}

func robocopy(src, dst string, dryRun bool) error {
	args := []string{src, dst, "/MIR"}
	if dryRun {
		args = append(args, "/L")
		fmt.Printf("[dry-run] robocopy %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.Command("robocopy", args...)
	cmd.Stdout = newPrefixWriter("[robocopy] ", os.Stdout)
	cmd.Stderr = newPrefixWriter("[robocopy] ", os.Stderr)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// robocopy exit codes 0–7 are success/info; >=8 are errors
			if exitErr.ExitCode() < 8 {
				return nil
			}
		}
		return err
	}
	return nil
}

func netUseDelete(driveLetter string, dryRun bool) error {
	args := []string{"use", driveLetter, "/delete"}
	if dryRun {
		fmt.Printf("[dry-run] net %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.Command("net", args...)
	cmd.Stdout = newPrefixWriter("[net] ", os.Stdout)
	cmd.Stderr = newPrefixWriter("[net] ", os.Stderr)
	return cmd.Run()
}
