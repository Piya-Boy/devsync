package transport_test

import (
	"testing"

	"github.com/Piya-Boy/devsync/config"
	"github.com/Piya-Boy/devsync/transport"
)

func TestDetect_UnreachableHost_ReturnsSMB(t *testing.T) {
	// Use a host that is guaranteed to be unreachable on port 22.
	server := &config.Server{Name: "test", Host: "192.0.2.1", User: "u", Port: 22}
	tr, err := transport.Detect(server)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Name() != "smb" {
		t.Errorf("expected smb fallback for unreachable host, got %s", tr.Name())
	}
}

func TestIsPortOpen_Localhost(t *testing.T) {
	// We can't guarantee any port is open in all test environments,
	// but we can verify that the function returns a bool without panic.
	open := transport.IsPortOpen("localhost", 9)
	_ = open // port 9 (discard) may or may not be open; just verify no panic
}
