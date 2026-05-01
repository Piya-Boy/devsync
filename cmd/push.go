package cmd

import (
	"fmt"
	"os"

	"github.com/Piya-Boy/devsync/config"
	"github.com/spf13/cobra"
)

var (
	serverFlag string
	dryRun     bool
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push local folders to remote server",
	Long:  "Loads .devsync.json, resolves the target server and folders, then syncs them.",
	RunE:  runPush,
}

func init() {
	pushCmd.Flags().StringVarP(&serverFlag, "server", "s", "", "Target server name (default: first server in config)")
	pushCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be synced without making changes")
}

func runPush(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return err
	}

	server, err := cfg.FindServer(serverFlag)
	if err != nil {
		return err
	}

	fmt.Printf("Resolved server: %s (%s@%s:%d)\n", server.Name, server.User, server.Host, server.Port)
	fmt.Printf("Folders to sync: %d\n", len(cfg.Folders))
	for _, f := range cfg.Folders {
		label := f.Name
		if label == "" {
			label = f.Local
		}
		fmt.Printf("  [%s] %s → %s\n", label, f.Local, f.Remote)
	}

	if dryRun {
		fmt.Println("\n[dry-run] No changes made.")
		return nil
	}

	fmt.Println("\nSync not yet implemented — transport layer pending (Phase 3).")
	return nil
}
