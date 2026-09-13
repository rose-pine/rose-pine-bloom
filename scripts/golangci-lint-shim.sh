#!/bin/sh

# Copied from https://github.com/jesseduffield/lazygit/blob/master/scripts/golangci-lint-shim.sh

set -e

# Must be kept in sync with the version in .github/workflows/ci.yml
version="v2.13.2"

go run "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$version" "$@"
