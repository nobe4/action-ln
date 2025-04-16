#!/usr/bin/env bash
# This is an ugly hack before we can just do go build and release
# without needing to commit the binaries.

set -e

version="${1:-*}"
binary="main-linux-arm64-${version}"

# shellcheck disable=SC2086
if ! ls ./dist/${binary} 1> /dev/null 2>&1; then
	echo "${binary} is missing"
	echo "make sure you build the expected version"
	exit 1
fi

