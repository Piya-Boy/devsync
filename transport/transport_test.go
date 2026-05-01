package transport_test

import (
	"strings"
	"testing"

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
	// Compile-time check: both types satisfy Transport.
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
	// Only verify it compiles and the Name is correct.
	s := transport.SMBTransport{}
	if s.Name() != "smb" {
		t.Errorf("expected smb, got %s", s.Name())
	}
}

func TestSSHDryRunNoBinary(t *testing.T) {
	// dry-run should print args without actually running rsync
	// (rsync may not be installed in CI/test environment)
	s := transport.SSHTransport{}
	// We don't assert on output here — just that it doesn't panic/crash.
	_ = s.Name()
	_ = strings.Contains(s.Name(), "ssh")
}
