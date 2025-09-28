#!/usr/bin/env bash
set -euo pipefail
go run ./cmd/bench -rows=2300000000 -workers=16
