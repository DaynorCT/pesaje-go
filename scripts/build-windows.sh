#!/usr/bin/env bash
# Genera build/bin/sistema-pesaje.exe (y opcionalmente el instalador NSIS).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export PATH="${HOME}/go/bin:${PATH}"

if ! command -v wails >/dev/null 2>&1; then
  echo "Wails no está en PATH. Instale con: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
  exit 1
fi

cd "$ROOT"
ARGS=(-platform windows/amd64)
if [[ "${1:-}" == "--nsis" ]]; then
  ARGS+=(-nsis)
fi

wails build "${ARGS[@]}"
echo ""
echo "Listo: $ROOT/build/bin/sistema-pesaje.exe"
if [[ "${1:-}" == "--nsis" ]]; then
  if compgen -G "$ROOT/build/bin/*-installer.exe" >/dev/null 2>&1; then
    ls -la "$ROOT/build/bin/"*-installer.exe
  else
    echo "Nota: el instalador NSIS requiere 'makensis' en PATH (típico en Windows con NSIS instalado)."
  fi
fi