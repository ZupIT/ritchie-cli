#!/bin/sh

curl -fsSL https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh | GH=mh-cbon/go-bin-deb sh -xe
go-bin-deb generate --file packaging/debian/deb-single.json --version ${RELEASE_VERSION} -o dist/ritchie-single.deb -a amd64
go-bin-deb generate --file packaging/debian/deb-team.json --version ${RELEASE_VERSION} -o dist/ritchie-team.deb -a amd64