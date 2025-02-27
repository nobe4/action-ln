#!/usr/bin/env bash

VERSION='dev'

build (){
	OS="${1}"
	ARCH="${2}"

	CGO_ENABLED=0 GOOS="${OS}" GOARCH="${ARCH}" \
		go build \
			-ldflags="-s -w" \
			-o "dist/main-${OS}-${ARCH}-${VERSION}" \
			./cmd/action-ln/
}

build linux arm64
build linux amd64
