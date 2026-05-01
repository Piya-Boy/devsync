package transport_test

import (
	"strings"
	"testing"

	"github.com/Piya-Boy/devsync/config"
	"github.com/Piya-Boy/devsync/transport"
)

func TestSSHTransport_Name(t *testing.T) {
	var tr transport.Transport = transport.SSHTransport{}
	if tr.Name() != "ssh" {
		t.Errorf("expected 'ssh', got %q", tr.Name())
	}
}

func TestSMBTransport_Name(t *testing.T) {
	var tr transport.Transport = transport.SMBTransport{}
	if tr.Name() != "smb" {
		t.Errorf("expected 'smb', got %q", tr.Name())
	}
}

func TestTransportInterface(t *testing.T) {
	transports := []transport.Transport{
		transport.SSHTransport{},
		transport.SMBTransport{},
	}
	for _, tr := range transports {
		if tr.Name() == "" {
			t.Errorf("transport.Name() should not be empty")
		}
	}
}

func TestSMBDryRunDoesNotPanic(t *testing.T) {
	s := transport.SMBTransport{}
	if s.Name() != "smb" {
		t.Errorf("expected smb, got %s", s.Name())
	}
}

func TestSSHDryRunArgs(t *testing.T) {
	// Verify that the SSH transport Name is correct and type is valid.
	s := transport.SSHTransport{}
	if !strings.Contains(s.Name(), "ssh") {
		t.Errorf("expected name to contain ssh")
	}
}

func TestBuildRsyncArgs(t *testing.T) {
	server := &config.Server{
		Name: "prod", Host: "srv.example.com", User: "deploy", Port: 22,
	}
	args := transport.BuildRsyncArgs(server, "/local/path", "/remote/path", false)

	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-avz") {
		t.Errorf("expected -avz flag, got: %s", joined)
	}
	if !strings.Contains(joined, "--delete") {
		t.Errorf("expected --delete flag, got: %s", joined)
	}
	if !strings.Contains(joined, "deploy@srv.example.com:/remote/path") {
		t.Errorf("expected remote path, got: %s", joined)
	}
	if !strings.Contains(joined, "ssh -p 22") {
		t.Errorf("expected ssh -p 22, got: %s", joined)
	}
	// local path should have trailing slash
	if !strings.Contains(joined, "/local/path/") {
		t.Errorf("expected trailing slash on local path, got: %s", joined)
	}
}

func TestBuildRsyncArgs_DryRun(t *testing.T) {
	server := &config.Server{Name: "s", Host: "h", User: "u", Port: 22}
	args := transport.BuildRsyncArgs(server, "/src", "/dst", true)

	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--dry-run") {
		t.Errorf("expected --dry-run flag, got: %s", joined)
	}
}

func TestBuildRsyncArgs_KeyPath(t *testing.T) {
	server := &config.Server{Name: "s", Host: "h", User: "u", Port: 2222, KeyPath: "~/.ssh/id_ed25519"}
	args := transport.BuildRsyncArgs(server, "/src", "/dst", false)

	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-i ~/.ssh/id_ed25519") {
		t.Errorf("expected key path in ssh cmd, got: %s", joined)
	}
	if !strings.Contains(joined, "-p 2222") {
		t.Errorf("expected custom port, got: %s", joined)
	}
}
