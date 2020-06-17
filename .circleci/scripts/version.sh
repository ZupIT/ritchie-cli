#!/bin/bash

if expr "$CIRCLE_BRANCH" : 'qa' >/dev/null; then
  export RELEASE_VERSION="qa-${CIRCLE_BUILD_NUM}"
elif expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then
  export RELEASE_VERSION=$(echo "$CIRCLE_BRANCH" | cut -d '-' -f 2-)
elif expr "$CIRCLE_BRANCH" : '^nightly' >/dev/null; then
  export RELEASE_VERSION="nightly"
elif expr "$CIRCLE_BRANCH" : '^beta' >/dev/null; then
  export RELEASE_VERSION="$(.circleci/scripts/beta.sh)"
elif expr "$CIRCLE_BRANCH" : '^feature/packaging' >/dev/null; then
  export RELEASE_VERSION="1.0.0-test.1"
else
  export RELEASE_VERSION=$(curl https://commons-repo.ritchiecli.io/stable.txt)
  echo ""
fi