#!/bin/sh

if expr "$VERSION" : 'qa-*' >/dev/null; then
  echo "ritchie-cli-bucket234376412767550"
elif expr "$VERSION" : '^[0-9]\+' >/dev/null; then
  echo "ritchie-cli-bucket152849730126474"
else
  echo ""
fi
