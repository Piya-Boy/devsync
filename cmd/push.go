package cmd

import (
	"fmt"
	"os"

	"github.com/Piya-Boy/devsync/config"
	"github.com/Piya-Boy/devsync/osadapter"
	syncengine "github.com/Piya-Boy/devsync/sync"
	"github.com/Piya-Boy/devsync/transport"
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

	// Auto-detect transport (SSH if port 22 open, else SMB).
	tr, err := transport.Detect(server)
	if err != nil {
		return err
	}
	fmt.Printf("Using transport: %s\n", tr.Name())

	// Validate platform support before starting.
	if err := osadapter.ValidatePlatformSupport(tr.Name()); err != nil {
		return err
	}

	if dryRun {
		fmt.Println("\n[dry-run] showing what would be synced:")
	}

	engine := syncengine.New(tr)
	return engine.SyncAll(server, cfg.Folders, syncengine.Options{DryRun: dryRun})
}
