#!/usr/bin/env bash
#/ usage: build.sh [VERSION]

VERSION="${1:-dev}"

build (){
	os="${1}"
	arch="${2}"
	bin="main-${os}-${arch}-${VERSION}"

	echo "building ${bin}..."
	CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" \
		go build \
			-ldflags="-s -w" \
			-o "dist/${bin}" \
			./cmd/action-ln/
}

build linux arm64
build linux amd64
