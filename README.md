# devsync

DevSync is a cross-platform deploy tool with a Go CLI core and a minimal Tauri + React desktop app.

## Build

```powershell
go build ./...
go test ./...
npm run build
cargo check --manifest-path src-tauri\Cargo.toml
```

Build release CLI binaries:

```powershell
.\scripts\build.ps1
```

Build the Windows desktop bundle:

```powershell
npm run tauri:build
```

## Release

Create a GitHub Release and upload CLI assets:

```powershell
.\scripts\release.ps1 -Tag v0.1.0
```

The release assets are named for the self-updater:

| OS | Asset |
| --- | --- |
| Windows | `devsync-windows.exe` |
| Linux | `devsync-linux` |
| macOS | `devsync-macos` |
