package transport

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Piya-Boy/devsync/config"
)

// SSHTransport syncs using rsync over SSH.
type SSHTransport struct{}

func (SSHTransport) Name() string { return "ssh" }

func (SSHTransport) Sync(server *config.Server, localPath, remotePath string, opts Options) error {
	args := buildRsyncArgs(server, localPath, remotePath, opts.DryRun)
	cmd := exec.Command("rsync", args...)
	cmd.Stdout = newPrefixWriter("[rsync]")
	cmd.Stderr = newPrefixWriter("[rsync]")

	if opts.DryRun {
		fmt.Printf("[dry-run] rsync %s\n", strings.Join(args, " "))
		return nil
	}
	return cmd.Run()
}

func buildRsyncArgs(server *config.Server, localPath, remotePath string, dryRun bool) []string {
	args := []string{"-avz", "--delete"}
	if dryRun {
		args = append(args, "--dry-run")
	}

	sshCmd := fmt.Sprintf("ssh -p %d", server.Port)
	if server.KeyPath != "" {
		sshCmd += " -i " + server.KeyPath
	}
	args = append(args, "-e", sshCmd)

	// Ensure trailing slash on local path so rsync copies contents, not directory.
	local := strings.TrimRight(localPath, `/\`) + "/"
	remote := fmt.Sprintf("%s@%s:%s", server.User, server.Host, remotePath)

	args = append(args, local, remote)
	return args
}

// portStr is a helper used in tests to verify port encoding.
func portStr(n int) string { return strconv.Itoa(n) }
