#!/bin/bash

if   expr "$CIRCLE_BRANCH" : 'qa' >/dev/null; then export RELEASE_VERSION="qa-${CIRCLE_BUILD_NUM}"
if   expr "$CIRCLE_BRANCH" : 'fix/changelog' >/dev/null; then export RELEASE_VERSION="1.0.0.1"
elif expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then export RELEASE_VERSION=$(echo "$CIRCLE_BRANCH"| cut -d '-' -f 2-)
else echo ""
fi