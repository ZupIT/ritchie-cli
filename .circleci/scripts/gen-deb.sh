#!/bin/sh

#apt update -y
#apt install curl -y
ls -lah
curl -fsSL https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh | GH=mh-cbon/go-bin-deb sh -xe
go-bin-deb generate --file packaging/debian/deb.json --version ${RELEASE_VERSION} -o ritchiecli.deb -a amd64
