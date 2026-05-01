package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Piya-Boy/devsync/config"
)

func writeConfig(t *testing.T, dir string, cfg config.Config) {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".devsync.json"), data, 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestLoad_Valid(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	writeConfig(t, dir, config.Config{
		Servers: []config.Server{{Name: "prod", Host: "srv.example.com", User: "admin", Port: 22}},
		Folders: []config.Folder{{Local: "./dist", Remote: "/var/www"}},
	})

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Servers) != 1 || cfg.Servers[0].Name != "prod" {
		t.Errorf("unexpected servers: %+v", cfg.Servers)
	}
	if len(cfg.Folders) != 1 {
		t.Errorf("expected 1 folder, got %d", len(cfg.Folders))
	}
}

func TestLoad_DefaultPort(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	writeConfig(t, dir, config.Config{
		Servers: []config.Server{{Name: "s1", Host: "h", User: "u"}},
		Folders: []config.Folder{{Local: ".", Remote: "/r"}},
	})

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Servers[0].Port != 22 {
		t.Errorf("expected default port 22, got %d", cfg.Servers[0].Port)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	_ = os.WriteFile(filepath.Join(dir, ".devsync.json"), []byte("{bad json}"), 0644)

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoad_NoServers(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	writeConfig(t, dir, config.Config{
		Servers: []config.Server{},
		Folders: []config.Folder{{Local: ".", Remote: "/r"}},
	})
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected validation error for no servers")
	}
}

func TestLoad_NoFolders(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	writeConfig(t, dir, config.Config{
		Servers: []config.Server{{Name: "s", Host: "h", User: "u"}},
		Folders: []config.Folder{},
	})
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected validation error for no folders")
	}
}

func TestFindServer_ByName(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{
			{Name: "prod", Host: "h1", User: "u"},
			{Name: "staging", Host: "h2", User: "u"},
		},
		Folders: []config.Folder{{Local: ".", Remote: "/r"}},
	}
	s, err := cfg.FindServer("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Host != "h2" {
		t.Errorf("expected h2, got %s", s.Host)
	}
}

func TestFindServer_Default(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{
			{Name: "prod", Host: "h1", User: "u"},
		},
		Folders: []config.Folder{{Local: ".", Remote: "/r"}},
	}
	s, err := cfg.FindServer("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "prod" {
		t.Errorf("expected prod, got %s", s.Name)
	}
}

func TestFindServer_NotFound(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.Server{{Name: "prod", Host: "h", User: "u"}},
		Folders: []config.Folder{{Local: ".", Remote: "/r"}},
	}
	_, err := cfg.FindServer("nope")
	if err == nil {
		t.Fatal("expected error for unknown server name")
	}
}
