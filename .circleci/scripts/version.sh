#!/bin/bash

if expr "$CIRCLE_BRANCH" : 'qa' >/dev/null; then
  export RELEASE_VERSION="qa-${CIRCLE_BUILD_NUM}"
elif expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then
  export RELEASE_VERSION=$(echo "$CIRCLE_BRANCH" | cut -d '-' -f 2-)
elif expr "$CIRCLE_BRANCH" : '^nightly' >/dev/null; then
  export RELEASE_VERSION="nightly"
elif expr "$CIRCLE_BRANCH" : '^beta' >/dev/null; then
  export RELEASE_VERSION="beta"
elif expr "$CIRCLE_BRANCH" : '^feature/packaging' >/dev/null; then
  export RELEASE_VERSION="2.0.0-beta.1"
else
  echo ""
fi
