#!/bin/sh

changelog init

mkdir -p pkg-build/SPECS

go-bin-rpm generate-spec --file packaging/rpm/rpm-team.json -a amd64 --version ${RELEASE_VERSION} > pkg-build/SPECS/ritchiecli.spec
go-bin-rpm generate --file packaging/rpm/rpm-team.json -a amd64 --version ${RELEASE_VERSION} -o ritchie-team.rpm

go-bin-rpm generate-spec --file packaging/rpm/rpm-single.json -a amd64 --version ${RELEASE_VERSION} > pkg-build/SPECS/ritchiecli.spec
go-bin-rpm generate --file packaging/rpm/rpm-single.json -a amd64 --version ${RELEASE_VERSION} -o ritchie-single.rpm