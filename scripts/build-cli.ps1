param(
    [string]$OutputDir = "src-tauri\bin"
)

$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$out = Join-Path $root $OutputDir
New-Item -ItemType Directory -Path $out -Force | Out-Null

$outputName = if ($IsWindows) { "devsync.exe" } else { "devsync" }
$output = Join-Path $out $outputName

Write-Host "Building CLI -> $output"
go build -o $output .
