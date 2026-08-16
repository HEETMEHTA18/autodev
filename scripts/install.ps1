# AutoDev Installer for Windows
# Usage: irm https://raw.githubusercontent.com/heetmehta18/autodev/main/scripts/install.ps1 | iex
# Optional: -Version 0.6.0, -InstallDir "C:\autodev"

param(
    [string]$Version = "latest",
    [string]$InstallDir = ""
)

$ErrorActionPreference = "Stop"

$Repo = "heetmehta18/autodev"
$Binary = "autodev.exe"

function Write-Info  { Write-Host "[autodev]" -ForegroundColor Blue -NoNewline; Write-Host " $args" }
function Write-Ok   { Write-Host "[autodev]" -ForegroundColor Green -NoNewline; Write-Host " $args" }
function Write-Warn { Write-Host "[autodev]" -ForegroundColor Yellow -NoNewline; Write-Host " $args" }

# ── Default install dir: user-local, no admin needed ─────────────────────────
if (-not $InstallDir) {
    $InstallDir = Join-Path $env:LOCALAPPDATA "autodev\bin"
}
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

# ── Resolve latest version ────────────────────────────────────────────────────
if ($Version -eq "latest") {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Headers @{ "User-Agent" = "autodev-installer" }
    $Version = $release.tag_name
    Write-Info "Latest release: $Version"
} elseif ($Version -notmatch "^v") {
    $Version = "v$Version"
}

# ── Detect arch ───────────────────────────────────────────────────────────────
$arch = if ($env:PROCESSOR_ARCHITECTURE -match "ARM64") { "arm64" } else { "amd64" }

$archive = "autodev_windows_$arch.zip"
$url = "https://github.com/$Repo/releases/download/$Version/$archive"
$tmp = Join-Path $env:TEMP "autodev-$([guid]::NewGuid().ToString('N'))"
New-Item -ItemType Directory -Force -Path $tmp | Out-Null

Write-Info "Downloading $url"
Invoke-WebRequest -Uri $url -OutFile (Join-Path $tmp $archive)

# ── Extract (handles nested folder layout) ───────────────────────────────────
Write-Info "Extracting..."
Expand-Archive -Path (Join-Path $tmp $archive) -DestinationPath $tmp -Force
$binary = Get-ChildItem -Path $tmp -Recurse -Filter $Binary | Select-Object -First 1
if (-not $binary) {
    throw "Could not locate $Binary inside the downloaded archive."
}

# ── Install ───────────────────────────────────────────────────────────────────
Copy-Item -Path $binary.FullName -Destination (Join-Path $InstallDir $Binary) -Force
Remove-Item -Path $tmp -Recurse -Force -ErrorAction SilentlyContinue
Write-Ok "Installed to $InstallDir"

# ── Add to user PATH (no admin required) ─────────────────────────────────────
$cur = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not $cur) { $cur = "" }
if ($cur -notlike "*$InstallDir*") {
    $newPath = $InstallDir + ";" + $cur
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Ok "Added $InstallDir to your user PATH"
} else {
    Write-Ok "$InstallDir is already on PATH"
}

Write-Host ""
Write-Host "Getting started:" -ForegroundColor Bold
Write-Host "  autodev              - open interactive installer"
Write-Host "  autodev install nodejs - install a specific package"
Write-Host "  autodev doctor       - check environment health"
Write-Host ""
Write-Warn "Open a new terminal (or run: setx Path `"%PATH%`") for PATH changes to apply."