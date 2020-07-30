#!/bin/sh

changelog init

mkdir -p pkg-build/SPECS

echo ${RELEASE_VERSION}

mkdir -p dist/installer

go-bin-rpm generate-spec --file packaging/rpm/rpm.json -a amd64 --version ${RELEASE_VERSION} > pkg-build/SPECS/ritchiecli.spec
go-bin-rpm generate --file packaging/rpm/rpm.json -a amd64 --version ${RELEASE_VERSION} -o dist/installer/ritchie.rpm
