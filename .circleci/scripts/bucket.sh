#!/bin/sh

if expr "$VERSION" : '^[0-9]\+.0.0-qa' >/dev/null; then
  echo "ritchie-13528094685555"
elif expr "$VERSION" : '^[0-9]\+.0.0-stg' >/dev/null; then
  echo "ritchie-216087623718649"
elif expr "$VERSION" : '^[0-9]\+' >/dev/null; then
  echo "ritchie-7395046262137"
elif expr "$VERSION" : '^nightly' >/dev/null; then
  echo "ritchie-7395046262137"
elif expr "$VERSION" : '^beta' >/dev/null; then
  echo "ritchie-7395046262137"
else
  echo ""
fi