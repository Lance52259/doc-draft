#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
go mod tidy
go test ./...
mkdir -p state .work
cp -n .env.example .env 2>/dev/null || true
echo "Setup done. Edit .env then: make detect"
