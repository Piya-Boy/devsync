package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Piya-Boy/devsync/version"
)

const defaultTimeout = 30 * time.Second

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

func LatestReleaseURL(owner, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
}

func FetchLatest(ctx context.Context, owner, repo string) (*Release, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, LatestReleaseURL(owner, repo), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "devsync/"+version.Version)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return nil, fmt.Errorf("GitHub release request failed: %s: %s", res.Status, strings.TrimSpace(string(body)))
	}

	var release Release
	if err := json.NewDecoder(res.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release response: %w", err)
	}
	if release.TagName == "" {
		return nil, fmt.Errorf("latest release response did not include tag_name")
	}
	return &release, nil
}

func IsUpdateAvailable(tagName string) bool {
	latest := strings.TrimPrefix(strings.TrimSpace(tagName), "v")
	return latest != "" && latest != version.Version
}

func AssetCandidates(goos string) []string {
	switch goos {
	case "windows":
		return []string{"devsync-windows.exe", "devsync-windows-amd64.exe", "devsync.exe"}
	case "linux":
		return []string{"devsync-linux", "devsync-linux-amd64", "devsync"}
	case "darwin":
		return []string{"devsync-macos", "devsync-darwin-amd64", "devsync"}
	default:
		return nil
	}
}

func SelectAsset(release *Release, goos string) (*Asset, error) {
	if release == nil {
		return nil, fmt.Errorf("release is nil")
	}

	candidates := AssetCandidates(goos)
	for _, candidate := range candidates {
		for i := range release.Assets {
			if release.Assets[i].Name == candidate {
				return &release.Assets[i], nil
			}
		}
	}

	return nil, fmt.Errorf("no %s asset found in release %s", goos, release.TagName)
}

func Download(ctx context.Context, asset *Asset) (string, error) {
	if asset == nil {
		return "", fmt.Errorf("asset is nil")
	}
	if asset.BrowserDownloadURL == "" {
		return "", fmt.Errorf("asset %q has no download URL", asset.Name)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.BrowserDownloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "devsync/"+version.Version)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", res.Status)
	}

	ext := filepath.Ext(asset.Name)
	tmp, err := os.CreateTemp("", "devsync-update-*"+ext)
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, res.Body); err != nil {
		_ = os.Remove(tmp.Name())
		return "", err
	}
	if err := tmp.Chmod(0755); err != nil {
		_ = os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}

func CurrentAsset() string {
	candidates := AssetCandidates(runtime.GOOS)
	if len(candidates) == 0 {
		return ""
	}
	return candidates[0]
}
