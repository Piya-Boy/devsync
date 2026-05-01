param(
    [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$out = Join-Path $root $OutputDir
New-Item -ItemType Directory -Path $out -Force | Out-Null

$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; Output = "devsync-windows.exe" },
    @{ GOOS = "linux";   GOARCH = "amd64"; Output = "devsync-linux" },
    @{ GOOS = "darwin";  GOARCH = "amd64"; Output = "devsync-macos" }
)

foreach ($target in $targets) {
    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    $output = Join-Path $out $target.Output
    Write-Host "Building $($target.GOOS)/$($target.GOARCH) -> $output"
    go build -o $output .
}

Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
