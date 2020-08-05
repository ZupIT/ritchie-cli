#!/bin/bash
GO_VERSION=$(go version)

if [[ $GO_VERSION != *"go1.14"* ]]; then
  echo "Go version is not compatible"
  exit 1
fi