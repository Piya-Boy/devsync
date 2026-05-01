// Package transport defines the Transport interface and related types for
// syncing folders from a local path to a remote server.
package transport

import "github.com/Piya-Boy/devsync/config"

// Options holds per-sync options shared by all transports.
type Options struct {
	DryRun bool
}

// Transport is the interface every sync backend must satisfy.
type Transport interface {
	// Name returns the human-readable name of this transport ("ssh", "smb", …).
	Name() string

	// Sync copies the contents of localPath to the remote directory remotePath
	// on the server described by server.
	Sync(server *config.Server, localPath, remotePath string, opts Options) error
}
