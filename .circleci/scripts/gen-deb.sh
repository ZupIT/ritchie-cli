#!/bin/sh


mkdir -p dist/installer

curl -fsSL https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh | GH=mh-cbon/go-bin-deb sh -xe
go-bin-deb generate --file packaging/debian/deb-single.json --version ${RELEASE_VERSION} -o dist/installer/ritchie-single.deb -a amd64
rm -rf pkg-build
go-bin-deb generate --file packaging/debian/deb-team.json --version ${RELEASE_VERSION} -o dist/installer/ritchie-team.deb -a amd64
ls -liah