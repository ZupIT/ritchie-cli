#!/bin/sh
mkdir -p bin
go test -v -coverprofile=bin/cov.out github.com/ZupIT/ritchie-cli/$1
go tool cover -html=bin/cov.out
