#!/usr/bin/env bash
# T12: pruebas automatizadas (Go + build frontend).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "== Go tests =="
go test ./... -count=1

echo ""
echo "== Frontend build =="
(cd frontend && npm run build)

echo ""
echo "OK — ver docs/T12-VERIFICACION.md para checklist visual manual."