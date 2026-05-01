param(
    [Parameter(Mandatory = $true)]
    [string]$Tag,
    [string]$Owner = "Piya-Boy",
    [string]$Repo = "devsync",
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"

if (-not $Tag.StartsWith("v")) {
    throw "Tag must start with 'v' (example: v0.1.0)"
}

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $root

if (-not $SkipBuild) {
    & "$PSScriptRoot\build.ps1"
}

$assets = @(
    "dist\devsync-windows.exe",
    "dist\devsync-linux",
    "dist\devsync-macos"
)

foreach ($asset in $assets) {
    if (-not (Test-Path $asset)) {
        throw "Missing release asset: $asset"
    }
}

git tag $Tag
git push origin $Tag

gh release create $Tag `
    --repo "$Owner/$Repo" `
    --title $Tag `
    --notes "DevSync $Tag release" `
    $assets
