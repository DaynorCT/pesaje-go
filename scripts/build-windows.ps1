# Genera build\bin\sistema-pesaje.exe (y opcionalmente instalador NSIS: -Nsis).
param(
    [switch]$Nsis
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$env:Path = "$env:USERPROFILE\go\bin;" + $env:Path

if (-not (Get-Command wails -ErrorAction SilentlyContinue)) {
    Write-Error "Wails no está en PATH. Ejecute: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
}

Set-Location $Root
$wailsArgs = @("build", "-platform", "windows/amd64")
if ($Nsis) {
    $wailsArgs += "-nsis"
}

& wails @wailsArgs
Write-Host ""
Write-Host "Listo: $Root\build\bin\sistema-pesaje.exe"
Get-ChildItem "$Root\build\bin\*-installer.exe" -ErrorAction SilentlyContinue