#!/bin/sh

if expr "$VERSION" : 'qa-*' >/dev/null; then
  echo "https://ritchie-server.itiaws.dev"
elif expr "$VERSION" : '^[0-9]\+' >/dev/null; then
  echo "https://ritchie-server.zup.io"
else
  echo ""
fi
