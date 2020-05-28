#!/bin/sh

if expr "$VERSION" : 'qa-*' >/dev/null; then
  echo "qa"
elif expr "$VERSION" : '^[0-9]\+' >/dev/null; then
  echo "prod"
elif expr "$VERSION" : '^nightly' >/dev/null; then
  echo "nightly"
elif expr "$VERSION" : '^beta' >/dev/null; then
  echo "beta"
else
  echo ""
fi
