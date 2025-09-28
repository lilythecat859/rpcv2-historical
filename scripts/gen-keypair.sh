#!/usr/bin/env bash
set -euo pipefail
go install github.com/rpcv2-historical/cmd/keygen@latest
keygen -alg=ed25519 -stdout | base64 -w0
