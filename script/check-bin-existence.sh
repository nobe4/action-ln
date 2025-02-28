#!/usr/bin/env bash

version="${1:-dev}"

# This is an ugly hack before we can just do go build and release
# without needing to commit the binaries.
binary="main-linux-arm64-${version}"
if [ ! -f "./dist/${binary}" ]; then
	echo "expected ${binary} to exist, but is missing"
	echo "commit it before making a release"
	exit 1
fi

