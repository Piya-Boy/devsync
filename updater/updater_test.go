package updater_test

import (
	"testing"

	"github.com/Piya-Boy/devsync/updater"
)

func TestAssetCandidates(t *testing.T) {
	tests := map[string]string{
		"windows": "devsync-windows.exe",
		"linux":   "devsync-linux",
		"darwin":  "devsync-macos",
	}

	for goos, want := range tests {
		got := updater.AssetCandidates(goos)
		if len(got) == 0 || got[0] != want {
			t.Fatalf("AssetCandidates(%q)[0] = %q, want %q", goos, got, want)
		}
	}
}

func TestSelectAsset(t *testing.T) {
	release := &updater.Release{
		TagName: "v0.2.0",
		Assets: []updater.Asset{
			{Name: "devsync-linux", BrowserDownloadURL: "https://example.com/linux"},
			{Name: "devsync-windows.exe", BrowserDownloadURL: "https://example.com/windows"},
		},
	}

	asset, err := updater.SelectAsset(release, "windows")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asset.Name != "devsync-windows.exe" {
		t.Fatalf("got %q", asset.Name)
	}
}

func TestSelectAsset_FallbackName(t *testing.T) {
	release := &updater.Release{
		TagName: "v0.2.0",
		Assets:  []updater.Asset{{Name: "devsync-windows-amd64.exe"}},
	}

	asset, err := updater.SelectAsset(release, "windows")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asset.Name != "devsync-windows-amd64.exe" {
		t.Fatalf("got %q", asset.Name)
	}
}

func TestSelectAsset_Missing(t *testing.T) {
	release := &updater.Release{TagName: "v0.2.0"}
	if _, err := updater.SelectAsset(release, "linux"); err == nil {
		t.Fatal("expected missing asset error")
	}
}
