// Package sync provides the SyncEngine which orchestrates folder syncing
// using a resolved Transport.
package sync

import (
	"fmt"

	"github.com/Piya-Boy/devsync/config"
	"github.com/Piya-Boy/devsync/osadapter"
	"github.com/Piya-Boy/devsync/transport"
)

// Options controls the behaviour of a sync run.
type Options struct {
	DryRun bool
}

// Engine orchestrates syncing all configured folders to a target server.
type Engine struct {
	tr transport.Transport
}

// New returns an Engine using the provided Transport.
func New(tr transport.Transport) *Engine {
	return &Engine{tr: tr}
}

// NewAuto detects the appropriate transport for server and returns an Engine.
func NewAuto(server *config.Server) (*Engine, error) {
	tr, err := transport.Detect(server)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Detected transport: %s\n", tr.Name())
	return New(tr), nil
}

// SyncAll iterates over cfg.Folders and syncs each to server using the engine's transport.
func (e *Engine) SyncAll(server *config.Server, folders []config.Folder, opts Options) error {
	topts := transport.Options{DryRun: opts.DryRun}
	var failed []string

	for _, folder := range folders {
		label := folder.Name
		if label == "" {
			label = folder.Local
		}
		remotePath := buildRemotePath(folder.Remote)
		fmt.Printf("→ Syncing [%s]: %s → %s\n", label, folder.Local, remotePath)

		if err := e.tr.Sync(server, folder.Local, remotePath, topts); err != nil {
			fmt.Printf("  ✗ failed: %v\n", err)
			failed = append(failed, label)
		} else {
			fmt.Printf("  ✓ done\n")
		}
	}

	if len(failed) > 0 {
		return fmt.Errorf("sync failed for folders: %v", failed)
	}
	return nil
}

// buildRemotePath ensures the remote path uses forward slashes (for SSH/rsync).
func buildRemotePath(remote string) string {
	return osadapter.ToForwardSlash(remote)
}
