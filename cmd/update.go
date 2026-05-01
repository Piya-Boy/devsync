package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Piya-Boy/devsync/updater"
	"github.com/Piya-Boy/devsync/version"
	"github.com/spf13/cobra"
)

var (
	updateOwner = "Piya-Boy"
	updateRepo  = "devsync"
	checkOnly   bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates from GitHub Releases",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().StringVar(&updateOwner, "owner", updateOwner, "GitHub repository owner")
	updateCmd.Flags().StringVar(&updateRepo, "repo", updateRepo, "GitHub repository name")
	updateCmd.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	release, err := updater.FetchLatest(ctx, updateOwner, updateRepo)
	if err != nil {
		return err
	}

	if !updater.IsUpdateAvailable(release.TagName) {
		fmt.Fprintf(cmd.OutOrStdout(), "devsync is up to date (v%s)\n", version.Version)
		return nil
	}

	asset, err := updater.SelectAsset(release, runtime.GOOS)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Update available: %s (asset: %s)\n", release.TagName, asset.Name)
	if checkOnly {
		return nil
	}

	tempPath, err := updater.Download(ctx, asset)
	if err != nil {
		return err
	}

	if err := installUpdate(tempPath); err != nil {
		_ = os.Remove(tempPath)
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Update installed.")
	return nil
}

func installUpdate(tempPath string) error {
	current, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to locate current executable: %w", err)
	}

	if runtime.GOOS == "windows" {
		return installWindows(tempPath, current)
	}

	if err := os.Chmod(tempPath, 0755); err != nil {
		return err
	}
	if err := os.Rename(tempPath, current); err != nil {
		return fmt.Errorf("failed to replace %s: %w", current, err)
	}
	return exec.Command(current, "version").Start()
}

func installWindows(tempPath, current string) error {
	script, err := os.CreateTemp("", "devsync-update-*.cmd")
	if err != nil {
		return err
	}
	defer script.Close()

	content := fmt.Sprintf(`@echo off
setlocal
:wait
tasklist /FI "PID eq %d" | find "%d" >nul
if not errorlevel 1 (
  timeout /t 1 /nobreak >nul
  goto wait
)
move /Y "%s" "%s" >nul
start "" "%s" version
del "%%~f0"
`, os.Getpid(), os.Getpid(), tempPath, current, current)

	if _, err := script.WriteString(content); err != nil {
		return err
	}
	return exec.Command("cmd", "/C", "start", "", script.Name()).Start()
}
