#!/bin/sh

mkdir -p bin
for i in $(go list ./pkg/... | grep -v vendor/); do
  go test -v -failfast -short -coverprofile=bin/cov.out $i || exit 1
  go tool cover -func=bin/cov.out
done

echo "\033[0;32m All tests run with success :D"