#!/bin/sh

docker-compose up -d

export REPO_URL=http://localhost:8882/formulas

mkdir -p bin
for i in $(go list ./pkg/... | grep -v vendor/); do
  go test -v -failfast -short -coverprofile=bin/cov.out $i || exit 1
done
go tool cover -func=bin/cov.out

docker-compose down