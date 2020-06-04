#!/bin/sh

docker-compose up -d

export REPO_URL=http://localhost:8882/formulas

mkdir -p bin
go test -v -short -coverprofile=bin/cov.out $(go list ./pkg/... | grep -v vendor/)
go tool cover -func=bin/cov.out

docker-compose down