package sync_test

import (
	"errors"
	"testing"

	"github.com/Piya-Boy/devsync/config"
	"github.com/Piya-Boy/devsync/sync"
	"github.com/Piya-Boy/devsync/transport"
)

// mockTransport is a test double for transport.Transport.
type mockTransport struct {
	name    string
	callLog []string
	errOn   string // folder label to fail on
}

func (m *mockTransport) Name() string { return m.name }

func (m *mockTransport) Sync(server *config.Server, localPath, remotePath string, opts transport.Options) error {
	m.callLog = append(m.callLog, localPath+"→"+remotePath)
	if m.errOn == localPath {
		return errors.New("mock sync error")
	}
	return nil
}

func TestEngine_SyncAll_AllPass(t *testing.T) {
	mock := &mockTransport{name: "ssh"}
	engine := sync.New(mock)

	server := &config.Server{Name: "prod", Host: "h", User: "u", Port: 22}
	folders := []config.Folder{
		{Name: "app", Local: "./dist", Remote: "/var/www/app"},
		{Name: "api", Local: "./api", Remote: "/var/www/api"},
	}

	if err := engine.SyncAll(server, folders, sync.Options{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(mock.callLog) != 2 {
		t.Errorf("expected 2 sync calls, got %d", len(mock.callLog))
	}
}

func TestEngine_SyncAll_PartialFailure(t *testing.T) {
	mock := &mockTransport{name: "ssh", errOn: "./dist"}
	engine := sync.New(mock)

	server := &config.Server{Name: "prod", Host: "h", User: "u", Port: 22}
	folders := []config.Folder{
		{Name: "app", Local: "./dist", Remote: "/var/www/app"},
		{Name: "api", Local: "./api", Remote: "/var/www/api"},
	}

	err := engine.SyncAll(server, folders, sync.Options{})
	if err == nil {
		t.Error("expected error from partial failure")
	}
	// Both folders should still be attempted (fail-all-then-report pattern)
	if len(mock.callLog) != 2 {
		t.Errorf("expected 2 sync calls even with failure, got %d", len(mock.callLog))
	}
}

func TestEngine_SyncAll_DryRun(t *testing.T) {
	mock := &mockTransport{name: "ssh"}
	engine := sync.New(mock)

	server := &config.Server{Name: "prod", Host: "h", User: "u", Port: 22}
	folders := []config.Folder{
		{Local: "./src", Remote: "/dest"},
	}

	if err := engine.SyncAll(server, folders, sync.Options{DryRun: true}); err != nil {
		t.Errorf("unexpected error in dry-run: %v", err)
	}
}

func TestEngine_SyncAll_RemotePathForwardSlash(t *testing.T) {
	mock := &mockTransport{name: "ssh"}
	engine := sync.New(mock)

	server := &config.Server{Name: "s", Host: "h", User: "u", Port: 22}
	folders := []config.Folder{
		{Local: `.\dist`, Remote: `C:\path\to\deploy`},
	}
	_ = engine.SyncAll(server, folders, sync.Options{})

	if len(mock.callLog) == 0 {
		t.Fatal("expected sync to be called")
	}
	if call := mock.callLog[0]; call != `.\dist→C:/path/to/deploy` {
		t.Errorf("unexpected call log: %q", call)
	}
}
