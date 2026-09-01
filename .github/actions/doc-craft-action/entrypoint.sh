#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$ROOT"
MODE="${1:-run}"
shift || true
go build -o bin/doc-craft ./cmd/doc-craft
exec ./bin/doc-craft "${MODE}" "$@"
